// Package middleware provides useful middleware for a [clir.Router].
package middleware

import (
	"errors"
	"flag"
	"io"
	"strconv"

	"maragu.dev/clir"
)

// Flags middleware allows you to set flags on a route.
func Flags(cb func(fs *flag.FlagSet)) clir.Middleware {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	cb(fs)

	return func(next clir.Runner) clir.Runner {
		return clir.RunnerFunc(func(ctx clir.Context) error {
			fs.SetOutput(ctx.Err)
			if err := fs.Parse(ctx.Args); err != nil {
				if errors.Is(err, flag.ErrHelp) {
					return nil
				}
				return err
			}
			ctx.Args = fs.Args()
			return next.Run(ctx)
		})
	}
}

// ArgSet is like [flag.FlagSet] but for positional arguments.
// The order of calls is significant.
type ArgSet struct {
	Usage  func()
	args   []string
	formal []*flag.Flag
	w      io.Writer
}

// String defines a string positional argument.
func (a *ArgSet) String(name string, value string, usage string) *string {
	p := new(string)
	a.StringVar(p, name, value, usage)
	return p
}

// StringVar defines a string positional argument with a pointer.
func (a *ArgSet) StringVar(p *string, name string, value string, usage string) {
	a.Var(newStringValue(value, p), name, usage)
}

// Int defines an int positional argument.
func (a *ArgSet) Int(name string, value int, usage string) *int {
	p := new(int)
	a.IntVar(p, name, value, usage)
	return p
}

// IntVar defines an int positional argument with a pointer.
func (a *ArgSet) IntVar(p *int, name string, value int, usage string) {
	a.Var(newIntValue(value, p), name, usage)
}

// Bool defines a bool positional argument.
func (a *ArgSet) Bool(name string, value bool, usage string) *bool {
	p := new(bool)
	a.BoolVar(p, name, value, usage)
	return p
}

// BoolVar defines a bool positional argument with a pointer.
func (a *ArgSet) BoolVar(p *bool, name string, value bool, usage string) {
	a.Var(newBoolValue(value, p), name, usage)
}

// Float64 defines a float64 positional argument.
func (a *ArgSet) Float64(name string, value float64, usage string) *float64 {
	p := new(float64)
	a.Float64Var(p, name, value, usage)
	return p
}

// Float64Var defines a float64 positional argument with a pointer.
func (a *ArgSet) Float64Var(p *float64, name string, value float64, usage string) {
	a.Var(newFloat64Value(value, p), name, usage)
}

func (a *ArgSet) Var(value flag.Value, name string, usage string) {
	f := &flag.Flag{Name: name, Usage: usage, Value: value, DefValue: value.String()}
	a.formal = append(a.formal, f)
}

func (a *ArgSet) Parse(args []string) error {
	a.args = args

	// Reset all positional arguments to their declared defaults before parsing.
	for _, f := range a.formal {
		if err := f.Value.Set(f.DefValue); err != nil {
			return err
		}
	}

	// Process all positional arguments based on formal definitions.
	for i, f := range a.formal {
		if i < len(args) {
			// Set the value from the args.
			if err := f.Value.Set(args[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *ArgSet) Args() []string {
	// Return any remaining args after the positional ones we consumed
	if len(a.formal) < len(a.args) {
		return a.args[len(a.formal):]
	}
	return []string{}
}

func (a *ArgSet) SetOutput(w io.Writer) {
	a.w = w
}

// Args middleware allows you to set positional arguments on a route.
func Args(cb func(as *ArgSet)) clir.Middleware {
	as := &ArgSet{}
	cb(as)

	return func(next clir.Runner) clir.Runner {
		return clir.RunnerFunc(func(ctx clir.Context) error {
			as.SetOutput(ctx.Err)
			if err := as.Parse(ctx.Args); err != nil {
				return err
			}
			ctx.Args = as.Args()
			return next.Run(ctx)
		})
	}
}

// Value implementations for different types

type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() any { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(val string) error {
	v, err := strconv.ParseInt(val, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*i = intValue(v)
	return nil
}

func (i *intValue) Get() any { return int(*i) }

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Set(val string) error {
	v, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*b = boolValue(v)
	return nil
}

func (b *boolValue) Get() any { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(val string) error {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	*f = float64Value(v)
	return nil
}

func (f *float64Value) Get() any { return float64(*f) }

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }
