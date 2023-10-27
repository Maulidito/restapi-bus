package request

type Bus struct {
	AgencyId    int    `json:"agencyId" form:"agencyId" binding:"required"`
	NumberPlate string `json:"numberPlate" form:"numberPlate" binding:"required"`
	TotalSeat   int    `default:"10" json:"totalSeat" form:"totalSeat,default=10" binding:"numeric,min=4"`
}

type BusFilter struct {
	FrontNumberPlate string `form:"frontNumberPlate" binding:"omitempty,alpha,len=1,uppercase"`
}
