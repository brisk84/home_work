package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInterfaceIsNotStruct = errors.New("interface is not a struct")
	ErrLenNotEqual          = errors.New("len not equal")
	ErrLessMin              = errors.New("less than min")
	ErrMoreMax              = errors.New("more than max")
	ErrNotMatchRegex        = errors.New("not match regex")
	ErrNotInRange           = errors.New("not it range")
	ErrArgFormatError       = errors.New("argument format error")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	ret := ""
	for _, e := range v {
		ret += e.Field + " - " + e.Err.Error() + "\n"
	}
	return ret
}

func CheckLen(field interface{}, param interface{}) error {
	strLen, err := strconv.Atoi(param.(string))
	if err != nil {
		return ErrArgFormatError
	}

	t := reflect.TypeOf(field)
	if t.String() == "string" {
		s := field.(string)
		if len(s) == strLen {
			return nil
		}
		return ErrLenNotEqual
	} else if t.String() == "[]string" {
		ss := field.([]string)
		for _, s := range ss {
			if len(s) != strLen {
				return ErrLenNotEqual
			}
		}
		return nil
	}
	return fmt.Errorf("not string or []string")
}

func CheckMin(field interface{}, param interface{}) error {
	val := field.(int)
	min, err := strconv.Atoi(param.(string))
	if err != nil {
		return err
	}
	if val < min {
		return ErrLessMin
	}
	return nil
}

func CheckMax(field interface{}, param interface{}) error {
	val := field.(int)
	max, err := strconv.Atoi(param.(string))
	if err != nil {
		return err
	}
	if val > max {
		return ErrMoreMax
	}
	return nil
}

func CheckRegexp(field interface{}, param interface{}) error {
	checkVal := field.(string)
	pattern := param.(string)

	matched, err := regexp.MatchString(pattern, checkVal)
	if err != nil {
		return err
	}
	if !matched {
		return ErrNotMatchRegex
	}
	return nil
}

func CheckIn(field interface{}, param interface{}) error {
	val := param.(string)
	vals := strings.Split(val, ",")

	if reflect.TypeOf(field).String() == "string" {
		checkVal := field.(string)
		for _, v := range vals {
			if v == checkVal {
				return nil
			}
		}
	}

	if reflect.TypeOf(field).String() == "[]string" {
		checkVals := field.([]string)
		for _, checkVal := range checkVals {
			for _, v := range vals {
				if v == checkVal {
					return nil
				}
			}
		}
	}

	if reflect.TypeOf(field).String() == "int" {
		checkVal := field.(int)
		for _, v := range vals {
			intV, err := strconv.Atoi(v)
			if err != nil {
				return ErrArgFormatError
			}
			if intV == checkVal {
				return nil
			}
		}
	}

	if reflect.TypeOf(field).String() == "[]int" {
		checkVals := field.([]int)
		for _, checkVal := range checkVals {
			for _, v := range vals {
				intV, err := strconv.Atoi(v)
				if err != nil {
					return ErrArgFormatError
				}
				if intV == checkVal {
					return nil
				}
			}
		}
	}

	return ErrNotInRange
}

func Validate(v interface{}) error {

	var errs ValidationErrors

	ff := map[string]func(field interface{}, param interface{}) error{
		"len":    CheckLen,
		"min":    CheckMin,
		"max":    CheckMax,
		"regexp": CheckRegexp,
		"in":     CheckIn,
	}

	vv := reflect.ValueOf(v)
	tt := reflect.TypeOf(v)

	if vv.Kind() != reflect.Struct {
		return fmt.Errorf("%w. struct != %T", ErrInterfaceIsNotStruct, v)
	}

	for i := 0; i < vv.NumField(); i++ {
		t := tt.Field(i)
		tag := t.Tag.Get("validate")
		if tag == "" {
			continue
		}

		vals := strings.Split(tag, "|")
		for _, curVal := range vals {
			funcArg := strings.Split(curVal, ":")
			if f, ok := ff[funcArg[0]]; ok {
				value := vv.Field(i).Interface()
				fieldType := vv.Field(i).Type().String()
				if (fieldType != "string") && (fieldType != "[]string") && (fieldType != "int") && (fieldType != "[]int") {
					value = vv.Field(i).String()
				}

				err := f(value, funcArg[1])
				if err != nil {
					e := ValidationError{
						Field: tt.Field(i).Name,
						Err:   err,
					}
					errs = append(errs, e)
				}
			}

		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}
