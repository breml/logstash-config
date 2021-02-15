package ast_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/breml/logstash-config/ast"
)

func TestNotFoundError(t *testing.T) {
	cases := []struct {
		formatVerb string

		wantErrorString string
	}{
		{
			formatVerb: "%s",

			wantErrorString: "not found: error",
		},
		{
			formatVerb: "%v",

			wantErrorString: "not found: error",
		},
		{
			formatVerb: "%+v",

			wantErrorString: "not found: error",
		},
		{
			formatVerb: "%q",

			wantErrorString: `"not found: error"`,
		},
	}

	for _, test := range cases {
		t.Run(test.formatVerb, func(t *testing.T) {
			err := ast.NotFoundErrorf("%s", "error")

			if err == nil {
				t.Fatal("Expect an error, but got none.")
			}

			if test.wantErrorString != fmt.Sprintf(test.formatVerb, err) {
				t.Errorf("Expect error to print %q with verb %q, but got: %q", test.wantErrorString, test.formatVerb, fmt.Sprintf(test.formatVerb, err))
			}

			if !ast.IsNotFoundError(err) {
				t.Fatalf("Expect err %v to implement NotFounder interface, but it does not.", err)
			}
		})
	}
}

func TestNotFoundErrorNotFound(t *testing.T) {
	ast.NewNotFoundError(errors.New("error")).(ast.NotFounder).NotFound()
}

func TestIsNotFoundError(t *testing.T) {
	cases := []struct {
		name string
		err  error

		want bool
	}{
		{
			name: "nil",
		},
		{
			name: "nil not founder",
			err:  ast.NewNotFoundError(nil),
		},
		{
			name: "error",
			err:  ast.NewNotFoundError(errors.New("error")),
			want: true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := ast.IsNotFoundError(test.err)
			if test.want != got {
				t.Fatalf("Expectation (%v) not met: %v", test.want, got)
			}
		})
	}
}
