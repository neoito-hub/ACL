// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: proxy/spaces_proxy.proto

package captain

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on SpacesRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SpacesRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SpacesRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SpacesRequestMultiError, or
// nil if none found.
func (m *SpacesRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *SpacesRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Body

	// no validation rules for Queryparams

	if len(errors) > 0 {
		return SpacesRequestMultiError(errors)
	}

	return nil
}

// SpacesRequestMultiError is an error wrapping multiple validation errors
// returned by SpacesRequest.ValidateAll() if the designated constraints
// aren't met.
type SpacesRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SpacesRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SpacesRequestMultiError) AllErrors() []error { return m }

// SpacesRequestValidationError is the validation error returned by
// SpacesRequest.Validate if the designated constraints aren't met.
type SpacesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SpacesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SpacesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SpacesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SpacesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SpacesRequestValidationError) ErrorName() string { return "SpacesRequestValidationError" }

// Error satisfies the builtin error interface
func (e SpacesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSpacesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SpacesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SpacesRequestValidationError{}

// Validate checks the field values on SpacesReply with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SpacesReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SpacesReply with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SpacesReplyMultiError, or
// nil if none found.
func (m *SpacesReply) ValidateAll() error {
	return m.validate(true)
}

func (m *SpacesReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Err

	// no validation rules for Msg

	// no validation rules for Data

	// no validation rules for Status

	if len(errors) > 0 {
		return SpacesReplyMultiError(errors)
	}

	return nil
}

// SpacesReplyMultiError is an error wrapping multiple validation errors
// returned by SpacesReply.ValidateAll() if the designated constraints aren't met.
type SpacesReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SpacesReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SpacesReplyMultiError) AllErrors() []error { return m }

// SpacesReplyValidationError is the validation error returned by
// SpacesReply.Validate if the designated constraints aren't met.
type SpacesReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SpacesReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SpacesReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SpacesReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SpacesReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SpacesReplyValidationError) ErrorName() string { return "SpacesReplyValidationError" }

// Error satisfies the builtin error interface
func (e SpacesReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSpacesReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SpacesReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SpacesReplyValidationError{}