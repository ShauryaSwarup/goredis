package resp

import (
	"fmt"
	"strconv"
)

func (v *Value) Marshal() []byte {
	fmt.Println("IN MARSHAL with v.Typ: ", v.Typ)
	switch v.Typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "simplestring":
		return v.marshalString()
	case "integer":
		return v.marshalInteger()
	case "null":
		return v.marshalNull()
	case "simpleerror":
		return v.marshalError()
	default:
		return []byte{}
	}
}
func (v *Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v *Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, "\r\n"...)
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, "\r\n"...)
	return bytes
}

func (v *Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.Err...)
	bytes = append(bytes, "\r\n"...)
	return bytes
}

func (v *Value) marshalInteger() []byte {
	var bytes []byte
	bytes = append(bytes, INTEGER)
	bytes = append(bytes, strconv.Itoa(v.Num)...)
	bytes = append(bytes, "\r\n"...)
	return bytes
}

func (v *Value) marshalArray() []byte {
	len := len(v.Array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.Array[i].Marshal()...)
	}

	return bytes
}

func (v *Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}
