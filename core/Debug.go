package core

/**
*   A class to do debug printfs
 */

 // Debugger type
type Debugger struct {
	Debug       bool
	ParserDebug bool
}
// Debug : Exported Object
var Debug = Debugger{true, false}

func (debugger *Debugger) print(s string) {
	if debugger.Debug {
		LogWrite.LogMessage(s)
	}
}

func (debugger *Debugger) println(s string) {
	if debugger.Debug {
		LogWrite.LogMessage(s + "\n")
	}
}
