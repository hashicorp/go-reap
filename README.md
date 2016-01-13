# go-reap

Provides a super simple set of functions for reaping child processes. This is
useful for running applications as PID 1 in a Docker container.

This should be supported on most UNIX flavors, but is not supported on Windows
or Solaris. Unsupported platforms have a stub implementation that's safe to call,
as well as an API to check if reaping is supported so that you can produce an
error in your application code.

Use care with `ReapChildren` if you also use Go's `exec` functions which internally
wait for subprocesses to complete. This will steal the results away from them. If
that's the case, then you may need to use `ReapOnce` during a safe time when you know
that no waits are occurring.

Documentation
=============

The full documentation is available on [Godoc](http://godoc.org/github.com/hashicorp/go-reap).

Example
=======

Below is a simple example of usage

```go
// Reap children with no control or feedback.
go ReapChildren(nil, nil, nil)

// Get feedback on reaped children and errors.
if reap.IsSupported() {
	pids := make(reap.PidCh, 1)
	errors := make(reap.ErrorCh, 1)
	done := make(chan struct{})
	go ReapChildren(pids, errors, done)
	// ...
	close(done)
} else {
	fmt.Println("Sorry, go-reap isn't supported on your platform.")
}

// Poll for children to reap and reap them all.
for {
	pid, err := ReapOnce()
        if err != nil {
		panic(err)
	}
        if pid > 0 {
		fmt.Printf("Reaped child process %d\n", pid)
		continue
	}
        break
}
```

