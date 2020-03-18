package tuple

func Query(query string, expression interface{}, next Next) {
	// TODO match query against expression
	match := true // TODO
	if match {
		next(expression)
	}
}
