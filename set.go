package gese

import (
	"errors"
	"reflect"
)

var (
	ErrPathNotExists       = errors.New("path not exists")
	ErrInvalidDest         = errors.New("destination must be pointer, map or slice")
	ErrBadValue            = errors.New("cannot set the value into dest with the path by two mismatched types")
	ErrMapKeyTypeUnsupport = errors.New("map key type only support: string, int")
)

func Set(dest, path, val interface{}) error {
	return set(dest, path, val)
}

func set(dest, path, val interface{}) error {
	dVal := reflect.ValueOf(dest)
	if !dVal.IsValid() ||
		(dVal.Type().Kind() != reflect.Ptr &&
			dVal.Type().Kind() != reflect.Map) {
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
		mVal := ptrVal.MapIndex(reflect.ValueOf(firstPInfo.val))
		if !mVal.IsValid() {
			goto InvalidDestLabel
		}

		if pLength != 1 {
			return set(mVal.Interface(), pInfo[1:pLength], val)
		}

		switch ptrVal.Type().Key().Kind() {
		case reflect.String:
			ptrVal.SetMapIndex(reflect.ValueOf(firstPInfo.val), reflect.ValueOf(val))
		case reflect.Int:
			if firstPInfo.num == nil {
				goto InvalidDestLabel
			}
			ptrVal.SetMapIndex(reflect.ValueOf(*firstPInfo.num), reflect.ValueOf(val))
		default:
			return ErrMapKeyTypeUnsupport
		}

		return nil
	case reflect.Ptr:
		if ptrVal.IsNil() {
			ptrVal.Set(reflect.New(ptrVal.Type().Elem()))
		}

		return set(ptrVal.Interface(), pInfo, val)
	case reflect.Interface:
		// TODO
	default:
		goto InvalidDestLabel
	}

BadValueLabel:
	return ErrBadValue
InvalidDestLabel:
	return ErrInvalidDest
}
