package response

type Bus struct {
	BusId       int    `json:"busId"`
	AgencyId    int    `json:"agencyId"`
	NumberPlate string `json:"numberPlate"`
}

type BusNoAgency struct {
	BusId       int    `json:"busId"`
	NumberPlate string `json:"numberPlate"`
}
