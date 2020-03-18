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

// The Grammar interface represents a particular language Grammar or Grammar or File Format.
//
// The print and parse method ought to be inverse functions of each other
// so the output of parse can be passed to print which in principle should be parsable by the parse function.
type Grammar interface {
	// A friendly name for the syntax
	Name() string

	// A standard suffix for source files.
	FileSuffix() string
	
	// Parses an input stream of characters into an internal representation (AST)
	// The output ought to be printable by the 'print' method.
	Parse(context * ParserContext) // , next func(tuple Tuple)) (interface{}, error)
	
	// Pretty prints the objects in the given syntax.
	// The output ought to be parsable by the 'parse' method.
	Print(token interface{}, next func(value string))
}

// A set of Grammars
type Grammars struct {
	all map[string]Grammar
}

// Returns a new empty set of syntaxes
func NewGrammars() Grammars{
	return Grammars{make(map[string]Grammar)}
}

func (syntaxes * Grammars) Add(syntax Grammar) {
	suffix := syntax.FileSuffix()
	syntaxes.all[suffix] = syntax
}

func (syntaxes * Grammars) FindBySuffix(suffix string) (*Grammar, bool) {
	syntax, ok := syntaxes.all[suffix]
	return &syntax, ok
}

func (syntaxes * Grammars) FindBySuffixOrPanic(suffix string) *Grammar {
	syntax, ok := syntaxes.FindBySuffix(suffix)
	if ! ok {
		panic("Unsupported file suffix: '" + suffix + "'")
	}
	return syntax
}

