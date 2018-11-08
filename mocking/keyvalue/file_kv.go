package keyvalue

import "fmt"

type fileKeyValue struct {
	path string
}

func (kv *fileKeyValue) Set(key, value string) {
	fmt.Println("fileKeyValue.Set()")
}

func (kv *fileKeyValue) Get(key string) (value string, found bool) {
	fmt.Println("fileKeyValue.Get()")
	return "", false
}

func NewFileKeyValue(path string) KeyValue {
	return &fileKeyValue{path}
}
