package repository

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"

	"github.com/einherij/the-platform/pkg/models"
)

type UsersRepository struct {
	pgClient *gorm.DB
}

func NewUsersRepository(pgClient *gorm.DB) *UsersRepository {
	return &UsersRepository{
		pgClient: pgClient,
	}
}

func (r *UsersRepository) GetUser(username string) (models.User, error) {
	var user models.User
	if err := r.pgClient.First(&user, "username=?", username).Error; err != nil {
		return user, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

func (r *UsersRepository) GetUserByToken(token string) (models.User, error) {
	var user models.User

	tokenInt, err := strconv.Atoi(token)
	if err != nil {
		return user, fmt.Errorf("bad token %s", token)
	}

	if err := r.pgClient.First(&user, "id=?", tokenInt).Error; err != nil {
		return user, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

func (r *UsersRepository) GetBalance(userID int) (models.Balance, error) {
	var balance models.Balance
	if err := r.pgClient.First(&balance, "user_id=?", userID).Error; err != nil {
		return balance, fmt.Errorf("error getting balance: %w", err)
	}
	return balance, nil
}
