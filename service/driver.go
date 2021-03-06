package service

import (
	"context"
	"database/sql"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/repository"
)

type ServiceDriverInterface interface {
	GetAllDriver(ctx context.Context, filter *request.DriverFilter) []response.Driver
	GetAllDriverOnSpecificAgency(ctx context.Context, agencyId int) []response.Driver
	GetOneDriverOnSpecificAgency(ctx context.Context, driverId int) response.Driver
	AddDriver(ctx context.Context, driver *request.Driver)
	DeleteDriver(ctx context.Context, driverId int) response.Driver
}

type ServiceDriverImplementation struct {
	Db         *sql.DB
	RepoDriver repository.DriverRepositoryInterface
	RepoAgency repository.AgencyRepositoryInterface
}

func NewServiceDriver(db *sql.DB, repoDriver repository.DriverRepositoryInterface, repoAgency repository.AgencyRepositoryInterface) ServiceDriverInterface {
	return &ServiceDriverImplementation{Db: db, RepoDriver: repoDriver, RepoAgency: repoAgency}
}

func (service *ServiceDriverImplementation) GetAllDriver(ctx context.Context, filter *request.DriverFilter) []response.Driver {

	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)

	listDriver := service.RepoDriver.GetAllDriver(tx, ctx, helper.RequestFilterDriverToString(filter))

	res := []response.Driver{}
	for _, val := range listDriver {
		res = append(res, helper.DriverEntityToResponse(&val))
	}

	return res

}
func (service *ServiceDriverImplementation) GetAllDriverOnSpecificAgency(ctx context.Context, agencyId int) []response.Driver {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)

	agency := entity.Agency{AgencyId: agencyId}
	listDriver := []entity.Driver{}
	listDriverResponse := []response.Driver{}
	chanErr := make(chan error, 1)
	go func() {
		defer func() {
			tempRecover := recover()

			if tempRecover != nil {
				chanErr <- tempRecover.(error)
			}
			close(chanErr)

		}()
		service.RepoAgency.GetOneAgency(ctx, tx, &agency)
		listDriver = service.RepoDriver.GetAllDriverOnSpecificAgency(tx, ctx, agency.AgencyId)
	}()

	helper.PanicIfError(<-chanErr)

	for _, val := range listDriver {
		listDriverResponse = append(listDriverResponse, helper.DriverEntityToResponse(&val))
	}

	return listDriverResponse

}
func (service *ServiceDriverImplementation) GetOneDriverOnSpecificAgency(ctx context.Context, driverId int) response.Driver {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)

	driver := entity.Driver{DriverId: driverId}

	service.RepoDriver.GetOneDriverOnSpecificAgency(tx, ctx, &driver)

	return helper.DriverEntityToResponse(&driver)

}
func (service *ServiceDriverImplementation) AddDriver(ctx context.Context, driver *request.Driver) {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)

	service.RepoAgency.GetOneAgency(ctx, tx, &entity.Agency{AgencyId: driver.AgencyId})
	driverEntity := helper.DriverRequestToEntity(driver)

	service.RepoDriver.AddDriver(tx, ctx, &driverEntity)
}
func (service *ServiceDriverImplementation) DeleteDriver(ctx context.Context, driverId int) response.Driver {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)

	driverEntity := entity.Driver{DriverId: driverId}
	service.RepoDriver.GetOneDriverOnSpecificAgency(tx, ctx, &driverEntity)
	service.RepoDriver.DeleteDriver(tx, ctx, &driverEntity)

	return helper.DriverEntityToResponse(&driverEntity)
}
