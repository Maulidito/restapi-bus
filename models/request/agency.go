package request

type Agency struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Place string `json:"place" form:"place" binding:"required"`
}
