package proposals

import (
	"crop_connect/business/proposals"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID                    primitive.ObjectID `bson:"_id"`
	Code                  primitive.ObjectID `bson:"code"`
	ValidatorID           primitive.ObjectID `bson:"validatorID,omitempty"`
	CommodityID           primitive.ObjectID `bson:"commodityID"`
	RegionID              primitive.ObjectID `bson:"regionID"`
	Name                  string             `bson:"name"`
	Description           string             `bson:"description"`
	Status                string             `bson:"status"`
	RejectReason          string             `bson:"rejectReason,omitempty"`
	EstimatedTotalHarvest float64            `bson:"estimatedTotalHarvest"`
	PlantingArea          float64            `bson:"plantingArea"`
	Address               string             `bson:"address"`
	IsAvailable           bool               `bson:"isAvailable"`
	CreatedAt             primitive.DateTime `bson:"createdAt"`
	UpdatedAt             primitive.DateTime `bson:"updatedAt,omitempty"`
	DeletedAt             primitive.DateTime `bson:"deletedAt,omitempty"`
}

func FromDomain(domain *proposals.Domain) *Model {
	return &Model{
		ID:                    domain.ID,
		Code:                  domain.Code,
		ValidatorID:           domain.ValidatorID,
		CommodityID:           domain.CommodityID,
		RegionID:              domain.RegionID,
		Name:                  domain.Name,
		Description:           domain.Description,
		Status:                domain.Status,
		RejectReason:          domain.RejectReason,
		EstimatedTotalHarvest: domain.EstimatedTotalHarvest,
		PlantingArea:          domain.PlantingArea,
		Address:               domain.Address,
		IsAvailable:           domain.IsAvailable,
		CreatedAt:             domain.CreatedAt,
		UpdatedAt:             domain.UpdatedAt,
		DeletedAt:             domain.DeletedAt,
	}
}

func (model *Model) ToDomain() proposals.Domain {
	return proposals.Domain{
		ID:                    model.ID,
		Code:                  model.Code,
		ValidatorID:           model.ValidatorID,
		CommodityID:           model.CommodityID,
		RegionID:              model.RegionID,
		Name:                  model.Name,
		Description:           model.Description,
		Status:                model.Status,
		EstimatedTotalHarvest: model.EstimatedTotalHarvest,
		PlantingArea:          model.PlantingArea,
		Address:               model.Address,
		IsAvailable:           model.IsAvailable,
		CreatedAt:             model.CreatedAt,
		UpdatedAt:             model.UpdatedAt,
		DeletedAt:             model.DeletedAt,
	}
}

func ToDomainArray(model []Model) []proposals.Domain {
	var result []proposals.Domain
	for _, proposal := range model {
		result = append(result, proposal.ToDomain())
	}
	return result
}
