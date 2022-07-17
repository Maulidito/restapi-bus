package response

type Driver struct {
	DriverId int    `json:"driverId"`
	AgencyId int    `json:"agencyId"`
	Name     string `json:"name"`
}

type AllDriverOnAgency struct {
	Agency *Agency   `json:"agency"`
	Driver *[]Driver `json:"driver"`
}

type DriverNoAgency struct {
	DriverId int    `json:"driverId"`
	Name     string `json:"name"`
}
