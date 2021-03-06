package request

type Agency struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Place string `json:"place" form:"place" binding:"required"`
}

type AgencyFilter struct {
	Name          string `form:"name" binding:"omitempty"`
	Limit         int    `form:"limit" binding:"omitempty,numeric"`
	Place         string `form:"place" binding:"omitempty,alpha"`
	BelowBusCount int    `form:"belowBusCount" binding:"omitempty,numeric"`
	AboveBusCount int    `form:"aboveBusCount" binding:"omitempty,numeric"`
}
