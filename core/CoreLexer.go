package core

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

/** A lexical analyzer that is used by all parsers in our implementation.
 */
// IMPORTANT - All keyword matches should be between START and END
const CORELEXER_START = 2048
const CORELEXER_END = CORELEXER_START + 2048
const CORELEXER_ID = CORELEXER_END - 1 // IMPORTANT -- This should be < END

// Individial token classes.
const CORELEXER_WHITESPACE = CORELEXER_END + 1
const CORELEXER_DIGIT = CORELEXER_END + 2
const CORELEXER_ALPHA = CORELEXER_END + 3
const CORELEXER_BACKSLASH = (int)('\\')
const CORELEXER_QUOTE = (int)('\'')
const CORELEXER_AT = (int)('@')
const CORELEXER_SP = (int)(' ')
const CORELEXER_HT = (int)('\t')
const CORELEXER_COLON = (int)(':')
const CORELEXER_STAR = (int)('*')
const CORELEXER_DOLLAR = (int)('$')
const CORELEXER_PLUS = (int)('+')
const CORELEXER_POUND = (int)('#')
const CORELEXER_MINUS = (int)('-')
const CORELEXER_DOUBLEQUOTE = (int)('"')
const CORELEXER_TILDE = (int)('~')
const CORELEXER_BACK_QUOTE = (int)('`')
const CORELEXER_NULL = (int)(0) //('\0')	;
const CORELEXER_EQUALS = (int)('=')
const CORELEXER_SEMICOLON = (int)(';')
const CORELEXER_SLASH = (int)('/')
const CORELEXER_L_SQUARE_BRACKET = (int)('[')
const CORELEXER_R_SQUARE_BRACKET = (int)(']')
const CORELEXER_R_CURLY = (int)('}')
const CORELEXER_L_CURLY = (int)('{')
const CORELEXER_HAT = (int)('^')
const CORELEXER_BAR = (int)('|')
const CORELEXER_DOT = (int)('.')
const CORELEXER_EXCLAMATION = (int)('!')
const CORELEXER_LPAREN = (int)('(')
const CORELEXER_RPAREN = (int)(')')
const CORELEXER_GREATER_THAN = (int)('>')
const CORELEXER_LESS_THAN = (int)('<')
const CORELEXER_PERCENT = (int)('%')
const CORELEXER_QUESTION = (int)('?')
const CORELEXER_AND = (int)('&')
const CORELEXER_UNDERSCORE = (int)('_')

type CoreLexer struct {
	StringTokenizer

	globalSymbolTable map[int]string
	lexerTables       map[string]LexerMap
	currentLexer      LexerMap
	currentLexerName  string
	currentMatch      *Token
}

func NewCoreLexer(lexerName string, buffer string) *CoreLexer {
	coreLexer := &CoreLexer{}

	coreLexer.StringTokenizer.super(buffer)

	coreLexer.globalSymbolTable = make(map[int]string)
	coreLexer.lexerTables = make(map[string]LexerMap)
	coreLexer.currentLexer = make(LexerMap)
	coreLexer.currentLexerName = lexerName

	return coreLexer
}

func (coreLexer *CoreLexer) Super(lexerName, buffer string) {
	coreLexer.StringTokenizer.super(buffer)

	coreLexer.globalSymbolTable = make(map[int]string)
	coreLexer.lexerTables = make(map[string]LexerMap)
	coreLexer.currentLexer = make(LexerMap)
	coreLexer.currentLexerName = lexerName
}

func (coreLexer *CoreLexer) SetLexerName(lexerName string) {
	coreLexer.currentLexerName = lexerName
}

func (coreLexer *CoreLexer) GetLexerName() string {
	return coreLexer.currentLexerName
}

func (coreLexer *CoreLexer) AddKeyword(name string, value int) {
	coreLexer.currentLexer[name] = value
	if _, ok := coreLexer.globalSymbolTable[value]; !ok {
		coreLexer.globalSymbolTable[value] = name
	}
}

func (coreLexer *CoreLexer) LookupToken(value int) string {
	if value > CORELEXER_START {
		return coreLexer.globalSymbolTable[value]
	} else {
		return strconv.Itoa(value)
	}
}

func (coreLexer *CoreLexer) AddLexer(lexerName string) LexerMap {
	var ok bool
	coreLexer.currentLexer, ok = coreLexer.lexerTables[lexerName]
	if !ok {
		coreLexer.currentLexer = make(LexerMap)
		coreLexer.lexerTables[lexerName] = coreLexer.currentLexer
	}
	return coreLexer.currentLexer
}

