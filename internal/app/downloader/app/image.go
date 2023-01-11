package app

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
)

func downloadFile(URL string) ([]byte, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("received non 200 response code")
	}
	//
	imageBytes := streamToByte(response.Body)
	//
	//encodedImage := make([]byte, base64.StdEncoding.EncodedLen(len(imageBytes)))
	//base64.StdEncoding.Encode(encodedImage, imageBytes)

	return imageBytes, nil
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(stream)
	if err != nil {
		return nil
	}
	log.Println("Bytes read from file: ", n)

	return buf.Bytes()
}
