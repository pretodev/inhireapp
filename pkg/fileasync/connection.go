package fileasync

import (
	"io"
	"log"
	"sync"
	"time"
)

type Connection[T io.Closer] struct {
	timeout   time.Duration
	instance  T
	refCount  int
	closeChan chan struct{}
	mutex     sync.Mutex
}

func (conn *Connection[T]) manageClosing() {
	for range conn.closeChan {
		time.Sleep(conn.timeout) // espera um tempo antes de fechar
		conn.mutex.Lock()
		if conn.refCount == 0 {
			conn.Close()
			conn.mutex.Unlock()
			return
		}
		conn.mutex.Unlock()
	}
}

func (conn *Connection[T]) Close() error {
	log.Println("Conexão fechada por inatividade.")
	return conn.instance.Close()
}

func (conn *Connection[T]) Exec(exec func(r T) error) error {
	conn.Acquire()
	defer conn.Release()
	return exec(conn.instance)
}

func (conn *Connection[T]) Acquire() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.refCount++
}

func (conn *Connection[T]) Release() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.refCount--
	if conn.refCount == 0 {
		conn.closeChan <- struct{}{}
	}
}
