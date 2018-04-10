package core

import (
	"strconv"
)

type Token struct {
	tokenValue string
	tokenType  int
}

func NewToken(tvalue string, ttype int) *Token {
	return &Token{tokenValue: tvalue, tokenType: ttype}
}

func (token *Token) GetTokenValue() string {
	return token.tokenValue
}
func (token *Token) GetTokenType() int {
	return token.tokenType
}
func (token *Token) SetTokenValue(tvalue string) {
	token.tokenValue = tvalue
}
func (token *Token) SetTokenType(ttype int) {
	token.tokenType = ttype
}

func (token *Token) String() string {
	return "tokenValue = " + token.tokenValue + " / tokenType = " + strconv.Itoa(token.tokenType)
}

func (token *Token) Clone() interface{} {
	retval := &Token{}

	retval.tokenType = token.tokenType
	retval.tokenValue = token.tokenValue

	return retval
}
