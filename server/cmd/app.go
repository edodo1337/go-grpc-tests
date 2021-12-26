package main

import (
	"fmt"
	"grpc-server/pkg/pb"
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const StorageCapacity = 1000

type App struct {
	serviceServer *pb.StorageService
	grpcServer    *grpc.Server
	Port          int
	logger        *logrus.Logger
}

func NewApp() *App {
	logger := logrus.New()

	return &App{
		serviceServer: pb.NewStorageService(logger, StorageCapacity),
		Port:          9000,
		logger:        logger,
	}
}

func (app *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", app.Port))
	if err != nil {
		app.logger.Errorf("failed to listen: %v", err)

		return err
	}

	defer lis.Close()

	app.grpcServer = grpc.NewServer()

	pb.RegisterKVStorageServiceServer(app.grpcServer, app.serviceServer)

	log.Printf("server listening at %v", lis.Addr())

	if err := app.grpcServer.Serve(lis); err != nil {
		app.logger.Errorf("failed to serve: %v", err)

		return err
	}

	return nil
}

func (app *App) Shutdown() {
	app.grpcServer.GracefulStop()
	fmt.Println("Server shutdown")
}
