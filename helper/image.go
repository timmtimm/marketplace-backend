package helper

import (
	"errors"
	"marketplace-backend/util"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UpdateImage struct {
	Image    *multipart.FileHeader
	IsChange bool
	IsDelete bool
}

func ValidateImage(image *multipart.FileHeader) (int, error) {
	if image.Size > 10*1024*1024 {
		return http.StatusRequestEntityTooLarge, errors.New("ukuran gambar maksimal 10MB")
	}

	checkImageContentType := util.CheckStringOnArray([]string{"image/jpg", "image/jpeg", "image/png"}, image.Header.Get("Content-Type"))
	if !checkImageContentType {
		return http.StatusUnsupportedMediaType, errors.New("image type not supported")
	}

	return http.StatusOK, nil
}

func GetCreateImageRequest(c echo.Context, keys []string) ([]*multipart.FileHeader, int, error) {
	images := []*multipart.FileHeader{}

	for _, key := range keys {
		image, _ := c.FormFile(key)
		if image != nil {
			if statusCode, err := ValidateImage(image); err != nil {
				return nil, statusCode, err
			}

			images = append(images, image)
		}
	}

	return images, http.StatusOK, nil
}

func GetUpdateImageRequest(c echo.Context, keys []string, imageURLs []string, isChange []bool, isDelete []bool) ([]*UpdateImage, int, error) {
	images := []*UpdateImage{}

	if len(imageURLs) != len(isChange) {
		return nil, http.StatusBadRequest, errors.New("array perubahan gambar tidak sesuai")
	} else if len(imageURLs) != len(isDelete) {
		return nil, http.StatusBadRequest, errors.New("array penghapusan gambar tidak sesuai")
	}

	for i := 0; i < len(imageURLs); i++ {
		if isChange[i] {
			image, _ := c.FormFile(keys[i])
			if image == nil {
				return nil, http.StatusBadRequest, errors.New("gambar yang ingin diperbarui kosong")
			}

			if statusCode, err := ValidateImage(image); err != nil {
				return nil, statusCode, err
			}

			images = append(images, &UpdateImage{
				Image:    image,
				IsChange: isChange[i],
				IsDelete: false,
			})
		} else {
			images = append(images, &UpdateImage{
				Image:    nil,
				IsChange: false,
				IsDelete: isDelete[i],
			})
		}
	}

	return images, http.StatusOK, nil
}
