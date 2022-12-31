package validator

import (
	"testing"
)

func TestCheckStringCommaDigit(t *testing.T) {

	s := "1,2,3,4,5,6,7,8,9"
	if !CheckStringCommaDigit(s) {
		t.Error("error")
	}
}

func TestCheckStringCommaAlpha(t *testing.T) {

	s := "a,b,c,d,e,f,g,h,i"
	if !CheckStringCommaAlpha(s) {
		t.Error("error")
	}
}

func TestCheckStringAlnum(t *testing.T) {

	s := "a1bcdefghi sdasd12432423324---"
	if !CheckStringAlnum(s) {
		t.Error("error")
	}
}

func TestIsVietnamesePhone(t *testing.T) {

	s := "0331019073"
	if !IsVietnamesePhone(s) {
		t.Error("invalid phone")
	}

	t.Log("valid phone")
}