func (coreLexer *CoreLexer) SelectLexer(lexerName string) {
	coreLexer.currentLexer = coreLexer.lexerTables[lexerName]
	coreLexer.currentLexerName = lexerName
}

func (coreLexer *CoreLexer) CurrentLexer() LexerMap {
	return coreLexer.currentLexer
}

/** Peek the next id but dont move the buffer pointer forward.
 */
func (coreLexer *CoreLexer) PeekNextId() string {
	oldPtr := coreLexer.ptr
	retval := coreLexer.Ttoken()
	coreLexer.savedPtr = coreLexer.ptr
	coreLexer.ptr = oldPtr
	return retval
}

/** Get the next id.
 */
func (coreLexer *CoreLexer) GetNextId() string {
	return coreLexer.Ttoken()
}

// call coreLexer after you call match
func (coreLexer *CoreLexer) GetNextToken() *Token {
	return coreLexer.currentMatch

}

/** Look ahead for one token.
 */
func (coreLexer *CoreLexer) PeekNextToken() (*Token, error) {
	tok, err := coreLexer.PeekNextTokenK(1)
	if err != nil {
		return nil, err
	} else {
		return tok[0], nil
	}
}

func (coreLexer *CoreLexer) PeekNextTokenK(ntokens int) ([]*Token, error) {
	old := coreLexer.ptr
	retval := make([]*Token, ntokens)
	var err error
	for i := 0; i < ntokens; i++ {
		tok := &Token{}
		if coreLexer.StartsId() {
			id := coreLexer.Ttoken()
			tok.tokenValue = id
			if _, ok := coreLexer.currentLexer[strings.ToUpper(id)]; ok {
				tok.tokenType = coreLexer.currentLexer[strings.ToUpper(id)]
			} else {
				tok.tokenType = CORELEXER_ID
			}
		} else {
			nextChar, err := coreLexer.GetNextChar()
			if err != nil {
				break
			}
			tok.tokenValue += string(nextChar)
			if coreLexer.IsAlpha(nextChar) {
				tok.tokenType = CORELEXER_ALPHA
			} else if coreLexer.IsDigit(nextChar) {
				tok.tokenType = CORELEXER_DIGIT
			} else {
				tok.tokenType = (int)(nextChar)
			}
		}
		retval[i] = tok
	}
	coreLexer.savedPtr = coreLexer.ptr
	coreLexer.ptr = old
	return retval, err
}

/** Match the given token or throw an exception if no such token
 * can be matched.
 */
func (coreLexer *CoreLexer) Match(tok int) (t *Token, ParseException error) {
	if Debug.ParserDebug {
		Debug.println("match " + strconv.Itoa(tok))
	}
	if tok > CORELEXER_START && tok < CORELEXER_END {
		if tok == CORELEXER_ID {
			// Generic ID sought.
			if !coreLexer.StartsId() {
				return nil, errors.New("ParseException: ID expected")
			}
			id := coreLexer.GetNextId()
			coreLexer.currentMatch = &Token{}
			coreLexer.currentMatch.tokenValue = id
			coreLexer.currentMatch.tokenType = CORELEXER_ID
		} else {
			nexttok := coreLexer.GetNextId()
			cur, ok := coreLexer.currentLexer[strings.ToUpper(nexttok)]
			if !ok || cur != tok {
				return nil, errors.New("ParseException: Unexpected Token")
			}
			coreLexer.currentMatch = &Token{}
			coreLexer.currentMatch.tokenValue = nexttok
			coreLexer.currentMatch.tokenType = tok
		}
	} else if tok > CORELEXER_END {
		// Character classes.
		next, err := coreLexer.LookAheadK(0)
		if err != nil {
			return nil, errors.New("ParseException: Expecting DIGIT")
		}
		if tok == CORELEXER_DIGIT {
			if !coreLexer.IsDigit(next) {
				return nil, errors.New("ParseException: Expecting DIGIT")
			}
			coreLexer.currentMatch = &Token{}
			coreLexer.currentMatch.tokenValue = string(next)
			coreLexer.currentMatch.tokenType = tok
			coreLexer.ConsumeK(1)
		} else if tok == CORELEXER_ALPHA {
			if !coreLexer.IsAlpha(next) {
				return nil, errors.New("ParseException: Expecting ALPHA")
			}
			coreLexer.currentMatch = &Token{}
			coreLexer.currentMatch.tokenValue = string(next)
			coreLexer.currentMatch.tokenType = tok
			coreLexer.ConsumeK(1)
		}
	} else {
		// This is a direct character spec.
		ch := byte(tok)
		next, err := coreLexer.LookAheadK(0)
		if err != nil {
			return nil, errors.New("ParseException: Expecting DIGIT")
		}
		if next == ch {
			coreLexer.currentMatch = &Token{}
			coreLexer.currentMatch.tokenValue = string(ch)
			coreLexer.currentMatch.tokenType = tok
			coreLexer.ConsumeK(1)
		} else {
			return nil, errors.New("ParseException: Expecting")
		}
	}
	return coreLexer.currentMatch, nil
}

