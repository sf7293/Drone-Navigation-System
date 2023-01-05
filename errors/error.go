package errors

// Error represents an error in the context of this service. Error implements go's error interface and can be used as one.
// When processing errors, use the getter functions to access data when possible (e.g. GetData(Err) instead of Err.Data)
// When creating new errors, either use E(...) convenience function, or create a new Error struct directly.
type Error struct {
	// Op is the operation that caused the error.
	Op Op
	// Err is the wrapped error, which was caused while Op was executing.
	Err error
	// Data is a map[string]interface{} under the hood, which allows for collection of any arbitrary data about Error.
	Data Data
}

// Cause is the root cause of the error, implements Error(), and should be the first error that is wrapped.
type Cause struct {
	Severity     Severity
	Exposable    bool
	TitleLabel   string
	MessageLabel string
}

// Op represents an operation. Each function in the code should define its own Op with a unique string.
type Op string

type Severity int

type Data map[string]interface{}

// E is a convenience function to initialize a new Error.
func E(op Op, wrappedError error, data ...Data) (Err *Error) {
	e := &Error{Op: op, Err: wrappedError}

	e.Data = make(map[string]interface{})

	for _, d := range data {
		e.Data.Merge(d)
	}

	return e
}

// Ops returns the stack of operations that lead to the error.
func (e *Error) Ops() (ops []Op) {
	ops = []Op{e.Op}

	// loop though Errors to the first error
	err := e
	for {
		var ok bool
		err, ok = err.Err.(*Error)
		if !ok {
			break
		}

		ops = append(ops, err.Op)
	}

	return ops
}

func GetSeverity(err error) (severity Severity) {
	cause, ok := err.(Cause)
	if ok {
		return cause.Severity
	}

	Err, ok := err.(*Error)
	if !ok {
		return SeverityCritical
	}

	// recurse through Errors
	return GetSeverity(Err.Err)
}

func GetCause(err error) (cause Cause) {
	cause, ok := err.(Cause)
	if ok {
		return cause
	}

	Err, ok := err.(*Error)
	if !ok {
		return ErrUnexpected
	}

	// recurse through Errors
	return GetCause(Err.Err)
}

// GetData iterates over the Errors stack and returns a merged Data object, containing the Data from all errors in the stack.
// If multiple errors in the stack contain the same keys, the most recent value is used for that key.
func GetData(err error) (data Data) {
	Err, ok := err.(*Error)
	if !ok {
		return make(map[string]interface{})
	}

	data = GetData(Err.Err)
	data.Merge(Err.Data)

	return data
}

func (e *Error) Error() string {
	if e.Err == nil {
		return ""
	}

	return e.Err.Error()
}

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityCritical
)

// Merge takes another Data and merges it with this one. If a value is present in both of this and other, the value from
// other will be used.
func (d Data) Merge(other Data) {
	for k, v := range other {
		d[k] = v
	}
}

func (c Cause) Error() string {
	return c.TitleLabel
}
