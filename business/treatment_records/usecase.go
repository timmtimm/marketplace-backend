package treatment_records

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/constant"
	"crop_connect/dto"
	"crop_connect/helper"
	"crop_connect/helper/cloudinary"
	"crop_connect/util"
	"errors"
	"mime/multipart"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TreatmentRecordUseCase struct {
	treatmentRecordRepository Repository
	batchRepository           batchs.Repository
	proposalRepository        proposals.Repository
	commodityRepository       commodities.Repository
	cloudinary                cloudinary.Function
}

func NewUseCase(trr Repository, br batchs.Repository, pr proposals.Repository, cr commodities.Repository, cldry cloudinary.Function) UseCase {
	return &TreatmentRecordUseCase{
		treatmentRecordRepository: trr,
		batchRepository:           br,
		proposalRepository:        pr,
		commodityRepository:       cr,
		cloudinary:                cldry,
	}
}

func (tru *TreatmentRecordUseCase) CheckFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (Domain, batchs.Domain, proposals.Domain, commodities.Domain, int, error) {
	treatmentRecord, err := tru.treatmentRecordRepository.GetByID(id)
	if err == mongo.ErrNoDocuments {
		return Domain{}, batchs.Domain{}, proposals.Domain{}, commodities.Domain{}, http.StatusNotFound, errors.New("riwayat perawatan tidak ditemukan")
	} else if err != nil {
		return Domain{}, batchs.Domain{}, proposals.Domain{}, commodities.Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan")
	}

	batch, err := tru.batchRepository.GetByID(treatmentRecord.BatchID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, batchs.Domain{}, proposals.Domain{}, commodities.Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
	} else if err != nil {
		return Domain{}, batchs.Domain{}, proposals.Domain{}, commodities.Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	proposal, err := tru.proposalRepository.GetByIDWithoutDeleted(batch.ProposalID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, batchs.Domain{}, proposals.Domain{}, commodities.Domain{}, http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return Domain{}, batchs.Domain{}, proposals.Domain{}, commodities.Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan proposal")
	}

	commodity, err := tru.commodityRepository.GetByIDWithoutDeleted(proposal.CommodityID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, batchs.Domain{}, proposals.Domain{}, commodities.Domain{}, http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return Domain{}, batchs.Domain{}, proposals.Domain{}, commodities.Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
	}

	if commodity.FarmerID != farmerID {
		return Domain{}, batchs.Domain{}, proposals.Domain{}, commodities.Domain{}, http.StatusForbidden, errors.New("riwayat perawatan tidak ditemukan")
	}

	return treatmentRecord, batch, proposal, commodity, http.StatusOK, nil
}

/*
Create
*/

func (tru *TreatmentRecordUseCase) RequestToFarmer(domain *Domain) (Domain, int, error) {
	batch, err := tru.batchRepository.GetByID(domain.BatchID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	if batch.Status != constant.BatchStatusPlanting {
		return Domain{}, http.StatusBadRequest, errors.New("batch tidak sedang dalam tahap tanam")
	} else if batch.CreatedAt >= domain.Date {
		return Domain{}, http.StatusBadRequest, errors.New("tanggal perawatan harus lebih besar dari tanggal tanam")
	}

	count, err := tru.treatmentRecordRepository.CountByBatchID(domain.BatchID)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan jumlah riwayat perawatan")
	}

	newestTreatmentRecord, err := tru.treatmentRecordRepository.GetNewestByBatchIDAndStatus(domain.BatchID, "")
	if err != mongo.ErrNoDocuments {
		if newestTreatmentRecord.Status != constant.TreatmentRecordStatusApproved {
			return Domain{}, http.StatusBadRequest, errors.New("riwayat perawatan terakhir belum selesai")
		} else if primitive.NewDateTimeFromTime(time.Now()) >= domain.Date {
			return Domain{}, http.StatusBadRequest, errors.New("tanggal perawatan harus lebih besar dari tanggal hari ini")
		} else if newestTreatmentRecord.Date >= domain.Date {
			return Domain{}, http.StatusBadRequest, errors.New("tanggal perawatan harus lebih besar dari tanggal perawatan terakhir")
		} else if domain.Date > batch.EstimatedHarvestDate {
			return Domain{}, http.StatusBadRequest, errors.New("tanggal perawatan harus lebih kecil dari tanggal perkiraan panen")
		}
	} else if err != mongo.ErrNoDocuments && err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan terakhir")
	}

	domain.ID = primitive.NewObjectID()
	domain.Number = count + 1
	domain.Status = constant.TreatmentRecordStatusWaitingResponse
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	treatmentRecord, err := tru.treatmentRecordRepository.Create(domain)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal membuat riwayat perawatan")
	}

	return treatmentRecord, http.StatusCreated, nil
}

