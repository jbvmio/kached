package kached

import (
	"bytes"
	"encoding/gob"
)

// EncodeFunc is used to encode data.
type EncodeFunc func(interface{}) []byte

// DecodeFunc is used to decode data.
type DecodeFunc func(raw []byte) interface{}

// D is used to encode / decode between cache and db.
type D struct {
	I interface{}
}

func encode(value interface{}) (raw []byte) {
	var err error
	d := D{I: value}
	raw, err = gobData(d)
	if err != nil {
		panic("encode error: " + err.Error())
	}
	return
}

func decode(raw []byte) (value interface{}) {
	var d D
	err := unGobData(raw, &d)
	if err != nil {
		panic("decode error: " + err.Error())
	}
	value = d.I
	return
}

func gobData(v interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(v)
	return buffer.Bytes(), err
}

func unGobData(data []byte, v interface{}) error {
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	err := dec.Decode(v)
	return err
}
