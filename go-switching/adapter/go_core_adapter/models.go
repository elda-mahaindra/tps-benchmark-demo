package go_core_adapter

// Account represents account information from go-core service
type Account struct {
	AccountID     int64
	AccountNumber string
	CustomerID    int64
	AccountType   string
	AccountStatus string
	Balance       string
	Currency      string
	OpenedDate    string
	ClosedDate    string
	CreatedAt     string
	UpdatedAt     string
}

// Customer represents customer information from go-core service
type Customer struct {
	CustomerNumber string
	FullName       string
	IDNumber       string
	PhoneNumber    string
	Email          string
	Address        string
	DateOfBirth    string
}
