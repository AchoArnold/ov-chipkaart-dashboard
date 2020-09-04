package ovchipkaart

import (
	"net/http"
	"strings"
	"time"
)

// HTTPClient is the class used to perform http requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// TimeInMilliSeconds represents time in milliseconds
type TimeInMilliSeconds int

// ToTime converts time in milliseconds to a time object
func (t TimeInMilliSeconds) ToTime() time.Time {
	return time.Unix(0, int64(t)*1000000)
}

// ToInt64 converts time in milliseconds to an int64 value
func (t TimeInMilliSeconds) ToInt64() int64 {
	return int64(t)
}

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

// RawRecord represents a transaction record
type RawRecord struct {
	CheckInInfo            string             `json:"checkInInfo" bson:"check_in_info"`
	CheckInText            string             `json:"checkInText" bson:"check_in_text"`
	Fare                   *float64           `json:"fare" bson:"fare"`
	FareCalculation        string             `json:"fareCalculation" bson:"fare_calculation"`
	FareText               string             `json:"fareText" bson:"fare_text"`
	ModalType              string             `json:"modalType" bson:"modal_type"`
	ProductInfo            string             `json:"productInfo" bson:"product_info"`
	ProductText            string             `json:"productText" bson:"product_text"`
	Pto                    string             `json:"pto" bson:"pto"`
	TransactionDateTime    TimeInMilliSeconds `json:"transactionDateTime" bson:"transaction_timestamp"`
	TransactionInfo        string             `json:"transactionInfo" bson:"transaction_info"`
	TransactionName        TransactionName    `json:"transactionName" bson:"transaction_name"`
	EPurseMut              *float64           `json:"ePurseMut" bson:"e_purse_mut"`
	EPurseMutInfo          string             `json:"ePurseMutInfo" bson:"e_purse_mut_info"`
	TransactionExplanation string             `json:"transactionExplanation" bson:"transaction_explanation"`
	TransactionPriority    string             `json:"transactionPriority" bson:"transaction_priority"`
}

// TransactionFetchOptions are the options needed when fetching a list of transactions
type TransactionFetchOptions struct {
	Username   string
	Password   string
	CardNumber string
	StartDate  time.Time
	EndDate    time.Time
}
