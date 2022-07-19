package request

type Driver struct {
	AgencyId int    `json:"agencyId" form:"agencyId" binding:"required"`
	Name     string `json:"name" form:"name" binding:"required"`
}

type DriverFilter struct {
	Name  string `form:"name" binding:"omitempty,alpha"`
	Limit int    `form:"limit" binding:"omitempty,number"`
}
