package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgxpool"

	"apps/api/internal/models"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

var userStruct = sqlbuilder.NewStruct(new(models.User)).
	For(sqlbuilder.PostgreSQL)

func (r *UserRepo) CreateUser(
	ctx context.Context,
	userCreate models.UserCreate,
) (*models.User, error) {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("users")
	ib.Cols("email", "password_hash")
	ib.Values(userCreate.Email, userCreate.PasswordHash)
	ib.Returning(strings.Join(userStruct.Columns(), ","))
	sql, args := ib.Build()

	var user models.User
	err := r.db.QueryRow(ctx, sql, args...).Scan(userStruct.Addr(&user)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) GetUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	return r.getUserByUniqField(ctx, "email", email)
}

func (r *UserRepo) GetUserById(
	ctx context.Context,
	id string,
) (*models.User, error) {
	return r.getUserByUniqField(ctx, "id", id)
}

func (r *UserRepo) getUserByUniqField(
	ctx context.Context,
	fieldName string,
	fieldValue any,
) (*models.User, error) {
	sb := userStruct.SelectFrom("users")
	sb.Where(sb.Equal(fieldName, fieldValue))
	query, args := sb.Build()

	var user models.User
	err := r.db.QueryRow(ctx, query, args...).Scan(userStruct.Addr(&user)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user by %s: %w", fieldName, err)
	}

	return &user, nil
}
