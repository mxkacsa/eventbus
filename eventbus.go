package eventbus

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	// ErrSubscribed is the error returned when a callback is already subscribed.
	ErrSubscribed = errors.New("already subscribed")
	// ErrNotPtr is the error returned when the provided pointer is not valid.
	ErrNotPtr = errors.New("not a pointer")
)

// CbFunc represents a callback function type.
type CbFunc[T any] func(T)

type subscriber[T any] struct {
	ptrId string
	cb    CbFunc[T]
}

// EventBus represents an event channel where subscribers can register to receive updates.
type EventBus[T any] struct {
	subscribers []subscriber[T]
	lock        sync.RWMutex
	lastValue   T
}

// SubscribeFunc adds a callback function to the subscribers list.
func (eb *EventBus[T]) SubscribeFunc(cb CbFunc[T]) error {
	return eb.subscribe(fmt.Sprintf("%p", cb), cb)
}

// SubscribeMethod adds a method as a callback function to the subscribers list.
func (eb *EventBus[T]) SubscribeMethod(ptr any, cb CbFunc[T]) error {
	if !isPtr(ptr) {
		return ErrNotPtr
	}

	return eb.subscribe(fmt.Sprintf("%p_%p", ptr, cb), cb)
}

func (eb *EventBus[T]) subscribe(id string, cb CbFunc[T]) error {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	for _, sub := range eb.subscribers {
		if sub.ptrId == id {
			return ErrSubscribed
		}
	}

	eb.subscribers = append(eb.subscribers, subscriber[T]{
		ptrId: id,
		cb:    cb,
	})

	return nil
}

// UnsubscribeFunc removes a callback function from the subscribers list.
func (eb *EventBus[T]) UnsubscribeFunc(cb CbFunc[T]) error {
	eb.unsubscribe(fmt.Sprintf("%p", cb))
	return nil
}

// UnsubscribeMethod removes a method from the subscribers list.
func (eb *EventBus[T]) UnsubscribeMethod(ptr any, cb CbFunc[T]) error {
	if !isPtr(ptr) {
		return ErrNotPtr
	}

	eb.unsubscribe(fmt.Sprintf("%p_%p", ptr, cb))
	return nil
}

func (eb *EventBus[T]) unsubscribe(id string) {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	var newSubscribers []subscriber[T]

	for _, sub := range eb.subscribers {
		if sub.ptrId != id {
			newSubscribers = append(newSubscribers, sub)
		}
	}

	eb.subscribers = newSubscribers
}

// Publish sends an event to all subscribers synchronously.
func (eb *EventBus[T]) Publish(value T) {
	eb.publish(value, false)
}

// PublishAsync sends an event to all subscribers asynchronously.
func (eb *EventBus[T]) PublishAsync(value T) {
	eb.publish(value, true)
}

func (eb *EventBus[T]) publish(value T, async bool) {
	eb.lock.Lock()
	eb.lastValue = value
	subs := eb.subscribers
	eb.lock.Unlock()

	for _, subscriber := range subs {
		if async {
			go subscriber.cb(value)
		} else {
			subscriber.cb(value)
		}
	}
}

// GetLastValue returns the last event published to the EventBus.
func (eb *EventBus[T]) GetLastValue() T {
	eb.lock.RLock()
	defer eb.lock.RUnlock()

	return eb.lastValue
}

// GetSubscribersCount returns the number of current subscribers.
func (eb *EventBus[T]) GetSubscribersCount() int {
	eb.lock.RLock()
	defer eb.lock.RUnlock()

	return len(eb.subscribers)
}

// isPtr checks if the provided parameter is a pointer.
func isPtr(param interface{}) bool {
	return reflect.ValueOf(param).Kind() == reflect.Ptr
}
