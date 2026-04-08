package books

import (
	"backend-challenge/internal/entities"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Insert(ctx context.Context, book *entities.Books) (*entities.Books, error)
	FindAll(ctx context.Context, params FindAllParams) ([]*entities.Books, int, error)
	FindById(ctx context.Context, id string) (*entities.Books, error)
	Update(ctx context.Context, id string, book *entities.Books) (*entities.Books, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Insert(ctx context.Context, book *entities.Books) (*entities.Books, error) {
	var b entities.Books
	query := `  
            INSERT INTO books (title, author, year) 
            VALUES ($1, $2, $3) 
            RETURNING id, title, author, year, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, book.Title, book.Author, book.Year).
		Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.CreatedAt, &b.UpdatedAt)

	return &b, err
}

func (r *repository) FindAll(ctx context.Context, params FindAllParams) ([]*entities.Books, int, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	offset := (params.Page - 1) * params.Limit

	where := []string{"deleted_at IS NULL"}
	args := []any{}
	i := 1

	if params.Author != "" {
		where = append(where, fmt.Sprintf("author ILIKE $%d", i))
		args = append(args, "%"+params.Author+"%")
		i++
	}
	if params.Title != "" {
		where = append(where, fmt.Sprintf("title ILIKE $%d", i))
		args = append(args, "%"+params.Title+"%")
		i++
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM books %s", whereClause)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, params.Limit, offset)
	dataQuery := fmt.Sprintf(`
        SELECT id, title, author, year, created_at, updated_at
        FROM books %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, i, i+1)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var books []*entities.Books
	for rows.Next() {
		b := new(entities.Books)
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, 0, err
		}
		books = append(books, b)
	}
	return books, total, nil
}

func (r *repository) FindById(ctx context.Context, id string) (*entities.Books, error) {
	query := `
        SELECT id, title, author, year, created_at, updated_at
        FROM books
        WHERE id = $1 AND deleted_at IS NULL
    `
	b := new(entities.Books)

	// Tambahkan &b.Year di sini
	err := r.db.QueryRow(ctx, query, id).
		Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return b, nil
}

func (r *repository) Update(ctx context.Context, id string, book *entities.Books) (*entities.Books, error) {
    query := `
        UPDATE books
        SET title = $1,
            author = $2,
            year = $3, -- Tambahkan ini
            updated_at = NOW()
        WHERE id = $4 AND deleted_at IS NULL
        RETURNING id, title, author, year, created_at, updated_at
    `

    b := new(entities.Books)

    // Update Scan untuk menyertakan b.Year
    err := r.db.QueryRow(ctx, query, book.Title, book.Author, book.Year, id).
        Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.CreatedAt, &b.UpdatedAt)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrBookNotFound
        }
        return nil, err
    }

    return b, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE books
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrBookNotFound
	}

	return nil
}
