package gese

import (
	"reflect"
)

// Get
// args1: type arg object need get value includes below types:
// 1. array, slice like []T{T1, T2}
// 2. struct, map like {A: T1, B: T2}
// 3. string "abc"
// 4. remain type like int, string, func, chan, ...
//
// args2: path arg path to direct value like below:
// ["A", "B"], "A.B", ["A", 2, "B", "C"]
//
// args3: defaultVal arg default value when the value at path not exists
//
// args4: replaceZeroVal arg when the value is zero will return defaultVal
// from, path, defaultVal interface{}, replaceZeroVal bool
func Get(args ...interface{}) interface{} {
	var (
		from, path, defaultVal        interface{}
		replaceZeroVal, isReplace, ok bool
		argsLength                    = len(args)
	)
	if argsLength == 0 {
		return nil
	}

	from = args[0]
	if argsLength <= 1 {
		goto GetLabel
	}

	path = args[1]
	if argsLength <= 2 {
		goto GetLabel
	}

	defaultVal = args[2]
	if argsLength <= 3 {
		goto GetLabel
	}

	isReplace, ok = args[3].(bool)
	if ok {
		replaceZeroVal = isReplace
	}

GetLabel:
	return get(from, path, defaultVal, replaceZeroVal, true)
}

func get(from, path, defaultVal interface{}, replaceZeroVal, isFirst bool) interface{} {
	fVal := reflect.ValueOf(from)
	if !fVal.IsValid() {
		return defaultVal
	}

	var pInfo []pathItem
	if isFirst {
		var valid bool
		pInfo, valid = detectPath(path)
		if !valid {
			return defaultVal
		}
	} else {
		pInfo = path.([]pathItem)
	}

	var (
		pLength    = len(pInfo)
		firstPInfo = pInfo[0]
		idxNum     = firstPInfo.num
		fType      = fVal.Type()
		fKind      = fType.Kind()
	)

	if fKind == reflect.Interface || fKind == reflect.Ptr {
		fVal = fVal.Elem()
		fType = fVal.Type()
		fKind = fType.Kind()
	}

	switch fKind {
	case reflect.String:
		fLength := fVal.Len()
		if !isPosNum(idxNum) || pLength > 1 || *idxNum > fLength {
			return defaultVal
		}
		val := fVal.Index(*idxNum)
		if val.IsZero() && replaceZeroVal {
			return defaultVal
		}
		return rune(val.Interface().(uint8))
	case reflect.Array, reflect.Slice:
		fLength := fVal.Len()
		switch {
		case !isPosNum(idxNum) || fLength-1 < *idxNum:
			return defaultVal
		case pLength == 1:
			val := fVal.Index(*idxNum)
			if val.IsZero() && replaceZeroVal {
				return defaultVal
			}
			return val.Interface()
		default:
			return get(fVal.Index(*idxNum).Interface(), pInfo[1:pLength], defaultVal, replaceZeroVal, false)
		}
	case reflect.Struct:
		if idxNum != nil {
			return defaultVal
		}

		val := fVal.FieldByName(firstPInfo.str)
		if !val.IsValid() {
			return defaultVal
		}

		if val.IsZero() && replaceZeroVal {
			return defaultVal
		}

		if pLength == 1 {
			return val.Interface()
		}
		return get(val.Interface(), pInfo[1:pLength], defaultVal, replaceZeroVal, false)
	case reflect.Map:
		var (
			mKeyType  = reflect.TypeOf(fVal.Interface()).Key().Kind()
			mPathType = reflect.TypeOf(firstPInfo.val).Kind()
			val       reflect.Value
		)

		if mKeyType == mPathType || mKeyType == reflect.Interface {
			val = fVal.MapIndex(reflect.ValueOf(firstPInfo.val))
		}

		if !val.IsValid() {
			if idxNum != nil && mKeyType == reflect.Interface {
				val = fVal.MapIndex(reflect.ValueOf(*idxNum))
				if !val.IsValid() {
					return defaultVal
				}
			} else if mKeyType == reflect.String {
				val = fVal.MapIndex(reflect.ValueOf(firstPInfo.str))
				if !val.IsValid() {
					return defaultVal
				}
			} else {
				return defaultVal
			}
		}

		if val.IsZero() && replaceZeroVal {
			return defaultVal
		}

		if pLength == 1 {
			return val.Interface()
		}

		return get(val.Interface(), pInfo[1:pLength], defaultVal, replaceZeroVal, false)
	}

	return defaultVal
}
