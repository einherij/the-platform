package main

import (
	"errors"
	"net"

	"google.golang.org/grpc"

	"github.com/einherij/enterprise"
	"github.com/einherij/enterprise/db"
	"github.com/einherij/enterprise/utils"
	"github.com/einherij/the-platform/pkg/transactions"
	transactions_protocol "github.com/einherij/the-platform/pkg/transactions/protocol"
	transactions_repository "github.com/einherij/the-platform/pkg/transactions/repository"
	"github.com/einherij/the-platform/pkg/users"
	users_protocol "github.com/einherij/the-platform/pkg/users/protocol"
	users_repository "github.com/einherij/the-platform/pkg/users/repository"
)

func main() {
	app := enterprise.NewApplication()

	pgClient := utils.Must(db.NewPostgresClient(db.PostgresConfig{
		Host:     "0.0.0.0",
		Port:     "5432",
		Username: "default_user",
		Password: "default_password",
		DBName:   "platform",
	}))

	grpcServer := grpc.NewServer()

	usersRepository := users_repository.NewUsersRepository(pgClient)
	usersService := users.NewService(usersRepository)
	users_protocol.RegisterUsersServer(grpcServer, usersService)

	transactionsRepository := transactions_repository.NewTransactionsRepository(pgClient)
	transactionsService := transactions.NewService(transactionsRepository)
	transactions_protocol.RegisterTransactionsServer(grpcServer, transactionsService)

	ln, err := net.Listen("tcp", "0.0.0.0:50051")
	utils.PanicOnError(err)
	app.RegisterOnRun(func() {
		err := grpcServer.Serve(ln)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			panic(err)
		}
	})
	app.RegisterOnShutdown(func() {
		grpcServer.GracefulStop()
	})
}
