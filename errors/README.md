# Errors demo

In this example we will demonstrate approach of using [errors library](https://github.com/pkg/errors) among with **temporary** property.

## Temporary errors
  
Temporary errors is errors which will be recovered eventually. You almost always get 'em in
production code: there are backends you're calling, and they can fail (connection failed, timeout, etc.).
These types of errors are recoverable - since it may happen next time you'll success.
And there are usually many other errors: parsing errors, failed tp encode something, bad state, etc.

You want to **distinct** these errors:
- Temporary: log with INFO, don't do anything else, just try one more time later.
- Others: log with ERROR severity, send to Sentry, alarm monitoring, etc.


Full article: https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully

## Installation
```
$ go get github.com/pkg/errors
$ go build
```

## Experiment
```
$ ./errors
```

Also see comments in main.go:requestBackend()