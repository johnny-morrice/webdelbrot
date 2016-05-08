package main

import (
    "sync"
    "time"
)

type debouncer struct {
    length time.Duration
    batch *filterbatch
    mut sync.RWMutex
}

type filterbatch struct {
    wg sync.WaitGroup
    once sync.Once
}

func newdebounce(length uint) *debouncer {
    tf := &debouncer{}
    tf.length = time.Duration(length)
    tf.batch = &filterbatch{}
    return tf
}

// Do f only once for a batch of time-close invocations.
func (tf *debouncer) do(f func ()) {
    tf.wait()
    tf.once(f)
}

    // once executes the function exactly once per batch.
func (tf *debouncer) once(f func ()) {
    tf.batch.wg.Done()

    tf.mut.Lock()
    batch := tf.batch
    tf.mut.Unlock() // defer useless in this situation
    batch.once.Do(func () {
        f()
        tf.mut.Lock()
        defer tf.mut.Unlock()
        tf.batch = &filterbatch{}
    })
}

// wait until the current batch has finished
func (tf *debouncer) wait() {
    tf.mut.RLock()
    defer tf.mut.RUnlock()
    batch := tf.batch

    batch.wg.Add(1)
    go func() {
        time.Sleep(tf.length)
        batch.wg.Done()
    }()
}