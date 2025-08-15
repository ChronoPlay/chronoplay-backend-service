package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func getCloudinaryURL() string {
	return fmt.Sprintf("cloudinary://%s:%s@%s", os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"), os.Getenv("CLOUDINARY_CLOUD_NAME"))
}

func UploadImageToCloudinary(ctx context.Context, image *multipart.FileHeader) (string, *helpers.CustomError) {
	src, err := image.Open()
	if err != nil {
		return "", helpers.System("Failed to open image file: " + err.Error())
	}
	defer src.Close()
	cld, err := cloudinary.NewFromURL(getCloudinaryURL())
	if err != nil {
		return "", helpers.System("Failed to create Cloudinary client: " + err.Error())
	}
	if cld == nil {
		return "", helpers.System("Cloudinary client is nil")
	}
	uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID: "cards/" + image.Filename,
	})
	if err != nil {
		return "", helpers.System("Failed to upload image to Cloudinary: " + err.Error())
	}
	return uploadResult.SecureURL, nil
}
