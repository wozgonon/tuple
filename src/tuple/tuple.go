package tuple

func check(e error) {
    if e != nil {
        panic(e)
    }
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
		subTuple, ok := token.(Tuple)
		var value string
		if ok {
			value = subTuple.PrettyPrint (newDepth);
			result = result + value + "\n"
		} else {
			value, ok = token.(string)
			result = result + newDepth + value + "\n"
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

