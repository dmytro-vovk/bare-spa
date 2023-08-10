//go:build sysboard

package memory

import (
	"github.com/Sergii-Kirichok/pr/pkg/binary"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

// memMutex make work with memory safe
var memMutex = &sync.Mutex{}

type Memory struct {
	file *os.File
	mem  []byte
}

func At(baseAddr, buffSize int) Memory {
	f, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0666) // todo: mb 600
	if err != nil {
		panic(err)
	}

	mem, err := syscall.Mmap(int(f.Fd()), int64(baseAddr), buffSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	memMutex.Lock()
	return Memory{
		file: f,
		mem:  mem,
	}
}

func (m *Memory) ReadWord(offset binary.Bits) binary.Bits {
	return *(*binary.Bits)(unsafe.Pointer(&m.mem[offset]))
}

func (m *Memory) WriteWord(offset, value binary.Bits) {
	*(*binary.Bits)(unsafe.Pointer(&m.mem[offset])) = value
}

func (m *Memory) Close() {
	if err := syscall.Munmap(m.mem); err != nil {
		panic(err)
	}

	if err := m.file.Close(); err != nil {
		panic(err)
	}

	m.file = nil
	m.mem = nil
	memMutex.Unlock()
}
