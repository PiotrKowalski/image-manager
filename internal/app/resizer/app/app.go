package app

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/PiotrKowalski/image-manager/pkg/api/gen/proto/go/downloader/v1"
	"github.com/PiotrKowalski/image-manager/pkg/config"
	"github.com/allegro/bigcache/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"math"
	"os"
	"time"
)

type App struct {
	cache *bigcache.BigCache
}

func NewApp() (*App, error) {
	var err error

	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		return nil, err
	}

	return &App{
		cache: cache,
	}, nil
}

// ResizeImage
// 1. Get image from downloader
// 2. Resize it
// 3. Stream to consumer
func (a App) DownloadImageAndSetInCache(ctx context.Context, imageId string) (string, error) {
	grpcPort, err := config.ReadEnvString("DAPR_GRPC_PORT")
	if err != nil {
		return "", err
	}

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", grpcPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return "", err
	}
	c := pb.NewDownloaderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	ctx = metadata.AppendToOutgoingContext(ctx, "dapr-app-id", "downloader")
	defer cancel()
	time.Sleep(1 * time.Second)
	defer conn.Close()
	imageCl, err := c.GetImage(ctx, &pb.GetImageRequest{ImageId: imageId})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	var fullImage []byte
	log.Println(1)

	for {
		log.Println(2)

		recv, err := imageCl.Recv()
		if err != nil {
			st := status.Convert(err)
			switch st.Code() {
			case codes.NotFound:
			}
		}

		if err == io.EOF {
			log.Println(21)
			err = nil
			break

		}
		if err != nil {
			log.Println(22, err)

			return "", err
		}

		fullImage = append(fullImage, recv.GetChunk()...)
		log.Println(len(fullImage))
	}

	err = a.cache.Set(imageId, fullImage)
	if err != nil {
		return "", err
	}

	err = os.WriteFile("/tmp/dat3.jpeg", fullImage, 0644)

	return imageId, nil
}

func (a App) ResizeImage(imageId string) ([]byte, error) {

	log.Println("1here", imageId)

	data, err := a.cache.Get(imageId)
	if err != nil {
		log.Println("2here")
		return nil, err
	}

	log.Println("before", len(data))

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	resImg := resize(img, 100, 100)

	return imgToBytes(resImg), nil

}

func resize(img image.Image, length int, width int) image.Image {
	//truncate pixel size
	minX := img.Bounds().Min.X
	minY := img.Bounds().Min.Y
	maxX := img.Bounds().Max.X
	maxY := img.Bounds().Max.Y
	for (maxX-minX)%length != 0 {
		maxX--
	}
	for (maxY-minY)%width != 0 {
		maxY--
	}
	scaleX := (maxX - minX) / length
	scaleY := (maxY - minY) / width

	imgRect := image.Rect(0, 0, length, width)
	resImg := image.NewRGBA(imgRect)
	draw.Draw(resImg, resImg.Bounds(), &image.Uniform{C: color.White}, image.ZP, draw.Src)
	for y := 0; y < width; y += 1 {
		for x := 0; x < length; x += 1 {
			averageColor := getAverageColor(img, minX+x*scaleX, minX+(x+1)*scaleX, minY+y*scaleY, minY+(y+1)*scaleY)
			resImg.Set(x, y, averageColor)
		}
	}
	return resImg
}

func getAverageColor(img image.Image, minX int, maxX int, minY int, maxY int) color.Color {
	var averageRed float64
	var averageGreen float64
	var averageBlue float64
	var averageAlpha float64
	scale := 1.0 / float64((maxX-minX)*(maxY-minY))

	for i := minX; i < maxX; i++ {
		for k := minY; k < maxY; k++ {
			r, g, b, a := img.At(i, k).RGBA()
			averageRed += float64(r) * scale
			averageGreen += float64(g) * scale
			averageBlue += float64(b) * scale
			averageAlpha += float64(a) * scale
		}
	}

	averageRed = math.Sqrt(averageRed)
	averageGreen = math.Sqrt(averageGreen)
	averageBlue = math.Sqrt(averageBlue)
	averageAlpha = math.Sqrt(averageAlpha)

	averageColor := color.RGBA{
		R: uint8(averageRed),
		G: uint8(averageGreen),
		B: uint8(averageBlue),
		A: uint8(averageAlpha)}

	return averageColor
}

func imgToBytes(img image.Image) []byte {
	var opt jpeg.Options
	opt.Quality = 80

	buff := bytes.NewBuffer(nil)
	err := jpeg.Encode(buff, img, &opt)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("/tmp/dat3_after.jpeg", buff.Bytes(), 0644)

	return buff.Bytes()
}
