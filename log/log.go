package log

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// These flags define which text to prefix to each log entry generated by the Logger.
const (
	// Bits or'ed together to control what's printed. There is no control over the
	// order they appear (the order listed here) or the format they present (as
	// described in the comments).  A colon appears after these items:
	//	2009/0123 01:23:23.123123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota                              // the date: 2009/01/23
	Ltime                                                  // the time: 01:23:23
	Lmicroseconds                                          // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                                              // full file name and line number: /a/b/c/d.go:23
	Lshortfile                                             // final file name element and line number: d.go:23. overrides Llongfile
	Llevel                                                 // preprend the message level
	Lshortlevel                                            // preprend the abbreviated message level (overrides Llevel)
	Lcolored                                               // uses colors around the level name
	LstdFlags     = Ldate | Ltime | Lshortlevel | Lcolored // initial values for the standard logger
	maxPoolCap    = 512
)

var (
	Std  = New(NewIOWriter(os.Stderr, LDebug), LstdFlags, LDefault)
	pool = make(chan []byte, 8)
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer.  Each logging operation makes a single call to
// the Writer's Write method.  A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	flags   int // properties
	level   LLevel
	writers []Writer // destination for output
}

// New creates a new Logger.   The out variable sets the
// destination to which log data will be written.
// The flag argument defines the logging properties.
func New(out Writer, flags int, level LLevel) *Logger {
	logger := &Logger{flags: flags, level: level}
	logger.AddWriter(out)
	return logger
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
// Knows the buffer has capacity.
func itoa(buf *[]byte, i int, wid int) {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		*buf = append(*buf, '0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}
	*buf = append(*buf, b[bp:]...)
}

func (l *Logger) formatHeader(level LLevel, buf *[]byte, t time.Time, file string, line int) {
	if l.flags&(Lshortlevel|Llevel) != 0 {
		var lev string
		if l.flags&Lshortlevel != 0 {
			lev = level.Initial()
		} else {
			lev = level.String()
		}
		*buf = append(*buf, ("[" + lev + "] ")...)
	}
	if l.flags&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flags&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flags&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flags&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flags&(Lshortfile|Llongfile) != 0 {
		if l.flags&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func (l *Logger) FormatMessage(level LLevel, calldepth int, s string) []byte {
	now := time.Now() // get this early.
	var file string
	var line int
	if l.flags&(Lshortfile|Llongfile) != 0 {
		// release lock while getting caller info - it's expensive.
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
	}
	var buf []byte
	select {
	case buf = <-pool:
		buf = buf[:0]
	default:
		buf = make([]byte, 0, maxPoolCap)
	}
	l.formatHeader(level, &buf, now, file, line)
	buf = append(buf, s...)
	return buf
}

func (l *Logger) AddWriter(w Writer) {
	l.writers = append(l.writers, w)
}

func (l *Logger) RemoveWriters() {
	l.writers = nil
}

// Output writes the output for a logging event.  The string s contains
// the text to print after the prefix specified by the flags of the
// Logger.  A newline is appended if the last character of s is not
// already a newline.  Calldepth is used to recover the PC.
func (l *Logger) Write(level LLevel, calldepth int, v ...interface{}) {
	if level >= l.level {
		s := fmt.Sprint(v...)
		msg := l.FormatMessage(level, calldepth, s)
		for _, w := range l.writers {
			if level >= w.Level() {
				w.Write(level, l.flags, msg)
			}
		}
		if cap(msg) <= maxPoolCap {
			select {
			case pool <- msg:
			default:
			}
		}
	}
}

func (l *Logger) Writef(level LLevel, calldepth int, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Write(level, calldepth+1, s)
}

func (l *Logger) Writeln(level LLevel, calldepth int, v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Write(level, calldepth+1, s)
}

func (l *Logger) Logf(level LLevel, format string, v ...interface{}) {
	l.Writef(level, 3, format, v...)
}

func (l *Logger) Log(level LLevel, v ...interface{}) {
	l.Write(level, 3, v...)
}

func (l *Logger) Logln(level LLevel, v ...interface{}) {
	l.Writeln(level, 3, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Writef(LDebug, 3, format, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Write(LDebug, 3, v...)
}

func (l *Logger) Debugln(v ...interface{}) {
	l.Writeln(LDebug, 3, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Writef(LInfo, 3, format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.Write(LInfo, 3, v...)
}

func (l *Logger) Infoln(v ...interface{}) {
	l.Writeln(LInfo, 3, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Writef(LWarning, 3, format, v)
}

func (l *Logger) Warning(v ...interface{}) {
	l.Write(LWarning, 3, v...)
}

func (l *Logger) Warningln(v ...interface{}) {
	l.Writeln(LWarning, 3, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Writef(LError, 3, format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Write(LError, 3, v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	l.Writeln(LError, 3, v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Writef(LDefault, 3, format, v...)
}

func (l *Logger) Print(v ...interface{}) {
	l.Write(LDefault, 3, v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.Writeln(LDefault, 3, v...)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Writef(LFatal, 3, format, v...)
	os.Exit(1)

}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	l.Write(LFatal, 3, v...)
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	l.Writeln(LFatal, 3, v...)
	os.Exit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Write(LFatal, 3, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Write(LFatal, 3, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Write(LFatal, 3, s)
	panic(s)
}

// Nil does nothing if the argument is nil. Otherwise, it's equivalent to Panic().
func (l *Logger) Nil(v interface{}) {
	if v != nil {
		s := fmt.Sprint(v)
		l.Write(LFatal, 3, s)
		panic(v)
	}
}

// Flags returns the output flags for the logger.
func (l *Logger) Flags() int {
	return l.flags
}

// SetFlags sets the output flags for the logger.
func (l *Logger) SetFlags(flags int) {
	l.flags = flags
}

func (l *Logger) Level() LLevel {
	return l.level
}

func (l *Logger) SetLevel(level LLevel) {
	l.level = level
}

// AddWriter adds a writer to the standard logger for the standard logger.
func SetOutput(out Writer) {
	Std.AddWriter(out)
}

// Flags returns the output flags for the standard logger.
func Flags() int {
	return Std.Flags()
}

// SetFlags sets the output flags for the standard logger.
func SetFlags(flag int) {
	Std.SetFlags(flag)
}

func Level() LLevel {
	return Std.Level()
}

func SetLevel(level LLevel) {
	Std.SetLevel(level)
}

// These functions write to the standard logger.

func Debug(v ...interface{}) {
	Std.Write(LDebug, 3, v...)
}

func Debugf(format string, v ...interface{}) {
	Std.Writef(LDebug, 3, format, v...)
}

func Debugln(v ...interface{}) {
	Std.Writeln(LDebug, 3, v...)
}

func Info(v ...interface{}) {
	Std.Writeln(LInfo, 3, v...)
}

func Infof(format string, v ...interface{}) {
	Std.Writef(LInfo, 3, format, v...)
}

func Infoln(v ...interface{}) {
	Std.Writeln(LInfo, 3, v...)
}

func Warning(v ...interface{}) {
	Std.Write(LWarning, 3, v...)
}

func Warningf(format string, v ...interface{}) {
	Std.Writef(LWarning, 3, format, v...)
}

func Warningln(v ...interface{}) {
	Std.Writeln(LWarning, 3, v...)
}

func Error(v ...interface{}) {
	Std.Write(LError, 3, v...)
}

func Errorf(format string, v ...interface{}) {
	Std.Writef(LError, 3, format, v...)
}

func Errorln(v ...interface{}) {
	Std.Writeln(LError, 3, v...)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	Std.Write(LFatal, 3, v...)
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	Std.Writef(LFatal, 3, format, v...)
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	Std.Writeln(LFatal, 3, v...)
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	Std.Write(LPanic, 3, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	Std.Write(LPanic, 3, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	Std.Write(LPanic, 3, s)
	panic(s)
}

// Nil does nothing if the argument is nil. Otherwise, it's equivalent to Panic().
func Nil(v interface{}) {
	if v != nil {
		s := fmt.Sprint(v)
		Std.Write(LPanic, 3, s)
		panic(v)
	}
}

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	Std.Write(LDefault, 3, v...)
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Std.Writef(LDefault, 3, format, v...)
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	Std.Writeln(LDefault, 3, v...)
}
