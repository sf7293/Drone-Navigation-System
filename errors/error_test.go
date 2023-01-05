package errors

import (
	"fmt"
	"testing"
)

func TestOps(t *testing.T) {
	e0 := &Error{}
	e1 := &Error{Op: "1"}
	e2 := &Error{Op: "2", Err: e1}
	e3 := &Error{Op: "3", Err: e2}

	res0 := []Op{""}
	res3 := []Op{"3", "2", "1"}

	tcs := []struct {
		input  *Error
		output []Op
	}{
		{
			input:  e0,
			output: res0,
		}, {
			input:  e3,
			output: res3,
		},
	}

	for _, tc := range tcs {
		res := tc.input.Ops()
		if (res == nil) != (tc.output == nil) {
			t.Errorf("expected nil, got %v", res)
			continue
		}

		if len(res) != len(tc.output) {
			t.Errorf("expected %v, got %v", tc.output, res)
			continue
		}

		for i := range res {
			if res[i] != tc.output[i] {
				t.Errorf("expected %v, got %v", tc.output, res)
				break
			}
		}
	}
}

func TestGetData(t *testing.T) {
	eNoData := &Error{}
	eEmpty := &Error{Data: Data{}}

	eUser10 := &Error{Data: Data{"UserID": 10}}
	eNameBlabla := &Error{Data: Data{"Name": "Blabla"}}

	withCause := E("test with cause", Cause{TitleLabel: "test cause"}, Data{"UserID": 123, "Name": "causy"})

	testCases := []struct {
		input        *Error
		outputUserID interface{}
		outputName   interface{}
	}{
		{
			// nil data
			input:        eNoData,
			outputUserID: nil,
			outputName:   nil,
		}, {
			// empty data
			input:        eEmpty,
			outputUserID: nil,
			outputName:   nil,
		}, {
			// overwrite 10 with 15
			input:        &Error{Err: eUser10, Data: Data{"UserID": 15}},
			outputUserID: 15,
			outputName:   nil,
		}, {
			// overwrite 10 with 0
			input:        &Error{Err: eUser10, Data: Data{"UserID": 0}},
			outputUserID: 0,
			outputName:   nil,
		}, {
			// should not overwrite 10 when UserID does not exist in data
			input:        &Error{Err: eUser10, Data: Data{}},
			outputUserID: 10,
			outputName:   nil,
		}, {
			// Name should also work
			input:        &Error{Err: &Error{Err: eNameBlabla}, Data: Data{}},
			outputUserID: nil,
			outputName:   "Blabla",
		}, {
			// Name should be merged with UserID
			input:        &Error{Err: eNameBlabla, Data: Data{"UserID": 20}},
			outputUserID: 20,
			outputName:   "Blabla",
		}, {
			// Only overwrite Name, and not UserID
			input:        &Error{Err: &Error{Data: Data{"UserID": 15, "Name": "Albalb"}}, Data: Data{"Name": "Guril Anguri"}},
			outputUserID: 15,
			outputName:   "Guril Anguri",
		}, {
			// with cause
			input:        withCause,
			outputUserID: 123,
			outputName:   "causy",
		},
	}

	for _, tc := range testCases {
		res := GetData(tc.input)
		if res == nil {
			t.Errorf("expected a value for Data, got nil")
			continue
		}

		if res["UserID"] != tc.outputUserID {
			t.Errorf("expected UserID %v, got %v", tc.outputUserID, res["UserID"])
			continue
		}
	}
}

func TestGetSeverity(t *testing.T) {
	eInfo := &Error{Err: Cause{Severity: SeverityInfo}}
	eWarn := &Error{Err: Cause{Severity: SeverityWarning}}
	eCritical := &Error{Err: Cause{Severity: SeverityCritical}}
	eUnwrapped := &Error{Err: fmt.Errorf("this should have been wrapped")}

	tcs := []struct {
		input  *Error
		output Severity
	}{
		{
			input:  &Error{Err: eWarn},
			output: SeverityWarning,
		}, {
			input:  eInfo,
			output: SeverityInfo,
		}, {
			input:  eCritical,
			output: SeverityCritical,
		}, {
			input:  eUnwrapped,
			output: SeverityCritical,
		},
	}

	for _, tc := range tcs {
		res := GetSeverity(tc.input)
		if res != tc.output {
			t.Errorf("expected %v, got %v", tc.output, res)
			continue
		}
	}
}

func TestGetCause(t *testing.T) {
	input1 := ErrInternal
	input2 := &Error{Err: input1}
	input3 := &Error{Err: input2}

	tcs := []struct {
		input  *Error
		output Cause
	}{
		{
			input:  &Error{Err: input1},
			output: ErrInternal,
		}, {
			input:  input2,
			output: ErrInternal,
		}, {
			input:  input3,
			output: ErrInternal,
		},
	}

	for _, tc := range tcs {
		res := GetCause(tc.input)
		if res != tc.output {
			t.Errorf("expected %v, got %v", tc.output, res)
			continue
		}
	}
}
