package commodities

import (
	"crop_connect/business/users"
	"crop_connect/constant"
	"crop_connect/helper"
	"crop_connect/helper/cloudinary"
	"errors"
	"mime/multipart"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommodityUseCase struct {
	commoditiesRepository Repository
	userRepository        users.Repository
	cloudinary            cloudinary.Function
}

func NewUseCase(cr Repository, ur users.Repository, cldry cloudinary.Function) UseCase {
	return &CommodityUseCase{
		commoditiesRepository: cr,
		userRepository:        ur,
		cloudinary:            cldry,
	}
}

/*
Create
*/

func (cu *CommodityUseCase) Create(domain *Domain, images []*multipart.FileHeader) (int, error) {
	_, err := cu.commoditiesRepository.GetByNameAndFarmerID(domain.Name, domain.FarmerID)
	if err == mongo.ErrNoDocuments {
		domain.ImageURLs = []string{}
		if len(images) != 0 {
			for _, image := range images {
				cloudinaryURL, err := cu.cloudinary.UploadOneWithGeneratedFilename(constant.CloudinaryFolderCommodities, image)
				if err != nil {
					return http.StatusInternalServerError, errors.New("gagal mengunggah gambar")
				}

				domain.ImageURLs = append(domain.ImageURLs, cloudinaryURL)
			}
		}

		domain.ID = primitive.NewObjectID()
		domain.Code = primitive.NewObjectID()
		domain.IsAvailable = true
		domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = cu.commoditiesRepository.Create(domain)
		if err != nil {
			err = cu.cloudinary.DeleteManyByURL(constant.CloudinaryFolderCommodities, domain.ImageURLs)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			return http.StatusInternalServerError, errors.New("gagal membuat komoditas")
		}

		return http.StatusCreated, nil
	}

	return http.StatusConflict, errors.New("komoditas sudah ada")
}

/*
Read
*/

func (cu *CommodityUseCase) GetByPaginationAndQuery(query Query) ([]Domain, int, int, error) {
	commodities, totalData, err := cu.commoditiesRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, err
	}

	return commodities, totalData, http.StatusOK, nil
}

func (cu *CommodityUseCase) GetByID(id primitive.ObjectID) (Domain, int, error) {
	commodity, err := cu.commoditiesRepository.GetByID(id)
	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	}

	return commodity, http.StatusOK, nil
}

func (cu *CommodityUseCase) GetByIDWithoutDeleted(id primitive.ObjectID) (Domain, int, error) {
	commodity, err := cu.commoditiesRepository.GetByIDWithoutDeleted(id)
	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	}

	return commodity, http.StatusOK, nil
}

func (cu *CommodityUseCase) GetByFarmerID(farmerID primitive.ObjectID) ([]Domain, int, error) {
	commodities, err := cu.commoditiesRepository.GetByFarmerID(farmerID)
	if err != nil {
		return []Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
	}

	return commodities, http.StatusOK, nil
}

func (cu *CommodityUseCase) CountTotalCommodity(year int) (int, int, error) {
	totalCommodity, err := cu.commoditiesRepository.CountTotalCommodity(year)
	if err != nil {
		return 0, http.StatusInternalServerError, errors.New("gagal mendapatkan total komoditas")
	}

	return totalCommodity, http.StatusOK, nil
}

func (cu *CommodityUseCase) CountTotalCommodityByFarmer(farmerID primitive.ObjectID) (int, int, error) {
	_, err := cu.userRepository.GetByID(farmerID)
	if err == mongo.ErrNoDocuments {
		return 0, http.StatusNotFound, errors.New("petani tidak ditemukan")
	} else if err != nil {
		return 0, http.StatusInternalServerError, errors.New("gagal mendapatkan petani")
	}

	totalCommodity, err := cu.commoditiesRepository.CountTotalCommodityByFarmer(farmerID)
	if err != nil {
		return 0, http.StatusInternalServerError, errors.New("gagal mendapatkan total komoditas")
	}

	return totalCommodity, http.StatusOK, nil
}

func (cu *CommodityUseCase) GetPerennialsByFarmerID(farmerID primitive.ObjectID) ([]Domain, int, error) {
	commodities, err := cu.commoditiesRepository.GetPerennialsByFarmerID(farmerID)
	if err != nil {
		return []Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
	}

	return commodities, http.StatusOK, nil
}

/*
Update
*/

func (cu *CommodityUseCase) Update(domain *Domain, updateImage []*helper.UpdateImage) (Domain, int, error) {
	commodity, err := cu.commoditiesRepository.GetByIDAndFarmerID(domain.ID, domain.FarmerID)
	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	}

	if commodity.Name != domain.Name {
		_, err = cu.commoditiesRepository.GetByNameAndFarmerID(domain.Name, domain.FarmerID)
		if err != mongo.ErrNoDocuments {
			return Domain{}, http.StatusConflict, errors.New("nama komoditas telah terdaftar")
		}
	}

	imageURLs, err := cu.cloudinary.UpdateArrayImage(constant.CloudinaryFolderCommodities, commodity.ImageURLs, updateImage)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mengupdate gambar")
	}

	domain.ImageURLs = imageURLs

	err = cu.commoditiesRepository.Delete(domain.ID)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal menghapus komoditas")
	}

	domain.ID = primitive.NewObjectID()
	domain.Code = commodity.Code
	domain.IsPerennials = commodity.IsPerennials
	domain.CreatedAt = commodity.CreatedAt
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	commodity, err = cu.commoditiesRepository.Create(domain)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mengupdate komoditas")
	}

	return commodity, http.StatusOK, nil
}

func (cu *CommodityUseCase) GetByIDAndFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (Domain, int, error) {
	commodity, err := cu.commoditiesRepository.GetByIDAndFarmerID(id, farmerID)
	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	}

	return commodity, http.StatusOK, nil
}

/*
Delete
*/

func (cu *CommodityUseCase) Delete(id primitive.ObjectID, farmerID primitive.ObjectID) (int, error) {
	commodity, err := cu.commoditiesRepository.GetByIDAndFarmerID(id, farmerID)
	if err != nil {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	}

	err = cu.commoditiesRepository.Delete(commodity.ID)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal menghapus komoditas")
	}

	return http.StatusOK, nil
}
