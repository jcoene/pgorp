package pgorp

import (
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"regexp"
	"strconv"
	"strings"
)

// Database type for struct definitions
type ArrayInt64 []int64

// Converter type to assign to gorp
type TypeConverter struct{}

var (
	// unquoted array values must not contain: (" , \ { } whitespace NULL)
	// and must be at least one char
	unquotedChar  = `[^",\\{}\s(NULL)]`
	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

	// quoted array values are surrounded by double quotes, can be any
	// character except " or \, which must be backslash escaped:
	quotedChar  = `[^"\\]|\\"|\\\\`
	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

	// an array value may be either quoted or unquoted:
	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

	// Array values are separated with a comma IF there is more than one value:
	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

	valueIndex int
)

// Find the index of the 'value' named expression
func init() {
	for i, subexp := range arrayExp.SubexpNames() {
		if subexp == "value" {
			valueIndex = i
			break
		}
	}
}

// Called by gorp to serialize custom types to database values
func (c TypeConverter) ToDb(val interface{}) (interface{}, error) {
	switch t := val.(type) {
	case ArrayInt64:
		return pgBuildIntArray(t), nil
	}

	return val, nil
}

// Called by gorp to deserialize database values to custom types
func (c TypeConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	switch target.(type) {
	case ArrayInt64:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(string)
			if !ok {
				return errors.New("unable to convert array int64 to string")
			}
			target = pgParseIntArray(s)
			return nil
		}
		return gorp.CustomScanner{new(string), target, binder}, true
	}

	return gorp.CustomScanner{}, false
}

// Builds a postgres array from an ArrayInt64
func pgBuildIntArray(value ArrayInt64) (result string) {
	strvals := make([]string, len(value))
	for i, v := range value {
		strvals[i] = fmt.Sprintf("%d", v)
	}

	return fmt.Sprintf("{%s}", strings.Join(strvals, ","))
}

// Builds an ArrayInt64 from a Postgres array string
func pgParseIntArray(value string) (result ArrayInt64) {
	result = make(ArrayInt64, 0)
	matches := arrayExp.FindAllStringSubmatch(value, -1)
	for _, match := range matches {
		s := match[valueIndex]
		s = strings.Trim(s, "\"")
		n, _ := strconv.ParseInt(s, 10, 64)
		result = append(result, n)
	}

	return
}
