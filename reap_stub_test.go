// +build windows solaris

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
	errors := make(ErrorCh, 1)
	ReapChildren(pids, errors, nil)
	select {
	case <-pids:
		t.Fatalf("should not report any pids")
	case <-errors:
		t.Fatalf("should not report any errors")
	default:
	}
}

func TestReap_ReapOnce(t *testing.T) {
	pid, err := ReapOnce()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if pid != 0 {
		f.Fatalf("bad: %d", pid)
	}
}
