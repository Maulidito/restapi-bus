package web

import "restapi-bus/models/response"

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type AllBusOnAgency struct {
	Agency *response.Agency `json:"agency"`
	Bus    *[]response.Bus  `json:"bus"`
}

type WebResponseAllBusOnAgency struct {
	Code   int            `json:"code"`
	Status string         `json:"status"`
	Data   AllBusOnAgency `json:"data"`
}
