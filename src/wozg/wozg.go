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
package main

import (
	"tuple"
	"os"
	"fmt"
)

import "flag"


//func check(e error) {
//    if e != nil {
//        panic(e)
//    }
//}

/*
	lisp := func (context * tuple.ParserContext) {
		parser := tuple.NewSExpressionParser(lispStyle, outputStyle, nextFunction(outputStyle))
		parser.ParseSExpression (context)
	}
	jml :=  func (context * tuple.ParserContext) {
		parser := tuple.NewSExpressionParser(jmlStyle, outputStyle, nextFunction(outputStyle))
		parser.ParseSExpression (context)
	}
	tcl := func (context * tuple.ParserContext) {
		parser := tuple.NewSExpressionParser(tclStyle, outputStyle, nextFunction(outputStyle))
		parser.ParseCommandShell (context)
	}
	tuple := func (context * tuple.ParserContext) {
		parser := tuple.NewSExpressionParser(tupleStyle, outputStyle, nextFunction(outputStyle))
		parser.ParseTuple (context)
	}
*/	

func style (value string) (tuple.Style) {


	tclStyle := tuple.Style{"", "", "  ", "{", "}", "", "\n", "true", "false", '#'}
	jmlStyle := tuple.Style{"\n", "", "  ", "{", "}", "", "\n", "true", "false", '#'}
	tupleStyle := tuple.Style{"", "", "  ", "(", ")", ",", "\n", "true", "false", '%'} // prolog, sql '--' for 
	lispStyle := tuple.Style{"", "", "  ", "(", ")", "", "\n", "true", "false", ';'}

	yamlStyle := tuple.Style{"---\n", "...\n", "  ", ":", "", "", "\n", "true", "false", '#'}

	// https://en.wikipedia.org/wiki/INI_file
	iniStyle := tuple.Style{"", "", "", ": ", "", "", "\n", "true", "false", '#'}
	// https://en.wikipedia.org/wiki/.properties
	propertiesStyle := tuple.Style{"", "", "", " = ", "", "", "\n", "true", "false", '#'}

	switch value {
	case ".l": return lispStyle
	case ".jml": return jmlStyle
	case ".tuple": return tupleStyle
	case ".fl.tcl": return tclStyle
	case ".tcl": return tclStyle
	case ".yaml": return yamlStyle
	case ".ini": return iniStyle
	case ".properties": return propertiesStyle
	case ".json": fallthrough
	case ".xml": fallthrough
	case ".jpost": fallthrough
	case ".tsv": fallthrough
	case ".csv":
	case ".init":
	case ".sql":
		return lispStyle
	default:
		return lispStyle
	}
		return lispStyle
}

func main() {

	var in = flag.String("in", ".l", "The format of the input.")
	var out = flag.String("out", ".l", "The format of the output.")
	var logger = flag.String("log", ".l", "The format of the error logging.")
	var verbose = flag.Bool("verbose", false, "Verbose logging.")
	var eval = flag.Bool("eval", false, "Run 'eval' interpretter.")
	var version = flag.Bool("version", false, "Print version of this software.")
	//var interactive = flag.Bool("interactive", false, "Runs in interactive code, as a CLI or REPL.  Set -in")

	flag.Parse()

	if *version {
		fmt.Printf("%s %s %s %s\n", os.Args[0], VERSION, COMMIT, BUILT)
		return
	}

	args := len(os.Args)
	numberOfFiles := flag.NArg()
	files := os.Args[args-numberOfFiles:]

	//if len(files) == 0 && !*interactive {
	//	return
	//}
	
	outputStyle := style(*out)
	logStyle := style(*logger)

	nextFunction := func(outputStyle tuple.Style) tuple.Next {
		pretty := func(tuple interface{}) {
			outputStyle.PrettyPrint(tuple, func(value string) {
				fmt.Printf ("%s", value)
			})
		}
		if *eval {
			return func(value interface{}) {
				tuple.SimpleEval(value, pretty)
			}
		} else {
			return pretty
		}
	}
	tuple.RunParser(files, logStyle, *verbose,
		func (context * tuple.ParserContext) {
			suffix := context.Suffix()
			if suffix == "" {
				suffix = *in
			}
			fmt.Print(outputStyle.StartDoc)
			
			context.Verbose("source [%s] suffix [%s]", context.SourceName, suffix)
			inputStyle := style(*in)
			parser := tuple.NewSExpressionParser(inputStyle, outputStyle, nextFunction(outputStyle))
			switch suffix {
			case ".l":
				parser.ParseSExpression (context)
			case ".jml":
				parser.ParseSExpression (context)
			case ".fl.tcl": 
			case ".tcl": 
				parser.ParseCommandShell (context)
			case ".tuple": 
				parser.ParseTuple (context)
			case ".yaml":
				fallthrough
			case ".json": fallthrough
			case ".xml": fallthrough
			case ".jpost": fallthrough
			case ".tsv": fallthrough
			case ".csv":
				context.Error("Not implemented file suffix: '%s'", suffix)
			default:
				context.Error("Unsupported file suffix: '%s', source: '%s'", suffix, context.SourceName)
			}
			//fmt.Printf("\n")
			fmt.Print(outputStyle.EndDoc)
		})
}
