// +build !windows,!solaris

package reap

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestReap_IsSupported(t *testing.T) {
	if !IsSupported() {
		t.Fatalf("reap should be supported on %s", runtime.GOOS)
	}
}

func TestReap_ReapChildren(t *testing.T) {
	pids := make(PidCh, 1)
	errors := make(ErrorCh, 1)
	done := make(chan struct{}, 1)

	didExit := make(chan struct{}, 1)
	go func() {
		ReapChildren(pids, errors, done)
		didExit <- struct{}{}
	}()

	killAndCheck := func() {
		cmd := exec.Command("sleep", "5")
		if err := cmd.Start(); err != nil {
			t.Fatalf("err: %v", err)
		}

		childPid := cmd.Process.Pid
		if err := cmd.Process.Kill(); err != nil {
			t.Fatalf("err: %v", err)
		}

		select {
		case pid := <-pids:
			if pid != childPid {
				t.Fatalf("unexpected pid: %d != %d", pid, childPid)
			}
		case err := <-errors:
			t.Fatalf("err: %v", err)
		case <-time.After(1 * time.Second):
			t.Fatalf("should have reaped %d", childPid)
		}
	}

	// Kill a child process and make sure it gets detected.
	killAndCheck()

	// Fire off a subprocess.
	cmd := exec.Command("sleep", "5")
	if err := cmd.Start(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Send a spurious SIGCHLD.
	if err := unix.Kill(os.Getpid(), unix.SIGCHLD); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Make sure the reaper didn't report anything.
	select {
	case pid := <-pids:
		t.Fatalf("unexpected pid: %d", pid)
	case err := <-errors:
		t.Fatalf("err: %v", err)
	case <-time.After(1 * time.Second):
		// Good - nothing was sent to the channels.
	}

	// Now kill the child subprocess.
	childPid := cmd.Process.Pid
	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Make sure the reaper sees it.
	select {
	case pid := <-pids:
		if pid != childPid {
			t.Fatalf("unexpected pid: %d != %d", pid, childPid)
		}
	case err := <-errors:
		t.Fatalf("err: %v", err)
	case <-time.After(1 * time.Second):
		t.Fatalf("should have reaped %d", childPid)
	}

	// Run a few more cycles to make sure things work.
	killAndCheck()
	killAndCheck()
	killAndCheck()

	// Shut it down.
	close(done)
	select {
	case <-didExit:
		// Good - the goroutine shut down.
	case <-time.After(1 * time.Second):
		t.Fatalf("should have shut down")
	}
}

func TestReap_ReapOnce(t *testing.T) {
	killAndCheck := func() {
		cmd := exec.Command("sleep", "5")
		if err := cmd.Start(); err != nil {
			t.Fatalf("err: %v", err)
		}

		childPid := cmd.Process.Pid
		if err := cmd.Process.Kill(); err != nil {
			t.Fatalf("err: %v", err)
		}

		start := time.Now()
		for {
			pid, err := ReapOnce()
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			if pid != 0 && pid != childPid {
				t.Fatalf("unexpected pid: %d != %d", pid, childPid)
			}

			if pid == childPid {
				break
			}

			if time.Now().Sub(start) > time.Second {
				t.Fatalf("should have reaped %d", childPid)
			}
		}
	}

	// Run a few cycles to make sure things work.
	killAndCheck()
	killAndCheck()
	killAndCheck()
}
