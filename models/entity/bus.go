package entity

import (
	"context"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type Bus struct {
	BusId       int
	AgencyId    int
	NumberPlate string
}

type BusServiceInterface interface {
	GetAllBus(ctx context.Context, filter *request.BusFilter) []response.Bus
	AddBus(ctx context.Context, bus *request.Bus)
	GetOneBusSpecificAgency(ctx context.Context, idBus int) response.Bus
	DeleteOneBus(ctx context.Context, idBus int) response.Bus
	GetAllBusOnSpecificAgency(ctx context.Context, idAgency int) []response.Bus
}

type BusRepositoryInterface interface {
	GetAllBus(ctx context.Context, filter *request.BusFilter) []Bus
	AddBus(ctx context.Context, bus *Bus)
	GetOneBus(ctx context.Context, bus *Bus)
	DeleteOneBus(ctx context.Context, bus *Bus)
	GetAllBusSpecificAgency(ctx context.Context, agencyId int) []Bus
}
