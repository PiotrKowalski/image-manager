package app

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/dapr/go-sdk/client"
	"time"
)

type App struct {
	daprClient client.Client
	cache      *bigcache.BigCache
}

func NewApp() (*App, error) {
	var err error

	daprClient, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		return nil, err
	}

	return &App{
		daprClient: daprClient,
		cache:      cache,
	}, nil
}

func (a *App) DownloadImage(ctx context.Context, url string) (string, error) {

	image, err := downloadFile(url)
	if err != nil {
		return "", err
	}

	err = a.cache.Set(url, image)
	if err != nil {
		return "", err
	}

	return url, err
}

func (a *App) GetImage(ctx context.Context, url string) ([]byte, error) {
	image, err := a.cache.Get(url)
	if err != nil {
		return nil, err
	}

	return image, nil
	//reader := bytes.NewReader(image)
	//
	//// read only 4 byte from our io.Reader
	//buf := make([]byte, 10000)
	//
	//fileChan := make(chan []byte, 1)
	//go func() {
	//	defer close(fileChan)
	//	i := 0
	//	for {
	//
	//		n, err := reader.Read(buf)
	//		if err != nil {
	//			break
	//		}
	//		log.Printf("Loop %d buf length %d %d", i, n, reader.Len())
	//		fileChan <- buf
	//
	//		i += 1
	//
	//		if n == 0 {
	//			break
	//		}
	//	}
	//
	//}()
	//
	//return fileChan, nil
}
