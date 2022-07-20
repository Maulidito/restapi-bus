package testing

import (
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllTicket(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)

	repo := repository.NewTicketRepository()
	listTikcet := repo.GetAllTicket(tx, ctx, "")

	fmt.Println(listTikcet)
	assert.NotEmpty(t, listTikcet)
	assert.NotZero(t, len(listTikcet))

}

func TestGetOneTicket(t *testing.T) {
	tx, err := db.Begin()

	helper.PanicIfError(err)

	ticket := entity.Ticket{TicketId: 1}
	repo := repository.NewTicketRepository()
	repo.GetOneTicket(tx, ctx, &ticket)

	fmt.Println(ticket)
	assert.NotEmpty(t, ticket)

}

func TestGetAllTicketOnDriver(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)

	repo := repository.NewTicketRepository()
	listTicket := repo.GetAllTicketOnDriver(tx, ctx, 1)

	fmt.Println(listTicket)
	assert.Empty(t, listTicket)
}
