package dbmodel

type (
	// AccountEntry - testing
	AccountEntry struct {
		ID           int
		PID          int
		FirstName    string
		LastName     string
		EmailID      string
		PasswdDigest string
		Role         int
	}

	// PaymentEntry - testing
	PaymentEntry struct {
		ID             int
		CCNumber       string
		BillingAddress string
		CCExpiry       string
		CVVCode        int
	}

	// PaymentHistoryEntry - testing
	PaymentHistoryEntry struct {
		ID            int
		LastPaidState string
		LastType      string
	}

	// SubscriptionEntry - testing
	SubscriptionEntry struct {
		ID               int
		ProductID        int
		SubscriptionCode int
		ProductType      string
		StoreLocation    string
		StartDate        string
		EndDate          string
		NumberOfAdmins   int
	}

	// SubscriptionAccountEntry - testing
	SubscriptionAccountEntry struct {
		ID  int
		PID int
	}

	// ProductEntry - testing
	ProductEntry struct {
		ProductID   int
		ProductType string
		StoreSize   int
		Duration    int
		Amount      int
	}

	// MediaTypeEntry - testing
	MediaTypeEntry struct {
		ID          int
		Catalog     string
		FileName    string
		Title       string
		Description string
		URL         string
		Poster      string
	}
)
