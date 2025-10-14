package service

// Account represents account information in the service layer
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

// Customer represents customer information in the service layer
type Customer struct {
	CustomerNumber string
	FullName       string
	IDNumber       string
	PhoneNumber    string
	Email          string
	Address        string
	DateOfBirth    string
}
