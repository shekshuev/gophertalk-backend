package repository

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Like struct {
	UserID    uint64
	PostID    uint64
	CreatedAt time.Time
}

type Dislike struct {
	UserID uint64
	PostID uint64
}

type LikeBuffer struct {
	likeBuffer    []Like
	dislikeBuffer []Dislike
	lock          sync.Mutex
	maxRecords    int
	timer         time.Duration
}

func (r *PostRepositoryImpl) LikePost(id, likedByID uint64) error {
	r.lb.lock.Lock()
	defer r.lb.lock.Unlock()
	newLike := Like{UserID: likedByID, PostID: id, CreatedAt: time.Now()}
	r.lb.likeBuffer = append(r.lb.likeBuffer, newLike)
	if len(r.lb.likeBuffer) >= r.lb.maxRecords {
		return r.flushLikes()
	}
	return nil
}

func (r *PostRepositoryImpl) flushLikes() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createQuery := "create temp table tmp_likes (post_id bigint, user_id bigint, created_at timestamp) on commit drop;"
	if _, err = tx.Exec(createQuery); err != nil {
		tx.Rollback()
		return err
	}

	insertToTmpQuery := `insert into tmp_likes (post_id, user_id, created_at) values `
	params := []interface{}{}
	values := []string{}
	for i, like := range r.lb.likeBuffer {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		params = append(params, like.PostID, like.UserID, like.CreatedAt)
	}
	insertToTmpQuery += strings.Join(values, ",")
	if _, err = tx.Exec(insertToTmpQuery, params...); err != nil {
		tx.Rollback()
		return err
	}

	insertFromTmpQuery := `
	insert into likes (post_id, user_id, created_at)
	select post_id, user_id, min(created_at) as created_at
	from tmp_likes group by post_id, user_id on conflict (user_id, post_id) do nothing;
	`
	if _, err = tx.Exec(insertFromTmpQuery); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	r.lb.likeBuffer = r.lb.likeBuffer[:0]
	return nil
}

func (r *PostRepositoryImpl) startLikesTimer() {
	ticker := time.NewTicker(r.lb.timer)
	go func() {
		for range ticker.C {
			r.lb.lock.Lock()
			if len(r.lb.likeBuffer) > 0 {
				r.flushLikes()
			}
			r.lb.lock.Unlock()
		}
	}()
}

func (r *PostRepositoryImpl) DislikePost(id, dislikedByID uint64) error {
	r.lb.lock.Lock()
	defer r.lb.lock.Unlock()
	newDislike := Dislike{UserID: dislikedByID, PostID: id}
	r.lb.dislikeBuffer = append(r.lb.dislikeBuffer, newDislike)
	if len(r.lb.dislikeBuffer) >= r.lb.maxRecords {
		return r.flushDislikes()
	}
	return nil
}

func (r *PostRepositoryImpl) flushDislikes() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createQuery := "create temp table tmp_dislikes (post_id bigint, user_id bigint) on commit drop;"
	if _, err = tx.Exec(createQuery); err != nil {
		tx.Rollback()
		return err
	}

	insertToTmpQuery := `insert into tmp_dislikes (post_id, user_id) values `
	params := []interface{}{}
	values := []string{}
	for i, like := range r.lb.dislikeBuffer {
		values = append(values, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		params = append(params, like.PostID, like.UserID)
	}
	insertToTmpQuery += strings.Join(values, ",")
	if _, err = tx.Exec(insertToTmpQuery, params...); err != nil {
		tx.Rollback()
		return err
	}

	deleteFromLikesQuery := `
	delete from likes
        where (post_id, user_id) in (
            select post_id, user_id from tmp_dislikes
        );
	`
	if _, err = tx.Exec(deleteFromLikesQuery); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	r.lb.dislikeBuffer = r.lb.dislikeBuffer[:0]
	return nil
}

func (r *PostRepositoryImpl) startDislikesTimer() {
	ticker := time.NewTicker(r.lb.timer)
	go func() {
		for range ticker.C {
			r.lb.lock.Lock()
			if len(r.lb.dislikeBuffer) > 0 {
				r.flushDislikes()
			}
			r.lb.lock.Unlock()
		}
	}()
}
