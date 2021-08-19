package apperrors

import "fmt"

type ExporterError struct {
	Type 	string
	Reason 	string
}

func (e *ExporterError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Reason)
}

func NewConfigError(reason string) *ExporterError {
	return &ExporterError{Type: "config_error", Reason: reason}
}

func NewRequestError(reason string) *ExporterError {
	return &ExporterError{Type: "request_error", Reason: reason}
}

func NewResponseFormatError(structure string) *ExporterError {
	return &ExporterError{Type: "response_format_error", Reason: fmt.Sprintf("can't parse %s", structure)}
}

type SelectelApiError struct {
	Code int
	Body string
}

func (e *SelectelApiError) Error() string {
	return fmt.Sprintf("selectel_api_error: status %d - %s", e.Code, e.Body)
}
