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

func (s *PostServiceImpl) GetAllPosts(limit, offset uint64) ([]models.ReadPostDTO, error) {
	return s.repo.GetAllPosts(limit, offset)
}

func (s *PostServiceImpl) GetPostByID(id uint64) (*models.ReadPostDTO, error) {
	return s.repo.GetPostByID(id)
}

func (s *PostServiceImpl) CreatePost(post models.CreatePostDTO) (*models.ReadPostDTO, error) {
	return s.repo.CreatePost(post)
}

func (s *PostServiceImpl) DeletePost(id uint64) error {
	return s.repo.DeletePost(id)
}
