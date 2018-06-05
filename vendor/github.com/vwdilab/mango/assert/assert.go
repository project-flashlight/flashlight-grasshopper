package assert

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode/utf8"
)

var defOut io.Writer

func init() {
	defOut = os.Stdout
}

// SetOutput sets output target for assertion printouts.
func SetOutput(w io.Writer) {
	defOut = w
}

// Len checks if v has the expected lenght expLen
func Len(t testing.TB, expLen int, v interface{}, msgs ...interface{}) {
	val := reflect.ValueOf(v)

	defer func() {
		if r := recover(); r != nil {
			actStr := fmt.Sprintf("<no length>\t%+v", v)
			fail(t, "Length can't be inferred", toStr(expLen), actStr, "", msgs...)
		}
	}()

	l := val.Len()
	if expLen != l {
		actStr := fmt.Sprintf("%v\t(%+v)", l, v)
		fail(t, "Length does not match", toStr(expLen), actStr, "", msgs...)
	}
}

// True fails the test if the condition is false.
func True(t testing.TB, condition bool, msgs ...interface{}) {
	if !condition {
		fail(t, "Expected TRUE but got FALSE", "", "", "", msgs...)
	}
}

// False fails the test if the condition is true.
func False(t testing.TB, condition bool, msgs ...interface{}) {
	if condition {
		fail(t, "Expected FALSE but got TRUE", "", "", "", msgs...)
	}
}

// Nil fails test if v is not nil.
func Nil(t testing.TB, v interface{}, msgs ...interface{}) {
	if !isNil(v) {
		fail(t, "Expected nil but got something", toStr(nil), toStr(v), "", msgs...)
	}
}

// NotNil fails test if v is nil.
func NotNil(t testing.TB, v interface{}, msgs ...interface{}) {
	if isNil(v) {
		fail(t, "Expected something but got nil", "<something>", toStr(v), "", msgs...)
	}
}

// Error fails the test if err is nil.
func Error(t testing.TB, err error, msgs ...interface{}) {
	if err == nil {
		fail(t, "Expected error didn't occur", toStr("<an error>"), toStr(err), "", msgs...)
	}
}

// NoError fails the test if err is not nil.
func NoError(t testing.TB, err error, msgs ...interface{}) {
	if err != nil {
		fail(t, "Unexpected error occured", toStr(nil), toStr(err.Error()), "", msgs...)
	}
}

// Equal fails the test if exp is not equal to act.
func Equal(t testing.TB, exp, act interface{}, msgs ...interface{}) {
	if !reflect.DeepEqual(exp, act) {
		d := diff(exp, act)
		fail(t, "Not equal", toStr(exp), toStr(act), d, msgs...)
	}
}

// NotEqual fails the test if exp is equal to act.
func NotEqual(t testing.TB, exp, act interface{}, msgs ...interface{}) {
	if reflect.DeepEqual(exp, act) {
		d := diff(exp, act)
		fail(t, "Not equal", toStr(exp), toStr(act), d, msgs...)
	}
}

// IsType fails the test if type of exp is not equal to type of act.
func IsType(t testing.TB, exp, act interface{}, msgs ...interface{}) {
	if reflect.TypeOf(exp) != reflect.TypeOf(act) {
		fail(t, "Types are not equal", reflect.TypeOf(exp).String(), reflect.TypeOf(act).String(), "", msgs...)
	}
}

// Contains fails the test if src doesn't contain elem.
func Contains(t testing.TB, src, elem interface{}, msgs ...interface{}) {
	found, ok := contains(src, elem)

	if !ok {
		fail(t, "Contains can not be applied to provided source", "<map,slice or string>", reflect.TypeOf(src).String(), "", msgs...)
		return
	}

	if !found {
		fail(t, "Contains source doesn't contain element", toStr(elem), toStr(src), "", msgs...)
	}

}

// NotContains fails the test if src contains elem.
func NotContains(t testing.TB, src, elem interface{}, msgs ...interface{}) {
	found, ok := contains(src, elem)

	if !ok {
		fail(t, "NotContains can not be applied to provided source", "<map,slice or string>", reflect.TypeOf(src).String(), "", msgs...)
		return
	}

	if found {
		fail(t, "Contains source contains element", toStr(elem), toStr(src), "", msgs...)
	}

}

// Fail fails the test and prints message msg.
func Fail(t testing.TB, msg string, msgs ...interface{}) {
	fail(t, msg, "", "", "", msgs...)
}

func fail(t testing.TB, msg, exp, act, diff string, msgs ...interface{}) {
	file, line := callerInfo()
	fmtr := func(pref string, m string, intends int) string {
		if m == "" {
			return ""
		}
		in := ""
		for j := 0; j < intends; j++ {
			in = fmt.Sprintf("  %s", in)
		}
		return fmt.Sprintf("\n%s%s: %s", in, pref, m)
	}

	str := fmtr("FAIL", msg, 1)
	str += fmtr(" exp", exp, 2)
	str += fmtr(" act", act, 2)
	str += fmtr(" diff", diff, 2)
	str += fmtr("Messages", concatMsg(msgs...), 1)

	fmt.Fprintf(defOut, "%s:%d:\n%s\n\n", file, line, str)
	t.FailNow()
}

func callerInfo() (string, int) {
	_, file, line, _ := runtime.Caller(1)
	return file, line
}

func concatMsg(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	m := ""
	for _, a := range args {
		m = fmt.Sprintf("%s %s", m, a)
	}
	return m
}

func diff(base, comp interface{}) string {
	strBase := fmt.Sprintf("%+v", base)
	strComp := fmt.Sprintf("%+v", comp)

	diff := ""
	start := -1
	maxRune := 30

	for i, r := range strBase {
		if len(strComp)-1 < i {
			break
		}
		s, _ := utf8.DecodeRune([]byte{strComp[i]})
		if r != s && start == -1 {
			start = i
		}
		if start != -1 {
			diff += string(s)
			maxRune--
		}
		if maxRune == 0 {
			break
		}
	}

	if start == -1 {
		return ""
	}

	if start+maxRune < len(strComp) {
		return fmt.Sprintf("[%v] %s...", start, diff)
	}
	return fmt.Sprintf("[%v] %s", start, diff)
}

func toStr(v interface{}) string {
	return fmt.Sprintf("%+v", v)
}

func isNil(v interface{}) (ok bool) {
	if v == nil {
		return true
	}

	// val.isNil might panic if the argument must be
	// a chan, func, interface, map, pointer, or slice value
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	val := reflect.ValueOf(v)

	if val.IsNil() {
		return true
	}
	return false
}

func contains(source interface{}, elem interface{}) (found, ok bool) {

	lVal := reflect.ValueOf(source)
	eVal := reflect.ValueOf(elem)
	defer func() {
		if e := recover(); e != nil {
			found = false
			ok = false
		}
	}()

	// if elem is string
	if reflect.TypeOf(source).Kind() == reflect.String {
		return strings.Contains(lVal.String(), eVal.String()), true
	}

	// if elem is map
	if reflect.TypeOf(source).Kind() == reflect.Map {
		k := lVal.MapKeys()
		for i := 0; i < len(k); i++ {
			if reflect.DeepEqual(k[i].Interface(), elem) {
				return true, true
			}
		}
		return false, true
	}

	// if elem is slice
	for i := 0; i < lVal.Len(); i++ {
		if reflect.DeepEqual(lVal.Index(i).Interface(), elem) {
			return true, true
		}
	}

	return false, true

}
