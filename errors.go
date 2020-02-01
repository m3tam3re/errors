// This package is based on the package "upspin.io/errors"
// Copyright 2016 The Upspin Authors. All rights reserved.
//
// for more details on the mentioned package please read this
// great article: https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html

package errors

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
)

// Error implements the error interface. Fields can be added easily. Please be aware
// that most fields have there own type to distinguish them from each other by a type switch
// in func E
type Error struct {
	Path Path
	Op   Op
	Kind Kind
	Msg  string
	Ref  Ref
	Err  error
}

type Path string

type Op string

type Kind uint8

type Ref string

// Kinds of errors
// new items can be added to the end
const (
	Other      Kind = iota //Unknown
	Invalid                //Operation not allowed for this item
	Permission             //No Permission to access this item
	IO                     //Error reading/writing, could be file network etc.
	Exists                 //Item already exists
	NotExist               //Item dows not exists
	IsDir                  //Item is a directory
	NotDir                 //Item is nor a directory
	NotEmpty               //Item is not empty
	Private                //Information requested is private
	Internal               //Internal Error
	BrokenLink             //Link target could not be found
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "unknown error"
	case Invalid:
		return "invalid operation"
	case Permission:
		return "permission denied"
	case IO:
		return "I/O error"
	case Exists:
		return "item already exists"
	case NotExist:
		return "item does not exist"
	case IsDir:
		return "item is a directoy"
	case NotDir:
		return "item is not a directory"
	case NotEmpty:
		return "directory is not empty"
	case Private:
		return "requested item is private"
	case Internal:
		return "internal error"
	case BrokenLink:
		return "link target could not be found"
	default:
		return "unknown error kind"
	}
}

// E populates the Error fields from the given arguments
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("errors.E call without arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Path:
			e.Path = arg
		case Op:
			e.Op = arg
		case Kind:
			e.Kind = arg
		case string:
			e.Msg = arg
		case Ref:
			e.Ref = arg
		case *Error:
			cp := *arg
			e.Err = &cp
		case error:
			e.Err = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: received broken call from %s:%d: %v", file, line, args)
			return fmt.Errorf("unknown type: %T, value: %v", arg, arg)
		}
	}
	return e
}

// pad appends a string to the buffer if the buffer is not empty
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

// Error() creates an error string using the fields of the Error struct
func (e *Error) Error() string {
	b := new(bytes.Buffer)
	if e.Op != "" {
		pad(b, "|")
		b.WriteString(string(e.Op))
	}
	if e.Path != "" {
		pad(b, "|")
		b.WriteString(string(e.Path))
	}
	if e.Kind != 0 {
		pad(b, "|")
		b.WriteString(e.Kind.String())
	}
	if e.Msg != "" {
		pad(b, "|")
		b.WriteString(e.Msg)
	}
	if e.Ref != "" {
		pad(b, "|")
		b.WriteString(e.Msg)
	}
	if e.Err != nil {
		pad(b, "|")
		b.WriteString(e.Err.Error())
	}
	return b.String()
}
