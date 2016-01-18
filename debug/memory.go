package debug

import (
	"os"
	"bufio"
	"bytes"
	"fmt"
)

// MemoryUsage Simple memory debugging base on http://stackoverflow.com/a/31881979
// will keep file open so when calling repeatedly wont need to close/open the file
// for every call, Should be able to get memory size of every pid this process has
// access to.
type MemoryUsage struct {
	// fd of open file
	file 	*os.File
	// prefix to search for, defaults:
	// []byte("Pss:") // (Proportional Set Size)
	prefix	[]byte
	// Result of last read
	result  uint64
}

// Close will close the opend file
func (m *MemoryUsage) Close() {
	m.file.Close()
}

// LastResult will return last queried result
func (m *MemoryUsage) LastResult() uint64 {
	return m.result
}

// GetSize returns current memory usage
func (m *MemoryUsage) GetSize() (uint64, error) {
	m.result = uint64(0)
	m.file.Seek(0,0)
	scanner := bufio.NewScanner(m.file)

	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, m.prefix) {
			var size uint64
			_, err := fmt.Sscanf(string(line[4:]), "%d", &size)
			if err != nil {
				return 0, err
			}
			m.result += size
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return m.result, nil
}

// NewSelfMemoryUsage creates a new MemoryUsage and monitoring
// memory usage of self by opening /proc/self/smaps
func NewSelfMemoryUsage() (*MemoryUsage, error) {
	fd, err := os.Open("/proc/self/smaps")
	if err != nil {
		return nil, err
	}
	return &MemoryUsage{file: fd, prefix: []byte("Pss:")}, nil
}

// NewPidMemoryUsage creates a new MemoryUsage and monitoring
// memory usage of given pid by opening /proc/<pind>/smaps
func NewPidMemoryUsage(pid int) (*MemoryUsage, error) {
	fd, err := os.Open(fmt.Sprintf("/proc/%d/smaps", pid))
	if err != nil {
		return nil, err
	}
	return &MemoryUsage{file: fd, prefix: []byte("Pss:")}, nil
}