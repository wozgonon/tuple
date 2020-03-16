package tuple

import "math"

//  A simple toy evaluator

type Next func(value interface{})

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
	case "Pi": return math.Pi
	case "Phi": return math.Phi
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
			case "-": return -aa
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