/*
Read
*/

func (tru *TreatmentRecordUseCase) GetByPaginationAndQuery(query Query) ([]Domain, int, int, error) {
	treatmentRecords, totalData, err := tru.treatmentRecordRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan")
	}

	return treatmentRecords, totalData, http.StatusOK, nil
}

func (tru *TreatmentRecordUseCase) GetByBatchID(batchID primitive.ObjectID) ([]Domain, int, error) {
	_, err := tru.batchRepository.GetByID(batchID)
	if err == mongo.ErrNoDocuments {
		return []Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
	} else if err != nil {
		return []Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	treatmentRecords, err := tru.treatmentRecordRepository.GetByBatchID(batchID)
	if err != nil {
		return []Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan")
	}

	return treatmentRecords, http.StatusOK, nil
}

func (tru *TreatmentRecordUseCase) CountByYear(year int) (int, int, error) {
	statistic, err := tru.treatmentRecordRepository.CountByYear(year)
	if err != nil {
		return 0, http.StatusInternalServerError, errors.New("gagal mendapatkan statistik riwayat perawatan")
	}

	return statistic, http.StatusOK, nil
}

func (tru *TreatmentRecordUseCase) StatisticByYear(year int) ([]dto.StatisticByYear, int, error) {
	statistic, err := tru.treatmentRecordRepository.StatisticByYear(year)
	if err != nil {
		return []dto.StatisticByYear{}, http.StatusInternalServerError, errors.New("gagal mendapatkan statistik riwayat perawatan")
	}

	if len(statistic) < 12 {
		statistic = util.FillNotAvailableMonth(statistic)
	}

	return statistic, http.StatusOK, nil
}

func (tru *TreatmentRecordUseCase) GetByIDAndFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (Domain, int, error) {
	var treatmentRecord Domain
	var statusCode int
	var err error

	if farmerID != primitive.NilObjectID {
		treatmentRecord, _, _, _, statusCode, err = tru.CheckFarmerID(id, farmerID)
		if err != nil {
			return Domain{}, statusCode, err
		}
	} else {
		treatmentRecord, err = tru.treatmentRecordRepository.GetByID(id)
		if err == mongo.ErrNoDocuments {
			return Domain{}, http.StatusNotFound, errors.New("riwayat perawatan tidak ditemukan")
		} else if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan")
		}
	}

	return treatmentRecord, http.StatusOK, nil
}

/*
Update
*/

func (tru *TreatmentRecordUseCase) FillTreatmentRecord(domain *Domain, farmerID primitive.ObjectID, images []*multipart.FileHeader, notes []string) (Domain, int, error) {
	treatmentRecord, _, _, _, statusCode, err := tru.CheckFarmerID(domain.ID, farmerID)
	if err != nil {
		return Domain{}, statusCode, err
	}

	if treatmentRecord.Date > primitive.NewDateTimeFromTime(time.Now().Add(7*time.Hour)) {
		return Domain{}, http.StatusBadRequest, errors.New("riwayat perawatan belum bisa diisi")
	}

	if treatmentRecord.Status == constant.TreatmentRecordStatusApproved {
		return Domain{}, http.StatusBadRequest, errors.New("riwayat perawatan sudah diterima")
	}

	var imageURLs []string

	if len(images) > 0 {
		imageURLs, err = tru.cloudinary.UploadManyWithGeneratedFilename(constant.CloudinaryFolderTreatmentRecords, images)
		if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mengunggah gambar")
		}

		for i := 0; i < len(imageURLs); i++ {
			treatmentRecord.Treatment = append(treatmentRecord.Treatment, dto.ImageAndNote{
				ImageURL: imageURLs[i],
				Note:     notes[i],
			})
		}
	} else {
		return Domain{}, http.StatusBadRequest, errors.New("gambar dan catatan tidak boleh kosong")
	}

	treatmentRecord.Status = constant.TreatmentRecordStatusPending
	treatmentRecord.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	treatmentRecord, err = tru.treatmentRecordRepository.Update(&treatmentRecord)
	if err != nil {
		err = tru.cloudinary.DeleteManyByURL(constant.CloudinaryFolderTreatmentRecords, imageURLs)
		if err != nil {
			return Domain{}, 0, err
		}
		return Domain{}, http.StatusInternalServerError, errors.New("gagal memperbarui riwayat perawatan")
	}

	return treatmentRecord, http.StatusOK, nil
}

