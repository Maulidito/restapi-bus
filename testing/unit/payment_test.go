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
	data := payment.MakeVirtualAccount(&c, "Ismail Rabbanii", "VA_fixed-1615935264", "MANDIRI", 200000)

	if data["id"] == "" {
		t.Error(fmt.Sprintf("id not found"))
	}
	if data["external_id"] == "" {
		t.Error(fmt.Sprintf("external id not found"))
	}

	err := recover()
	if err != nil {
		t.Fatalf("got err %s", err.(error).Error())
	}
	t.Logf("got id :%s, external_id : %s ", data["id"], data["external_id"])

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
	data := payment.MakeVirtualAccount(context.Background(), "Ismail Rabbnii", "VA_fixed-1715931461", "MANDIRI", amount)
	data1 := payment.MakeVirtualAccount(context.Background(), "bAMBANG", "VA_fixed-1923123423", "MANDIRI", amount)

	err := payment.PayVirtualAccount(context.Background(), data["external_id"].(string), amount)

	if err != nil {
		t.Errorf("got error : %s", err.Error())
	}
	payment.GetDataVirtualAccount(context.Background(), data["id"].(string))
	err = payment.PayVirtualAccount(context.Background(), data1["external_id"].(string), amount)
	if err != nil {
		t.Errorf("got error : %s", err.Error())
	}
	if recover() != nil {
		t.Errorf("got error : %s", err.Error())
	}
	t.Logf("SUCCESS %v %s", data, data["id"])
}
