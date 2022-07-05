package request

type Driver struct {
	AgencyId int    `json:"agencyId" form:"agencyId" binding:"required"`
	Name     string `json:"name" form:"name" binding:"required"`
}
