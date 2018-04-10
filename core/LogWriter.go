package core

import (
	"os"
)

/**
*  Log System Errors. Also used for debugging log.
 */
/** Dont trace
 */
const TRACE_NONE = 0

/** Trace message processing
 */
const TRACE_MESSAGES = 16

/** Trace exception processing
 */
const TRACE_EXCEPTION = 17

/** Debug trace level (all tracing enabled).
 */
const TRACE_DEBUG = 32

/** Name of the log file in which the trace is written out
 * (default is /tmp/sipserverlog.txt)
 */
type LogWriter struct {

	/** Print writer that is used to write out the log file.
	 */
	printWriter *os.File

	/** Flag to indicate that logging is enabled. logWriter needs to be
	 * static and public in order to globally turn logging on or off.
	 * logWriter is static for efficiency reasons (the java compiler will not
	 * generate the logging code if logWriter is set to false).
	 */
	traceLevel int

	logFileName string

	needsLogging bool

	lineCount int
}

var LogWrite = LogWriter{nil, TRACE_NONE, "debug.log", false, 0}

/** Set the log file name
*@param name is the name of the log file to set.
 */
func (logWriter *LogWriter) SetLogFileName(name string) {
	logWriter.logFileName = name
}

func (logWriter *LogWriter) LogMessageToFile(message, logFileName string) {
	var err error
	logWriter.printWriter, err = os.OpenFile(logFileName, os.O_APPEND, 0)
	if err != nil {
		println("Can't open file in LogMessageToFile")
		return
	}
	defer logWriter.printWriter.Close()

	logWriter.printWriter.WriteString(" ---------------------------------------------- ")
	logWriter.printWriter.WriteString(message)
}

func (logWriter *LogWriter) checkLogFile() {
	if logWriter.printWriter != nil {
		return
	}
	if logWriter.logFileName == "" {
		return
	}

	var err error
	logWriter.printWriter, err = os.OpenFile(logWriter.logFileName, os.O_APPEND, 0)
	if err != nil {
		println("Can't open file in checkLogFile")
		return
	}
}

func (logWriter *LogWriter) println(message string) {
	for i := 0; i < len(message); i++ {
		if message[i] == '\n' {
			logWriter.lineCount++
		}
	}
	logWriter.checkLogFile()
	if logWriter.printWriter != nil {
		logWriter.printWriter.WriteString(message)
	}
	logWriter.lineCount++
}

/** Log a message into the log file.
 * @param message message to log into the log file.
 */
func (logWriter *LogWriter) LogMessage(message string) {
	if !logWriter.needsLogging {
		return
	}

	logWriter.checkLogFile()
	logWriter.println(message)
}

/** Set the trace level for the stack.
 */
func (logWriter *LogWriter) SetTraceLevel(level int) {
	logWriter.traceLevel = level
}

/** Get the trace level for the stack.
 */
func (logWriter *LogWriter) GetTraceLevel() int {
	return logWriter.traceLevel
}
