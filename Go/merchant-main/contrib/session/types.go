package session

import (
	"time"
	"github.com/valyala/fasthttp"
)

// Config config struct
type Config struct {

	// cookie name
	CookieName string

	// cookie domain
	Domain string

	// If you want to delete the cookie when the browser closes, set it to -1.
	//
	//  0 means no expire, (24 years)
	// -1 means when browser closes
	// >0 is the time.Duration which the session cookies should expire.
	Expires time.Duration

	// gc life time to execute it
	GCLifetime time.Duration

	// set whether to pass this bar cookie only through HTTPS
	Secure bool

	// IsSecureFunc should return whether the communication channel is secure
	// in order to set the secure flag to true according to Secure flag.
	IsSecureFunc func(*fasthttp.RequestCtx) bool

	// value cookie length
	//cookieLen uint32

	Addr string
	Auth string
}


// Cookie cookie struct
type Cookie struct{}

