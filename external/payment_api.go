package external

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"restapi-bus/helper"
	"time"
)

type InterfacePayment interface {
	MakeVirtualAccount(c context.Context, name string, id string, bank_code string, expected_amount int) map[string]interface{}
	GetDataVirtualAccount(c context.Context, id string) map[string]interface{}
	PayVirtualAccount(c context.Context, externalId string, amount int) error
	SetXenditWebhookUrl(typeUrl string, urlPublic string) error
}

type Payment struct {
	passwordPayment string
	usernamePayment string
}

var paymentSingleton *Payment

func NewPayment() InterfacePayment {

	if paymentSingleton == nil {
		payment := Payment{
			passwordPayment: os.Getenv("PASSWORD_XENDIT"),
			usernamePayment: os.Getenv("SECRET_KEY_XENDIT"),
		}
		paymentSingleton = &payment
	}

	return paymentSingleton

}

func (p *Payment) MakeVirtualAccount(c context.Context, name string, id string, bank_code string, expected_amount int) map[string]interface{} {

	if name == "" || id == "" || bank_code == "" {
		panic(errors.New("something went wrong"))
	}

	client := http.DefaultClient
	dataBodyReq := struct {
		Name           string `json:"name"`
		Id             string `json:"external_id"`
		BankCode       string `json:"bank_code"`
		IsClosed       bool   `json:"is_closed"`
		ExpectedAmount int    `json:"expected_amount"`
		ExpirationDate string `json:"expiration_date"`
		IsSingleUse    bool   `json:"is_single_use"`
	}{Name: name, Id: id, BankCode: bank_code, IsClosed: true,
		ExpectedAmount: expected_amount,
		ExpirationDate: time.Now().Local().Add(1 * time.Hour).Format(time.RFC3339),
		IsSingleUse:    false,
	}
	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(dataBodyReq)
	req, err := http.NewRequestWithContext(c, http.MethodPost, "https://api.xendit.co/callback_virtual_accounts", &buffer)
	helper.PanicIfError(err)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(p.usernamePayment, p.passwordPayment)

	resp, err := client.Do(req)
	helper.PanicIfError(err)
	if resp.StatusCode != 200 {
		p, _ := io.ReadAll(resp.Body)
		fmt.Println(string(p))
		panic(fmt.Errorf("something went wrong, with error : %s", resp.Status))

	}

	defer resp.Body.Close()

	helper.PanicIfError(err)
	dataBodyResp := map[string]interface{}{}
	json.NewDecoder(resp.Body).Decode(&dataBodyResp)

	return dataBodyResp
}

func (p *Payment) GetDataVirtualAccount(c context.Context, id string) map[string]interface{} {
	client := http.DefaultClient
	url := fmt.Sprintf("https://api.xendit.co/callback_virtual_accounts/%s", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.SetBasicAuth(p.usernamePayment, p.passwordPayment)
	helper.PanicIfError(err)

	resp, err := client.Do(req)
	helper.PanicIfError(err)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		p, _ := io.ReadAll(resp.Body)
		fmt.Println(string(p))
		panic(fmt.Errorf("something went wrong, with error : %s", resp.Status))
	}

	dataStruct := map[string]interface{}{}

	err = json.NewDecoder(resp.Body).Decode(&dataStruct)
	helper.PanicIfError(err)

	return dataStruct

}

func (p *Payment) PayVirtualAccount(c context.Context, externalId string, amount int) error {
	var bodyBuffer bytes.Buffer
	err := json.NewEncoder(&bodyBuffer).Encode(map[string]int{"amount": amount})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, fmt.Sprintf("https://api.xendit.co/callback_virtual_accounts/external_id=%s/simulate_payment", externalId), &bodyBuffer)
	req.SetBasicAuth(p.usernamePayment, p.passwordPayment)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}
	client := http.DefaultClient

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		p, _ := io.ReadAll(resp.Body)
		fmt.Println(string(p))
		panic(fmt.Errorf("something went wrong, with error : %s", resp.Status))
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("something went wrong")
	}
	return nil
}

func (p *Payment) SetXenditWebhookUrl(typeUrl string, urlPublic string) error {

	body := map[string]string{
		"url": urlPublic,
	}
	var bodyBuffer bytes.Buffer
	err := json.NewEncoder(&bodyBuffer).Encode(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.xendit.co/callback_urls/%s", typeUrl), &bodyBuffer)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(p.usernamePayment, p.passwordPayment)
	if err != nil {
		return err
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("something went wrong")
	}

	return nil
}
