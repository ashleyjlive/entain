package db_test

import (
	"regexp"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestOrderBySingle(t *testing.T) {
	out, err := ValidateOrderBy("foo asc")
	want := regexp.MustCompile("foo")
	if !want.MatchString(out[0]) || err != nil {
		t.Fatalf(`validateOrderBy("foo asc") = %q, %v, want match for %#q, nil`, out, err, want)
	}
}
