/*
    This file is part of WOZG.

    WOZG is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    WOZG is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with WOZG.  If not, see <https://www.gnu.org/licenses/>.
*/
package eval

import "math"
import "strings"
import "tuple"
import "reflect"
import "fmt"

// These functions are harmless in the sense that they just do basic functions and do not provide any access to resources
// such as operating system or even allocating memory.
func AddBooleanAndArithmeticFunctions(table * SymbolTable) {

	table.Add("exp", math.Exp)
	table.Add("log", math.Log)
	table.Add("sin", math.Sin)
	table.Add("cos", math.Cos)
	table.Add("tan", math.Tan)
	table.Add("acos", math.Acos)
	table.Add("asin", math.Asin)
	table.Add("atan", math.Atan)
	table.Add("atan2", math.Atan2)
	table.Add("round", math.Round)
	table.Add("round2", func(value float64) float64 { return math.Round(value*100)/100 }) // Not needed
	table.Add("**", math.Pow)
	table.Add("+", func (aa float64) float64 { return aa })
	table.Add("-", func (aa float64) float64 { return -aa })
	table.Add("+", func (aa float64, bb float64) float64 { return aa+bb })
	table.Add("-", func (aa float64, bb float64) float64 { return aa-bb })
	table.Add("*", func (aa float64, bb float64) float64 { return aa*bb })
	table.Add("/", func (aa float64, bb float64) float64 { return aa/bb })
	table.Add("eq", func (aa string, bb string) bool { return aa==bb })  // Should take Value as argument
	table.Add("==", func (aa float64, bb float64) bool { return aa==bb })  // Should take Value as argument
	table.Add("!=", func (aa float64, bb float64) bool { return aa!=bb })
	table.Add(">=", func (aa float64, bb float64) bool { return aa>=bb })
	table.Add("<=", func (aa float64, bb float64) bool { return aa<=bb })
	table.Add(">", func (aa float64, bb float64) bool { return aa>bb })
	table.Add("<", func (aa float64, bb float64) bool { return aa<bb })
	table.Add("&&", func (aa bool, bb bool) bool { return aa&&bb })
	table.Add("||", func (aa bool, bb bool) bool { return aa||bb })
	table.Add("!", func (aa bool) bool { return ! aa })
	table.Add("PI", func () float64 { return math.Pi })
	table.Add("PHI", func () float64 { return math.Phi })
	table.Add("E", func () float64 { return math.E })
	table.Add("true", func () bool { return true })
	table.Add("false", func () bool { return false })

}

// String functions that do not allocate any memory
func AddHarmlessStringFunctions(table * SymbolTable) {
	table.Add("len", func(value string) int64 { return int64(len(value)) })
	table.Add("lower", strings.ToLower)
	table.Add("upper", strings.ToUpper)
}

// Tuple functions that do not allocate any memory
func AddHarmlessTupleFunctions(table * SymbolTable)  {

	table.Add("nth", func(index int64, value Tuple) Value {
		if index < 0 || index >= int64(value.Length()) {
			return tuple.EMPTY
		}
		return value.List[index]
	})

	table.Add("istuple", func (context EvalContext, value Value) bool {
		evaluated := Eval(context, value)
		_, ok := evaluated.(Tuple)
		return ok
	})
}

/////////////////////////////////////////////////////////////////////////////

type ErrorIfFunctionNotFound struct {}

func (function * ErrorIfFunctionNotFound) Find(context EvalContext, name Atom, args [] Value) (*SymbolTable, reflect.Value) {
	return nil, reflect.ValueOf(func(args... Value) bool {
		fmt.Printf("ERROR: function not found: '%s' %s\n", name.Name, args)  // TODO ought to use context logger
		return false
	})
}

/////////////////////////////////////////////////////////////////////////////

// These functions are harmless
// One can execute them from a script without any worry they will access something they ought not to or use up resources.
func NewHarmlessSymbolTable(notFound CallHandler) SymbolTable {
	table := NewSymbolTable(notFound)
	AddBooleanAndArithmeticFunctions(&table)
	AddHarmlessStringFunctions(&table)
	AddHarmlessTupleFunctions(&table)
	
	return table
}
