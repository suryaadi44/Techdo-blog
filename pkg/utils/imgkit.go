package utils

import (
	"context"
	"log"
	"os"

	"github.com/codedius/imagekit-go"
)

func UploadImage(ctx context.Context, name string, data string, folder string) (*imagekit.UploadResponse, error) {
	opts := imagekit.Options{
		PrivateKey: os.Getenv("IMGKIT_PRIVKEY"),
		PublicKey:  os.Getenv("IMGKIT_PUBKEY"),
	}

	ik, err := imagekit.NewClient(&opts)
	if err != nil {
		log.Println("[IMGKIT] Error creating imgkit client -> error:", err)
		return nil, err
	}

	ur := imagekit.UploadRequest{
		File:              data,
		FileName:          name,
		UseUniqueFileName: true,
		Tags:              []string{},
		Folder:            folder,
		IsPrivateFile:     false,
		CustomCoordinates: "",
		ResponseFields:    nil,
	}

	upr, err := ik.Upload.ServerUpload(ctx, &ur)
	if err != nil {
		log.Println("[IMGKIT] Error uploading -> error:", err)
		return nil, err
	}

	return upr, nil
}
