package resp

import (
	"bufio"
	"fmt"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Resp struct {
	r *bufio.Reader
}

func NewResp(r *bufio.Reader) *Resp {
	return &Resp{
		r: r,
	}
}

type Value struct {
	Typ   string
	Err   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

func (resp *Resp) ReadValue() (Value, error) {
	dataType, err := resp.r.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch dataType {
	case ARRAY: // Array
		return resp.readArray()
	case BULK: // Bulk string
		return resp.readBulk()
	case STRING: // Simple string
		return resp.readSimpleString()
	case INTEGER: // Integers
		return resp.readNumber()
	case ERROR: // Simple Errors
		return resp.readSimpleError()
	default:
		return Value{}, ErrUnexpectedType
	}
}

func (resp *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := resp.r.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (resp *Resp) readInteger() (int, int, error) {
	str, n, err := resp.readLine()
	if err != nil {
		return 0, 0, err
	}
	num, err := strconv.ParseInt(string(str), 10, 64)
	if err != nil {
		return 0, n, fmt.Errorf("%w: %v", ErrInvalidSyntax, err)
	}
	return int(num), n, nil
}

func (resp *Resp) readArray() (Value, error) {
	v := Value{}
	v.Typ = "array"
	n, _, err := resp.readInteger()
	if err != nil {
		return v, err
	}
	v.Array = make([]Value, n)
	for i := 0; i < n; i++ {
		val, err := resp.ReadValue()
		if err != nil {
			return v, err
		}
		v.Array[i] = val
	}
	return v, nil
}

func (resp *Resp) readBulk() (Value, error) {
	v := Value{Typ: "bulk"}
	len, _, err := resp.readInteger()
	if err != nil {
		return v, err
	}
	if len < 0 {
		return Value{Typ: "null"}, nil
	}
	buf := make([]byte, len)
	_, err = resp.r.Read(buf)
	if err != nil {
		return v, err
	}
	v.Bulk = string(buf)
	resp.r.ReadLine()
	return v, nil
}

func (resp *Resp) readSimpleString() (Value, error) {
	v := Value{}
	v.Typ = "simplestring"
	str, _, err := resp.readLine()
	if err != nil {
		return v, err
	}
	v.Str = string(str)
	return v, err
}

func (resp *Resp) readSimpleError() (Value, error) {
	v := Value{}
	v.Typ = "simpleerror"
	str, _, err := resp.readLine()
	if err != nil {
		return v, err
	}
	v.Err = string(str)
	return v, err
}

func (resp *Resp) readNumber() (Value, error) {
	num, _, err := resp.readInteger()
	if err != nil {
		return Value{}, err
	}
	return Value{Typ: "integer", Num: int(num)}, nil
}
