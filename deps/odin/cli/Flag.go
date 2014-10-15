package cli

import "github.com/Sam-Izdat/pogo/deps/odin/cli/values"

// A Flag represents the state of a flag.
type Flag struct {
	Name     string // name as it appears on command line
	Usage    string // help message
	DefValue string // default value (as text); for usage message

	value values.Value // value as set
}

type boolFlag interface {
	values.Value
	IsBoolValue() bool
}

var noString string
var noValue = values.NewString("", &noString)
var noFlag = &Flag{
	Name:     "",
	Usage:    "",
	DefValue: noValue.String(),
	value:    noValue,
}
