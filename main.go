package main

import (
	"errors"
	"fmt"
	"time"
)

type Semaphore interface {
	Acquire() error
	Release() error
	Len() int
}

type semaphore struct {
	ch chan struct{}
}

func NewSemaphore(count int) Semaphore {
	return &semaphore{ch: make(chan struct{}, count)}
}

func (s *semaphore) Acquire() error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-time.After(3000 * time.Millisecond):
		return errors.New("fail acquire")
	}
}

func (s *semaphore) Release() error {
	select {
	case <-s.ch:
		return nil
	case <-time.After(3000 * time.Millisecond):
		return errors.New("fail release")
	}
}

func (s *semaphore) Len() int {
	return len(s.ch)
}

func runTicket(i int, s Semaphore) {
	defer func() {
		if err := s.Release(); err != nil {
			fmt.Printf("error %s\n", err)
		}
	}()

	time.Sleep(time.Millisecond * 1000)
	fmt.Println("work", i)
}

func main() {
	s := NewSemaphore(10)
	for i := 0; i < 100; i++ {
		if err := s.Acquire(); err != nil {
			fmt.Printf("error: %s\n", err)
		}
		go runTicket(i, s)
	}

	for s.Len() > 0 {
		time.Sleep(10 * time.Millisecond)
	}

}
