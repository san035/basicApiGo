package logger

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"path"
	"runtime"
	"strings"
)

const DeepStackDefault = 2

type ExtendedError struct {
	Err      error
	Stack    string        // функция в которой ошибка
	Addition []interface{} // доп информация
}

// New Новая ошибка с текстом textError
func New(textError string, addition ...interface{}) *ExtendedError {

	resp := &ExtendedError{Err: errors.New(textError), Stack: StackFunc()}
	if len(addition) != 0 {
		resp.Addition = addition
	}
	return resp
}

// New Новая ошибка с текстом textError и глубиной стека 1
func NewWithDeep1(textError string, addition ...interface{}) *ExtendedError {

	resp := &ExtendedError{Err: errors.New(textError), Stack: StackWithDeep(2, 2)}
	if len(addition) != 0 {
		resp.Addition = addition
	}
	return resp
}

func (e *ExtendedError) Error() string {
	return e.Err.Error()
}

func (e *ExtendedError) GetAddition() interface{} {
	return e.Addition
}

// Wrap возврашает новывй экземпляр ExtendedError с заполненым Addition
func Wrap(err *error, addition ...interface{}) *ExtendedError {
	switch errExt := (*err).(type) {
	case *ExtendedError:
		errExt.Addition = append(errExt.Addition, addition)
		return errExt
	}

	return &ExtendedError{Err: *err, Addition: addition, Stack: StackWithDeep(2, DeepStackDefault+1)}
}

// Wrap возврашает новывй экземпляр ExtendedError с заполненым Addition
func WrapWithDeep1(err *error, addition ...interface{}) *ExtendedError {
	switch errExt := (*err).(type) {
	case *ExtendedError:
		if addition != nil {
			if errExt.Addition == nil {
				errExt.Addition = addition
			} else {
				errExt.Addition = append(errExt.Addition, addition)
			}
		}
		return errExt
	}

	return &ExtendedError{Err: *err, Addition: addition, Stack: StackWithDeep(2, 2)}
}

// StackFunc возаращает стек возова функции, например: main.go main.x 50;y 55;z 58;
func StackFunc() (listNameFunc string) {
	return StackWithDeep(2, DeepStackDefault+1)
}

// StackFunc возаращает стек возова функции c глубиной от levelDeepFrom до levelDeepTo
func StackWithDeep(levelDeepFrom, levelDeepTo int) (listNameFunc string) {
	var lastModul string
	for i := levelDeepTo; i >= levelDeepFrom; i-- {
		pc, file, lineNumber, ok := runtime.Caller(i)
		file = path.Base(file)

		if !ok {
			return
		}
		nameFunc := file + " " + runtime.FuncForPC(pc).Name()
		indexLastPoint := strings.LastIndex(nameFunc, ".")
		if indexLastPoint > 0 {
			newModule := nameFunc[0:indexLastPoint]
			if lastModul == newModule {
				nameFunc = nameFunc[indexLastPoint+1:]
			} else {
				lastModul = newModule
			}
		}

		listNameFunc += fmt.Sprintf("%s %d;", nameFunc, lineNumber)
	}

	return
}

// AddAddition возвращает zerolog.Event с записанными полями Stack и Addition
func AddAddition(err *error, e *zerolog.Event) *zerolog.Event {
	switch errExt := (*err).(type) {
	case *ExtendedError:
		respEvent := e.Str("Stack", errExt.Stack)
		if len(errExt.Addition) != 0 {
			respEvent = respEvent.Interface("Addition", errExt.Addition)
		}
		return respEvent
	}

	return e
}

// Error Возвращает log.Error() с записаннами данными из err
func Error(err *error) (respEvent *zerolog.Event) {
	errExt, ok := (*err).(*ExtendedError)
	if ok {
		respEvent = log.Error().Err(errExt.Err).Str("Stack", errExt.Stack)
		if len(errExt.Addition) != 0 {
			respEvent = respEvent.Interface("Addition", errExt.Addition)
		}
		return
	}

	respEvent = log.Error().Err(*err)
	return
}
