package core

import (
	"bytes"
	"errors"
)

/** SIPParser for host names.
 */
type HostNameParser struct {
	CoreParser
}

func NewHostNameParser(hname string) *HostNameParser {
	hostNameParser := &HostNameParser{}

	hostNameParser.lexer = NewCoreLexer("charLexer", hname)

	return hostNameParser
}

/** The lexer is initialized with the buffer.
 */
func NewHostNameParserFromLexer(lexer Lexer) *HostNameParser {
	hostNameParser := &HostNameParser{}

	hostNameParser.CoreParser.SetLexer(lexer)
	hostNameParser.CoreParser.GetLexer().SelectLexer("charLexer")

	return hostNameParser
}

func (hostNameParser *HostNameParser) DomainLabel() (s string, ParseException error) {
	var retval bytes.Buffer
	if Debug.ParserDebug {
		hostNameParser.Dbg_enter("domainLabel")
		defer hostNameParser.Dbg_leave("domainLabel")
	}

	for hostNameParser.lexer.HasMoreChars() {
		la, err := hostNameParser.lexer.LookAheadK(0)
		if err != nil {
			return retval.String(), err
		}
		if hostNameParser.lexer.IsAlpha(la) {
			hostNameParser.lexer.ConsumeK(1)
			retval.WriteByte(la)
		} else if hostNameParser.lexer.IsDigit(la) {
			hostNameParser.lexer.ConsumeK(1)
			retval.WriteByte(la)
		} else if la == '-' {
			hostNameParser.lexer.ConsumeK(1)
			retval.WriteByte(la)
		} else {
			break
		}
	}

	return retval.String(), nil
}

func (hostNameParser *HostNameParser) Ipv6Reference() (s string, ParseException error) {
	var retval bytes.Buffer
	if Debug.ParserDebug {
		hostNameParser.Dbg_enter("ipv6Reference")
		defer hostNameParser.Dbg_leave("ipv6Reference")
	}

	for hostNameParser.lexer.HasMoreChars() {
		la, err := hostNameParser.lexer.LookAheadK(0)
		if err != nil {
			return retval.String(), err
		}
		if hostNameParser.lexer.IsHexDigit(la) {
			hostNameParser.lexer.ConsumeK(1)
			retval.WriteByte(la)
		} else if la == '.' ||
			la == ':' ||
			la == '[' {
			hostNameParser.lexer.ConsumeK(1)
			retval.WriteByte(la)
		} else if la == ']' {
			hostNameParser.lexer.ConsumeK(1)
			retval.WriteByte(la)
			return retval.String(), nil
		} else {
			break
		}
	}

	return retval.String(), errors.New("ParseException: Illegal Host name")
}

func (hostNameParser *HostNameParser) GetHost() (h *Host, err error) {
	if Debug.ParserDebug {
		hostNameParser.Dbg_enter("host")
		defer hostNameParser.Dbg_leave("host")
	}

	var hname bytes.Buffer
	var next byte
	var nextToks string

	//IPv6 referene
	if next, err = hostNameParser.lexer.LookAheadK(0); err == nil && next == '[' {
		if nextToks, err = hostNameParser.Ipv6Reference(); err != nil {
			return nil, err
		}
		hname.WriteString(nextToks)
	} else { //IPv4 address or hostname
		if nextToks, err = hostNameParser.DomainLabel(); err != nil {
			return nil, err
		}
		hname.WriteString(nextToks)
		// Bug reported by Stuart Woodsford (used to barf on
		// more than 4 components to the name).
		for hostNameParser.lexer.HasMoreChars() {
			// Reached the end of the buffer.
			if next, err = hostNameParser.lexer.LookAheadK(0); err == nil && next == '.' {
				hostNameParser.lexer.ConsumeK(1)
				if nextToks, err = hostNameParser.DomainLabel(); err != nil {
					return nil, err
				}
				hname.WriteString(".")
				hname.WriteString(nextToks)
			} else {
				break
			}
		}
	}

	hostname := hname.String()

	if hostname == "" {
		return nil, errors.New("ParseException: Illegal Host name")
	} else {
		return NewHost(hostname), nil
	}
}

func (hostNameParser *HostNameParser) GetHostPort() (hp *HostPort, ParseException error) {
	if Debug.ParserDebug {
		hostNameParser.Dbg_enter("hostPort")
		defer hostNameParser.Dbg_leave("hostPort")
	}

	host, err := hostNameParser.GetHost()
	if err != nil {
		return nil, err
	}
	hp = &HostPort{host: host, port: -1}
	// Has a port?
	if hostNameParser.lexer.HasMoreChars() {
		if next, err := hostNameParser.lexer.LookAheadK(0); err == nil && next == ':' {
			hostNameParser.lexer.ConsumeK(1)

			port, err := hostNameParser.lexer.Number()
			if err != nil {
				return nil, err
			}
			hp.SetPort(port)

		}
	}
	return hp, nil
}
