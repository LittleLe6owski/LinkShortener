//一个简单的生成错误的帮助类
package goerr

import (
	"bytes"
	"fmt"
	"strconv"
)

// new errorContext with error and string
//
// return error
func Errorf(err error, format string, p ...interface{}) error {
	return &errorContext{text: fmt.Sprintf(format, p...), err: err, line: -1}
}

// new errorContext with error
//
//  return  error
func Error(err error) *errorContext {
	if e, ok := err.(*errorContext); ok {
		return e
	} else {
		return &errorContext{err: err, line: -1}
	}
}

// new errorContext with string
func String(format string, p ...interface{}) *errorContext {
	return &errorContext{text: fmt.Sprintf(format, p...), line: -1}
}

//new error
type errorContext struct {
	code int
	text string
	file string
	line int
	err  error
}

func (e *errorContext) Format(format string, p ...interface{}) *errorContext {
	e.text = fmt.Sprintf(format, p...)
	return e
}
func (e *errorContext) Code(code int) *errorContext {
	e.code = code
	return e
}
func (e *errorContext) Line(line int) *errorContext {
	e.line = line
	return e
}
func (e *errorContext) File(file string) *errorContext {
	e.file = file
	return e
}
func (e *errorContext) E(err error) *errorContext {
	e.err = err
	return e
}
func (e *errorContext) AttachE(err error) {
	if er, ok := e.err.(*errorContext); er != nil && ok {
		er.AttachE(err)
	} else {
		e.err = err
	}
}

// trace error message
func (e *errorContext) Trace() string {
	var buffer bytes.Buffer
	if e.file != "" {
		buffer.WriteString("File: ")
		buffer.WriteString(e.file)
		buffer.WriteRune('\t')
	}
	if e.line != -1 {
		bs := strconv.AppendInt(nil, int64(e.line), 10)
		buffer.WriteString("Line: ")
		buffer.Write(bs)
		buffer.WriteRune('\t')
	}
	if e.code != 0 {
		bs := strconv.AppendInt(nil, int64(e.code), 10)
		buffer.WriteString("Error Code: ")
		buffer.Write(bs)
		buffer.WriteRune('\t')
	}
	if e.text != "" {
		buffer.WriteString("Text: ")
		buffer.WriteString(e.text)
		buffer.WriteRune('\t')
	}
	buffer.WriteRune('\n')
	if e.err != nil {
		if err, ok := e.err.(*errorContext); ok {
			buffer.WriteString(err.Trace())
		} else {
			buffer.WriteString("Trace:")
			buffer.WriteString(e.err.Error())
			buffer.WriteRune('\n')
		}
	}
	return buffer.String()
}

//实现error接口
func (e *errorContext) Error() string {
	if e.text != "" {
		return e.text
	} else if e.err != nil {
		return e.err.Error()
	} else {
		return e.Trace()
	}
}
