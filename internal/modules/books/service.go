package books

import (
	"backend-challenge/internal/configs"
	"backend-challenge/internal/entities"
	"context"
	"errors"
)

var ErrBookNotFound = errors.New("book not found")

type Service interface {
	Create(ctx context.Context, req CreateBookRequest) (*BookResponse, error)
	FindAll(ctx context.Context, params FindAllParams) ([]BookResponse, int, error) // ← update
	FindByID(ctx context.Context, id string) (*BookResponse, error)
	Update(ctx context.Context, id string, req UpdateBookRequest) (*BookResponse, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
	cfg  *configs.Setting
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateBookRequest) (*BookResponse, error) {
	book := &entities.Books{
		Title:  req.Title,
		Author: req.Author,
		Year:   req.Year,
	}

	res, err := s.repo.Insert(ctx, book)
	if err != nil {
		return nil, err
	}

	return &BookResponse{
		ID: res.ID, Title: res.Title, Author: res.Author, Year: res.Year,
		CreatedAt: res.CreatedAt, UpdatedAt: res.UpdatedAt,
	}, nil
}

func (s *service) FindAll(ctx context.Context, params FindAllParams) ([]BookResponse, int, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	books, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	var res []BookResponse
	for _, b := range books {
		res = append(res, BookResponse{
			ID:        b.ID,
			Title:     b.Title,
			Author:    b.Author,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}

	return res, total, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*BookResponse, error) {
	b, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrBookNotFound
	}

	return &BookResponse{
		ID:        b.ID,
		Title:     b.Title,
		Author:    b.Author,
		Year:      b.Year, // Pastikan ini dipassing
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}, nil
}

func (s *service) Update(ctx context.Context, id string, req UpdateBookRequest) (*BookResponse, error) {
    book := &entities.Books{
        Title:  req.Title,
        Author: req.Author,
        Year:   req.Year,
    }

    res, err := s.repo.Update(ctx, id, book)
    if err != nil {
        return nil, err
    }

    return &BookResponse{
        ID:        res.ID,
        Title:     res.Title,
        Author:    res.Author,
        Year:      res.Year, // Masukkan ke response
        CreatedAt: res.CreatedAt,
        UpdatedAt: res.UpdatedAt,
    }, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
