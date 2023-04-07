package utils

import (
	"bytes"
	"log"
	"reflect"
	"strings"
)

func ObjectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

// containsElement try loop over the list check if the list includes the element.
// return (false, false) if impossible.
// return (true, false) if element was not found.
// return (true, true) if element was found.
func includeElement(list interface{}, element interface{}) (ok, found bool) {

	listValue := reflect.ValueOf(list)
	listKind := reflect.TypeOf(list).Kind()
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if listKind == reflect.String {
		elementValue := reflect.ValueOf(element)
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if listKind == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if ObjectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if ObjectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}

	return true, false
}

func Contains(s, contains interface{}) bool {
	ok, found := includeElement(s, contains)
	if !ok {
		log.Printf("%#v could not be applied builtin len()", s)
		return false
	}
	if !found {
		log.Printf("%#v does not contain %#v", s, contains)
		return false
	}

	return true
}

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    assert.NotContains(t, "Hello World", "Earth")
//    assert.NotContains(t, ["Hello", "World"], "Earth")
//    assert.NotContains(t, {"Hello": "World"}, "Earth")
func NotContains(s, contains interface{}) bool {
	ok, found := includeElement(s, contains)
	if !ok {
		log.Printf("%#v could not be applied builtin len()", s)
		return false
	}
	if found {
		log.Printf("\"%s\" should not contain \"%s\"", s, contains)
		return false
	}

	return true
}

func OneOf[T comparable](e T, l []T) bool {
	for _, i := range l {
		if i == e {
			return true
		}
	}

	return false
}
