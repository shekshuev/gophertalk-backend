package repository

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type View struct {
	UserID    uint64
	PostID    uint64
	CreatedAt time.Time
}

type ViewBuffer struct {
	buffer     []View
	lock       sync.Mutex
	maxRecords int
	timer      time.Duration
}

func (r *PostRepositoryImpl) ViewPost(id, viewedByID uint64) error {
	r.vb.lock.Lock()
	defer r.vb.lock.Unlock()
	newView := View{UserID: viewedByID, PostID: id, CreatedAt: time.Now()}
	r.vb.buffer = append(r.vb.buffer, newView)
	if len(r.vb.buffer) >= r.vb.maxRecords {
		return r.flushViews()
	}
	return nil
}

func (r *PostRepositoryImpl) flushViews() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createQuery := "create temp table tmp_views (post_id bigint, user_id bigint, created_at timestamp) on commit drop;"
	if _, err = tx.Exec(createQuery); err != nil {
		tx.Rollback()
		return err
	}

	insertToTmpQuery := `insert into tmp_views (post_id, user_id, created_at) values `
	params := []interface{}{}
	values := []string{}
	for i, view := range r.vb.buffer {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		params = append(params, view.PostID, view.UserID, view.CreatedAt)
	}
	insertToTmpQuery += strings.Join(values, ",")
	if _, err = tx.Exec(insertToTmpQuery, params...); err != nil {
		tx.Rollback()
		return err
	}

	insertFromTmpQuery := `
	insert into views (post_id, user_id, created_at) 
	select post_id, user_id, min(created_at) as created_at
	from tmp_views group by post_id, user_id on conflict (user_id, post_id) do nothing;
	`
	if _, err = tx.Exec(insertFromTmpQuery); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	r.vb.buffer = r.vb.buffer[:0]
	return nil
}

func (r *PostRepositoryImpl) startViewsTimer() {
	ticker := time.NewTicker(r.vb.timer)
	go func() {
		for range ticker.C {
			r.vb.lock.Lock()
			if len(r.vb.buffer) > 0 {
				r.flushViews()
			}
			r.vb.lock.Unlock()
		}
	}()
}
