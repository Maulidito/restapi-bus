package entity

import (
	"context"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type Agency struct {
	AgencyId int
	Name     string
	Place    string
	Username string
	Password string
	Salt     string
}

type AgencyServiceInterface interface {
	GetAllAgency(ctx context.Context, filter *request.AgencyFilter) []response.Agency
	RegisterAgency(ctx context.Context, agency *request.Agency)
	GetOneAgency(ctx context.Context, id int) response.Agency
	DeleteOneAgency(ctx context.Context, id int) response.Agency
	LoginAgency(ctx context.Context, agencyAuth *request.AgencyAuth) (string, int, response.Agency)
}

type AgencyRepositoryInterface interface {
	GetAllAgency(ctx context.Context, filter *request.AgencyFilter) []Agency
	RegisterAgency(ctx context.Context, agency *Agency)
	GetOneAgency(ctx context.Context, agency *Agency)
	IsUsenameAgencyExist(ctx context.Context, agencyUsername string) bool
	DeleteOneAgency(ctx context.Context, agency *Agency)
	GetOneAgencyAuth(ctx context.Context, agency *Agency)
	GetSaltAgencyWithUsername(ctx context.Context, agencyUsername string) string
}
