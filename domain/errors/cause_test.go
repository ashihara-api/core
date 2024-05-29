package errors_test

import (
	er "errors"
	"fmt"
	"testing"

	"github.com/ashihara-api/core/domain/errors"
)

func TestCause_Append(t *testing.T) {
	type fields struct {
		err    error
		domain string
		c      errors.ErrCase
	}
	type args struct {
		d string
		e error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		isCause bool
	}{
		{
			name: "append_cause",
			fields: fields{
				err:    er.New("unknown error"),
				domain: "foobar",
				c:      errors.CaseBackendError,
			},
			args: args{
				d: errors.ServiceDomainGlobal,
				e: errors.NewCause(
					er.New("maintainance"),
					errors.ServiceDomainGlobal,
					errors.CaseUnavailable,
				),
			},
			isCause: true,
		},
		{
			name: "append_non_cause",
			fields: fields{
				err:    er.New("unknown error"),
				domain: "foobar",
				c:      errors.CaseBackendError,
			},
			args: args{
				d: errors.ServiceDomainGlobal,
				e: er.New("test error"),
			},
			isCause: true,
		},
		{
			name: "append_non_cause_with_custom_domain",
			fields: fields{
				err:    er.New("unknown error"),
				domain: "foobar",
				c:      errors.CaseBackendError,
			},
			args: args{
				d: "fizzbizz",
				e: er.New("test error"),
			},
			isCause: true,
		},
	}
	for _, tt := range tests {
		errors.ServiceDomain = tt.args.d
		t.Run(tt.name, func(t *testing.T) {
			e := errors.NewCause(
				tt.fields.err,
				tt.fields.domain,
				tt.fields.c,
			)
			var c *errors.Cause
			ok := er.As(e, &c)
			if ok != tt.isCause {
				t.Errorf("NewCause() ok = %v, isCause %v", ok, tt.isCause)
			}
			if ok {
				c.Append(tt.args.e)
				// if *c != *tt.want {
				// 	t.Errorf("NewCause() got %#v, want %#v", c, tt.want)
				// }
				// spew.Dump(c)
				fmt.Printf("%#v\n", c)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	s := `{
	"error": {
		"code": 500,
		"status": "INTERNAL",
		"message": "missing destination name source in *[]mysql.product",
		"details": [
		{
			"Domain": "products",
			"Reason": "backendError",
			"Message": "missing destination name source in *[]mysql.product"
		}
		]
	}
}`
	e := errors.NewCause(
		er.New("unknown error"),
		"foobar",
		errors.CaseBackendError,
	)
	var c *errors.Cause
	er.As(e, &c)
	if err := c.UnmarshalJSON([]byte(s)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
