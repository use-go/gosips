package core

import "bytes"

/** Generic parser class.
* All parsers inherit coreParser class.
 */
type CoreParser struct {
	nesting_level int

	lexer Lexer //*CoreLexer
}

func NewCoreParser(buffer string) *CoreParser {
	coreParser := &CoreParser{}

	coreParser.lexer = NewCoreLexer("CharLexer", buffer)

	return coreParser
}

func (coreParser *CoreParser) Super(buffer string) {
	coreParser.lexer = NewCoreLexer("CharLexer", buffer)
}

func (coreParser *CoreParser) GetLexer() Lexer {
	return coreParser.lexer
}
func (coreParser *CoreParser) SetLexer(lexer Lexer) {
	coreParser.lexer = lexer
}

func (coreParser *CoreParser) NameValue(separator byte) (nv *NameValue, ParseException error) {
	if Debug.ParserDebug {
		coreParser.Dbg_enter("nameValue")
		defer coreParser.Dbg_leave("nameValue")
	}

	if _, ParseException = coreParser.lexer.Match(CORELEXER_ID); ParseException != nil {
		return nil, ParseException
	}
	name := coreParser.lexer.GetNextToken()

	// eat white space.
	coreParser.lexer.SPorHT()

	quoted := false
	la, err := coreParser.lexer.LookAheadK(0)
	if la == separator && err == nil {
		coreParser.lexer.ConsumeK(1)
		coreParser.lexer.SPorHT()

		var str string

		if la, err = coreParser.lexer.LookAheadK(0); la == '"' && err == nil {
			str, _ = coreParser.lexer.QuotedString()
			quoted = true
		} else {
			if _, ParseException = coreParser.lexer.Match(CORELEXER_ID); ParseException != nil {
				return nil, ParseException
			}
			value := coreParser.lexer.GetNextToken()
			str = value.tokenValue
		}
		nv := NewNameValue(name.tokenValue, str)
		if quoted {
			nv.SetQuotedValue()
		}

		return nv, nil
	} else {
		return NewNameValue(name.tokenValue, ""), nil
	}
}

func (coreParser *CoreParser) Dbg_enter(rule string) {
	var stringBuffer bytes.Buffer
	for i := 0; i < coreParser.nesting_level; i++ {
		stringBuffer.WriteString(">")
	}
	if Debug.ParserDebug {
		println(stringBuffer.String() + rule + "\nlexer buffer = \n" + coreParser.lexer.GetRest())
	}
	coreParser.nesting_level++
}

func (coreParser *CoreParser) Dbg_leave(rule string) {
	var stringBuffer bytes.Buffer
	for i := 0; i < coreParser.nesting_level; i++ {
		stringBuffer.WriteString("<")
	}
	if Debug.ParserDebug {
		println(stringBuffer.String() + rule + "\nlexer buffer = \n" + coreParser.lexer.GetRest())
	}
	coreParser.nesting_level--
}

func (coreParser *CoreParser) PeekLine(rule string) {
	if Debug.ParserDebug {
		Debug.println(rule + " " + coreParser.lexer.PeekLine())
	}
}
