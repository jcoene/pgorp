package pgorp

import (
	"reflect"
	"testing"
)

func TestPgBuildIntArray(t *testing.T) {
	expect := "{1,2,3,4,5,6,1234567890}"
	data := ArrayInt64{1, 2, 3, 4, 5, 6, 1234567890}
	got := pgBuildIntArray(data)

	if got != expect {
		t.Errorf("expected '%s', got '%s'", expect, got)
	}
}

func TestPgParseIntArray(t *testing.T) {
	expect := ArrayInt64{1, 2, 3, 4, 5, 6, 1234567890}
	data := "{1,2,3,4,5,6,1234567890}"
	got := pgParseIntArray(data)

	if !reflect.DeepEqual(got, expect) {
		t.Errorf("expected %+v, got %+v", expect, got)
	}
}
