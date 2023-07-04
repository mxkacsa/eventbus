package eventbus

import (
	"testing"
	"time"
)

func TestEventBus(t *testing.T) {
	t.Run("Subscribe and Unsubscribe", func(t *testing.T) {
		localValue := 0
		cb := func(e int) {
			localValue = e
		}

		var eh EventBus[int]
		err := eh.SubscribeFunc(cb)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		eh.Publish(10)

		if localValue != 10 {
			t.Errorf("Expected event %d, but got %d", 10, localValue)
		}

		err = eh.UnsubscribeFunc(cb)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		eh.Publish(20)

		if localValue != 10 {
			t.Errorf("Expected event %d, but got %d", 10, localValue)
		}
	})

	t.Run("Subscribe and Unsubscribe With Async", func(t *testing.T) {
		localValue := 0
		cb := func(e int) {
			localValue = e
		}

		var eh EventBus[int]
		err := eh.SubscribeFunc(cb)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		eh.PublishAsync(10)
		if localValue != 0 {
			t.Errorf("Expected event %d, but got %d", 0, localValue)
		}

		time.Sleep(20 * time.Millisecond)

		if localValue != 10 {
			t.Errorf("Expected event %d, but got %d", 10, localValue)
		}

		err = eh.UnsubscribeFunc(cb)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		eh.PublishAsync(20)

		time.Sleep(20 * time.Millisecond)

		if localValue != 10 {
			t.Errorf("Expected event %d, but got %d", 10, localValue)
		}
	})

	t.Run("Publish and GetLastValue", func(t *testing.T) {
		var eh EventBus[string]
		event := "item 1"
		eh.Publish(event)
		lastEvent := eh.GetLastValue()
		if lastEvent != event {
			t.Errorf("Expected last event %s, but got %s", event, lastEvent)
		}
	})

	t.Run("Get subscribers count", func(t *testing.T) {
		cb := func(e string) {}

		var eh EventBus[string]
		err := eh.SubscribeFunc(cb)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		count := eh.GetSubscribersCount()
		if count != 1 {
			t.Errorf("Expected 1 subscriber, but got %d", count)
		}

		err = eh.UnsubscribeFunc(cb)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		count = eh.GetSubscribersCount()
		if count != 0 {
			t.Errorf("Expected 0 subscriber, but got %d", count)
		}
	})

	// check if same subscribing
	t.Run("Subscribe same function", func(t *testing.T) {
		cb := func(e string) {}

		var eh EventBus[string]
		err := eh.SubscribeFunc(cb)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		err = eh.SubscribeFunc(cb)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})

	t.Run("Subscribe same method", func(t *testing.T) {
		cb := func(e string) {}

		var eh EventBus[string]
		err := eh.SubscribeFunc(cb)
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		err = eh.SubscribeFunc(cb)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})

	t.Run("Subscribe and Unsubscribe with different methods", func(t *testing.T) {
		test1 := new(Test)
		test2 := new(Test)

		var eh EventBus[string]
		err := eh.SubscribeMethod(test1, test1.cb)

		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		err = eh.SubscribeMethod(test2, test2.cb)

		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		count := eh.GetSubscribersCount()
		if count != 2 {
			t.Errorf("Expected 2 subscribers, but got %d", count)
		}

		eh.Publish("test")

		if test1.Val != "test" {
			t.Errorf("Expected event %s, but got %s", "test", test1.Val)
		}

		if test2.Val != "test" {
			t.Errorf("Expected event %s, but got %s", "test", test2.Val)
		}

		err = eh.UnsubscribeMethod(test1, test1.cb)

		count = eh.GetSubscribersCount()
		if count != 1 {
			t.Errorf("Expected 1 subscriber, but got %d", count)
		}

		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		eh.Publish("test2")

		if test1.Val != "test" {
			t.Errorf("Expected event %s, but got %s", "test", test1.Val)
		}

		if test2.Val != "test2" {
			t.Errorf("Expected event %s, but got %s", "test2", test2.Val)
		}
	})

	t.Run("Subscribe to method with not a pointer", func(t *testing.T) {
		test := Test{}
		var eh EventBus[string]
		err := eh.SubscribeMethod(test, test.cb)

		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})

	t.Run("Unsubscribe to method with not a pointer", func(t *testing.T) {
		test := Test{}
		var eh EventBus[string]
		err := eh.SubscribeMethod(&test, test.cb)

		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		err = eh.UnsubscribeMethod(test, test.cb)

		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})
}

type Test struct {
	Val string
}

func (t *Test) cb(e string) {
	t.Val = e
}
