package request

type Customer struct {
	Name        string `json:"name" form:"name" binding:"required" `
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" binding:"required"`
}
