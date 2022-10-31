package test

import (
	"fmt"
	"net/http"
)

const (
	AuthzPath = "/data/rbac/manage/allow"
)

type imap map[string]interface{}

func AuthzCookie(tenant, subject string) []*http.Cookie {
	return []*http.Cookie{
		{
			Name:  "user",
			Value: fmt.Sprintf("%s / %s", tenant, subject),
		},
	}
}

func AuthzMock(tenant, subject string, allowed bool) *Mock {
	authzRequest := imap{
		"input": imap{
			"tenant":  tenant,
			"subject": subject,
		},
	}

	authzResponse := imap{
		"result": allowed,
	}

	code := 0
	switch allowed {
	case true:
		code = http.StatusOK
	case false:
		code = http.StatusForbidden
	}

	return &Mock{
		Checks: []CheckRequest{
			CheckRequestMethod(http.MethodPost),
			CheckRequestContentType(ApplicationJson),
			CheckRequestBody(authzRequest),
		},
		Response: DefaultResponse(code, authzResponse),
	}
}
