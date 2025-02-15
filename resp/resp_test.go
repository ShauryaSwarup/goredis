package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestReadValue(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected Value
		wantErr  error
	}{
		{
			name:  "Simple String",
			input: []byte("+OK\r\n"),
			expected: Value{
				Typ: "simplestring",
				Str: "OK",
			},
		},
		{
			name:  "Empty Simple String",
			input: []byte("+\r\n"),
			expected: Value{
				Typ: "simplestring",
				Str: "",
			},
		},
		{
			name:  "Bulk String",
			input: []byte("$5\r\nhello\r\n"),
			expected: Value{
				Typ:  "bulk",
				Bulk: "hello",
			},
		},
		{
			name:  "Empty Bulk String",
			input: []byte("$0\r\n\r\n"),
			expected: Value{
				Typ:  "bulk",
				Bulk: "",
			},
		},
		{
			name:     "Null Bulk String",
			input:    []byte("$-1\r\n"),
			expected: Value{Typ: "null"},
		},
		{
			name:  "integer",
			input: []byte(":1000\r\n"),
			expected: Value{
				Typ: "integer",
				Num: 1000,
			},
		},
		{
			name:  "Negative integer",
			input: []byte(":-123\r\n"),
			expected: Value{
				Typ: "integer",
				Num: -123,
			},
		},
		{
			name:  "Simple Error",
			input: []byte("-ERR something\r\n"),
			expected: Value{
				Typ: "simpleerror",
				Err: "ERR something",
			},
		},
		{
			name:  "Array",
			input: []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
			expected: Value{
				Typ: "array",
				Array: []Value{
					{Typ: "bulk", Bulk: "hello"},
					{Typ: "bulk", Bulk: "world"},
				},
			},
		},
		{
			name:  "Empty Array",
			input: []byte("*0\r\n"),
			expected: Value{
				Typ:   "array",
				Array: []Value{},
			},
		},
		{
			name:  "Nested Arrays",
			input: []byte("*2\r\n*2\r\n:1\r\n:2\r\n*2\r\n+a\r\n-b\r\n"),
			expected: Value{
				Typ: "array",
				Array: []Value{
					{
						Typ: "array",
						Array: []Value{
							{Typ: "integer", Num: 1},
							{Typ: "integer", Num: 2},
						},
					},
					{
						Typ: "array",
						Array: []Value{
							{Typ: "simplestring", Str: "a"},
							{Typ: "simpleerror", Err: "b"},
						},
					},
				},
			},
		},
		{
			name:    "Invalid Type",
			input:   []byte("?invalid\r\n"),
			wantErr: ErrUnexpectedType,
		},
		{
			name:    "Malformed integer",
			input:   []byte(":notanum\r\n"),
			wantErr: ErrInvalidSyntax,
		},
		{
			name:    "Incomplete Array",
			input:   []byte("*2\r\n$5\r\nhello\r\n"),
			wantErr: io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(bytes.NewReader(tt.input))
			resp := NewResp(r)
			val, err := resp.ReadValue()

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(val, tt.expected) {
				t.Errorf("\nexpected: %+v\nreceived: %+v", tt.expected, val)
			}
		})
	}
}
