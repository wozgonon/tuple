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
	//space := ""
	result :=  depth + "(\n"
	newDepth := depth + "  "
	for _, token := range tuple.list {
		switch token.(type) {
		case Tuple: result = result + token.(Tuple).PrettyPrint (newDepth) + "\n"
		case Atom:  result = result + newDepth + token.(Atom).Name + "\n"
		case string: result = result + newDepth + "\"" + token.(string) + "\"" + "\n"  // TODO Escape
		case int64: result = result + newDepth + strconv.FormatInt(int64(token.(int64)), 10) + "\n"
		case float64: result = result + newDepth + fmt.Sprint(token.(float64)) + "\n"
		default:
			log.Printf("Type not recognised: %s", token);
		}
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

