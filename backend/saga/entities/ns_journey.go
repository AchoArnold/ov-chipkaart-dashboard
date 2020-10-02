package entities

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	internalTime "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/time"
)

// NSJourney are options for fetching the price of a journey
type NSJourney struct {
	Year            string `bson:"year"`
	FromStationCode string `bson:"from_station_code"`
	ToStationCode   string `bson:"to_station_code"`
	date            time.Time
}

// NewNSJourney creates a new NSJourney instance
func NewNSJourney(timestamp time.Time, fromStationCode, toStationCode string) NSJourney {
	return NSJourney{
		Year:            timestamp.Format(internalTime.YearFormat),
		FromStationCode: fromStationCode,
		ToStationCode:   toStationCode,
		date:            timestamp,
	}
}

// ToMap converts  the JS journey struct to a `map[string]string` map
func (journey NSJourney) ToMap() map[string]string {
	date := journey.date.Format(internalTime.DateFormat)
	if journey.date.Year() < time.Now().Year() {
		parsed, err := time.Parse(internalTime.DateFormat, strconv.Itoa(time.Now().Year()-1)+"-12-30")
		if err == nil {
			date = parsed.Format(internalTime.DateFormat)
		}
	}

	return map[string]string{
		"date":        date,
		"toStation":   journey.FromStationCode,
		"fromStation": journey.ToStationCode,
	}
}

// NSPriceHash gets the hash for an ns journey used to determine the price of the journey
func (journey NSJourney) NSPriceHash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(journey.FromStationCode+"-"+journey.ToStationCode+"-"+journey.Year)))
}
