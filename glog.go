// Modified from:
// log.go, Copyright 2009 The Go Authors, BSD-style.

// InLogger = Indexed Input Logger (i.e. has an index nr, 00-99)
// SGLogger = Standard Single(ton) Global Logger

// https://dave.cheney.net/2015/11/05/lets-talk-about-logging
// I believe that there are only two things you should log:
// - Things that developers care about when developing or debugging SW.
// - Things that users care about when using your SW.
// Obviously these are debug and info levels, respectively.

// Package glog implements a simple logging package.
// It defines a type `InLogger` with methods for formatting output.
//
// It also has a predefined "standard singleton global" `SGLogger`
// accessible thru helper functions `Print[f|ln], Fatal[f|ln]`
//
// All loggers are thread-safe for concurrent use. Every log
// message is output on a separate line: if the message being
// printed does not end in a newline, the logger will add one.
//
package miscutils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	S "strings"
	"sync"
	"time"

	FP "path/filepath"

	"github.com/fatih/color"
	SU "github.com/fbaube/stringutils"
	WU "github.com/fbaube/wasmutils"
	"github.com/jimlawless/whereami"
)

type Writer interface {
	Write(p []byte) (n int, err error)
}

func TracedError(e error) error {
	return errors.New(
		e.Error() +
			"\n\t" + whereami.WhereAmI(2) +
			"\n\t" + whereami.WhereAmI(3) +
			"\n\t" + whereami.WhereAmI(4))
}

var SessionLogger *Logger
var UserHomeDir string

func init() {
	// So is non-nil
	// func New(out io.Writer, prefix string, flag int) *Logger
	SessionLogger = New(os.Stdout, "log> ", 0)

	if !WU.IsWasm() {
		UserHomeDir, _ = os.UserHomeDir()
		if S.HasSuffix(UserHomeDir, "/") {
			println("--> Trimming trailing slash from UserHomeDir:", UserHomeDir)
			UserHomeDir = S.TrimSuffix(UserHomeDir, "/")
			println("--> UserHomeDir:", UserHomeDir)
		}
	}
}

func ErrorTrace(w io.Writer, e error) {
	fmt.Fprintf(w, SU.Rfg(SU.Ybg(" ** ERROR ** ")))
	color.Set(color.FgHiRed)
	fmt.Fprintf(w, "\n"+e.Error()+"\n")
	color.Unset()
}

// These flags define which text to prefix to each log entry generated by the Logger.
// Bits are or'ed together to control what's printed.
// There is no control over the order they appear (the order listed
// here) or the format they present (as described in the comments).
// The prefix is followed by a colon only when Llongfile or Lshortfile
// is specified.
// For example, flags Ldate | Ltime (or LstdFlags) produce,
//	2009/01/23 01:23:23 message
// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// A Logger is an active logging object that generates output lines for an
// `io.Writer`. Each log call makes one call to the Writer's `Write(..)`.
// A Logger has a mutex and can be called simultaneously from multiple
// goroutines; it is guaranteed to serialize access to the Writer.
type Logger struct {
	// The fields below are from stdlib logger.Logger
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix string     // prefix to write at beginning of each line
	flag   int        // properties
	out    io.Writer  // destination for output
	buf    []byte     // for accumulating text to write
	// The fields below are unique to `GLogger`
	creatime time.Time
	path     string // absolute path; can be "" for `os.Stdout os.Stderr`
}

func (L *Logger) Close() {
	if F, ok := L.out.(*os.File); ok {
		F.Sync()
		F.Close()
	} else if C, ok := L.out.(io.Closer); ok {
		C.Close()
	}
}

// NewAtPath creates a new Logger.
// - arg `path` is the destination for writes.
// - arg `prefix` will start each log line.
func NewAtPath(path string, prefix string) *Logger {
	var afp string
	var e error
	afp, e = FP.Abs(path)
	if e != nil {
		panic("L nu Absol")
	}
	F, e := os.Create(afp)
	if e != nil {
		panic("L nu Fopen")
	}
	_ = F.Truncate(0)
	if e != nil {
		panic("L nu Trunc")
	}
	L := New(F, prefix, 0)
	L.path = afp
	_, _ = L.out.Write([]byte("## as " + afp + "\n##\n"))
	return L
}

// New creates a new Logger.
// It is up to the caller to set `Logger.Path` if desired.
// - arg `out` is the destination for writes.
// - arg `prefix` will start each log line.
// - arg `flag` is logging properties (0, but keep for compatibility)
func New(out io.Writer, prefix string, flag int) *Logger {
	t := time.Now()
	if flag&LUTC != 0 {
		t = t.UTC()
	}
	// s := "Test message" + "\n"
	L := &Logger{out: out, prefix: prefix, flag: flag}
	L.creatime = t
	// Write header line
	// func (l *Logger) Output(calldepth int, s string) error {
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "F??L??"
		line = 0
	}
	L.buf = L.buf[:0]
	// formatHeader(&L.buf, now, file, line)
	buf := &L.buf
	*buf = append(*buf, SU.B("## LOGFILE created at ")...)
	// if l.flag&Ldate != 0 {
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')
	// if l.flag&(Ltime|Lmicroseconds) != 0 {
	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
	// if l.flag&Lmicroseconds != 0 {
	*buf = append(*buf, '.')
	itoa(buf, t.Nanosecond()/1e6, 3)
	// }
	*buf = append(*buf, SU.B(" local time\n## by ")...)
	// }
	// }
	// if l.flag&(Lshortfile|Llongfile) != 0 {
	*buf = append(*buf, file...)
	*buf = append(*buf, ':')
	itoa(buf, line, -1)
	// *buf = append(*buf, ": "...)
	// }
	// L.buf = append(L.buf, s...)
	L.buf = append(L.buf, '\n')
	_, err := L.out.Write(L.buf)
	if err != nil {
		println("LOG WRITE ERROR:", L.path, err.Error())
		// panic("LOG WRITE ERROR")
		L.Close()
	}
	return L
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided).
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	*buf = append(*buf, l.prefix...)
	if l.flag&LUTC != 0 {
		t = t.UTC()
	}
	var elapsed time.Duration
	elapsed = time.Since(l.creatime)
	*buf = append(*buf, []byte(elapsed.String())...)
	*buf = append(*buf, ' ')

	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
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

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(calldepth int, s string) error {
	now := time.Now() // get this early.
	var file string
	var line int
	/*
		if l.mu == nil {
			l.mu = new(sync.Mutex)
		}
	*/
	println("No mu.lock")
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) { l.Output(2, fmt.Sprint(v...)) }

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...interface{}) { l.Output(2, fmt.Sprintln(v...)) }

// Prefix returns the output prefix for the logger.
func (l *Logger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// Writer returns the output destination for the logger.
func (l *Logger) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}

// Blare is used for errors that will stop processing.
func Blare(s string) {
	SessionLogger.Printf(s)
	fmt.Fprintf(os.Stderr, s)
}