package util

import (
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	// listenFdsStart corresponds to `SD_LISTEN_FDS_START`.
	listenFdsStart = 3
)

var ErrSocketNotFound = errors.New("socket not found")
var socketFiles []*os.File

// socketFiles returns a slice containing a `os.File` object for each file
// descriptor passed to this process via systemd fd-passing protocol. I stole
// this code from the go-systemd project
func getSocketFiles() []*os.File {
	if socketFiles != nil {
		return socketFiles
	}

	pid, err := strconv.Atoi(os.Getenv("LISTEN_PID"))
	if err != nil || pid != os.Getpid() {
		return nil
	}

	nfds, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
	if err != nil || nfds == 0 {
		return nil
	}

	names := strings.Split(os.Getenv("LISTEN_FDNAMES"), ":")

	files := make([]*os.File, 0, nfds)
	for fd := listenFdsStart; fd < listenFdsStart+nfds; fd++ {
		syscall.CloseOnExec(fd)
		name := "LISTEN_FD_" + strconv.Itoa(fd)
		offset := fd - listenFdsStart
		if offset < len(names) && len(names[offset]) > 0 {
			name = names[offset]
		}
		files = append(files, os.NewFile(uintptr(fd), name))
	}

	socketFiles = files

	return files
}

// SystemdSocketByName returns a net.Listener if there is a systemd socket with
// that name available
func SystemdListenerByName(name string) (net.Listener, error) {
	for _, f := range getSocketFiles() {
		if f.Name() == name {
			return net.FileListener(f)
		}
	}
	return nil, ErrSocketNotFound
}

// SystemdFileByName returns a *os.File if there is a systemd socket with that
// name available
func SystemdFileByName(name string) (*os.File, error) {
	for _, f := range getSocketFiles() {
		if f.Name() == name {
			return f, nil
		}
	}
	return nil, ErrSocketNotFound
}
