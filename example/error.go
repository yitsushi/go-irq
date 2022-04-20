package main

import "fmt"

type UnknownAPICallError struct {
	Name    string
	Payload interface{}
}

func (e UnknownAPICallError) Error() string {
	return fmt.Sprintf("unknown api call: %s", e.Name)
}

func newUnknownAPICallError(name string, payload interface{}) UnknownAPICallError {
	return UnknownAPICallError{Name: name, Payload: payload}
}
