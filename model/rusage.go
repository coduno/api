package model

import "syscall"

// Rusage is is a copy of syscall.Rusage for Linux.
// It is needed as syscall.Rusage built on Windows cannot be saved in Datastore
// because it holds unsigned integers. Runs are executed inside Docker and
// therefore will always generate this version of syscall.Rusage.
// See https://godoc.org/google.golang.org/cloud/datastore#Property
// See https://golang.org/src/syscall/syscall_windows.go
// See https://golang.org/src/syscall/ztypes_linux_amd64.go
type Rusage struct {
	Utime,
	Stime syscall.Timeval
	Maxrss,
	Ixrss,
	Idrss,
	Isrss,
	Minflt,
	Majflt,
	Nswap,
	Inblock,
	Oublock,
	Msgsnd,
	Msgrcv,
	Nsignals,
	Nvcsw,
	Nivcsw int64
}
