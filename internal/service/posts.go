package service

import (
	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/models"
	"github.com/shekshuev/gophertalk-backend/internal/repository"
)

type PostServiceImpl struct {
	repo repository.PostRepository
	cfg  *config.Config
}

func NewPostServiceImpl(repo repository.PostRepository, cfg *config.Config) *PostServiceImpl {
	return &PostServiceImpl{repo: repo, cfg: cfg}
}

func (s *PostServiceImpl) GetAllPosts(dto models.FilterPostDTO) ([]models.ReadPostDTO, error) {
	return s.repo.GetAllPosts(dto)
}

func (s *PostServiceImpl) GetPostByID(id, userID uint64) (*models.ReadPostDTO, error) {
	return s.repo.GetPostByID(id, userID)
}

func (s *PostServiceImpl) CreatePost(post models.CreatePostDTO) (*models.ReadPostDTO, error) {
	return s.repo.CreatePost(post)
}

func (s *PostServiceImpl) DeletePost(id, ownerID uint64) error {
	return s.repo.DeletePost(id, ownerID)
}

func (s *PostServiceImpl) ViewPost(id, viewedByID uint64) error {
	return s.repo.ViewPost(id, viewedByID)
}

func (s *PostServiceImpl) LikePost(id, likedByID uint64) error {
	return s.repo.LikePost(id, likedByID)
}

func (s *PostServiceImpl) DislikePost(id, dislikedByID uint64) error {
	return s.repo.DislikePost(id, dislikedByID)
}
