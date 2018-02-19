// Package errors Modified from Upspin for generic use case.
package errors

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/wirepair/ewserver/types"
)

// Error is the type that implements the error interface.
// It contains a number of fields, each of different type.
// An Error value may leave some values unset.
type Error struct {
	// Path is the path name of the item being accessed.
	Path types.PathName
	// User is the username of the user attempting the operation.
	User types.UserName
	// Op is the operation being performed, usually the name of the method
	// being invoked (Get, Put, etc.). It should not contain an at sign @.
	Op Op
	// Kind is the class of error, such as permission failure,
	// or "Other" if its class is unknown or irrelevant.
	Kind Kind
	// The underlying error that triggered this one, if any.
	Err error
}

func (e *Error) isZero() bool {
	return e.Path == "" && e.User == "" && e.Op == "" && e.Kind == 0 && e.Err == nil
}

// Op describes an operation, usually as the package and method,
// such as "key/server.Lookup".
type Op string

// Separator is the string used to separate nested errors. By
// default, to make errors easier on the eye, nested errors are
// indented on a new line. A server may instead choose to keep each
// error on a single line by modifying the separator string, perhaps
// to ":: ".
var Separator = ":\n\t"

// Kind defines the kind of error this is, mostly for use by systems
// such as FUSE that must act differently depending on the error.
type Kind uint8

// Kinds of errors.
//
// The values of the error kinds new items must be added only to the end.
const (
	Other         Kind = iota // Unclassified error. This value is not printed in the error message.
	Invalid                   // Invalid operation for this type of item.
	Permission                // Permission denied.
	IO                        // External I/O error such as network failure.
	Exist                     // Item already exists.
	NotExist                  // Item does not exist.
	IsDir                     // Item is a directory.
	NotDir                    // Item is not a directory.
	NotEmpty                  // Directory not empty.
	Private                   // Information withheld.
	Internal                  // Internal error or inconsistency.
	CannotDecrypt             // No wrapped key for user with read access.
	CannotDecode              // Unable to decode an item
	CannotEncode              // Unable to encode an item
	Transient                 // A transient error.
	BrokenLink                // Link target does not exist.
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case Invalid:
		return "invalid operation"
	case Permission:
		return "permission denied"
	case IO:
		return "I/O error"
	case Exist:
		return "item already exists"
	case NotExist:
		return "item does not exist"
	case BrokenLink:
		return "link target does not exist"
	case IsDir:
		return "item is a directory"
	case NotDir:
		return "item is not a directory"
	case NotEmpty:
		return "directory not empty"
	case Private:
		return "information withheld"
	case Internal:
		return "internal error"
	case CannotDecrypt:
		return `no wrapped key for user; owner must "upspin share -fix"`
	case Transient:
		return "transient error"
	}
	return "unknown error kind"
}

