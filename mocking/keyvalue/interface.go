package keyvalue

type KeyValue interface {
	Set(key, value string)
	Get(key string) (value string, found bool)
}
