package transactions

import (
	"context"
	"fmt"

	"github.com/einherij/the-platform/pkg/transactions/protocol"
	"github.com/einherij/the-platform/pkg/transactions/repository"
)

type TransactionsService struct {
	transactionsRepository *repository.TransactionsRepository

	protocol.UnimplementedTransactionsServer
}

var _ = (protocol.TransactionsServer)(new(TransactionsService))

func NewService(transactionsRepository *repository.TransactionsRepository) *TransactionsService {
	return &TransactionsService{
		transactionsRepository: transactionsRepository,
	}
}

func (s *TransactionsService) Up(ctx context.Context, balanceChange *protocol.BalanceChange) (*protocol.Balance, error) {
	user, err := s.transactionsRepository.GetUserByToken(balanceChange.GetToken())
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	balance, err := s.transactionsRepository.AddBalance(user.ID, int(balanceChange.GetAmount()))
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &protocol.Balance{Amount: int32(balance.Balance)}, nil
}

func (s *TransactionsService) Down(ctx context.Context, balanceChange *protocol.BalanceChange) (*protocol.Balance, error) {
	user, err := s.transactionsRepository.GetUserByToken(balanceChange.GetToken())
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	balance, err := s.transactionsRepository.AddBalance(user.ID, int(-balanceChange.GetAmount()))
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &protocol.Balance{Amount: int32(balance.Balance)}, nil
}

func (s *TransactionsService) Transaction(ctx context.Context, transactionMessage *protocol.TransactionMessage) (*protocol.Balance, error) {
	user, err := s.transactionsRepository.GetUserByToken(transactionMessage.GetToken())
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	balance, err := s.transactionsRepository.Transfer(user.ID, int(transactionMessage.GetToUser()), int(transactionMessage.GetAmount()))
	if err != nil {
		return nil, fmt.Errorf("error transfering balance: %w", err)
	}
	return &protocol.Balance{Amount: int32(balance.Balance)}, nil
}
