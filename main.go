package main

import (
	"fmt"

	_route "marketplace-backend/app/route"
	_driver "marketplace-backend/driver"
	_mongo "marketplace-backend/driver/mongo"
	_util "marketplace-backend/util"

	_batchUseCase "marketplace-backend/business/batchs"
	_commodityUseCase "marketplace-backend/business/commodities"
	_proposalUseCase "marketplace-backend/business/proposals"
	_transactionUseCase "marketplace-backend/business/transactions"
	_userUseCase "marketplace-backend/business/users"

	_batchController "marketplace-backend/controller/batchs"
	_commodityController "marketplace-backend/controller/commodities"
	_proposalController "marketplace-backend/controller/proposals"
	_transactionController "marketplace-backend/controller/transactions"
	_userController "marketplace-backend/controller/users"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	database := _mongo.Init(_util.GetConfig("DB_NAME"))

	userRepository := _driver.NewUserRepository(database)
	commodityRepository := _driver.NewCommodityRepository(database)
	proposalRepository := _driver.NewProposalRepository(database)
	transactionRepository := _driver.NewTransactionRepository(database)
	batchRepository := _driver.NewBatchRepository(database)

	userUseCase := _userUseCase.NewUserUseCase(userRepository)
	commodityUsecase := _commodityUseCase.NewCommodityUseCase(commodityRepository)
	proposalUseCase := _proposalUseCase.NewProposalUseCase(proposalRepository, commodityRepository)
	transactionUseCase := _transactionUseCase.NewTransactionUseCase(transactionRepository, commodityRepository, proposalRepository)
	batchUseCase := _batchUseCase.NewBatchUseCase(batchRepository, transactionRepository, proposalRepository, commodityRepository)

	userController := _userController.NewUserController(userUseCase)
	commodityController := _commodityController.NewCommodityController(commodityUsecase, userUseCase, proposalUseCase)
	proposalController := _proposalController.NewProposalController(proposalUseCase, commodityUsecase)
	transactionController := _transactionController.NewTransactionController(transactionUseCase, proposalUseCase, commodityUsecase, userUseCase, batchUseCase)
	batchController := _batchController.NewBatchController(batchUseCase)

	routeController := _route.ControllerList{
		UserController:        userController,
		CommodityController:   commodityController,
		ProposalController:    proposalController,
		TransactionController: transactionController,
		BatchController:       batchController,
	}

	routeController.Init(e)

	appPort := fmt.Sprintf(":%s", _util.GetConfig("APP_PORT"))
	e.Logger.Fatal(e.Start(appPort))
}
