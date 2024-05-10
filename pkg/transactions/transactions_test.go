package transactions

import (
	"context"
	"github.com/einherij/enterprise/db"
	"github.com/einherij/the-platform/pkg/transactions/protocol"
	"github.com/einherij/the-platform/pkg/transactions/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type UsersServiceSuite struct {
	pgClient        *gorm.DB
	usersRepository *repository.TransactionsRepository
	usersService    *TransactionsService

	suite.Suite
}

func TestUsersService(t *testing.T) {
	suite.Run(t, new(UsersServiceSuite))
}

func (s *UsersServiceSuite) SetupTest() {
	var err error
	s.pgClient, err = db.NewPostgresClient(db.PostgresConfig{
		Host:     "0.0.0.0",
		Port:     "5432",
		Username: "default_user",
		Password: "default_password",
		DBName:   "platform",
	})
	s.NoError(err)
	s.usersRepository = repository.NewTransactionsRepository(s.pgClient)
	s.usersService = NewService(s.usersRepository)
}

func (s *UsersServiceSuite) TestUpDown() {
	balance, err := s.usersService.Up(context.Background(), &protocol.BalanceChange{
		Token:  "1",
		Amount: 100,
	})
	s.NoError(err)
	s.Equal(int32(200), balance.GetAmount())

	balance, err = s.usersService.Down(context.Background(), &protocol.BalanceChange{
		Token:  "1",
		Amount: 100,
	})
	s.Equal(int32(100), balance.GetAmount())
}

func (s *UsersServiceSuite) TestTransaction() {
	balance, err := s.usersService.Transaction(context.Background(), &protocol.TransactionMessage{
		Token:  "2",
		ToUser: 1,
		Amount: 100,
	})
	s.NoError(err)
	s.Equal(int32(100), balance.GetAmount())

	balance, err = s.usersService.Transaction(context.Background(), &protocol.TransactionMessage{
		Token:  "1",
		ToUser: 2,
		Amount: 100,
	})
	s.NoError(err)
	s.Equal(int32(100), balance.GetAmount())
}
