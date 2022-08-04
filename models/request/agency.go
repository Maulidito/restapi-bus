package request

type Agency struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Place string `json:"place" form:"place" binding:"required"`
	Auth  *AgencyAuth
}

type AgencyFilter struct {
	Name          string `form:"name" binding:"omitempty"`
	Limit         int    `form:"limit" binding:"omitempty,numeric"`
	Place         string `form:"place" binding:"omitempty,alpha"`
	BelowBusCount int    `form:"belowBusCount" binding:"omitempty,numeric"`
	AboveBusCount int    `form:"aboveBusCount" binding:"omitempty,numeric"`
}

type AgencyAuth struct {
	Username string `form:"username" binding:"required,email"`
	Password string `form:"password" binding:"required,min=8,validateoneuppercase"`
}
