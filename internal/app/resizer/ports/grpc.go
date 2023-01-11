package ports

import (
	"context"
	"github.com/PiotrKowalski/image-manager/internal/app/resizer/app"
	pb "github.com/PiotrKowalski/image-manager/pkg/api/gen/proto/go/resizer/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
)

type GRPCServer struct {
	pb.UnimplementedResizerServiceServer
	app *app.App
}

func NewGRPCServer(server *grpc.Server) error {
	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(server, healthCheck)

	serviceApp, err := app.NewApp()
	if err != nil {
		return err
	}

	service := GRPCServer{app: serviceApp}

	pb.RegisterResizerServiceServer(server, service)

	return nil
}

func (s GRPCServer) ResizeImage(in *pb.ResizeImageRequest, server pb.ResizerService_ResizeImageServer) error {
	//TODO implement me

	ctx := context.Background()

	imageId, err := s.app.DownloadImageAndSetInCache(ctx, in.GetUrl())
	if err != nil {
		return err
	}

	log.Println(imageId)

	image, err := s.app.ResizeImage(imageId)
	if err != nil {
		return err
	}

	log.Println("after", len(image), image)

	return nil

	//return &pb.ResizeImageResponse{}, nil
}

//func (s GRPCServer) ResizeImage(ctx context.Context, in *pb.ResizeImageRequest) (*pb.ResizeImageResponse, error) {
//	s.app.ResizeImage(ctx)
//
//	grpcPort, err := config.ReadEnvString("DAPR_GRPC_PORT")
//	if err != nil {
//		return nil, err
//	}
//
//	conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", grpcPort), grpc.WithInsecure(), grpc.WithBlock())
//	if err != nil {
//		log.Fatalf("did not connect: %v", err)
//	}
//	c := downloaderv1.NewDownloaderServiceClient(conn)
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
//	ctx = metadata.AppendToOutgoingContext(ctx, "dapr-app-id", "downloader")
//	defer cancel()
//	go func() {
//		time.Sleep(10 * time.Second)
//		defer conn.Close()
//		imageCl, err := c.GetImage(ctx, &pb.GetImageRequest{ImageId: "https://i.imgur.com/b062vl1.jpeg"})
//		if err != nil {
//			log.Fatalf("could not greet: %v", err)
//		}
//	loop:
//		for {
//			select {
//			case <-time.After(2 * time.Millisecond):
//				recv, err := imageCl.Recv()
//				log.Print(recv)
//				if err == io.EOF {
//					break loop
//				}
//				if err != nil {
//					log.Print(err)
//				}
//
//			}
//		}
//
//	}()
//
//	return &pb.ResizeImageResponse{}, nil
//}
