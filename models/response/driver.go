package response

type Driver struct {
	DriverId int    `json:"driverId"`
	AgencyId int    `json:"agencyId"`
	Name     string `json:"name"`
}
