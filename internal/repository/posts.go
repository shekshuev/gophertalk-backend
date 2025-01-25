package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/models"
)

type PostRepositoryImpl struct {
	db  *sql.DB
	cfg *config.Config
	vb  *ViewBuffer
}

func NewPostRepositoryImpl(cfg *config.Config) *PostRepositoryImpl {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("Error connecting to database", err)
		return nil
	}
	viewsBufferSize := 100
	viewsBufferTimer := time.Second
	vb := &ViewBuffer{
		buffer:     make([]View, 0, viewsBufferSize),
		maxRecords: viewsBufferSize,
		timer:      viewsBufferTimer,
	}
	repository := &PostRepositoryImpl{cfg: cfg, db: db, vb: vb}
	go repository.startViewsTimer()
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
	query := `
		select 
			p.id AS post_id,
			p.text,
			p.reply_to_id,
			p.created_at,
			u.id AS user_id,
			u.user_name,
			u.first_name,
			u.last_name,
			p.likes_count,
			p.views_count,
			p.replies_count,
		    case 
		        when l.user_id is not null then true
		        else false
		    end as user_liked,
			case 
		        when v.user_id is not null then true
		        else false
		    end as user_viewed
		from posts p
		join users u ON p.user_id = u.id
		left join likes l on l.post_id = p.id and l.user_id = $1
		left join views v on v.post_id = p.id and v.user_id = $1
		where p.deleted_at is null
	`
	params := []interface{}{dto.UserID}
	if dto.Search != "" {
		query += fmt.Sprintf(" and p.text ilike $%d", len(params)+1)
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

	var readDTOs []models.ReadPostDTO = make([]models.ReadPostDTO, 0)
	rows, err := r.db.Query(query, params...)
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
			&postDTO.RepliesCount,
			&postDTO.UserLiked,
			&postDTO.UserViewed,
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
		select 
			p.id AS post_id,
			p.text,
			p.reply_to_id,
			p.created_at,
			u.id AS user_id,
			u.user_name,
			u.first_name,
			u.last_name,
			p.likes_count,
			p.views_count,
			p.replies_count,
		    case 
		        when l.user_id is not null then true
		        else false
		    end as user_liked,
			case 
		        when v.user_id is not null then true
		        else false
		    end as user_viewed
		from posts p
		join users u ON p.user_id = u.id
		left join likes l on l.post_id = p.id and l.user_id = $1
		left join views v on v.post_id = p.id and v.user_id = $1
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
		&postDTO.RepliesCount,
		&postDTO.UserLiked,
		&postDTO.UserViewed,
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
