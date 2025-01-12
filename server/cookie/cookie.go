package cookie

import (
	"net/http"
	"net/url"

	"github.com/authorizerdev/authorizer/server/constants"
	"github.com/authorizerdev/authorizer/server/parsers"
	"github.com/gin-gonic/gin"
)

// SetSession sets the session cookie in the response
func SetSession(gc *gin.Context, sessionID string) {
	secure := true
	httpOnly := true
	hostname := parsers.GetHost(gc)
	host, _ := parsers.GetHostParts(hostname)
	domain := parsers.GetDomainName(hostname)
	if domain != "localhost" {
		domain = "." + domain
	}

	// TODO allow configuring from dashboard
	year := 60 * 60 * 24 * 365

	gc.SetSameSite(http.SameSiteNoneMode)
	gc.SetCookie(constants.AppCookieName+"_session", sessionID, year, "/", host, secure, httpOnly)
	gc.SetCookie(constants.AppCookieName+"_session_domain", sessionID, year, "/", domain, secure, httpOnly)
}

// DeleteSession sets session cookies to expire
func DeleteSession(gc *gin.Context) {
	secure := true
	httpOnly := true
	hostname := parsers.GetHost(gc)
	host, _ := parsers.GetHostParts(hostname)
	domain := parsers.GetDomainName(hostname)
	if domain != "localhost" {
		domain = "." + domain
	}

	gc.SetSameSite(http.SameSiteNoneMode)
	gc.SetCookie(constants.AppCookieName+"_session", "", -1, "/", host, secure, httpOnly)
	gc.SetCookie(constants.AppCookieName+"_session_domain", "", -1, "/", domain, secure, httpOnly)
}

// GetSession gets the session cookie from context
func GetSession(gc *gin.Context) (string, error) {
	var cookie *http.Cookie
	var err error
	cookie, err = gc.Request.Cookie(constants.AppCookieName + "_session")
	if err != nil {
		cookie, err = gc.Request.Cookie(constants.AppCookieName + "_session_domain")
		if err != nil {
			return "", err
		}
	}

	decodedValue, err := url.PathUnescape(cookie.Value)
	if err != nil {
		return "", err
	}
	return decodedValue, nil
}
