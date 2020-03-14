package tuple

import "fmt"
import "log"
import "strconv"

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type Atom struct {
	Name string
}

/////////////////////////////////////////////////////////////////////////////
// Tuple
/////////////////////////////////////////////////////////////////////////////

type Tuple struct {
	list []interface{}
}

func (tuple Tuple) PrettyPrint(depth string) string {
	result :=  depth + "(\n"
	newDepth := depth + "  "
	for _, token := range tuple.list {
		var value string
		switch token.(type) {
		case Tuple: value = token.(Tuple).PrettyPrint (newDepth)
		case Atom:  value = newDepth + token.(Atom).Name
		case string: value = newDepth + "\"" + token.(string) + "\""  // TODO Escape
		case int64: value = newDepth + strconv.FormatInt(int64(token.(int64)), 10)
		case float64: value = newDepth + fmt.Sprint(token.(float64))
		default:
			value = "???"
			log.Printf("Type not recognised: %s", token);
		}
		result = result + value + "\n"
	}
	result = result + depth + ")"
	return result
}


func (tuple *Tuple) Append(token interface{}) {
	tuple.list = append(tuple.list, token)
}

func NewTuple() Tuple {
	return Tuple{make([]interface{}, 0)}
}

