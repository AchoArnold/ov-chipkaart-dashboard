package middlewares

import internalContext "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/context"

const (
	// ContextKeyUserID is the get to fetch the user id from the context
	ContextKeyUserID = internalContext.Key("user-id")

	// ContextKeyLanguageTag is the key for to get the language tag from the context
	ContextKeyLanguageTag = internalContext.Key("language-tag")
)

// Client provides the collection of middlewares.
type Client struct {
}

// New creates a middle ware client
func New() Client {
	return Client{}
}
