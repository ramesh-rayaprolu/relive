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

// CreateProductReq - used to create product
type CreateProductReq struct {
	ProductID   uint32 `json:"ProductID"`
	ProductType string `json:"ProductType"`
	StoreSize   uint32 `json:"StoreSize,omitempty"`
	Duration    uint32 `json:"Duration,omitempty"`
	Amount      uint32 `json:"Amount,omitempty"`
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
