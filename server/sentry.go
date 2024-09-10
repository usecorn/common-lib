package server

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// AddTagToRootScope adds tags to the root scope of the current request's Sentry hub.
func AddTagToRootScope(c *gin.Context, tags map[string]string) {
	hub := sentrygin.GetHubFromContext(c)
	if hub == nil {
		return
	}
	for k, v := range tags {
		hub.Scope().SetTag(k, v)
	}
}
