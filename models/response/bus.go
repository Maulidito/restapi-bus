package response

type Bus struct {
	BusId       int    `json:"busId"`
	AgencyId    int    `json:"agencyId"`
	NumberPlate string `json:"numberPlate"`
	TotalSeat   int    `json:"totalSeat"`
}

type BusNoAgency struct {
	BusId       int    `json:"busId"`
	NumberPlate string `json:"numberPlate"`
	TotalSeat   int    `json:"totalSeat"`
}

type AllBusOnAgency struct {
	Agency *Agency `json:"agency"`
	Bus    *[]Bus  `json:"bus"`
}
