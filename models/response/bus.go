package response

type Bus struct {
	BusId       int    `json:"busId"`
	AgencyId    int    `json:"agencyId"`
	NumberPlate string `json:"numberPlate"`
}
