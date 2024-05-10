package users

import (
	"context"
	"github.com/einherij/enterprise/db"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"

	"github.com/einherij/the-platform/pkg/users/protocol"
	"github.com/einherij/the-platform/pkg/users/repository"
)

type UsersServiceSuite struct {
	pgClient        *gorm.DB
	usersRepository *repository.UsersRepository
	usersService    *UsersService

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
	s.usersRepository = repository.NewUsersRepository(s.pgClient)
	s.usersService = NewService(s.usersRepository)
}

func (s *UsersServiceSuite) TestLogin() {
	token, err := s.usersService.Login(context.Background(), &protocol.Auth{
		Username: "test_name_1",
		Password: "",
	})
	s.NoError(err)
	s.Equal("1", token.GetToken())
}

func (s *UsersServiceSuite) Test() {
	for _, tt := range []struct {
		name    string
		token   string
		balance int32
	}{
		{
			name:    "first_user",
			token:   "1",
			balance: 100,
		},
		{
			name:    "second_user",
			token:   "2",
			balance: 200,
		},
	} {
		s.Run(tt.name, func() {
			balance, err := s.usersService.Balance(context.Background(), &protocol.Token{
				Token: tt.token,
			})
			s.NoError(err)
			s.Equal(tt.balance, balance.GetAmount())
		})
	}
}
