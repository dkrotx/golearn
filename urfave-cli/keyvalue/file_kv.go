package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

type fileKeyValue struct {
	path string
}

type ValueAndOptions struct {
	Value string `json:"value"`
}

type kvData struct {
	Mapping map[string]ValueAndOptions `json:"content"`
}

func (fkv *fileKeyValue) Get(key string) (value string, err error) {
	data, err := fkv.load()
	if err != nil {
		return
	}

	valOpts, found := data.Mapping[key]
	if !found {
		err = KeyNotFoundErr
	}

	value = valOpts.Value
	return
}

func (fkv *fileKeyValue) Set(key, value string, opts ...Option) error {
	data, err := fkv.load()
	if err != nil {
		return err
	}

	valOpts := ValueAndOptions{ Value: value }
	data.Mapping[key] = valOpts

	return fkv.save(data)
}

func (fkv *fileKeyValue) Delete(key string) error {
	data, err := fkv.load()
	if err != nil {
		return err
	}

	delete(data.Mapping, key)
	return fkv.save(data)
}

func (fkv *fileKeyValue) load() (*kvData, error) {
	raw, err := ioutil.ReadFile(fkv.path)

	if err != nil {
		return nil, err
	}

	var data kvData
	if err = json.Unmarshal(raw, &data); err != nil {
		return nil, errors.Wrap(err, "failed to parse KV file")
	}

	if data.Mapping == nil {
		data.Mapping = make(map[string]ValueAndOptions)
	}
	return &data, nil
}

func (fkv *fileKeyValue) save(data *kvData) error {
	raw, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n') // make it pretty as human-entered JSON

	ioutil.WriteFile(fkv.path, raw, 0666)
	return nil
}

func (fkv *fileKeyValue) init() error {
	_, err := os.Stat(fkv.path)

	if err != nil {
		if os.IsNotExist(err) {
			return fkv.save(&kvData{})
		}
		return err
	}

	return nil
}

func NewFileKeyValue(path string) (KeyValue, error) {
	fkv := &fileKeyValue{path: path}
	if err := fkv.init(); err != nil {
		return nil, err
	}
	return fkv, nil
}