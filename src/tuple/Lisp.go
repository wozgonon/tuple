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

/////////////////////////////////////////////////////////////////////////////
// Lisp Grammar
/////////////////////////////////////////////////////////////////////////////

type Lisp struct {
	parser SExpressionParser
}

func (syntax Lisp) Name() string {
	return "Lisp"
}

func (syntax Lisp) FileSuffix() string {
	return ".l"
}

func (syntax Lisp) Parse(context * ParserContext) {
	syntax.parser.ParseSExpression(context)
}

func (syntax Lisp) Print(token interface{}, next func(value string)) {
	syntax.parser.outputStyle.PrettyPrint(token, next)
}

func NewLispGrammar() Grammar {
	style := Style{"", "", "  ", "(", ")", "", "\n", "true", "false", ';'}
	return Lisp{NewSExpressionParser(style, style)}
}

/////////////////////////////////////////////////////////////////////////////
// Tcl Grammar
/////////////////////////////////////////////////////////////////////////////

type Tcl struct {
	parser SExpressionParser
}

func (syntax Tcl) Name() string {
	return "Tcl"
}

func (syntax Tcl) FileSuffix() string {
	return ".tcl"
}

func (syntax Tcl) Parse(context * ParserContext) {
	syntax.parser.ParseCommandShell (context)
	//syntax.parser.ParseSExpression(context)
}

func (syntax Tcl) Print(token interface{}, next func(value string)) {
	syntax.parser.outputStyle.PrettyPrint(token, next)
}

func NewTclGrammar() Grammar {
	style := Style{"", "", "  ", "{", "}", "", "\n", "true", "false", '#'}
	return Tcl{NewSExpressionParser(style, style)}
}

/////////////////////////////////////////////////////////////////////////////
// JML Grammar
/////////////////////////////////////////////////////////////////////////////

type Jml struct {
	parser SExpressionParser
}

func (syntax Jml) Name() string {
	return "Jml"
}

func (syntax Jml) FileSuffix() string {
	return ".jml"
}

func (syntax Jml) Parse(context * ParserContext) {
	syntax.parser.ParseSExpression(context)
}

func (syntax Jml) Print(token interface{}, next func(value string)) {
	syntax.parser.outputStyle.PrettyPrint(token, next)
}

func NewJmlGrammar() Grammar {
	style := Style{"\n", "", "  ", "{", "}", "", "\n", "true", "false", '#'}
	return Jml{NewSExpressionParser(style, style)}
}

/////////////////////////////////////////////////////////////////////////////
// Tuple Grammar
/////////////////////////////////////////////////////////////////////////////

type TupleGrammar struct {
	parser SExpressionParser
}

func (syntax TupleGrammar) Name() string {
	return "TupleGrammar"
}

func (syntax TupleGrammar) FileSuffix() string {
	return ".tuple"
}

func (syntax TupleGrammar) Parse(context * ParserContext) {
	syntax.parser.ParseTuple (context)
	//syntax.parser.ParseSExpression(context)
}

func (syntax TupleGrammar) Print(token interface{}, next func(value string)) {
	syntax.parser.outputStyle.PrettyPrint(token, next)
}

func NewTupleGrammar() Grammar {
	style := Style{"", "", "  ", "(", ")", ",", "\n", "true", "false", '%'} // prolog, sql '--' for 
	return TupleGrammar{NewSExpressionParser(style, style)}
}

/////////////////////////////////////////////////////////////////////////////
// Yaml Grammar
/////////////////////////////////////////////////////////////////////////////

type Yaml struct {
	parser SExpressionParser
}

func (syntax Yaml) Name() string {
	return "Yaml"
}

func (syntax Yaml) FileSuffix() string {
	return ".yaml"
}

func (syntax Yaml) Parse(context * ParserContext) {
	context.Error("Not implemented file suffix: '%s'", syntax.FileSuffix())
}

func (syntax Yaml) Print(token interface{}, next func(value string)) {
	next(syntax.parser.outputStyle.StartDoc)
	syntax.parser.outputStyle.PrettyPrint(token, next)
	next(syntax.parser.outputStyle.EndDoc)
}

func NewYamlGrammar() Grammar {
	style := Style{"---\n", "...\n", "  ", ":", "", "", "\n", "true", "false", '#'}
	return Yaml{NewSExpressionParser(style, style)}
}

/////////////////////////////////////////////////////////////////////////////
// Ini Grammar
/////////////////////////////////////////////////////////////////////////////

type Ini struct {
	parser SExpressionParser
}

func (syntax Ini) Name() string {
	return "Ini"
}

func (syntax Ini) FileSuffix() string {
	return ".ini"
}

func (syntax Ini) Parse(context * ParserContext) {
	context.Error("Not implemented file suffix: '%s'", syntax.FileSuffix())
}

func (syntax Ini) Print(token interface{}, next func(value string)) {
	syntax.parser.outputStyle.PrettyPrint(token, next)
}

func NewIniGrammar() Grammar {
	// https://en.wikipedia.org/wiki/INI_file
	style := Style{"", "", "ini", "= ", "", "", "\n", "true", "false", '#'}
	return Ini{NewSExpressionParser(style, style)}
}

/////////////////////////////////////////////////////////////////////////////
// PropertyGrammar Grammar
/////////////////////////////////////////////////////////////////////////////

type PropertyGrammar struct {
	parser SExpressionParser
}

func (syntax PropertyGrammar) Name() string {
	return "PropertyGrammar"
}

func (syntax PropertyGrammar) FileSuffix() string {
	return ".properties"
}

func (syntax PropertyGrammar) Parse(context * ParserContext) {
	context.Error("Not implemented file suffix: '%s'", syntax.FileSuffix())
}

func (syntax PropertyGrammar) Print(token interface{}, next func(value string)) {
	syntax.parser.outputStyle.PrettyPrint(token, next)
}

func NewPropertyGrammar() Grammar {
	// https://en.wikipedia.org/wiki/.properties
	style := Style{"", "", "", " = ", "", "", "\n", "true", "false", '#'}
	return PropertyGrammar{NewSExpressionParser(style, style)}
}

// TODO json xml postfix

