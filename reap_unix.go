// +build !windows,!solaris

package reap

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

// IsSupported returns true if child process reaping is supported on this
// platform.
func IsSupported() bool {
	return true
}

// ReapChildren is a long-running routine that blocks waiting for child
// processes to exit and reaps them, reporting reaped process IDs to the
// optional pids channel and any errors to the optional errors channel.
//
// Be careful using this if your process uses Go's exec module, which waits
// for processes to complete internally. This may steal results from that and
// make exec think that the process doesn't exist.
func ReapChildren(pids PidCh, errors ErrorCh, done chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, unix.SIGCHLD)

	for {
	WAIT:
		// Block for an incoming signal that a child has exited.
		select {
		case <-c:
			// Got a child signal, drop out and reap.
		case <-done:
			return
		}

		// Try to reap children until there aren't any more. We never
		// block in here so that we are always responsive to signals, at
		// the expense of possibly leaving a child behind if we get
		// here too quickly. Any stragglers should get reaped the next
		// time we see a signal, so we won't leak in the long run.
	POLL:
		pid, err := ReapOnce()
		if err != nil {
			if errors != nil {
				errors <- err
			}
			goto WAIT
		}

		// Got a child, clean this up and poll again to clean up
		// any other child, if there is one.
		if pid > 0 {
			if pids != nil {
				pids <- pid
			}
			goto POLL
		}

		// No child found, wait for another signal.
		goto WAIT
	}
}

// ReapOnce looks for an unclaimed child process and reaps it if one is found.
// This will never block. If the returned pid > 0, then a child process was
// reaped and had the given pid.
func ReapOnce() (int, error) {
	var status unix.WaitStatus
	pid, err := unix.Wait4(-1, &status, unix.WNOHANG, nil)
	switch err {
	case unix.ECHILD:
		// No child was found.
		return 0, nil

	case unix.EINTR:
		// We got interrupted. This likely can't happen since we are
		// calling Wait4 in a non-blocking fashion, but it's good to be
		// complete and handle this case rather than fail.
		return 0, nil

	default:
		return pid, err
	}
}
