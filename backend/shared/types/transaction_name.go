package types

import "strings"

// TransactionName  represents the various transaction names
type TransactionName string

// String returns the transaction name as a string
func (name TransactionName) String() string {
	return string(name)
}

// IsTheSameAs is used to compare 2 transaction names
func (name TransactionName) IsTheSameAs(comp TransactionName) bool {
	return strings.ToLower(name.String()) == strings.ToLower(comp.String())
}

const (
	// TransactionNameCheckIn represents a check in transaction
	TransactionNameCheckIn = TransactionName("Check-in")

	// TransactionNameCheckOut represents a check out transaction
	TransactionNameCheckOut = TransactionName("Check-uit")

	// TransactionNameIntercityDirectSurcharge represents an intercity direct transaction
	TransactionNameIntercityDirectSurcharge = TransactionName("Toeslag Intercity Direct")
)
