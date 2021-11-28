package gese

import (
	"errors"
	"reflect"
)

var (
	ErrPathNotExists = errors.New("path not exists")
	ErrInvalidDest   = errors.New("destination must be pointer, map or slice")
	ErrBadValue      = errors.New("cannot set the value into dest with the path by two mismatched types")
)

func Set(dest, path, val interface{}) error {
	return set(dest, path, val)
}

func set(dest, path, val interface{}) error {
	dVal := reflect.ValueOf(dest)
	if !dVal.IsValid() || dVal.Type().Kind() != reflect.Ptr {
		return ErrInvalidDest
	}

	pInfo, ok := path.([]pathItem)
	if !ok {
		var valid bool
		pInfo, valid = detectPath(path)
		if !valid {
			return ErrPathNotExists
		}
	}

	var (
		pLength    = len(pInfo)
		firstPInfo = pInfo[0]
		idxNum     = firstPInfo.num
	)

	ptrVal := reflect.Indirect(dVal)
	switch ptrVal.Type().Kind() {
	case reflect.String:
		str := ptrVal.String()
		if idxNum == nil || len(str) < *idxNum || pLength > 1 {
			return ErrPathNotExists
		}

		if val == nil || reflect.TypeOf(val).Kind() != reflect.Int32 {
			goto BadValueLabel
		}

		r := []rune(str)
		r[*idxNum] = rune(reflect.ValueOf(val).Int())

		reflect.Indirect(dVal).Set(reflect.ValueOf(string(r)))
		return nil
	case reflect.Struct:
		if idxNum != nil {
			return ErrPathNotExists
		}

		fVal := ptrVal.FieldByName(firstPInfo.str)
		if !fVal.IsValid() {
			return ErrPathNotExists
		}

		if pLength != 1 {
			return set(fVal.Addr().Interface(), pInfo[1:pLength], val)
		}

		var sVal = reflect.ValueOf(val)
		if sVal.Type() != fVal.Type() {
			return ErrBadValue
		}

		dVal.Elem().FieldByName(firstPInfo.str).Set(sVal)

		return nil
	case reflect.Map:

	case reflect.Ptr:

		if ptrVal.IsNil() {
			ptrVal.Set(reflect.New(ptrVal.Type().Elem()))
		}
		return set(ptrVal.Interface(), pInfo, val)
	case reflect.Interface:
	// case reflect.Bool,
	// 	reflect.Complex128, reflect.Complex64,
	// 	reflect.Float32, reflect.Float64,
	// 	reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
	// 	reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
	// 	if reflect.ValueOf(val).Type() != dVal.Type() {
	// 		goto BadValueLabel
	// 	}

	default:
		goto InvalidLabel

	}
	// return dest, Set(dest, pInfo, setVal)

BadValueLabel:
	return ErrBadValue
InvalidLabel:
	return ErrInvalidDest
}
