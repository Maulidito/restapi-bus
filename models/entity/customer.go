package entity

import (
	"context"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type Customer struct {
	CustomerId  int
	Name        string
	PhoneNumber string
	Email       string
}

type CustomerServiceInterface interface {
	GetAllCustomer(ctx context.Context, filter *request.CustomerFilter) []response.Customer
	AddCustomer(ctx context.Context, customer *request.Customer)
	GetOneCustomer(ctx context.Context, id int) response.Customer
	DeleteOneCustomer(ctx context.Context, id int) response.Customer
}

type CustomerRepositoryInterface interface {
	GetAllCustomer(ctx context.Context, filter *request.CustomerFilter) []Customer
	AddCustomer(ctx context.Context, customer *Customer)
	GetOneCustomer(ctx context.Context, customer *Customer)
	DeleteOneCustomer(ctx context.Context, customer *Customer)
}
