# EventBus

The EventBus is a simple Go package that offers a pub/sub pattern implementation in Go. This means that it allows parts of your application to subscribe to events and publish events.
## Features

- Type-Safe Generics: EventBus uses generics, which enforce type safety at compile-time, preventing potential runtime errors due to incorrect data types. This means the publisher, the event, and the listing function are all of the same type, reducing the risk of type mismatch or assertion errors.
- Subscribe/Unsubscribe: You can add callback functions or methods and remove them from the subscriber list easily with `Subscribe*` and `Unsubscribe*` functions.
- Synchronous and asynchronous events: You can publish events either synchronously (`Publish`) or asynchronously (`PublishAsync`).
- Get last event and subscriber count: You can use `GetLastValue` to get the last published event and `GetSubscribersCount` to get the current number of subscribers.

## Installation

To install eventbus, simply run:

```bash
go get github.com/mxkacsa/eventbus
```

## Usage

To use this package, you first create an `EventBus` object:

```go
eb := new(eventbus.EventBus[string]) // string, int, or any other type what you want
```

Then, you can subscribe functions to this EventBus:

```go
func callback(message string) {
    log.Println(message)
}

err := eb.SubscribeFunc(callback)
```

Or you can subscribe methods:

```go
type MyStruct struct {}

func (m *MyStruct) callback(message string) {
    log.Println(message)
}

...
sample := new(MyStruct)
err := eb.SubscribeMethod(sample, sample.callback)
```

To publish an event, use the Publish or PublishAsync method:

```go
eb.Publish("Hello World!")
// or asynchronously
eb.PublishAsync("Hello World!")
```

### For examples, check out the [examples](examples) folder.

## Errors

The EventBus package has two errors that can be returned when subscribing and unsubscribing:
- ErrSubscribed: Returned when a callback is already subscribed.
- ErrNotPtr: Returned when the provided pointer is not valid.

## License

This code is available under the MIT license.