package gese

import (
	"reflect"
	"strconv"
	"strings"
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
		from, path, defaultVal interface{}
		replaceZeroVal         bool
		argsLength             = len(args)
	)
	if argsLength > 0 {
		from = args[0]
		if argsLength > 1 {
			path = args[1]
			if argsLength > 2 {
				defaultVal = args[2]
				if argsLength > 3 {
					isReplace, ok := args[3].(bool)
					if ok {
						replaceZeroVal = isReplace
					}
				}
			}
		}
	}

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
		mKeyType := reflect.TypeOf(fVal.Interface()).Key().Kind()
		mPathType := reflect.TypeOf(firstPInfo.val).Kind()
		var val reflect.Value
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

// isPosNum Check the number is positive
func isPosNum(numPtr *int) bool {
	return numPtr != nil && *numPtr >= 0
}

// detectPath detect and parse path to string, real value and number if value is string and can parse to number.
// if detectPath invalid = true => list path item detail at least a item
func detectPath(path interface{}) ([]pathItem, bool) {
	pVal := reflect.ValueOf(path)
	if !pVal.IsValid() {
		return nil, false
	}

	var (
		pType    = pVal.Type()
		pKind    = pType.Kind()
		pathInfo []pathItem
	)
	if pKind == reflect.Interface || pKind == reflect.Ptr {
		pVal = pVal.Elem()
		pKind = pVal.Type().Kind()
	}

	switch pKind {
	case reflect.String: // "A.B.C", "A.B.6"
		list := strings.Split(pVal.String(), ".")
		pathInfo = make([]pathItem, len(list))
		for idx := range list {
			if list[idx] == "" {
				return nil, false
			}
			var numPtr *int
			num, err := strconv.Atoi(list[idx])
			if err == nil {
				numPtr = &num
			}
			pathInfo[idx] = pathItem{
				val: list[idx],
				num: numPtr,
				str: list[idx],
			}
		}
	case reflect.Array, reflect.Slice: // ["A", "B", 5]
		var length = pVal.Len()
		if length == 0 {
			return nil, false
		}
		pathInfo = make([]pathItem, length)
		for i := 0; i < pVal.Len(); i++ {
			var (
				iVal   = pVal.Index(i)
				iKind  = iVal.Type().Kind()
				val    interface{}
				numPtr *int
				str    string
			)
			if iKind == reflect.Interface {
				iVal = reflect.ValueOf(iVal.Interface())
				iKind = iVal.Type().Kind()
			}

			switch iKind {
			case reflect.String:
				val = iVal.String()
				str = val.(string)
				num, err := strconv.Atoi(str)
				if err == nil {
					numPtr = &num
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				iInt := int(iVal.Int())
				val = iInt
				numPtr = &iInt
				str = strconv.Itoa(iInt)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				iInt := int(iVal.Uint())
				val = iInt
				numPtr = &iInt
				str = strconv.Itoa(iInt)
			default:
				return nil, false
			}
			pathInfo[i] = pathItem{
				val: val,
				num: numPtr,
				str: str,
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: // 10
		var iInt = int(pVal.Int())
		pathInfo = append(pathInfo, pathItem{
			val: iInt,
			num: &iInt,
			str: strconv.Itoa(iInt),
		})
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: // 11 :) :) :)
		var iInt = int(pVal.Uint())
		pathInfo = append(pathInfo, pathItem{
			val: iInt,
			num: &iInt,
			str: strconv.Itoa(iInt),
		})
	default:
		pathInfo = append(pathInfo, pathItem{
			val: pVal.Interface(),
		})
	}

	return pathInfo, true
}

type pathItem struct {
	str string
	val interface{}
	num *int
}
