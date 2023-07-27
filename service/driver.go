package service

import (
	"context"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type ServiceDriverImplementation struct {
	RepoDriver entity.DriverRepositoryInterface
	RepoAgency entity.AgencyRepositoryInterface
}

func NewServiceDriver(repoDriver entity.DriverRepositoryInterface, repoAgency entity.AgencyRepositoryInterface) entity.ServiceDriverInterface {
	return &ServiceDriverImplementation{RepoDriver: repoDriver, RepoAgency: repoAgency}
}

func (service *ServiceDriverImplementation) GetAllDriver(ctx context.Context, filter *request.DriverFilter) []response.Driver {

	listDriver := service.RepoDriver.GetAllDriver(ctx, filter)

	res := []response.Driver{}
	for _, val := range listDriver {
		res = append(res, helper.DriverEntityToResponse(&val))
	}

	return res

}
func (service *ServiceDriverImplementation) GetAllDriverOnSpecificAgency(ctx context.Context, agencyId int) []response.Driver {

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
		service.RepoAgency.GetOneAgency(ctx, &agency)
		listDriver = service.RepoDriver.GetAllDriverOnSpecificAgency(ctx, agency.AgencyId)
	}()

	helper.PanicIfError(<-chanErr)

	for _, val := range listDriver {
		listDriverResponse = append(listDriverResponse, helper.DriverEntityToResponse(&val))
	}

	return listDriverResponse

}
func (service *ServiceDriverImplementation) GetOneDriverOnSpecificAgency(ctx context.Context, driverId int) response.Driver {

	driver := entity.Driver{DriverId: driverId}

	service.RepoDriver.GetOneDriverOnSpecificAgency(ctx, &driver)

	return helper.DriverEntityToResponse(&driver)

}
func (service *ServiceDriverImplementation) AddDriver(ctx context.Context, driver *request.Driver) {

	service.RepoAgency.GetOneAgency(ctx, &entity.Agency{AgencyId: driver.AgencyId})
	driverEntity := helper.DriverRequestToEntity(driver)

	service.RepoDriver.AddDriver(ctx, &driverEntity)
}
func (service *ServiceDriverImplementation) DeleteDriver(ctx context.Context, driverId int) response.Driver {

	driverEntity := entity.Driver{DriverId: driverId}
	service.RepoDriver.GetOneDriverOnSpecificAgency(ctx, &driverEntity)
	service.RepoDriver.DeleteDriver(ctx, &driverEntity)

	return helper.DriverEntityToResponse(&driverEntity)
}
