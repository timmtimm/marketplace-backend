package users

import (
	"crop_connect/business/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID          primitive.ObjectID `bson:"_id"`
	RegionID    primitive.ObjectID `bson:"regionID"`
	Name        string             `bson:"name"`
	Email       string             `bson:"email"`
	Description string             `bson:"description"`
	PhoneNumber string             `bson:"phoneNumber"`
	Password    string             `bson:"password"`
	Role        string             `bson:"role"`
	CreatedAt   primitive.DateTime `bson:"createdAt"`
	UpdatedAt   primitive.DateTime `bson:"updatedAt,omitempty"`
}

func FromDomain(domain *users.Domain) *Model {
	return &Model{
		ID:          domain.ID,
		RegionID:    domain.RegionID,
		Name:        domain.Name,
		Email:       domain.Email,
		Description: domain.Description,
		PhoneNumber: domain.PhoneNumber,
		Password:    domain.Password,
		Role:        domain.Role,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}

func (model *Model) ToDomain() users.Domain {
	return users.Domain{
		ID:          model.ID,
		RegionID:    model.RegionID,
		Name:        model.Name,
		Email:       model.Email,
		Description: model.Description,
		PhoneNumber: model.PhoneNumber,
		Password:    model.Password,
		Role:        model.Role,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func ToDomainArray(models []Model) []users.Domain {
	var result []users.Domain
	for _, value := range models {
		result = append(result, value.ToDomain())
	}
	return result
}
