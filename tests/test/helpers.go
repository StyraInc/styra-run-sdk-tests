package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

const (
	ContentTypeHeader = "Content-Type"
	ApplicationJson   = "application/json"
)

var (
	DefaultContentTypeHeader = map[string]string{
		ContentTypeHeader: ApplicationJson,
	}
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
		expectedBytes, expectedValue, err := remarshal(body)
		if err != nil {
			return err
		}

		gotBytes, gotValue, err := remarshal(response.Body)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(expectedValue, gotValue) {
			return fmt.Errorf("response: expected:\n\n%s\n\ngot:\n\n%s",
				string(expectedBytes), string(gotBytes),
			)
		}

		return nil
	}
}

func CheckRequestMethod(method string) CheckRequest {
	return func(w http.ResponseWriter, r *MockRequest) error {
		if r.Request.Method != method {
			MethodNotAllowed(w)

			return fmt.Errorf("request: expected method %s, got %s", method, r.Request.Method)
		}

		return nil
	}
}

func CheckRequestContentType(contentType string) CheckRequest {
	return func(w http.ResponseWriter, r *MockRequest) error {
		if headers, ok := r.Request.Header[ContentTypeHeader]; !ok {
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
	return func(w http.ResponseWriter, r *MockRequest) error {
		expectedBytes, expectedValue, err := remarshal(body)
		if err != nil {
			InternalServerError(w)
			return err
		}

		gotBytes := r.Body

		var gotValue interface{}
		err = json.Unmarshal(gotBytes, &gotValue)
		if err != nil {
			BadRequest(w)
			return err
		}

		if !reflect.DeepEqual(expectedValue, gotValue) {
			BadRequest(w)

			return fmt.Errorf("request: expected:\n\n%s\n\ngot:\n\n%s",
				string(expectedBytes), string(gotBytes),
			)
		}

		return nil
	}
}

func DefaultResponse(code int, body interface{}) EmitResponse {
	return func(w http.ResponseWriter, r *MockRequest) error {
		if bytes, err := json.Marshal(body); err != nil {
			InternalServerError(w)

			return err
		} else {
			w.WriteHeader(code)
			w.Header().Set(ContentTypeHeader, ApplicationJson)

			if _, err := w.Write(bytes); err != nil {
				InternalServerError(w)

				return err
			}
		}

		return nil
	}
}

// To use reflect.DeepEqual safely, both values must be of the
// same nested types. That is, even though a struct and a map
// of keys may be "the same", they won't be when comparing. The
// solution is to serialize to json, then serialize back out.
func remarshal(x interface{}) ([]byte, interface{}, error) {
	bytes, err := json.Marshal(x)
	if err != nil {
		return nil, nil, err
	}

	var value interface{}
	if err := json.Unmarshal(bytes, &value); err != nil {
		return bytes, nil, err
	}

	return bytes, value, nil
}
