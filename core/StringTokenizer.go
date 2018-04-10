package core

import (
	"bytes"
	"errors"
	"strings"
)


// StringTokenizer Base string token splitter.
type StringTokenizer struct {
	buffer   string
	ptr      int
	savedPtr int
}

// NewStringTokenizer for string
func NewStringTokenizer(buffer string) *StringTokenizer {
	stringTokenizer := &StringTokenizer{}
	stringTokenizer.buffer = buffer
	stringTokenizer.ptr = 0

	return stringTokenizer
}

func (stringtokenizer *StringTokenizer) super(buffer string) {
	stringtokenizer.buffer = buffer
	stringtokenizer.ptr = 0
}

// NextToken for string
func (stringtokenizer *StringTokenizer) NextToken() string {
	var retval bytes.Buffer

	for stringtokenizer.ptr < len(stringtokenizer.buffer) {
		if stringtokenizer.buffer[stringtokenizer.ptr] == '\n' {
			retval.WriteByte(stringtokenizer.buffer[stringtokenizer.ptr])
			stringtokenizer.ptr++
			break
		} else {
			retval.WriteByte(stringtokenizer.buffer[stringtokenizer.ptr])
			stringtokenizer.ptr++
		}
	}

	return retval.String()
}

func (stringtokenizer *StringTokenizer) HasMoreChars() bool {
	return stringtokenizer.ptr < len(stringtokenizer.buffer)
}

func (stringtokenizer *StringTokenizer) IsHexDigit(ch byte) bool {
	if stringtokenizer.IsDigit(ch) {
		return true
	}
	ch1 := strings.ToUpper(string(ch))[0]
	return ch1 == 'A' || ch1 == 'B' || ch1 == 'C' ||
		ch1 == 'D' || ch1 == 'E' || ch1 == 'F'

}

func (stringtokenizer *StringTokenizer) IsAlpha(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func (stringtokenizer *StringTokenizer) IsDigit(ch byte) bool {
	return ch == '0' || ch == '1' || ch == '2' || ch == '3' || ch == '4' ||
		ch == '5' || ch == '6' || ch == '7' || ch == '8' || ch == '9'
}

func (stringtokenizer *StringTokenizer) GetLine() string {
	var retval bytes.Buffer
	for stringtokenizer.ptr < len(stringtokenizer.buffer) && stringtokenizer.buffer[stringtokenizer.ptr] != '\n' {
		retval.WriteByte(stringtokenizer.buffer[stringtokenizer.ptr])
		stringtokenizer.ptr++
	}
	if stringtokenizer.ptr < len(stringtokenizer.buffer) && stringtokenizer.buffer[stringtokenizer.ptr] == '\n' {
		retval.WriteString("\n")
		stringtokenizer.ptr++
	}
	return retval.String()
}

func (stringtokenizer *StringTokenizer) PeekLine() string {
	curPos := stringtokenizer.ptr
	retval := stringtokenizer.GetLine()
	stringtokenizer.ptr = curPos
	return retval
}

func (stringtokenizer *StringTokenizer) LookAhead() (byte, error) {
	return stringtokenizer.LookAheadK(0)
}

func (stringtokenizer *StringTokenizer) LookAheadK(k int) (byte, error) {
	if stringtokenizer.ptr+k < len(stringtokenizer.buffer) {
		return stringtokenizer.buffer[stringtokenizer.ptr+k], nil
	}
	return 0, errors.New("StringTokenizer::LookAheadK: End of buffer")
}

func (stringtokenizer *StringTokenizer) GetNextChar() (byte, error) {
	if stringtokenizer.ptr >= len(stringtokenizer.buffer) {
		return 0, errors.New("StringTokenizer::GetNextChar: End of buffer")
	}
	ch := stringtokenizer.buffer[stringtokenizer.ptr]
	stringtokenizer.ptr++
	return ch, nil
}

func (stringtokenizer *StringTokenizer) Consume() {
	stringtokenizer.ptr = stringtokenizer.savedPtr
}

func (stringtokenizer *StringTokenizer) ConsumeK(k int) {
	stringtokenizer.ptr += k
}

// GetLines Get a Vector of the buffer tokenized by lines
func (stringtokenizer *StringTokenizer) GetLines() map[int]string {
	result := make(map[int]string)
	for stringtokenizer.HasMoreChars() {
		line := stringtokenizer.GetLine()
		result[len(result)] = line
	}
	return result
}

// GetNextTokenByDelim  Get the next token from the buffer.
func (stringtokenizer *StringTokenizer) GetNextTokenByDelim(delim byte) (string, error) {
	var retval bytes.Buffer
	for {
		la, err := stringtokenizer.LookAheadK(0)
		if err != nil {
			return "", err
		}
		if la == delim {
			break
		}
		retval.WriteByte(stringtokenizer.buffer[stringtokenizer.ptr])
		stringtokenizer.ConsumeK(1)
	}
	return retval.String(), nil
}

// GetSDPFieldName  String
func (stringtokenizer *StringTokenizer) GetSDPFieldName(line string) string {
	if line == "" {
		return ""
	}

	begin := strings.Index(line, "=")
	if begin != -1 {
		return line[0:begin]
	} else {
		return ""
	}
}
