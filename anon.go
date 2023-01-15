package anon

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/mitchellh/copystructure"
)

const tagName = "anon"

const (
	tagStars  = "stars"
	tagEmpty  = "empty"
	tagLen    = "stars_with_len"
	tagInfo   = "with_info"
	tagSHA512 = "sha512"

	emptyErr = "empty"
)

// Marshal anonymises fields that have anon tags, and then runs the given
// function (e.g. json.Marshal()) on the resulting struct.
// Marshal takes a pointer to a struct.
func Marshal(v any, m func(any) ([]byte, error)) ([]byte, error) {
	res, err := Anonymise(v)
	if err != nil {
		return nil, err
	}

	return m(res)
}

// Anonymise anonymises data according to the "anon" tags it has in it.
// It returns a copy of the given value. Use AnonymiseByRef if you'd
// prefer to pass by reference.
func Anonymise(v any) (res any, err error) {
	defer func() {
		// I don't trust reflection.
		if r := recover(); r != nil {
			err = fmt.Errorf("anon: %v", r)
		}
	}()

	// Deep copy, otherwise we risk changing maps/slices inside the struct.
	cp, err := copystructure.Copy(v)
	if err != nil {
		return nil, err
	}

	val := reflect.ValueOf(cp)
	tmp := reflect.New(val.Type())
	tmp.Elem().Set(val)

	err = anonymise("", tmp)
	return tmp.Interface(), err
}

// AnonymiseByRef anonymises the struct that the given pointer points
// to, according to the "anon" tags it has in it.
func AnonymiseByRef(v any) (err error) {
	defer func() error {
		// I don't trust reflection.
		if r := recover(); r != nil {
			err = fmt.Errorf("anon: %v", r)
		}
		return nil
	}()

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("value must be pointer")
	}

	return anonymise("", val)
}

func anonymise(tag string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return anonymise(tag, v.Elem())
	case reflect.Struct:
		for i := 0; i < v.Type().NumField(); i++ {
			tag := v.Type().Field(i).Tag.Get(tagName)

			err := anonymise(tag, v.Field(i))
			if err != nil {
				return err
			}
		}
	case reflect.String:
		err := obfuscate(tag, v)
		if err != nil {
			return err
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			err := anonymise(tag, v.Index(i))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func obfuscate(tag string, v reflect.Value) error {
	if tag == "-" || tag == "" {
		return nil
	}

	var a anonymiser

	switch tag {
	case tagStars:
		a = Stars
	case tagEmpty:
		a = Empty
	case tagLen:
		a = StarsWithLen
	case tagInfo:
		a = WithInfo
	case tagSHA512:
		a = SHA512
	default:
		return fmt.Errorf("no tag for %s", tag)
	}

	v.SetString(a(v.String()))
	return nil
}

type anonymiser func(string) string

func Stars(string) string {
	return "****"
}

// Returns a number of asterisks equal to string length.
func StarsWithLen(s string) string {
	return strings.Repeat("*", len(s))
}

func Empty(string) string {
	return ""
}

// Returns length and existence of non-ASCII characters in string.
// Do keep in mind that a non-ASCII character will have a length greater than 1.
//
// len("á") == 2
//
// len("a") == 1
func WithInfo(s string) string {
	return fmt.Sprintf("len:%d,is_ascii:%t", len(s), isASCII(s))
}

func SHA512(s string) string {
	sum := sha512.Sum512([]byte(s))
	return string(sum[:])
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		// if a byte is greater than RuneSelf, then the following byte is still
		// part of the same rune character, meaning it's not ASCII.
		if s[i] > utf8.RuneSelf {
			return false
		}
	}
	return true
}
