package request

type Customer struct {
	Name        string `json:"name" form:"name" binding:"required" `
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" binding:"required"`
}

type CustomerFilter struct {
	Name        string `form:"name" binding:"alpha"`
	Limit       int    `form:"limit" binding:"number"`
	FrontNumber string `form:"frontNumber" binding:"alpha, len=4"`
}
