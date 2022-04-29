package portscan

import "syscall"

// GetMaxFD returns the maximum number of open file descriptors.
func GetMaxFD() (int, error) {

	rlimit := syscall.Rlimit{}

	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		return 0, err
	}

	return int(rlimit.Cur), nil
}
