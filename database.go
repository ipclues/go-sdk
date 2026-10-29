package ipclues

// Database identifies an IPClues lookup database.
type Database string

// Valida reports whether the database ID is supported by the SDK.
func (d Database) Valid() bool {
	switch d {
	case DatabaseGeoI3Lite, DatabaseGeoI3Country, DatabaseNetI3Identity:
		return true
	default:
		return false
	}
}

const (
	DatabaseGeoI3Lite     Database = "geoi3-lite"
	DatabaseGeoI3Country  Database = "geoi3-country"
	DatabaseNetI3Identity Database = "neti3-identity"
)
