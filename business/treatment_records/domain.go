package treatment_records

import (
	"crop_connect/dto"
	"crop_connect/helper"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID           primitive.ObjectID
	RequesterID  primitive.ObjectID
	AccepterID   primitive.ObjectID
	BatchID      primitive.ObjectID
	Number       int
	Date         primitive.DateTime
	Status       string
	Description  string
	Treatment    []dto.ImageAndNote
	RevisionNote string
	WarningNote  string
	CreatedAt    primitive.DateTime
	UpdatedAt    primitive.DateTime
}

type Query struct {
	Skip      int64
	Limit     int64
	Sort      string
	Order     int
	Farmer    string
	FarmerID  primitive.ObjectID
	Commodity string
	BatchID   primitive.ObjectID
	Batch     string
	Number    int
	Status    string
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetNewestByBatchIDAndStatus(batchID primitive.ObjectID, status string) (Domain, error)
	CountByBatchID(batchID primitive.ObjectID) (int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByBatchID(batchID primitive.ObjectID) ([]Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	CountByYear(year int) (int, error)
	StatisticByYear(year int) ([]dto.StatisticByYear, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
}

type UseCase interface {
	// Create
	RequestToFarmer(domain *Domain) (Domain, int, error)
	// Read
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	GetByBatchID(batchID primitive.ObjectID) ([]Domain, int, error)
	StatisticByYear(year int) ([]dto.StatisticByYear, int, error)
	GetByIDAndFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (Domain, int, error)
	// Update
	FillTreatmentRecord(domain *Domain, farmerID primitive.ObjectID, images []*multipart.FileHeader, notes []string) (Domain, int, error)
	UpdateTreatmentRecord(domain *Domain, farmerID primitive.ObjectID, updateImages []*helper.UpdateImage, notes []string) (Domain, int, error)
	Validate(domain *Domain, validatorID primitive.ObjectID) (Domain, int, error)
	UpdateNotes(domain *Domain) (Domain, int, error)
	CountByYear(year int) (int, int, error)
	// Delete
}
