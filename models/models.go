package models

type Response struct {
	Code     int
	Msg      string
	Response interface{} `json:"omitempty"`
}
type Record struct {
	Address     string `json:"address"`
	City        string `json:"city"`
	CompanyName string `json:"company_name"`
	County      string `json:"county"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Postal      string `json:"postal"`
	Web         string `json:"web"`
}

type EditRequest struct {
	Id          *int    `json:"id"`
	Address     *string `json:"address,omitempty"`
	City        *string `json:"city,omitempty"`
	CompanyName *string `json:"company_name,omitempty"`
	County      *string `json:"county,omitempty"`
	Email       *string `json:"email,omitempty"`
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Postal      *string `json:"postal,omitempty"`
	Web         *string `json:"web,omitempty"`
}
type OrderedRecord struct {
	Key  string `json:"key"`
	Data Record `json:"data"`
}
