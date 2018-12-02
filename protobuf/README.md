# protobuf
Demonstration of use protobuf packge in Go.

## Installation
Install protoc following instruction [here](https://github.com/golang/protobuf)

## Prepare
You need generate .pb-files first:
```make proto```

## Demo
First you need to generate binary file:  
```go run main.go write /tmp/test.bin```

Now you can see documents from it:  
```go run main.go read /tmp/test.bin```
