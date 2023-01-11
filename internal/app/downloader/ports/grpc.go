package ports

import (
	"bytes"
	"context"
	"errors"
	"github.com/PiotrKowalski/image-manager/internal/app/downloader/app"
	pb "github.com/PiotrKowalski/image-manager/pkg/api/gen/proto/go/downloader/v1"
	"github.com/allegro/bigcache/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"io"
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
	ctx := context.Background()
	image, err := s.app.GetImage(ctx, req.GetImageId())
	switch {
	case errors.Is(err, bigcache.ErrEntryNotFound):
		return status.Error(codes.NotFound, "image not found")
	case err != nil:
		return err
	}

	reader := bytes.NewReader(image)
	buf := make([]byte, 10000)

	for {
		_, err := reader.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		err = server.Send(&pb.GetImageResponse{Chunk: buf})
		if err != nil {
			return err

		}

	}
	return nil
}
