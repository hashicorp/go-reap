// +build windows

package reap

import (
	"runtime"
	"testing"
)

func TestReap_IsSupported(t *testing.T) {
	if IsSupported() {
		t.Fatalf("reap should not be supported on %s", runtime.GOOS)
	}
}

func TestReap_ReapChildren(t *testing.T) {
	pids := make(PidCh, 1)
	errors := make(ErrCh, 1)
	ReapChildren(pids, errors)
	select {
	case <-pids:
		t.Fatalf("should not report any pids")
	case <-errors:
		t.Fatalf("should not report any errors")
	default:
	}
}
