package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

// Add adds `func() error` callback to the globalCloser
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait is a function which waits until all closer functions are executed
func Wait() {
	globalCloser.Wait()
}

// CloseAll is a function that calls all closer functions
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer is an object that has methods for adding and calling closer functions
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New returns new Closer, if []os.Signal is specified Closer will automatically call CloseAll when one of signals is received from OS
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

// Add func to closer
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait blocks until all closer functions are done
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll calls all closer functions
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		// call all Closer funcs async
		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer")
			}
		}
	})
}
