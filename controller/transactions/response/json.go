package response

import (
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/transactions"
	"crop_connect/business/users"
	"net/http"

	commodityResponse "crop_connect/controller/commodities/response"
	proposalResponse "crop_connect/controller/proposals/response"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Buyer struct {
	ID         primitive.ObjectID          `json:"_id"`
	Commodity  commodityResponse.Commodity `json:"commodity"`
	Proposal   proposalResponse.Buyer      `json:"proposal"`
	Address    string                      `json:"address"`
	Status     string                      `json:"status"`
	TotalPrice float64                     `json:"totalPrice"`
	CreatedAt  primitive.DateTime          `json:"createdAt"`
}

func FromDomainToBuyer(domain *transactions.Domain, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) (Buyer, int, error) {
	proposal, statusCode, err := proposalUC.GetByIDWithoutDeleted(domain.ProposalID)
	if err != nil {
		return Buyer{}, statusCode, err
	}

	commodityDomain, statusCode, err := commodityUC.GetByIDWithoutDeleted(proposal.CommodityID)
	if err != nil {
		return Buyer{}, statusCode, err
	}

	commodity, statusCode, err := commodityResponse.FromDomain(commodityDomain, userUC, regionUC)
	if err != nil {
		return Buyer{}, statusCode, err
	}

	return Buyer{
		ID:         domain.ID,
		Commodity:  commodity,
		Proposal:   proposalResponse.FromDomainToBuyer(&proposal),
		Address:    domain.Address,
		Status:     domain.Status,
		TotalPrice: domain.TotalPrice,
		CreatedAt:  domain.CreatedAt,
	}, http.StatusOK, nil
}

func FromDomainArrayToBuyer(domain []transactions.Domain, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) ([]Buyer, int, error) {
	var buyers []Buyer
	for _, value := range domain {
		buyer, statusCode, err := FromDomainToBuyer(&value, proposalUC, commodityUC, userUC, regionUC)
		if err != nil {
			return []Buyer{}, statusCode, err
		}

		buyers = append(buyers, buyer)
	}

	return buyers, http.StatusOK, nil
}
