package sdks

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SDKError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (m *SDKError) Error() string {
	return fmt.Sprintf("request failed with code %d: %s", m.Code, m.Message)
}

func errorHook(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}
	var errSDK SDKError
	err := json.NewDecoder(resp.Body).Decode(&errSDK)
	if err != nil {
		return &SDKError{Code: resp.StatusCode, Message: err.Error()}
	}
	errSDK.Code = resp.StatusCode
	return &errSDK
}
