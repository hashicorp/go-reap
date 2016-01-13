// +build windows solaris

package reap

// IsSupported returns true if child process reaping is supported on this
// platform. This stub version always returns false.
func IsSupported() bool {
	return false
}

// ReapChildren is not supported so this always returns right away.
func ReapChildren(pids PidCh, errors ErrorCh, done chan struct{}) {
}

// ReapOnce is not supported so this always returns as if nothing was reaped.
func ReapOnce() (int, error) {
	return 0, nil
}
