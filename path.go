package gese

import (
	"reflect"
	"strconv"
	"strings"
)

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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var iInt = int(pVal.Int())
		pathInfo = append(pathInfo, pathItem{
			val: iInt,
			num: &iInt,
			str: strconv.Itoa(iInt),
		})
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
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
