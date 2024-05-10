package users

import (
	"context"
	"fmt"
	"strconv"

	"github.com/einherij/the-platform/pkg/users/protocol"
	"github.com/einherij/the-platform/pkg/users/repository"
)

type UsersService struct {
	usersRepository *repository.UsersRepository

	protocol.UnimplementedUsersServer
}

var _ = (protocol.UsersServer)(new(UsersService))

func NewService(usersRepository *repository.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

func (s *UsersService) Login(ctx context.Context, auth *protocol.Auth) (*protocol.Token, error) {
	user, err := s.usersRepository.GetUser(auth.GetUsername())
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &protocol.Token{
		Token: strconv.Itoa(user.ID),
	}, nil
}

func (s *UsersService) Balance(ctx context.Context, token *protocol.Token) (*protocol.BalanceMessage, error) {
	user, err := s.usersRepository.GetUserByToken(token.GetToken())
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	balance, err := s.usersRepository.GetBalance(user.ID)
	if err != nil {
		return &protocol.BalanceMessage{
			Amount: 0,
		}, nil
	}
	return &protocol.BalanceMessage{
		Amount: int32(balance.Balance),
	}, nil
}
