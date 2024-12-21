package repository

import (
	"database/sql"
	"log"

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
		returning id, text, user_id, created_at;
	`
	var post models.ReadPostDTO
	err := r.db.QueryRow(
		query, dto.Text, dto.UserID).Scan(&post.ID, &post.Text, &post.UserID, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepositoryImpl) GetAllPosts(limit, offset uint64) ([]models.ReadPostDTO, error) {
	query := `
		select id, text, user_id, repost_of_id, created_at from posts 
		where deleted_at is null offset $1 limit $2;
	`
	var readDTOs []models.ReadPostDTO = make([]models.ReadPostDTO, 0)
	rows, err := r.db.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var postDTO models.ReadPostDTO
		err := rows.Scan(&postDTO.ID, &postDTO.Text, &postDTO.UserID, &postDTO.RepostOfID, &postDTO.CreatedAt)
		if err != nil {
			return nil, err
		}
		readDTOs = append(readDTOs, postDTO)
	}
	return readDTOs, nil
}

func (r *PostRepositoryImpl) GetPostByID(id uint64) (*models.ReadPostDTO, error) {
	query := `
		select 
			id, text, user_id, repost_of_id, created_at
		from posts where id = $1 and deleted_at is null;
	`
	var postDTO models.ReadPostDTO
	err := r.db.QueryRow(query, id).Scan(&postDTO.ID, &postDTO.Text, &postDTO.UserID, &postDTO.RepostOfID, &postDTO.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &postDTO, nil
}

func (r *PostRepositoryImpl) DeletePost(id uint64) error {
	query := `
        update posts set deleted_at = now() where id = $1 and deleted_at is null;
    `
	result, err := r.db.Exec(query, id)
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
