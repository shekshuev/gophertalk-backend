package repository

import (
	"database/sql"
	"log"
	"strings"

	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/models"
)

type PostRepositoryImpl struct {
	db  *sql.DB
	cfg *config.Config
}

func NewPostRepositoryImpl(cfg *config.Config) *PostRepositoryImpl {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("Error connecting to database", err)
		return nil
	}
	repository := &PostRepositoryImpl{cfg: cfg, db: db}
	return repository
}

func (r *PostRepositoryImpl) CreatePost(dto models.CreatePostDTO) (*models.ReadPostDTO, error) {
	query := `
		insert into posts (text, user_id) values ($1, $2)
		returning id, text, created_at;
	`
	var post models.ReadPostDTO
	err := r.db.QueryRow(
		query, dto.Text, dto.UserID).Scan(&post.ID, &post.Text, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepositoryImpl) GetAllPosts(limit, offset, userID uint64) ([]models.ReadPostDTO, error) {
	query := `
		with likes_count AS (
			select post_id, count(*) as likes_count
			from likes group by post_id
		),
		views_count as (
			select post_id, count(*) AS views_count
			from views group by post_id
		)
		select 
			p.id AS post_id,
			p.text,
			p.reply_to_id,
			p.created_at,
			u.id AS user_id,
			u.user_name,
			u.first_name,
			u.last_name,
			coalesce(lc.likes_count, 0) AS likes_count,
			coalesce(vc.views_count, 0) AS views_count,
		    case 
		        when l.user_id is not null then true
		        else false
		    end as user_liked
		from posts p
		join users u ON p.user_id = u.id
		left join likes_count lc ON p.id = lc.post_id
		left join views_count vc ON p.id = vc.post_id
		left join likes l on l.post_id = p.id and l.user_id = $1
		where p.deleted_at is null
		offset $2 limit $3;
	`
	var readDTOs []models.ReadPostDTO = make([]models.ReadPostDTO, 0)
	rows, err := r.db.Query(query, userID, offset, limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var postDTO models.ReadPostDTO
		var userDTO models.ReadPostUserDTO
		err := rows.Scan(
			&postDTO.ID,
			&postDTO.Text,
			&postDTO.ReplyToID,
			&postDTO.CreatedAt,
			&userDTO.ID,
			&userDTO.UserName,
			&userDTO.FirstName,
			&userDTO.LastName,
			&postDTO.LikesCount,
			&postDTO.ViewsCount,
			&postDTO.UserLiked,
		)
		if err != nil {
			return nil, err
		}
		postDTO.User = &userDTO
		readDTOs = append(readDTOs, postDTO)
	}
	return readDTOs, nil
}

func (r *PostRepositoryImpl) GetPostByID(id, userID uint64) (*models.ReadPostDTO, error) {
	query := `
		with likes_count AS (
			select post_id, count(*) as likes_count
			from likes group by post_id
		),
		views_count as (
			select post_id, count(*) AS views_count
			from views group by post_id
		)
		select 
			p.id AS post_id,
			p.text,
			p.reply_to_id,
			p.created_at,
			u.id AS user_id,
			u.user_name,
			u.first_name,
			u.last_name,
			coalesce(lc.likes_count, 0) AS likes_count,
			coalesce(vc.views_count, 0) AS views_count,
		    case 
		        when l.user_id is not null then true
		        else false
		    end as user_liked
		from posts p
		join users u ON p.user_id = u.id
		left join likes_count lc on p.id = lc.post_id
		left join views_count vc on p.id = vc.post_id
		left join likes l on l.post_id = p.id and l.user_id = $1
		where p.id = $2 and p.deleted_at is null;
	`
	var postDTO models.ReadPostDTO
	var userDTO models.ReadPostUserDTO
	err := r.db.QueryRow(query, userID, id).Scan(
		&postDTO.ID,
		&postDTO.Text,
		&postDTO.ReplyToID,
		&postDTO.CreatedAt,
		&userDTO.ID,
		&userDTO.UserName,
		&userDTO.FirstName,
		&userDTO.LastName,
		&postDTO.LikesCount,
		&postDTO.ViewsCount,
		&postDTO.UserLiked,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	postDTO.User = &userDTO
	return &postDTO, nil
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

func (r *PostRepositoryImpl) ViewPost(id, viewedByID uint64) error {
	query := `
        insert into views (post_id, user_id) values ($1, $2);
    `
	result, err := r.db.Exec(query, id, viewedByID)
	if err != nil {
		if strings.Contains(err.Error(), "pk__views") {
			return ErrAlreadyViewed
		} else {
			return err
		}
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

func (r *PostRepositoryImpl) LikePost(id, likedByID uint64) error {
	query := `
        insert into likes (post_id, user_id) values ($1, $2);
    `
	result, err := r.db.Exec(query, id, likedByID)
	if err != nil {
		if strings.Contains(err.Error(), "pk__likes") {
			return ErrAlreadyLiked
		} else {
			return err
		}
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

func (r *PostRepositoryImpl) DislikePost(id, dislikedByID uint64) error {
	query := `
        delete from likes where post_id = $1 and user_id = $2;
    `
	result, err := r.db.Exec(query, id, dislikedByID)
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
