package service

import (
	"context"
	"restapi-bus/helper"
	"restapi-bus/models/database"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type ServiceDriverImplementation struct {
	RepoDriver entity.DriverRepositoryInterface
	RepoAgency entity.AgencyRepositoryInterface
	Tx         database.TrInterface
}

func NewServiceDriver(repoDriver entity.DriverRepositoryInterface, repoAgency entity.AgencyRepositoryInterface, tx database.TrInterface) entity.ServiceDriverInterface {
	return &ServiceDriverImplementation{RepoDriver: repoDriver, RepoAgency: repoAgency, Tx: tx}
}

func (service *ServiceDriverImplementation) GetAllDriver(ctx context.Context, filter *request.DriverFilter) []response.Driver {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	listDriver := service.RepoDriver.GetAllDriver(ctx, filter)

	res := []response.Driver{}
	for _, val := range listDriver {
		res = append(res, helper.DriverEntityToResponse(&val))
	}

	return res

}
func (service *ServiceDriverImplementation) GetAllDriverOnSpecificAgency(ctx context.Context, agencyId int) []response.Driver {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	agency := entity.Agency{AgencyId: agencyId}

	listDriverResponse := []response.Driver{}

	service.RepoAgency.GetOneAgency(ctx, &agency)
	listDriver := service.RepoDriver.GetAllDriverOnSpecificAgency(ctx, agency.AgencyId)

	for _, val := range listDriver {
		listDriverResponse = append(listDriverResponse, helper.DriverEntityToResponse(&val))
	}

	return listDriverResponse

}
func (service *ServiceDriverImplementation) GetOneDriverOnSpecificAgency(ctx context.Context, driverId int) response.Driver {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	driver := entity.Driver{DriverId: driverId}

	service.RepoDriver.GetOneDriverOnSpecificAgency(ctx, &driver)

	return helper.DriverEntityToResponse(&driver)

}
func (service *ServiceDriverImplementation) AddDriver(ctx context.Context, driver *request.Driver) {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	service.RepoAgency.GetOneAgency(ctx, &entity.Agency{AgencyId: driver.AgencyId})
	driverEntity := helper.DriverRequestToEntity(driver)

	service.RepoDriver.AddDriver(ctx, &driverEntity)
}
func (service *ServiceDriverImplementation) DeleteDriver(ctx context.Context, driverId int) response.Driver {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	driverEntity := entity.Driver{DriverId: driverId}
	service.RepoDriver.GetOneDriverOnSpecificAgency(ctx, &driverEntity)
	service.RepoDriver.DeleteDriver(ctx, &driverEntity)

	return helper.DriverEntityToResponse(&driverEntity)
}
