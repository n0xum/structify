package repository

import (
	"context"
	"github.com/ak/structify/internal/domain/validator"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *validator.ValidatedEntity) (*T, error)
	FindByID(ctx context.Context, id int64) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*T, error)
}

type EntityRepository interface {
	Create(ctx context.Context, entity *validator.ValidatedEntity) (interface{}, error)
	FindByID(ctx context.Context, id int64) (interface{}, error)
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]interface{}, error)
}

type QueryRepository interface {
	FindByID(ctx context.Context, id int64) (interface{}, error)
	List(ctx context.Context) ([]interface{}, error)
}

type CommandRepository interface {
	Create(ctx context.Context, entity *validator.ValidatedEntity) (interface{}, error)
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id int64) error
}
