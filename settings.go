package jayson

import (
	"net/http"
)

// DefaultSettings returns default settings for jayson instance
func DefaultSettings() Settings {
	return Settings{
		DefaultErrorStatus:        http.StatusInternalServerError,
		DefaultErrorDetailKey:     "errorDetail",
		DefaultErrorMessageKey:    "message",
		DefaultErrorStatusCodeKey: "code",
		DefaultErrorStatusTextKey: "status",
		DefaultResponseStatus:     http.StatusOK,
	}
}

// Settings for jayson instance
type Settings struct {
	DefaultErrorStatus        int
	DefaultErrorDetailKey     string
	DefaultErrorMessageKey    string
	DefaultErrorStatusCodeKey string
	DefaultErrorStatusTextKey string
	DefaultResponseStatus     int
}

func (s *Settings) Validate() {
	if s.DefaultErrorStatus == 0 {
		s.DefaultErrorStatus = http.StatusInternalServerError
	}
	if s.DefaultErrorDetailKey == "" {
		s.DefaultErrorDetailKey = "errorDetail"
	}
	if s.DefaultErrorMessageKey == "" {
		s.DefaultErrorMessageKey = "message"
	}
	if s.DefaultErrorStatusCodeKey == "" {
		s.DefaultErrorStatusCodeKey = "code"
	}
	if s.DefaultErrorStatusTextKey == "" {
		s.DefaultErrorStatusTextKey = "status"
	}
	if s.DefaultResponseStatus == 0 {
		s.DefaultResponseStatus = http.StatusOK
	}
}
