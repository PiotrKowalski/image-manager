package main

import (
	"fmt"
	"github.com/PiotrKowalski/image-manager/internal/app/downloader/ports"
	"github.com/PiotrKowalski/image-manager/pkg/config"
	"github.com/PiotrKowalski/image-manager/pkg/config/remote"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	configProvider := remote.NewConsulRemoteConfigProvider()
	err := configProvider.LoadStoreConfig()
	if err != nil {
		log.Print(err)
		return
	}
	err = configProvider.StartRemoteWatch()
	if err != nil {
		return
	}

	appPort, err := config.ReadEnvString("APP_PORT")
	if err != nil {
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", appPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	err = ports.NewGRPCServer(s)
	if err != nil {
		return
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