// E builds an error value from its arguments.
// There must be at least one argument or E panics.
// The type of each argument determines its meaning.
// If more than one argument of a given type is presented,
// only the last one is recorded.
//
// The types are:
//	PathName
//		The path name of the item being accessed.
//	UserName
//		The name of the user attempting the operation.
//	errors.Op
//		The operation being performed, usually the method
//		being invoked (Get, Put, etc.).
//	string
//		Treated as an error message and assigned to the
//		Err field after a call to errors.Str. To avoid a common
//		class of misuse, if the string contains an @, it will be
//		treated as a PathName or UserName, as appropriate. Use
//		errors.Str explicitly to avoid this special-casing.
//	errors.Kind
//		The class of error, such as permission failure.
//	error
//		The underlying error that triggered this one.
//
// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.
//
// If Kind is not specified or Other, we set it to the Kind of
// the underlying error.
//
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case types.PathName:
			e.Path = arg
		case types.UserName:
			e.User = arg
		case Op:
			e.Op = arg
		case string:
			// Someone might accidentally call us with a user or path name
			// that is not of the right type. Take care of that and log it.
			if strings.Contains(arg, "@") {
				_, file, line, _ := runtime.Caller(1)
				log.Printf("errors.E: unqualified type for %q from %s:%d", arg, file, line)
				if strings.Contains(arg, "/") {
					if e.Path == "" { // Don't overwrite a valid path.
						e.Path = types.PathName(arg)
					}
				} else {
					if e.User == "" { // Don't overwrite a valid user.
						e.User = types.UserName(arg)
					}
				}
				continue
			}
			e.Err = Str(arg)
		case Kind:
			e.Kind = arg
		case *Error:
			// Make a copy
			copy := *arg
			e.Err = &copy
		case error:
			e.Err = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			return Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}

	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}

	// The previous error was also one of ours. Suppress duplications
	// so the message won't contain the same kind, file name or user name
	// twice.
	if prev.Path == e.Path {
		prev.Path = ""
	}
	if prev.User == e.User {
		prev.User = ""
	}
	if prev.Kind == e.Kind {
		prev.Kind = Other
	}
	// If this error has Kind unset or Other, pull up the inner one.
	if e.Kind == Other {
		e.Kind = prev.Kind
		prev.Kind = Other
	}
	return e
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	if e.Op != "" {
		pad(b, ": ")
		b.WriteString(string(e.Op))
	}
	if e.Path != "" {
		pad(b, ": ")
		b.WriteString(string(e.Path))
	}
	if e.User != "" {
		if e.Path == "" {
			pad(b, ": ")
		} else {
			pad(b, ", ")
		}
		b.WriteString("user ")
		b.WriteString(string(e.User))
	}
	if e.Kind != 0 {
		pad(b, ": ")
		b.WriteString(e.Kind.String())
	}
	if e.Err != nil {
		// Indent on new line if we are cascading non-empty Upspin errors.
		if prevErr, ok := e.Err.(*Error); ok {
			if !prevErr.isZero() {
				pad(b, Separator)
				b.WriteString(e.Err.Error())
			}
		} else {
			pad(b, ": ")
			b.WriteString(e.Err.Error())
		}
	}
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

// Recreate the errors.New functionality of the standard Go errors package
// so we can create simple text errors when needed.

// Str returns an error that formats as the given text. It is intended to
// be used as the error-typed argument to the E function.
func Str(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// Errorf is equivalent to fmt.Errorf, but allows clients to import only this
// package for all error handling.
func Errorf(format string, args ...interface{}) error {
	return &errorString{fmt.Sprintf(format, args...)}
}

// Match compares its two error arguments. It can be used to check
// for expected errors in tests. Both arguments must have underlying
// type *Error or Match will return false. Otherwise it returns true
// iff every non-zero element of the first error is equal to the
// corresponding element of the second.
// If the Err field is a *Error, Match recurs on that field;
// otherwise it compares the strings returned by the Error methods.
// Elements that are in the second argument but not present in
// the first are ignored.
//
// For example,
//	Match(errors.E(upspin.UserName("joe@schmoe.com"), errors.Permission), err)
// tests whether err is an Error with Kind=Permission and User=joe@schmoe.com.
func Match(err1, err2 error) bool {
	e1, ok := err1.(*Error)
	if !ok {
		return false
	}
	e2, ok := err2.(*Error)
	if !ok {
		return false
	}
	if e1.Path != "" && e2.Path != e1.Path {
		return false
	}
	if e1.User != "" && e2.User != e1.User {
		return false
	}
	if e1.Op != "" && e2.Op != e1.Op {
		return false
	}
	if e1.Kind != Other && e2.Kind != e1.Kind {
		return false
	}
	if e1.Err != nil {
		if _, ok := e1.Err.(*Error); ok {
			return Match(e1.Err, e2.Err)
		}
		if e2.Err == nil || e2.Err.Error() != e1.Err.Error() {
			return false
		}
	}
	return true
}

// Is reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == kind
	}
	if e.Err != nil {
		return Is(kind, e.Err)
	}
	return false
}
