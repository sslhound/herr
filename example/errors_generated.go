// Code generated by go generate; DO NOT EDIT.
// This file was generated by herr at 2020-01-14 13:07:59.9248125 -0500 EST m=+0.006999601
package errors

import (
    "fmt"
)

type CodedError interface {
    Code() int
    Description() string
    Prefix() string
    error
}


type DebugErrorOneError struct {
    Err error
}

type DebugErrorTwoError struct {
    Err error
}

type InvalidAndroidVersionError struct {
    Err error
}

type InvalidAndroidDeviceError struct {
    Err error
}

type InvalidIOSDeviceError struct {
    Err error
}

var _ CodedError = DebugErrorOneError{}
var _ CodedError = DebugErrorTwoError{}
var _ CodedError = InvalidAndroidVersionError{}
var _ CodedError = InvalidAndroidDeviceError{}
var _ CodedError = InvalidIOSDeviceError{}

// ErrorFromCode returns the CodedError for a serialized coded error string. 
func ErrorFromCode(code string) (bool, error) {
    switch code {
    case "DBGAAAAAAAB":
        return true, DebugErrorOneError{}
    case "DBGAAAAAAAC":
        return true, DebugErrorTwoError{}
    case "MOBANDAAAAAAAE":
        return true, InvalidAndroidVersionError{}
    case "MOBANDAAAAAAAF":
        return true, InvalidAndroidDeviceError{}
    case "MOBIOSAAAAAAAD":
        return true, InvalidIOSDeviceError{}
    default:
        return false, fmt.Errorf("unknown error code: %s", code)
    }
}

func (e DebugErrorOneError) Error() string {
    return "DBGAAAAAAAB"
}

func (e DebugErrorOneError) Unwrap() error {
	return e.Err
}

func (e DebugErrorOneError) Is(target error) bool {
    t, ok := target.(DebugErrorOneError)
    if !ok {
        return false
    }
    return t.Prefix() == "DBG" && t.Code() == 1
}

func (e DebugErrorOneError) Code() int {
    return 1
}

func (e DebugErrorOneError) Description() string {
    return "The first debug error"
}

func (e DebugErrorOneError) Prefix() string {
    return "DBG"
}

func (e DebugErrorOneError) String() string {
    return "DBGAAAAAAAB The first debug error"
}

func (e DebugErrorTwoError) Error() string {
    return "DBGAAAAAAAC"
}

func (e DebugErrorTwoError) Unwrap() error {
	return e.Err
}

func (e DebugErrorTwoError) Is(target error) bool {
    t, ok := target.(DebugErrorTwoError)
    if !ok {
        return false
    }
    return t.Prefix() == "DBG" && t.Code() == 2
}

func (e DebugErrorTwoError) Code() int {
    return 2
}

func (e DebugErrorTwoError) Description() string {
    return "The second debug error"
}

func (e DebugErrorTwoError) Prefix() string {
    return "DBG"
}

func (e DebugErrorTwoError) String() string {
    return "DBGAAAAAAAC The second debug error"
}

func (e InvalidAndroidVersionError) Error() string {
    return "MOBANDAAAAAAAE"
}

func (e InvalidAndroidVersionError) Unwrap() error {
	return e.Err
}

func (e InvalidAndroidVersionError) Is(target error) bool {
    t, ok := target.(InvalidAndroidVersionError)
    if !ok {
        return false
    }
    return t.Prefix() == "MOBAND" && t.Code() == 4
}

func (e InvalidAndroidVersionError) Code() int {
    return 4
}

func (e InvalidAndroidVersionError) Description() string {
    return "The Android version is invalid."
}

func (e InvalidAndroidVersionError) Prefix() string {
    return "MOBAND"
}

func (e InvalidAndroidVersionError) String() string {
    return "MOBANDAAAAAAAE The Android version is invalid."
}

func (e InvalidAndroidDeviceError) Error() string {
    return "MOBANDAAAAAAAF"
}

func (e InvalidAndroidDeviceError) Unwrap() error {
	return e.Err
}

func (e InvalidAndroidDeviceError) Is(target error) bool {
    t, ok := target.(InvalidAndroidDeviceError)
    if !ok {
        return false
    }
    return t.Prefix() == "MOBAND" && t.Code() == 5
}

func (e InvalidAndroidDeviceError) Code() int {
    return 5
}

func (e InvalidAndroidDeviceError) Description() string {
    return "The Android device is invalid."
}

func (e InvalidAndroidDeviceError) Prefix() string {
    return "MOBAND"
}

func (e InvalidAndroidDeviceError) String() string {
    return "MOBANDAAAAAAAF The Android device is invalid."
}

func (e InvalidIOSDeviceError) Error() string {
    return "MOBIOSAAAAAAAD"
}

func (e InvalidIOSDeviceError) Unwrap() error {
	return e.Err
}

func (e InvalidIOSDeviceError) Is(target error) bool {
    t, ok := target.(InvalidIOSDeviceError)
    if !ok {
        return false
    }
    return t.Prefix() == "MOBIOS" && t.Code() == 3
}

func (e InvalidIOSDeviceError) Code() int {
    return 3
}

func (e InvalidIOSDeviceError) Description() string {
    return "The iOS device is invalid."
}

func (e InvalidIOSDeviceError) Prefix() string {
    return "MOBIOS"
}

func (e InvalidIOSDeviceError) String() string {
    return "MOBIOSAAAAAAAD The iOS device is invalid."
}


