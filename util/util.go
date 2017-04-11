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

// CreateSubscriptionReq - used to create subscription
type CreateSubscriptionReq struct {
	ID               uint32 // ID is constrained to Account ID
	ProductID        uint32 // constrained
	SubscriptionCode uint32
	ProductType      string
	StoreLocation    string
	StartDate        string
	EndDate          string
	NumberOfAdmins   uint32
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

// SubscrDetails struct
type SubscrDetails struct {
	ID          int
	ProductID   int
	SubscrCode  int
	ProductType string
}
