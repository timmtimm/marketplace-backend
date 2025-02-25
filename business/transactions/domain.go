package transactions

import (
	"crop_connect/business/commodities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID              primitive.ObjectID
	BuyerID         primitive.ObjectID
	TransactionType string
	ProposalID      primitive.ObjectID
	RegionID        primitive.ObjectID
	BatchID         primitive.ObjectID
	Address         string
	Status          string
	TotalPrice      float64
	CreatedAt       primitive.DateTime
	UpdatedAt       primitive.DateTime
}

type Statistic struct {
	Month            int
	TotalAccepted    int
	TotalTransaction int
	TotalIncome      float64
	TotalWeight      float64
	TotalUniqueBuyer int
}

type TotalTransactionByProvince struct {
	Province         string
	TotalAccepted    int
	TotalTransaction int
}

type StatisticTopCommodity struct {
	Commodity commodities.Domain
	Total     int
}

type ModelStatisticTopCommodity struct {
	CommodityCode primitive.ObjectID `bson:"_id"`
	Total         int                `bson:"total"`
}

type Query struct {
	Skip      int64
	Limit     int64
	Sort      string
	Order     int
	Commodity string
	Proposal  string
	Batch     string
	FarmerID  primitive.ObjectID
	BuyerID   primitive.ObjectID
	Status    string
	StartDate primitive.DateTime
	EndDate   primitive.DateTime
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByBuyerIDProposalIDAndStatus(buyerID primitive.ObjectID, proposalID primitive.ObjectID, status string) (Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	GetByIDAndBuyerID(id primitive.ObjectID, buyerID primitive.ObjectID) (Domain, error)
	StatisticByYear(farmerID primitive.ObjectID, year int) ([]Statistic, error)
	StatisticTopProvince(year int, limit int) ([]TotalTransactionByProvince, error)
	StatisticTopCommodity(farmerID primitive.ObjectID, year int, limit int) ([]ModelStatisticTopCommodity, error)
	CountByCommodityCode(Code primitive.ObjectID) (int, float64, error)
	GetByBuyerIDBatchIDAndStatus(buyerID primitive.ObjectID, batchID primitive.ObjectID, status string) (Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	RejectPendingByProposalID(proposalID primitive.ObjectID) error
	RejectPendingByBatchID(batchID primitive.ObjectID) error
	// Delete
}

type UseCase interface {
	// Create
	Create(domain *Domain) (int, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, int, error)
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	GetByIDAndBuyerIDOrFarmerID(id primitive.ObjectID, buyerID primitive.ObjectID, farmerID primitive.ObjectID) (Domain, int, error)
	StatisticByYear(farmerID primitive.ObjectID, year int) ([]Statistic, int, error)
	StatisticTopProvince(year int, limit int) ([]TotalTransactionByProvince, int, error)
	StatisticTopCommodity(farmerID primitive.ObjectID, year int, limit int) ([]StatisticTopCommodity, int, error)
	CountByCommodityID(commodityID primitive.ObjectID) (int, float64, int, error)
	// Update
	MakeDecision(domain *Domain, farmerID primitive.ObjectID) (int, error)
	CancelOnPending(id primitive.ObjectID, buyerID primitive.ObjectID) (int, error)
	// Delete
}