func (tru *TreatmentRecordUseCase) UpdateTreatmentRecord(domain *Domain, farmerID primitive.ObjectID, updateImages []*helper.UpdateImage, notes []string) (Domain, int, error) {
	treatmentRecord, _, _, _, statusCode, err := tru.CheckFarmerID(domain.ID, farmerID)
	if err != nil {
		return Domain{}, statusCode, err
	}

	if treatmentRecord.Date > primitive.NewDateTimeFromTime(time.Now().Add(7*time.Hour)) {
		return Domain{}, http.StatusBadRequest, errors.New("riwayat perawatan belum bisa diisi")
	} else if treatmentRecord.Status == constant.TreatmentRecordStatusApproved {
		return Domain{}, http.StatusBadRequest, errors.New("riwayat perawatan sudah diterima")
	}

	if len(updateImages) > 0 && len(notes) > 0 {
		imageURLs := []string{}
		for _, imageAndNote := range treatmentRecord.Treatment {
			imageURLs = append(imageURLs, imageAndNote.ImageURL)
		}

		newImageURLs, err := tru.cloudinary.UpdateArrayImage(constant.CloudinaryFolderCommodities, imageURLs, updateImages)
		if err != nil {
			return Domain{}, http.StatusInternalServerError, err
		}

		tempImageAndNotes := []dto.ImageAndNote{}
		for i := 0; i < len(newImageURLs); i++ {
			tempImageAndNotes = append(tempImageAndNotes, dto.ImageAndNote{
				ImageURL: newImageURLs[i],
				Note:     notes[i],
			})
		}

		treatmentRecord.Treatment = tempImageAndNotes
	} else {
		return Domain{}, http.StatusBadRequest, errors.New("gambar dan catatan tidak boleh kosong")
	}

	treatmentRecord.Status = constant.TreatmentRecordStatusPending
	treatmentRecord.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	treatmentRecord, err = tru.treatmentRecordRepository.Update(&treatmentRecord)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal memperbarui riwayat perawatan")
	}

	return treatmentRecord, http.StatusOK, nil
}

func (tru *TreatmentRecordUseCase) Validate(domain *Domain, validatorID primitive.ObjectID) (Domain, int, error) {
	isStatusAvailable := util.CheckStringOnArray([]string{constant.TreatmentRecordStatusRevision, constant.TreatmentRecordStatusApproved}, domain.Status)
	if !isStatusAvailable {
		return Domain{}, http.StatusBadRequest, errors.New("status proposal hanya tersedia approved dan revision")
	}

	treatmentRecord, err := tru.treatmentRecordRepository.GetByID(domain.ID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("riwayat perawatan tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan")
	}

	if treatmentRecord.Status != constant.TreatmentRecordStatusPending {
		return Domain{}, http.StatusBadRequest, errors.New("riwayat perawatan tidak dalam status menunggu validasi")
	}

	if domain.Status == constant.TreatmentRecordStatusRevision && domain.RevisionNote == "" {
		return Domain{}, http.StatusBadRequest, errors.New("catatan revisi tidak boleh kosong")
	}

	treatmentRecord.Status = domain.Status
	treatmentRecord.WarningNote = domain.WarningNote
	treatmentRecord.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	if domain.Status == constant.TreatmentRecordStatusRevision {
		treatmentRecord.RevisionNote = domain.RevisionNote
	} else {
		treatmentRecord.AccepterID = validatorID
		treatmentRecord.RevisionNote = ""
	}

	_, err = tru.treatmentRecordRepository.Update(&treatmentRecord)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal memperbarui riwayat perawatan")
	}

	return treatmentRecord, http.StatusOK, nil
}

func (tru *TreatmentRecordUseCase) UpdateNotes(domain *Domain) (Domain, int, error) {
	treatmentRecord, err := tru.treatmentRecordRepository.GetByID(domain.ID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("riwayat perawatan tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan")
	}

	if treatmentRecord.Status != constant.TreatmentRecordStatusRevision && domain.RevisionNote != "" {
		return Domain{}, http.StatusBadRequest, errors.New("riwayat perawatan tidak dalam status revisi")
	}

	treatmentRecord.RevisionNote = domain.RevisionNote
	treatmentRecord.WarningNote = domain.WarningNote
	treatmentRecord.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	treatmentRecord, err = tru.treatmentRecordRepository.Update(&treatmentRecord)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal memperbarui riwayat perawatan")
	}

	return treatmentRecord, http.StatusOK, nil
}

/*
Delete
*/
