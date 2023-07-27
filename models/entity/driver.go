package entity

import (
	"context"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type Driver struct {
	DriverId int
	AgencyId int
	Name     string
}

type ServiceDriverInterface interface {
	GetAllDriver(ctx context.Context, filter *request.DriverFilter) []response.Driver
	GetAllDriverOnSpecificAgency(ctx context.Context, agencyId int) []response.Driver
	GetOneDriverOnSpecificAgency(ctx context.Context, driverId int) response.Driver
	AddDriver(ctx context.Context, driver *request.Driver)
	DeleteDriver(ctx context.Context, driverId int) response.Driver
}

type DriverRepositoryInterface interface {
	GetAllDriver(ctx context.Context, filter *request.DriverFilter) []Driver
	GetAllDriverOnSpecificAgency(ctx context.Context, agencyId int) []Driver
	GetOneDriverOnSpecificAgency(ctx context.Context, driver *Driver)
	AddDriver(ctx context.Context, driver *Driver)
	DeleteDriver(ctx context.Context, driver *Driver)
}
