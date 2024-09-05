package fileasync

import (
	"io"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type DatasourceBuilder[T io.Closer] func(path string) (T, error)

type Manager[T io.Closer] struct {
	conns map[string]*Connection[T]
	mutex sync.Mutex
}

func NewConnectionManager[T io.Closer]() *Manager[T] {
	return &Manager[T]{
		conns: make(map[string]*Connection[T]),
	}
}

type poolOptions struct {
	Timeout time.Duration
}

type poolOption func(*poolOptions)

func WithTimeout(d time.Duration) poolOption {
	return func(opts *poolOptions) {
		opts.Timeout = d
	}
}

func (m *Manager[T]) NewPool(ctx context.Context, path string, dtb DatasourceBuilder[T], opts ...poolOption) (*Connection[T], error) {
	options := &poolOptions{
		Timeout: 30 * time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if conn, exists := m.conns[path]; exists {
		return conn, nil
	}

	instance, err := dtb(path)
	if err != nil {
		return nil, err
	}

	conn := &Connection[T]{
		instance:  instance,
		refCount:  0,
		closeChan: make(chan struct{}),
		timeout:   options.Timeout,
	}

	m.conns[path] = conn
	go conn.manageClosing()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	return conn, nil
}
