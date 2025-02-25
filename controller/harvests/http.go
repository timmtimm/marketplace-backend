package harvests

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/harvests"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/transactions"
	"crop_connect/business/users"
	"crop_connect/constant"
	"crop_connect/controller/harvests/request"
	"crop_connect/controller/harvests/response"
	"crop_connect/helper"
	"crop_connect/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	harvestUC     harvests.UseCase
	batchUC       batchs.UseCase
	transactionUC transactions.UseCase
	proposalUC    proposals.UseCase
	commodityUC   commodities.UseCase
	userUC        users.UseCase
	regionUC      regions.UseCase
}

func NewController(harvestUC harvests.UseCase, batchUC batchs.UseCase, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) *Controller {
	return &Controller{
		harvestUC:     harvestUC,
		batchUC:       batchUC,
		transactionUC: transactionUC,
		proposalUC:    proposalUC,
		commodityUC:   commodityUC,
		userUC:        userUC,
		regionUC:      regionUC,
	}
}

/*
Create
*/

func (hc *Controller) SubmitHarvest(c echo.Context) error {
	batchID, err := primitive.ObjectIDFromHex(c.Param("batch-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id batch tidak valid",
		})
	}

	userInput := request.SubmitHarvest{}
	c.Bind(&userInput)

	validationErr := userInput.Validate()
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	images, statusCode, err := helper.GetCreateImageRequest(c, []string{"image1", "image2", "image3", "image4", "image5"})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	userID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	inputDomain, err := userInput.ToDomain()
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	inputDomain.BatchID = batchID

	_, statusCode, err = hc.harvestUC.SubmitHarvest(inputDomain, userID, images, []string{userInput.Note1, userInput.Note2, userInput.Note3, userInput.Note4, userInput.Note5})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mengajukan hasil panen",
	})
}

/*
Read
*/

func (hc *Controller) GetByPaginationAndQuery(c echo.Context) error {
	queryPagination, err := helper.PaginationToQuery(c, []string{"status", "totalPrice", "createdAt"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	token, err := helper.GetPayloadFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	queryParam, err := request.QueryParamValidation(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	harvestQuery := harvests.Query{
		Skip:        queryPagination.Skip,
		Limit:       queryPagination.Limit,
		Sort:        queryPagination.Sort,
		Order:       queryPagination.Order,
		CommodityID: queryParam.CommodityID,
		Commodity:   queryParam.Commodity,
		BatchID:     queryParam.BatchID,
		Batch:       queryParam.Batch,
		Status:      queryParam.Status,
	}

	if token.Role == constant.RoleFarmer {
		FarmerID, err := primitive.ObjectIDFromHex(token.UID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:  http.StatusBadRequest,
				Message: "token tidak valid",
			})
		}

		harvestQuery.FarmerID = FarmerID
	}

	harvests, totalData, statusCode, err := hc.harvestUC.GetByPaginationAndQuery(harvestQuery)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	harvestResponses, statusCode, err := response.FromDomainArrayToResponse(harvests, hc.batchUC, hc.transactionUC, hc.proposalUC, hc.commodityUC, hc.userUC, hc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:     statusCode,
		Message:    "berhasil mendapatkan panen",
		Data:       harvestResponses,
		Pagination: helper.ConvertToPaginationResponse(queryPagination, totalData),
	})
}

func (hc *Controller) GetByBatchID(c echo.Context) error {
	batchID, err := primitive.ObjectIDFromHex(c.QueryParam("batch-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id batch tidak valid",
		})
	}

	harvest, statusCode, err := hc.harvestUC.GetByBatchIDAndStatus(batchID, constant.HarvestStatusApproved)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	harvestResponses, statusCode, err := response.FromDomain(&harvest, hc.batchUC, hc.transactionUC, hc.proposalUC, hc.commodityUC, hc.userUC, hc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan panen",
		Data:    harvestResponses,
	})
}

func (hc *Controller) CountByYear(c echo.Context) error {
	year, err := request.QueryParamValidationYear(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	count, statusCode, err := hc.harvestUC.CountByYear(year)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan jumlah panen",
		Data:    count,
	})
}

func (hc *Controller) GetByID(c echo.Context) error {
	harvestID, err := primitive.ObjectIDFromHex(c.Param("harvest-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id panen tidak valid",
		})
	}

	harvest, statusCode, err := hc.harvestUC.GetByID(harvestID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	harvestResponses, statusCode, err := response.FromDomain(&harvest, hc.batchUC, hc.transactionUC, hc.proposalUC, hc.commodityUC, hc.userUC, hc.regionUC)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil mendapatkan panen",
		Data:    harvestResponses,
	})
}

/*
Update
*/

func (hc *Controller) Validate(c echo.Context) error {
	harvestID, err := primitive.ObjectIDFromHex(c.Param("harvest-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id panen tidak valid",
		})
	}

	userInput := request.Validate{}
	c.Bind(&userInput)

	validationErr := userInput.Validate()
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	validatorID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	inputDomain := userInput.ToDomain()
	inputDomain.ID = harvestID

	_, statusCode, err := hc.harvestUC.Validate(inputDomain, validatorID)
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil memvalidasi panen",
	})
}

func (hc *Controller) Update(c echo.Context) error {
	harvestID, err := primitive.ObjectIDFromHex(c.Param("harvest-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "id panen tidak valid",
		})
	}

	farmerID, err := helper.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	userInput := request.SubmitHarvest{}
	c.Bind(&userInput)

	validationErr := userInput.Validate()
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validasi gagal",
			Error:   validationErr,
		})
	}

	inputDomain, err := userInput.ToDomain()
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	inputDomain.ID = harvestID

	updateImages, statusCode, err := helper.GetUpdateImageRequest(c, []string{"image1", "image2", "image3", "image4", "image5"}, util.ConvertArrayStringToBool(userInput.IsChange), util.ConvertArrayStringToBool(userInput.IsDelete))
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	_, statusCode, err = hc.harvestUC.UpdateHarvest(inputDomain, farmerID, updateImages, []string{userInput.Note1, userInput.Note2, userInput.Note3, userInput.Note4, userInput.Note5})
	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
		})
	}

	return c.JSON(statusCode, helper.BaseResponse{
		Status:  statusCode,
		Message: "berhasil memperbarui riwayat perawatan",
	})
}

/*
Delete
*/
