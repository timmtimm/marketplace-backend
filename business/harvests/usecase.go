package harvests

import (
	"errors"
	"marketplace-backend/business/batchs"
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/proposals"
	"marketplace-backend/business/transactions"
	treatmentRecords "marketplace-backend/business/treatment_records"
	"marketplace-backend/constant"
	"marketplace-backend/dto"
	"marketplace-backend/helper/cloudinary"
	"marketplace-backend/util"
	"mime/multipart"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HarvestUseCase struct {
	harvestRepository         Repository
	treatmentRecordRepository treatmentRecords.Repository
	batchRepository           batchs.Repository
	transactionRepository     transactions.Repository
	proposalRepository        proposals.Repository
	commodityRepository       commodities.Repository
	cloudinary                cloudinary.Function
}

func NewHarvestUseCase(hr Repository, br batchs.Repository, trr treatmentRecords.Repository, tr transactions.Repository, pr proposals.Repository, cr commodities.Repository, cldry cloudinary.Function) UseCase {
	return &HarvestUseCase{
		harvestRepository:         hr,
		treatmentRecordRepository: trr,
		batchRepository:           br,
		transactionRepository:     tr,
		proposalRepository:        pr,
		commodityRepository:       cr,
		cloudinary:                cldry,
	}
}

/*
Create
*/

func (hu *HarvestUseCase) SubmitHarvest(domain *Domain, farmerID primitive.ObjectID, images []*multipart.FileHeader, notes []string) (Domain, int, error) {
	newestTreatmentRecord, err := hu.treatmentRecordRepository.GetNewestByBatchID(domain.BatchID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("riwayat perawatan tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan")
	}

	if newestTreatmentRecord.Date > domain.Date {
		return Domain{}, http.StatusBadRequest, errors.New("tanggal panen tidak boleh lebih awal dari tanggal perawatan terakhir")
	} else if domain.Date > primitive.NewDateTimeFromTime(time.Now()) {
		return Domain{}, http.StatusBadRequest, errors.New("tanggal panen tidak boleh lebih dari tanggal hari ini")
	}

	batch, err := hu.batchRepository.GetByID(domain.BatchID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	checkHarvest, err := hu.harvestRepository.GetByBatchID(domain.BatchID)
	if err == mongo.ErrNoDocuments {
		transaction, err := hu.transactionRepository.GetByID(batch.TransactionID)
		if err == mongo.ErrNoDocuments {
			return Domain{}, http.StatusNotFound, errors.New("transaksi tidak ditemukan")
		} else if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan transaksi")
		}

		proposal, err := hu.proposalRepository.GetByID(transaction.ProposalID)
		if err == mongo.ErrNoDocuments {
			return Domain{}, http.StatusNotFound, errors.New("proposal tidak ditemukan")
		} else if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan proposal")
		}

		commodity, err := hu.commodityRepository.GetByID(proposal.CommodityID)
		if err == mongo.ErrNoDocuments {
			return Domain{}, http.StatusNotFound, errors.New("komoditas tidak ditemukan")
		} else if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
		}

		if commodity.FarmerID != farmerID {
			return Domain{}, http.StatusForbidden, errors.New("anda tidak memiliki akses")
		}

		var imageURLs []string
		notes = util.RemoveNilStringInArray(notes)

		if len(images) != len(notes) {
			return Domain{}, http.StatusBadRequest, errors.New("jumlah gambar dan catatan tidak sama")
		}

		if len(images) > 0 {
			imageURLs, err = hu.cloudinary.UploadManyWithGeneratedFilename(constant.CloudinaryFolderHarvests, images)
			if err != nil {
				return Domain{}, http.StatusInternalServerError, errors.New("gagal mengunggah gambar")
			}

			for i := 0; i < len(imageURLs); i++ {
				domain.Harvest = append(domain.Harvest, dto.ImageAndNote{
					ImageURL: imageURLs[i],
					Note:     notes[i],
				})
			}
		} else {
			return Domain{}, http.StatusBadRequest, errors.New("gambar dan catatan tidak boleh kosong")
		}

		domain.ID = primitive.NewObjectID()
		domain.Status = constant.HarvestStatusPending
		domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = hu.harvestRepository.Create(domain)
		if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mengajukan hasi panen")
		}

		return *domain, http.StatusCreated, nil
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan hasil panen")
	}

	if checkHarvest.Status == constant.HarvestStatusPending {
		return Domain{}, http.StatusBadRequest, errors.New("hasil panen sedang dalam proses verifikasi")
	} else if checkHarvest.Status == constant.HarvestStatusApproved {
		return Domain{}, http.StatusBadRequest, errors.New("hasil panen sudah diterima")
	} else {
		return Domain{}, http.StatusBadRequest, errors.New("hasil panen sedang dalam proses revisi")
	}
}

/*
Read
*/

/*
Update
*/

/*
Delete
*/
