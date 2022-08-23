package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ApplicationJson = "application/json"
)

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func UnsupportedMediaType(w http.ResponseWriter) {
	http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
}

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "bad request", http.StatusBadRequest)
}

func AuthzCookie(tenant, subject string) []*http.Cookie {
	return []*http.Cookie{
		{
			Name:  "user",
			Value: fmt.Sprintf("%s / %s", tenant, subject),
		},
	}
}

func CheckResponseCode(code int) CheckResponse {
	return func(response *Response) error {
		if response.Rest.Code != code {
			return fmt.Errorf("response: expected code %d, got %d", code, response.Rest.Code)
		}

		return nil
	}
}

func CheckResponseBody(body interface{}) CheckResponse {
	return func(response *Response) error {
		if expected, err := json.Marshal(body); err != nil {
			return err
		} else if got, err := json.Marshal(response.Body); err != nil {
			return err
		} else if string(expected) != string(got) {
			return fmt.Errorf("response: expected body:\n\n%s\n\ngot:\n\n%s",
				string(expected), string(got),
			)
		}

		return nil
	}
}

func CheckRequestMethod(method string) CheckRequest {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != method {
			MethodNotAllowed(w)

			return fmt.Errorf("request: expected method %s, got %s", method, r.Method)
		}

		return nil
	}
}

func CheckRequestContentType(contentType string) CheckRequest {
	return func(w http.ResponseWriter, r *http.Request) error {
		if headers, ok := r.Header["Content-Type"]; !ok {
			UnsupportedMediaType(w)

			return fmt.Errorf("request: missing content type %s", contentType)
		} else {
			for _, header := range headers {
				if header == contentType {
					return nil
				}
			}

			UnsupportedMediaType(w)

			return fmt.Errorf("request: missing content type %s", contentType)
		}
	}
}

func CheckRequestBody(body interface{}) CheckRequest {
	return func(w http.ResponseWriter, r *http.Request) error {
		responseData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			InternalServerError(w)
			return err
		}

		var responseBody interface{}
		if err := json.Unmarshal(responseData, &responseBody); err != nil {
			BadRequest(w)
			return err
		}

		if expected, err := json.Marshal(body); err != nil {
			InternalServerError(w)
			return err
		} else if got, err := json.Marshal(responseBody); err != nil {
			InternalServerError(w)
			return err
		} else if string(expected) != string(got) {
			BadRequest(w)

			return fmt.Errorf("request: expected body:\n\n%s\n\ngot:\n\n%s",
				string(expected), string(got),
			)
		}

		return nil
	}
}
