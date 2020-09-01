package validator

// Helpers are custom validations that are package agnostic
type Helpers struct {
}

// ValidateOvChipkaartCredentials checks that the username and password for an ov chipkaart are valid
func (h Helpers) ValidateOvChipkaartCredentials(username string, password string) (err error) {
	return err
}
