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
package tuple

import "math"
import "strings"

//  A simple toy evaluator.
//
//  Treats the data produced by parsing and treats it as an executable program
//  (code can be accessed and treated as if it is data and data can be treated as code).
//  This is normal for LISP but will work for any of the input grammars.
//
//  See:
//  * [Homoiconic](https://en.wikipedia.org/wiki/Homoiconicity) language treats "code as data".
//  * [Eval](https://en.wikipedia.org/wiki/Eval)
//  * [Meta-circular_evaluator](https://en.wikipedia.org/wiki/Meta-circular_evaluator)
//  
func SimpleEval(expression Value, next Next) {
	result := eval(expression)
	next(result)
}

func eval(expression Value) Value {

	switch val := expression.(type) {
	case Tuple:
		ll := len(val.List)
		if ll == 0 {
			return val
		}
		head := val.List[0]
		if ll == 1 {
			return eval(head)
		}
		atom, ok := head.(Atom)
		if ! ok {
			// TODO Handle case of list: (1 2 3)
		}
		evaluated := make([]Value, ll-1)
		for k, v := range val.List[1:] {
			evaluated[k] = eval(v)
		}

		name := atom.Name
		switch ll {
		case 2:
			if str, ok := evaluated[0].(String); ok {
				switch name {
				case "len": return Int64(len(str))
				case "lower": return String(strings.ToLower(string(str)))
				case "upper": return String(strings.ToUpper(string(str)))
				default: return NAN // Float64math.NaN()
				}
			}
			
			aa := toFloat64(evaluated[0])
			switch name {
			case "log":
				return Float64(math.Log(aa))
			case "exp": return Float64(math.Exp(aa))
			case "sin": return Float64(math.Sin(aa))
			case "cos": return Float64(math.Cos(aa))
			case "tan": return Float64(math.Tan(aa))
			case "acos": return Float64(math.Acos(aa))
			case "asin": return Float64(math.Asin(aa))
			case "atan": return Float64(math.Atan(aa))
			case "round": return Float64(math.Round(aa))
			case "-": return Float64(-aa)
			case "+": return Float64(aa)
			//case "_unary_minus": return -aa
			//case "_unary_plus": return aa
			default: return NAN // math.NaN()
			}
		case 3:
			aa := toFloat64(evaluated[0])
			bb := toFloat64(evaluated[1])
			switch name {
				//case ":=": 
			case "^": return Float64(math.Pow(aa,bb))
			case "+": return Float64(aa+bb)
			case "-": return Float64(aa-bb)
			case "*": return Float64(aa*bb)
			case "/": return Float64(aa/bb)
			case "==": return Bool(aa==bb)
			case "!=": return Bool(aa!=bb)
			case ">=": return Bool(aa>=bb)
			case "<=": return Bool(aa<=bb)
			case ">": return Bool(aa>bb)
			case "<": return Bool(aa<bb)
			case "atan2": return Float64(math.Atan2(aa,bb))
			default: return NAN
			}
		default:
			return NAN
			// TODO
		}
	case Int64:
		return expression
	case Float64:
		return expression
	case String:
		return expression
	case Bool:
		return expression
	case Atom:
		return Nullary(expression.(Atom))
	default:
		return NAN
	}
		
}
func toString(value Value) string {
	switch val := value.(type) {
	case Atom: return val.Name
	case String: return string(val)
	default:
		return ""
	}
}

func toFloat64(value Value) float64 {
	switch val := value.(type) {
	case Int64: return float64(val)
	case Float64: return float64(val)
	case Atom: return math.NaN() // TODO Nullary(val)
	//case bool: return float64(val)
	default:
		return math.NaN()
	}
}

func Nullary(val Atom) Value {
	switch val.Name {
	case "PI": return Float64(math.Pi)
	case "PHI": return Float64(math.Phi)
	case "E": return Float64(math.E)
	case "true": return Bool(true)
	case "false": return Bool(false)
	default: return NAN   // TODO look up variable
	}
}
