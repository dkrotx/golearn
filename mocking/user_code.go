package main

import (
	"fmt"
	"golearn/mocking/keyvalue"
)

func WriteAndGet(kv keyvalue.KeyValue) {
	kv.Set("test-key", "test-value")
	value, found := kv.Get("test-key")
	fmt.Printf("value: %q, found=%v", value, found)
}

func main() {
	WriteAndGet(keyvalue.NewFileKeyValue("/tmp/xyz.txt"))
}
