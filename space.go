// Package spacelang implements a simple QBasic like language runtime in go
// Spacelang is a VM that only implements go function calls from the language's
// syntax.
package spacelang

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"text/scanner"
)

var rxp = regexp.MustCompile("[a-zA-Z_][a-zA-Z0-9_]{0,31}")

func scanTokens(s string) []string {
	ss := make([]string, 0)

	scan := &scanner.Scanner{}
	scan.Init(strings.NewReader(s))

	for tok := scan.Scan(); tok != scanner.EOF; tok = scan.Scan() {
		temp := scan.TokenText()
		ss = append(ss, temp)
	}
	return ss
}

type VM struct {
	istack int64
	Vars   map[string]interface{}

	Funcs map[string]func(a ...*Token) error
}

func NewVM() *VM {
	vm := &VM{
		Vars:   make(map[string]interface{}),
		Funcs:  make(map[string]func(a ...*Token) error),
		istack: 1,
	}

	vm.Funcs["print"] = func(a ...*Token) error {
		for i, val := range a {
			if val.Type == VALUE {
				fmt.Print(val)
			} else {
				vx, ok := vm.Vars[val.String()]
				if !ok {
					return fmt.Errorf("'%s' : variable not found in scope", val.String())
				}
				fmt.Print(vx)
			}
			if i < len(a)-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
		return nil
	}

	vm.Funcs["goto"] = func(a ...*Token) error {
		if len(a) != 1 {
			return fmt.Errorf("not a valid syntax of goto")
		}

		if a[0].Type == REFERENCE {
			vx, ok := vm.Vars[a[0].String()]
			if !ok {
				return fmt.Errorf("'%s' : variable not found in scope", a[0].String())
			}
			val, ok := vx.(int64)
			if !ok {
				return fmt.Errorf("'%T' : variable of that type cannot be used for goto", val)
			}
			vm.istack = val - 1
			return nil
		}

		if a[0].ValueType != INT {
			return fmt.Errorf("goto requires only an integer argument. Got %s", a[0])
		}

		vm.istack = a[0].Value.(int64) - 1

		return nil
	}

	vm.Funcs["set"] = func(a ...*Token) error {
		if len(a) != 2 {
			return fmt.Errorf("not a valid syntax of set")
		}

		if a[0].Type != REFERENCE {
			return fmt.Errorf("can only assign value to a variable only")
		}

		vm.Vars[a[0].String()] = a[1].Value
		return nil
	}

	return vm
}

func (vm *VM) GetVariable(x string) (*Token, error) {
	num, err := strconv.ParseInt(x, 10, 64)
	if err == nil {
		return &Token{Type: VALUE, Value: num, ValueType: INT}, nil
	}

	numf, err := strconv.ParseFloat(x, 64)
	if err == nil {
		return &Token{Type: VALUE, Value: numf, ValueType: FLOAT}, nil
	}

	if x[0] == '"' && x[len(x)-1] == '"' {
		return &Token{Type: VALUE, Value: x[1 : len(x)-1], ValueType: STRING}, nil
	}

	if !rxp.MatchString(x) {
		return nil, fmt.Errorf("'%s' : not a valid variable name", x)
	}

	return &Token{Type: REFERENCE, Value: x}, nil
}

func (vm *VM) Eval(s string) error {
	lines := strings.Split(s, "\n")

	for ; vm.istack <= int64(len(lines)); vm.istack++ {
		if vm.istack < 1 {
			return fmt.Errorf("line %d is not an accessible line")
		}

		line := lines[vm.istack-1]
		if strings.HasPrefix(line, "#") {
			continue
		}
		toks := scanTokens(line)
		if len(toks) == 0 {
			continue
		}

		fn, ok := vm.Funcs[toks[0]]
		if !ok || fn == nil {
			return fmt.Errorf("In line : %d, function '%s' not defined in the vm", vm.istack, toks[0])
		}

		Vars := make([]*Token, 0)
		for _, val := range toks[1:] {
			vx, err := vm.GetVariable(val)
			if err != nil {
				return fmt.Errorf("In line : %d, %s", vm.istack, err)
			}
			Vars = append(Vars, vx)
		}

		err := fn(Vars...)
		if err != nil {
			return err
		}
	}
	return nil
}
