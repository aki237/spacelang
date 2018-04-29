package spacelang

import "fmt"

type TokenType int

const (
	REFERENCE TokenType = 0
	VALUE     TokenType = 1
)

type ValueType int

const (
	INT    ValueType = 0
	FLOAT  ValueType = 1
	STRING ValueType = 2
)

type Token struct {
	Type      TokenType
	ValueType ValueType
	Value     interface{}
}

func (t Token) String() string {
	return fmt.Sprint(t.Value)
}
