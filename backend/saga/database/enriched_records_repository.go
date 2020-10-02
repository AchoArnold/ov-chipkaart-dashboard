package database

// EnrichedRecordsRepository fetches enriched records.
type EnrichedRecordsRepository interface {
	StoreMany(records []EnrichedRecord) (err error)
}
