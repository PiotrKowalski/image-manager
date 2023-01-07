package ports

import (
	"context"
	"github.com/PiotrKowalski/image-manager/internal/app/downloader/app"
	pb "github.com/PiotrKowalski/image-manager/pkg/api/gen/proto/go/downloader/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
)

type GRPCServer struct {
	pb.UnimplementedDownloaderServiceServer
	app *app.App
}

func NewGRPCServer(server *grpc.Server) error {
	var err error

	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(server, healthCheck)

	serviceApp, err := app.NewApp()
	if err != nil {
		return err
	}

	service := GRPCServer{app: serviceApp}
	pb.RegisterDownloaderServiceServer(server, service)

	return nil
}

func (s GRPCServer) DownloadImage(ctx context.Context, in *pb.DownloadImageRequest) (*pb.DownloadImageResponse, error) {

	id, err := s.app.DownloadImage(ctx, in.GetUrl())
	if err != nil {
		return nil, err
	}

	return &pb.DownloadImageResponse{ImageId: id}, nil
}

func (s GRPCServer) GetImage(req *pb.GetImageRequest, server pb.DownloaderService_GetImageServer) error {
	//TODO implement me
	ctx := context.Background()
	imageChan, err := s.app.GetImage(ctx, req.GetImageId())
	if err != nil {
		return err
	}

	imageByte := make([]byte, 10000)
	var ok bool
loop:
	for {
		select {
		case imageByte, ok = <-imageChan:
			if !ok {
				log.Println("here breaks", imageByte)
				break loop
			}

			err = server.Send(&pb.GetImageResponse{Chunk: imageByte})

			if err != nil {
				log.Println("2 here breaks", imageByte)
				return err

			}

		}
	}
	return nil
}
