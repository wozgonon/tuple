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

//  A simple toy evaluator

func toFloat64(value interface{}) float64 {
	switch val := value.(type) {
	case int64: return float64(val)
	case float64: return val
	case Atom: return math.NaN() // TODO Nullary(val)
	//case bool: return float64(val)
	default:
		return math.NaN()
	}
}

func Nullary(val Atom) interface{} {
	switch val.Name {
	case "PI": return math.Pi
	case "PHI": return math.Phi
	case "E": return math.E
	case "true": return true
	case "false": return false
	default: return math.NaN()   // TODO look up variable
	}
}

func SimpleEval(expression interface{}, next Next) {
	result := eval(expression)
	next(result)
}

func eval(expression interface{}) interface{} {

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
		evaluated := make([]interface{}, ll-1)
		for k, v := range val.List[1:] {
			evaluated[k] = eval(v)
		}

		name := head.(Atom).Name
		switch ll {
		case 2:
			aa := toFloat64(evaluated[0])
			switch name {
			case "log": return math.Log(aa)
			case "sin": return math.Sin(aa)
			case "cos": return math.Cos(aa)
			case "tan": return math.Tan(aa)
			case "acos": return math.Cos(aa)
			case "asin": return math.Sin(aa)
			case "atan": return math.Tan(aa)
			case "_unary_-": return -aa
			case "_unary_+": return aa
			default: return math.NaN()
			}
		case 3:
			aa := toFloat64(evaluated[0])
			bb := toFloat64(evaluated[1])
			switch name {
			//case ":=": 
			case "+": return aa+bb
			case "-": return aa-bb
			case "*": return aa*bb
			case "/": return aa/bb
			case "==": return aa==bb
			case "!=": return aa!=bb
			case ">=": return aa>=bb
			case "<=": return aa<=bb
			case ">": return aa>bb
			case "<": return aa<bb
			default: return math.NaN()
			}
		default:
			return math.NaN()
			// TODO
		}
	case int64:
		return expression
	case float64:
		return expression
	case string:
		return expression
	case bool:
		return expression
	case Atom:
		return Nullary(expression.(Atom))
	default:
		return math.NaN()
	}
		
}

