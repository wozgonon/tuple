package tuple

func Query(query string, expression interface{}, next Next) {
	panic("TODO match query against expression.")
	match := true // TODO
	if match {
		next(expression)
	}
}
