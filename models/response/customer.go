package response

type Customer struct {
	CustomerId  int    `json:"customerId"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}
