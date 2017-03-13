package util

// CreateAccountReq - used to create account
type CreateAccountReq struct {
	UserName  string
	Email     string
	FirstName string
	LastName  string
	PWD       string
	CompanyID uint32 `json:"CompanyID,omitempty"`
	Role      uint32 `json:"Role,omitempty"`
}

// LoginReq - Login Account
type LoginReq struct {
	UserName string
	PWD      string
}

// SearchAccountReq - used to Search account
type SearchAccountReq struct {
	UserName  string
	Email     string
	FirstName string
	LastName  string
	PWD       string
	CompanyID uint32 `json:"CompanyID,omitempty"`
	Role      uint32 `json:"Role,omitempty"`
}