func (coreLexer *CoreLexer) SPorHT() {
	var ch byte

	for ch, _ = coreLexer.LookAheadK(0); ch == ' ' || ch == '\t'; ch, _ = coreLexer.LookAheadK(0) {
		coreLexer.ConsumeK(1)
	}
}

func (coreLexer *CoreLexer) StartsId() bool {
	nextChar, err := coreLexer.LookAheadK(0)
	if err != nil {
		return false
	}
	return (coreLexer.IsAlpha(nextChar) ||
		coreLexer.IsDigit(nextChar) ||
		nextChar == '-' ||
		nextChar == '.' ||
		nextChar == '!' ||
		nextChar == '%' ||
		nextChar == '*' ||
		nextChar == '_' ||
		nextChar == '+' ||
		nextChar == '`' ||
		nextChar == '\'' ||
		nextChar == '~')
}

func (coreLexer *CoreLexer) Ttoken() string {
	var nextId bytes.Buffer

	for coreLexer.HasMoreChars() {
		nextChar, err := coreLexer.LookAheadK(0)
		if err != nil {
			break
		}

		if coreLexer.IsAlpha(nextChar) ||
			coreLexer.IsDigit(nextChar) ||
			nextChar == '-' ||
			nextChar == '.' ||
			nextChar == '!' ||
			nextChar == '%' ||
			nextChar == '*' ||
			nextChar == '_' ||
			nextChar == '+' ||
			nextChar == '`' ||
			nextChar == '\'' ||
			nextChar == '~' {
			coreLexer.ConsumeK(1)
			nextId.WriteByte(nextChar)
		} else {
			break
		}
	}
	return nextId.String()
}

func (coreLexer *CoreLexer) TtokenAllowSpace() string {
	var nextId bytes.Buffer

	for coreLexer.HasMoreChars() {
		nextChar, err := coreLexer.LookAheadK(0)
		if err != nil {
			break
		}

		if coreLexer.IsAlpha(nextChar) ||
			coreLexer.IsDigit(nextChar) ||
			nextChar == '-' ||
			nextChar == '.' ||
			nextChar == '!' ||
			nextChar == '%' ||
			nextChar == '*' ||
			nextChar == '_' ||
			nextChar == '+' ||
			nextChar == '`' ||
			nextChar == '\'' ||
			nextChar == '~' ||
			nextChar == ' ' ||
			nextChar == '\t' {

			nextId.WriteByte(nextChar)
			coreLexer.ConsumeK(1)
		} else {
			break
		}
	}
	return nextId.String()
}

// Assume the cursor is at a quote.
func (coreLexer *CoreLexer) QuotedString() (s string, err error) {
	var retval bytes.Buffer
	var next byte

	if next, err = coreLexer.LookAheadK(0); next != '"' || err != nil {
		return "", err
	}
	coreLexer.ConsumeK(1)
	for {
		if next, err = coreLexer.GetNextChar(); err != nil {
			break
		}
		if next == '"' {
			// Got to the terminating quote.
			break
			// } else if next == 0 { //'\0' {
			// 	return "", errors.New("ParseException: unexpected EOL")
		} else if next == '\\' {
			retval.WriteByte(next)
			next, _ = coreLexer.GetNextChar()
			retval.WriteByte(next)
		} else {
			retval.WriteByte(next)
		}
	}
	return retval.String(), err
}

// Assume the cursor is at a "("
func (coreLexer *CoreLexer) Comment() (s string, err error) {
	var retval bytes.Buffer
	var next byte

	if next, err = coreLexer.LookAheadK(0); next != '(' || err != nil {
		return "", err
	}
	coreLexer.ConsumeK(1)
	for {
		if next, err = coreLexer.GetNextChar(); err != nil {
			break
		}
		if next == ')' {
			break
			// } else if next == 0 { //'\0' {
			// 	return "", errors.New("ParseException: unexpected EOL")
		} else if next == '\\' {
			retval.WriteByte(next)
			if next, err = coreLexer.GetNextChar(); err != nil {
				break
			}
			// if next == 0 { //'\0'{
			// 	return "", errors.New("ParseException: unexpected EOL")
			// }
			retval.WriteByte(next)
		} else {
			retval.WriteByte(next)
		}
	}
	return retval.String(), err
}

