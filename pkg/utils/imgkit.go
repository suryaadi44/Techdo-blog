package utils

import (
	"context"
	"log"
	"os"

	"github.com/codedius/imagekit-go"
)

func UploadImage(ctx context.Context, name string, data interface{}, folder string) (*imagekit.UploadResponse, error) {
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

func DeleteImage(ctx context.Context, id string) error {
	opts := imagekit.Options{
		PrivateKey: os.Getenv("IMGKIT_PRIVKEY"),
		PublicKey:  os.Getenv("IMGKIT_PUBKEY"),
	}

	ik, err := imagekit.NewClient(&opts)
	if err != nil {
		log.Println("[IMGKIT] Error creating imgkit client -> error:", err)
		return err
	}

	err = ik.Media.DeleteFile(ctx, id)
	if err != nil {
		log.Println("[IMGKIT] Error deleting -> error:", err)
		return err
	}

	return nil
}

func DeleteFolder(ctx context.Context, path string) error {
	opts := imagekit.Options{
		PrivateKey: os.Getenv("IMGKIT_PRIVKEY"),
		PublicKey:  os.Getenv("IMGKIT_PUBKEY"),
	}

	ik, err := imagekit.NewClient(&opts)
	if err != nil {
		log.Println("[IMGKIT] Error creating imgkit client -> error:", err)
		return err
	}

	r := imagekit.DeleteFolderRequest{
		FolderPath: path,
	}

	err = ik.Media.DeleteFolder(ctx, &r)
	if err != nil {
		log.Println("[IMGKIT] Error deleting -> error:", err)
		return err
	}

	return nil
}

func GetPictureUrl(ctx context.Context, id string) (string, error) {
	opts := imagekit.Options{
		PrivateKey: os.Getenv("IMGKIT_PRIVKEY"),
		PublicKey:  os.Getenv("IMGKIT_PUBKEY"),
	}

	ik, err := imagekit.NewClient(&opts)
	if err != nil {
		log.Println("[IMGKIT] Error creating imgkit client -> error:", err)
		return "", err
	}

	res, err := ik.Media.GetFileDetails(ctx, id)
	if err != nil {
		log.Println("[IMGKIT] Error getting file details -> error:", err)
		return "", err
	}

	return res.URL, nil
}
