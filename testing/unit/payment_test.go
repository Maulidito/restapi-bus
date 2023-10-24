package testing

import (
	"context"
	"fmt"
	"restapi-bus/external"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func TestMakeVirtualAccount(t *testing.T) {
	godotenv.Load("../../.env")
	payment := external.NewPayment()
	c := gin.Context{}
	id, external_id := payment.MakeVirtualAccount(&c, "Ismail Rabbanii", "VA_fixed-1615935264", "MANDIRI", 200000, "1234567890")

	if id == "" {
		t.Error(fmt.Sprintf("id not found"))
	}
	if external_id == "" {
		t.Error(fmt.Sprintf("external id not found"))
	}

	err := recover()
	if err != nil {
		t.Fatalf("got err %s", err.(error).Error())
	}
	t.Logf("got id :%s, external_id : %s ", id, external_id)

}

func TestGetVirtualAccount(t *testing.T) {
	godotenv.Load("../../.env")
	payment := external.NewPayment()
	data := payment.GetDataVirtualAccount(context.Background(), "59444945-73b3-4846-95e7-1b4eb1ded0b1")

	err := recover()

	if err != nil {
		t.Errorf("got error %s", err)
	}
	if data == nil {
		t.Errorf("data got %s", data)
	}
	t.Logf("got data %s", data["is_closed"])

}

func TestPayVirtualAccount(t *testing.T) {
	godotenv.Load("../../.env")
	payment := external.NewPayment()
	amount := 200000
	id, external_id := payment.MakeVirtualAccount(context.Background(), "Ismail Rabbnii", "VA_fixed-1715931461", "MANDIRI", amount, "0034761800")
	_, external_id1 := payment.MakeVirtualAccount(context.Background(), "bAMBANG", "VA_fixed-1923123423", "MANDIRI", amount, "0034761800")
	err := payment.PayVirtualAccount(context.Background(), external_id, amount)

	if err != nil {
		t.Errorf("got error : %s", err.Error())
	}
	data := payment.GetDataVirtualAccount(context.Background(), id)
	err = payment.PayVirtualAccount(context.Background(), external_id1, amount)
	if err != nil {
		t.Errorf("got error : %s", err.Error())
	}
	if recover() != nil {
		t.Errorf("got error : %s", err.Error())
	}
	t.Logf("SUCCESS %v %s", data, id)
}
