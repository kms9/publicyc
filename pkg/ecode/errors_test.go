package ecode

import (
	"errors"
	"fmt"
	"testing"

	pkgErr "github.com/pkg/errors"
)

func TestErrorsMatch(t *testing.T) {
	s := &StatusError{Code: 1}
	st := &StatusError{Code: 2}
	if errors.Is(s, st) {
		t.Errorf("errs is not match: %+v -> %+v", s, st)
	}
	s.Code = 1
	st.Code = 1
	if !errors.Is(s, st) {
		t.Errorf("errs is not match: %+v -> %+v", s, st)
	}
	s.WithDebugs(&ErrorInfo{Reason: "debug1"})
	st.WithDebugs(&ErrorInfo{Reason: "debug2"})
	if !errors.Is(s, st) {
		t.Errorf("errs is not match: %+v -> %+v", s, st)
	}

	if Reason(s).Reason != UnknownReason {
		t.Errorf("errs is not match: %+v -> %+v", s, st)
	}
	s.ErrorInfo = &ErrorInfo{Reason: "test_reason"}
	if Reason(s).Reason != "test_reason" {
		t.Errorf("errs is not match: %+v -> %+v", s, st)
	}
}

func TestErrorIs(t *testing.T) {
	err1 := errors.New("client timeout")
	t.Log(err1)
	err2 := pkgErr.Wrap(err1, "this is err 2")
	t.Log(err2)
	t.Log(pkgErr.Cause(err2))
	// 打印堆栈
	t.Log(fmt.Printf("%+v", err2))
	if !(pkgErr.Is(err2, err1)) {
		t.Errorf("errs is not match: a: %v b: %v ", err2, err1)
	}
}

func TestErrorAs(t *testing.T) {
	err1 := &StatusError{Code: 1}
	err2 := pkgErr.Wrap(err1, "this is err 2")
	err3 := new(StatusError)
	if !errors.As(err2, &err3) {
		t.Errorf("errs is not match: %v", err2)
	}
}

func TestHTTPCodeAndReason(t *testing.T) {
	err1 := InvalidArgument("required", "userId err", "")
	code, reason := HTTPCodeAndReason(err1)
	t.Log(code, reason)
	if code != 400 || reason.Reason != "required" {
		t.Errorf("errs code reason err %v", err1)
	}
}
