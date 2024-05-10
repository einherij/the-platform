package repository

import (
	"fmt"
	"github.com/einherij/the-platform/pkg/models"
	"gorm.io/gorm"
	"strconv"
)

type TransactionsRepository struct {
	pgClient *gorm.DB
}

func NewTransactionsRepository(pgClient *gorm.DB) *TransactionsRepository {
	return &TransactionsRepository{
		pgClient: pgClient,
	}
}

func (r *TransactionsRepository) GetUserByToken(token string) (models.User, error) {
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

func (r *TransactionsRepository) AddBalance(userID int, amount int) (models.Balance, error) {
	var balance models.Balance
	if err := r.pgClient.Model(&balance).Where("user_id=?", userID).Update("balance", gorm.Expr("balance + ?", amount)).First(&balance, "user_id=?", userID).Error; err != nil {
		return balance, fmt.Errorf("error adding balance: %w", err)
	}
	return balance, nil
}

func (r *TransactionsRepository) Transfer(from, to int, amount int) (models.Balance, error) {
	var balance models.Balance
	if err := r.pgClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&balance).Where("user_id=?", from).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return fmt.Errorf("error updating source balance: %w", err)
		}
		if err := tx.Model(&balance).Where("user_id=?", to).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return fmt.Errorf("error update destination balance: %w", err)
		}
		if err := tx.First(&balance, "user_id=?", from).Error; err != nil {
			return fmt.Errorf("error getting balance: %w", err)
		}
		return nil
	}); err != nil {
		return balance, fmt.Errorf("error transfering balance: %w", err)
	}

	return balance, nil
}