func (coreLexer *CoreLexer) ByteStringNoSemicolon() string {
	var retval bytes.Buffer

	for {
		next, err := coreLexer.LookAheadK(0)
		if err != nil {
			break
		}
		if /*next == 0*/ /*'\0'*/ /*||*/ next == '\n' || next == ';' {
			break
		} else {
			coreLexer.ConsumeK(1)
			retval.WriteByte(next)
		}
	}

	return retval.String()
}

func (coreLexer *CoreLexer) ByteStringNoComma() string {
	var retval bytes.Buffer

	for {
		next, err := coreLexer.LookAheadK(0)
		if err != nil {
			break
		}
		if next == '\n' || next == ',' {
			break
		} else {
			coreLexer.ConsumeK(1)
			retval.WriteByte(next)
		}
	}

	return retval.String()
}

func (coreLexer *CoreLexer) CharAsString(ch byte) string {
	var retval bytes.Buffer
	retval.WriteByte(ch)
	return retval.String()
}

/** Lookahead in the inputBuffer for n chars and return as a string.
 * Do not consume the input.
 */
func (coreLexer *CoreLexer) NCharAsString(nchars int) string {
	var retval bytes.Buffer

	for i := 0; i < nchars; i++ {
		next, err := coreLexer.LookAheadK(i)
		if err != nil {
			break
		}
		retval.WriteByte(next)
	}

	return retval.String()
}

/** Get and consume the next number.
 */
func (coreLexer *CoreLexer) Number() (n int, ParseException error) {
	var retval bytes.Buffer

	next, err := coreLexer.LookAheadK(0)
	if err != nil {
		return -1, err
	}
	if !coreLexer.IsDigit(next) {
		return -1, errors.New("ParseException: unexpected token \"" + string(next) + "\"")
	}

	retval.WriteByte(next)
	coreLexer.ConsumeK(1)
	for {
		next, err := coreLexer.LookAheadK(0)
		if err == nil && coreLexer.IsDigit(next) {
			retval.WriteByte(next)
			coreLexer.ConsumeK(1)
		} else {
			break
		}
	}

	if n, err = strconv.Atoi(retval.String()); err != nil {
		return -1, err
	} else {
		return n, nil
	}
}

/** Mark the position for backtracking.
 */
func (coreLexer *CoreLexer) MarkInputPosition() int {
	return coreLexer.ptr
}

/** Rewind the input ptr to the marked position.
 */
func (coreLexer *CoreLexer) RewindInputPosition(position int) {
	coreLexer.ptr = position
}

/** Get the rest of the String
 * @return String
 */
func (coreLexer *CoreLexer) GetRest() string {
	if coreLexer.ptr >= len(coreLexer.buffer) {
		return ""
	} else {
		return coreLexer.buffer[coreLexer.ptr:]
	}
}

/** Get the sub-String until the character is encountered.
 * Acknowledgement - Sylvian Corre submitted a bug fix for coreLexer
 * method.
 * @param char c the character to match
 * @return matching string.
 */
func (coreLexer *CoreLexer) GetString(c byte) (s string, err error) {
	var savedPtr int = coreLexer.ptr
	var retval bytes.Buffer
	var next byte

	for {
		next, err = coreLexer.LookAheadK(0)

		if err != nil /*next == 0*/ { //'\0'
			coreLexer.ptr = savedPtr
			break //return "", errors.New("ParseException: unexpected EOL")
		} else if next == c {
			coreLexer.ConsumeK(1)
			break
		} else if next == '\\' {
			coreLexer.ConsumeK(1)
			next, err = coreLexer.LookAheadK(0)
			if err != nil /*nextchar == 0*/ { //'\0'  {
				coreLexer.ptr = savedPtr
				break //return "", errors.New("ParseException: unexpected EOL")
			} else {
				coreLexer.ConsumeK(1)
				retval.WriteByte(next)
			}
		} else {
			coreLexer.ConsumeK(1)
			retval.WriteByte(next)
		}
	}
	return retval.String(), err
}

/** Get the read pointer.
 */
func (coreLexer *CoreLexer) GetPtr() int {
	return coreLexer.ptr
}

/** Get the buffer.
 */
func (coreLexer *CoreLexer) GetBuffer() string {
	return coreLexer.buffer
}
