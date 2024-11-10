package interfaces

import (
	"awesomeProject4/user-auth-service/internal/domains/models"
	"context"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	IsUserExists(ctx context.Context, user models.User) (bool, error)
	FindisUsernewByEmail(ctx context.Context, email string) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUserStatus(ctx context.Context, userID int, status string) error
	UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) error
}
