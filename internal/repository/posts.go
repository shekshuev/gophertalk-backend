package repository

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/models"
)

type PostRepositoryImpl struct {
	db  *sql.DB
	cfg *config.Config
	vb  *ViewBuffer
	lb  *LikeBuffer
}

func NewPostRepositoryImpl(cfg *config.Config) *PostRepositoryImpl {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("Error connecting to database", err)
		return nil
	}
	viewsBufferSize := 100
	viewsBufferTimer := time.Second
	likesBufferSize := 10
	likesBufferTimer := 5 * time.Second
	vb := &ViewBuffer{
		buffer:     make([]View, 0, viewsBufferSize),
		maxRecords: viewsBufferSize,
		timer:      viewsBufferTimer,
	}
	lb := &LikeBuffer{
		likeBuffer:    make([]Like, 0, 100),
		dislikeBuffer: make([]Dislike, 0, 100),
		maxRecords:    likesBufferSize,
		timer:         likesBufferTimer,
	}
	repository := &PostRepositoryImpl{cfg: cfg, db: db, vb: vb, lb: lb}
	go repository.startViewsTimer()
	go repository.startLikesTimer()
	go repository.startDislikesTimer()
	return repository
}

func (r *PostRepositoryImpl) CreatePost(dto models.CreatePostDTO) (*models.ReadPostDTO, error) {
	query := `
		insert into posts (text, user_id, reply_to_id) values ($1, $2, $3)
		returning id, text, created_at, reply_to_id;
	`
	var post models.ReadPostDTO
	err := r.db.QueryRow(
		query, dto.Text, dto.UserID, dto.ReplyToID).Scan(&post.ID, &post.Text, &post.CreatedAt, &post.ReplyToID)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepositoryImpl) GetAllPosts(dto models.FilterPostDTO) ([]models.ReadPostDTO, error) {
	var (
		posts    []models.ReadPostDTO
		likesMap map[uint64]bool
		viewsMap map[uint64]bool
		wg       sync.WaitGroup
		mu       sync.Mutex
		errChan  = make(chan error, 3)
		doneChan = make(chan struct{})
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		p, err := r.fetchPosts(dto)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		posts = p
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		lm, err := r.fetchLikesMap(dto.UserID)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		likesMap = lm
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		vm, err := r.fetchViewsMap(dto.UserID)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		viewsMap = vm
		mu.Unlock()
	}()

	go func() {
		wg.Wait()
		close(doneChan)
	}()
	select {
	case <-doneChan:
		for i, post := range posts {
			if likesMap[post.ID] {
				posts[i].UserLiked = true
			}
			if viewsMap[post.ID] {
				posts[i].UserViewed = true
			}
		}
		return posts, nil
	case err := <-errChan:
		return nil, err
	}
}

func (r *PostRepositoryImpl) fetchPosts(dto models.FilterPostDTO) ([]models.ReadPostDTO, error) {
	query := `
		select 
			p.id as post_id,
			p.text,
			p.reply_to_id,
			p.created_at,
			u.id as user_id,
			u.user_name,
			u.first_name,
			u.last_name,
			p.likes_count,
			p.views_count,
			p.replies_count
		from posts p
		join users u on p.user_id = u.id
		where p.deleted_at is null
	`
	params := []interface{}{}
	if dto.Search != "" {
		query += " and p.text ilike $1"
		params = append(params, "%"+dto.Search+"%")
	}

	if dto.OwnerID > 0 {
		query += fmt.Sprintf(" and p.user_id = $%d", len(params)+1)
		params = append(params, dto.OwnerID)
	}

	if dto.ReplyToID > 0 {
		query += fmt.Sprintf(" and p.reply_to_id = $%d order by p.created_at asc", len(params)+1)
		params = append(params, dto.ReplyToID)
	} else {
		query += " and p.reply_to_id is null order by p.created_at desc"
	}

	query += fmt.Sprintf(" offset $%d limit $%d", len(params)+1, len(params)+2)
	params = append(params, dto.Offset, dto.Limit)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.ReadPostDTO
	for rows.Next() {
		var post models.ReadPostDTO
		var user models.ReadPostUserDTO
		err := rows.Scan(
			&post.ID,
			&post.Text,
			&post.ReplyToID,
			&post.CreatedAt,
			&user.ID,
			&user.UserName,
			&user.FirstName,
			&user.LastName,
			&post.LikesCount,
			&post.ViewsCount,
			&post.RepliesCount,
		)
		if err != nil {
			return nil, err
		}
		post.User = &user
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepositoryImpl) fetchLikesMap(userID uint64) (map[uint64]bool, error) {
	query := `
		select post_id
		from likes
		where user_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	likesMap := make(map[uint64]bool)
	for rows.Next() {
		var postID uint64
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}
		likesMap[postID] = true
	}
	return likesMap, nil
}

func (r *PostRepositoryImpl) fetchViewsMap(userID uint64) (map[uint64]bool, error) {
	query := `
		select post_id
		from views
		where user_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	viewsMap := make(map[uint64]bool)
	for rows.Next() {
		var postID uint64
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}
		viewsMap[postID] = true
	}
	return viewsMap, nil
}

func (r *PostRepositoryImpl) DeletePost(id, ownerID uint64) error {
	query := `
        update posts set deleted_at = now() where id = $1 and user_id = $2 and deleted_at is null;
    `
	result, err := r.db.Exec(query, id, ownerID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
