package boot

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Container struct {
	items      sync.Map
	shutdownFn []shutdownFn
	m          sync.Mutex
}

type shutdownFn struct {
	name string
	fn   func()
}

func NewContainer() (container *Container, shutdown func()) {
	container = &Container{}
	container.arm(syscall.SIGINT, syscall.SIGTERM)

	return container, container.shutdown
}

func (c *Container) arm(signals ...os.Signal) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, signals...)

	go func() {
		s := <-sc

		log.Printf("Got %v, shutting down...", s)

		c.shutdown()

		log.Printf("Shutdown complete")

		os.Exit(0)
	}()
}

func (c *Container) Set(name string, item interface{}, fn func()) *Container {
	c.items.Store(name, item)

	if fn != nil {
		c.shutdownFn = append(c.shutdownFn, shutdownFn{
			name: name,
			fn:   fn,
		})
	}

	log.Printf("Initialised %s", name)

	return c
}

func (c *Container) Get(name string) interface{} {
	if it, ok := c.items.Load(name); ok {
		return it
	}

	return nil
}

func (c *Container) shutdown() {
	c.m.Lock()

	for i := len(c.shutdownFn) - 1; i >= 0; i-- {
		log.Printf("Shutting down %s...", c.shutdownFn[i].name)

		c.shutdownFn[i].fn()

		log.Printf("Shutting down %s complete", c.shutdownFn[i].name)
	}

	c.shutdownFn = c.shutdownFn[0:0]

	c.m.Unlock()
}
