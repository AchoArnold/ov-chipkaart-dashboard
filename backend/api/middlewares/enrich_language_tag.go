package middlewares

import (
	"context"
	"net/http"

	"golang.org/x/text/language"
)

// AddLanguageTag adds the language tag to the the context
func (middleware Client) AddLanguageTag() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			acceptLanguage := r.Header.Get("Accept-Language")

			matcher := language.NewMatcher([]language.Tag{
				language.English, // The first language is used as fallback.
				language.MustParse("en-NL"),
				language.Dutch,
			})

			tag, _ := language.MatchStrings(matcher, acceptLanguage)

			// put it in context
			ctx := context.WithValue(r.Context(), ContextKeyLanguageTag, &tag)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
