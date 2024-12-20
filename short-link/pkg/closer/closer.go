package closer

import (
	"context"
	"fmt"
	"sync"
)

// Closer - struct for graceful closing all services.
type Closer struct {
	mx         sync.Mutex
	closeFuncs []CloseFunc
}

// New - constructor for Closer.
func New() *Closer {
	return &Closer{
		mx:         sync.Mutex{},
		closeFuncs: make([]CloseFunc, 0),
	}
}

type CloseFunc func(ctx context.Context) error

// Add function adds closer func for each service, that closer will invoke in Close func.
// If you want to shut down your service gracefully, you need to add closer func for your service.
// It's not garantee that your service will be closed gracefully, if your service will be closing
// longer than context timeout from Close func param.
func (c *Closer) Add(f CloseFunc) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.closeFuncs = append(c.closeFuncs, f)
}

// Close function calls all services CloseFunc in parallel, that have been added by Add func.
func (c *Closer) Close(ctx context.Context) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	var (
		msgs     = newMessages(len(c.closeFuncs))
		complete = make(chan struct{}, 1)
		wg       = &sync.WaitGroup{}
	)

	wg.Add(len(c.closeFuncs))
	go func() {
		for _, f := range c.closeFuncs {
			f := f
			go func() {
				if err := f(ctx); err != nil {
					msgs.add(fmt.Sprintf("[!] %v", err))
				}
				wg.Done()
			}()
		}
		wg.Wait()
		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("graceful shutdown cancelled: %v", ctx.Err())
	}

	if msgs.len() > 0 {
		return fmt.Errorf(
			"shutdown finished with error(s): \n%s",
			msgs.String(),
		)
	}

	return nil
}
