package request

type Customer struct {
	Name        string `json:"name" form:"name" binding:"required" `
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" binding:"required"`
}

type CustomerFilter struct {
	Name        string `form:"name" binding:"omitempty,alpha"`
	Limit       int    `form:"limit" binding:"omitempty,number"`
	FrontNumber string `form:"frontNumber" binding:"omitempty,number,len=4"`
}
