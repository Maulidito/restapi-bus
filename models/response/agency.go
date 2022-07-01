package response

type Agency struct {
	AgencyId int    `json:"agencyId"`
	Name     string `json:"name"`
	Place    string `json:"place"`
}
