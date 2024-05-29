package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	// ServiceDomainGlobal ...
	ServiceDomainGlobal = "ashihara"
)

// Status
const (
	// StatusBadRequest ...
	StatusBadRequest = "INVALID_ARGUMENT"
	// StatusUnauthenticated ...
	StatusUnauthenticated = "UNAUTHENTICATED"
	// StatusPermissionDenied ...
	StatusPermissionDenied = "PERMISSION_DENIED"
	// StatusNotFound ...
	StatusNotFound = "NOT_FOUND"
	// StatusAborted ...
	StatusAborted = "ABORTED"
	// StatusAlreadyExists ...
	StatusAlreadyExists = "ALREADY_EXISTS"
	// StatusResourceExhausted ...
	StatusResourceExhausted = "RESOURCE_EXHAUSTED"
	// StatusUnavailable ...
	StatusUnavailable = "UNAVAILABLE"
	// StatusBackendError ...
	StatusBackendError = "INTERNAL"
)

// Reason
const (
	// ReasonBadRequest ...
	ReasonBadRequest = "badRequest"
	// ReasonUnauthenticated ...
	ReasonUnauthenticated = "unauthenticated"
	// ReasonPermissionDenied ...
	ReasonPermissionDenied = "permissionDenied"
	// ReasonNotFound ...
	ReasonNotFound = "notFound"
	// ReasonAborted ...
	ReasonAborted = "abourtedRequest"
	// ReasonAlreadyExists ...
	ReasonAlreadyExists = "alreadyExists"
	// ReasonResourceExhausted ...
	ReasonResourceExhausted = "userRateLimitExceeded"
	// ReasonUnavailable ...
	ReasonUnavailable = "unavailable"
	// ReasonBackendError ...
	ReasonBackendError = "backendError"
)

var (
	// ServiceDomain is default service domain name
	ServiceDomain = ServiceDomainGlobal

	// CaseBadRequest ...
	CaseBadRequest = ErrCase{
		Code:   http.StatusBadRequest,
		Status: StatusBadRequest,
		Reason: ReasonBadRequest,
	}

	// CaseUnauthenticated ...
	CaseUnauthenticated = ErrCase{
		Code:   http.StatusUnauthorized,
		Status: StatusUnauthenticated,
		Reason: ReasonUnauthenticated,
	}

	// CasePermissionDenied ...
	CasePermissionDenied = ErrCase{
		Code:   http.StatusForbidden,
		Status: StatusPermissionDenied,
		Reason: ReasonPermissionDenied,
	}

	// CaseNotFound ...
	CaseNotFound = ErrCase{
		Code:   http.StatusNotFound,
		Status: StatusNotFound,
		Reason: ReasonNotFound,
	}

	// CaseAborted ...
	CaseAborted = ErrCase{
		Code:   http.StatusConflict,
		Status: StatusAborted,
		Reason: ReasonAborted,
	}

	// CaseAlreadyExists ...
	CaseAlreadyExists = ErrCase{
		Code:   http.StatusConflict,
		Status: StatusBackendError,
		Reason: ReasonBackendError,
	}

	// CaseResourceExhausted ...
	CaseResourceExhausted = ErrCase{
		Code:   http.StatusTooManyRequests,
		Status: StatusResourceExhausted,
		Reason: ReasonResourceExhausted,
	}

	// CaseUnavailable ...
	CaseUnavailable = ErrCase{
		Code:   http.StatusServiceUnavailable,
		Status: StatusUnavailable,
		Reason: ReasonUnavailable,
	}

	// CaseBackendError ...
	CaseBackendError = ErrCase{
		Code:   http.StatusInternalServerError,
		Status: StatusBackendError,
		Reason: ReasonBackendError,
	}
)

// Cause ...
type Cause struct {
	code    int
	domain  string
	status  string
	message error
	details []Detail
}

// Detail ...
type Detail struct {
	Domain  string
	Reason  string
	Message string
}

// ErrCase ...
type ErrCase struct {
	Code   int
	Status string
	Reason string
}

// jsonResponse ...
type jsonResponse struct {
	Error jsonResponseElement `json:"error"`
}

// jsonResponseElement
type jsonResponseElement struct {
	Code    int      `json:"code"`
	Domain  string   `json:"domain"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Details []Detail `json:"details"`
}

// NewCause ...
func NewCause(err error, domain string, c ErrCase) error {
	return NewCauseWithStatus(err, domain, c.Code, c.Status, c.Reason)
}

// NewCauseWithStatus ...
func NewCauseWithStatus(err error, domain string, code int, status string, reason string) error {
	return &Cause{
		code:    code,
		domain:  domain,
		status:  status,
		message: err,
		details: []Detail{
			{
				Domain:  domain,
				Reason:  reason,
				Message: err.Error(),
			},
		},
	}
}

// Error return error message
func (c *Cause) Error() string {
	reason := ""
	if len(c.details) > 0 {
		reason = c.details[0].Reason
	}
	return fmt.Errorf("%s : %w", reason, c.message).Error()
}

// Unwrap implements errors.Unwrap method
func (c *Cause) Unwrap() error {
	return c.message
}

// Append one or more elements onto the end of details
func (c *Cause) Append(e error) {
	if c == nil || c.IsZero() {
		c.set(e)
		return
	}

	var v *Cause
	if errors.As(e, &v) {
		c.details = append(c.details, v.details...)
		return
	}
	c.details = append(c.details, Detail{
		Domain:  c.domain,
		Reason:  StatusBackendError,
		Message: e.Error(),
	})
}

// MarshalJSON implements the json.Marshaler interface
func (c *Cause) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonResponse{
		Error: jsonResponseElement{
			Code:    c.code,
			Domain:  c.domain,
			Status:  c.status,
			Message: c.message.Error(),
			Details: c.details,
		},
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (c *Cause) UnmarshalJSON(b []byte) (err error) {
	var je jsonResponse
	if err = json.Unmarshal(b, &je); err != nil {
		return
	}
	ce := &Cause{
		code:    je.Error.Code,
		domain:  je.Error.Domain,
		status:  je.Error.Status,
		message: errors.New(je.Error.Message),
		details: je.Error.Details,
	}
	c.Append(ce)
	return nil
}

// IsZero checks empty
func (c *Cause) IsZero() bool {
	return (c.code == 0 || c.status == "" || c.message == nil || len(c.details) == 0)
}

// set overwrites error
func (c *Cause) set(e error) {
	var v *Cause
	if !errors.As(e, &v) {
		errors.As(NewCause(e, c.domain, CaseBackendError), &v)
	}

	// c = v <== fail staticcheck. SA4006: this value of `c` is never used
	c.code = v.code
	c.domain = v.domain
	c.status = v.status
	c.message = v.message
	c.details = v.details
}

// Code get http status code
func (c *Cause) Code() int {
	return c.code
}
