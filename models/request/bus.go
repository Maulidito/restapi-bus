package request

type Bus struct {
	AgencyId    int    `json:"agencyId" form:"agencyId" binding:"required"`
	NumberPlate string `json:"numberPlate" form:"numberPlate" binding:"required"`
}

type BusFilter struct {
	FrontNumberPlate string `form:"frontNumberPlate" binding:"omitempty,alpha,len=1,uppercase"`
}
