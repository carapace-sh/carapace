package jsonc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Jsonc is the structure for parsing json with comments
type Jsonc struct {
	index    int    // current index position in source data
	comment  uint   // the type of comment: eithher 1 for // OR 2 for /*
	len      int    // the length of source data
	objDepth uint   // the depth of nested objects
	arrDepth uint   // the depth of nested arrays
	inStr    bool   // if inside string
	inArr    bool   // if inside array notation: [
	inObj    bool   // if inside object notation: {
	last     string // last significant non whitespace char
	strDelim string // string delimeter: either ' or "
}

// New creates Jsonc struct
func New() *Jsonc {
	return &Jsonc{}
}

// Strip strips comments and trailing commas from input byte array
func (j *Jsonc) Strip(jsonb []byte) []byte {
	s := j.StripS(string(jsonb))
	return []byte(s)
}

var sq = `'`  // single quote
var dq = `"`  // double quote
var esc = `\` // escape
var comma = regexp.MustCompile(`(?:,+)(\s*)$`)

// StripS strips comments and trailing commas from input string
func (j *Jsonc) StripS(data string) string {
	var oldprev, prev, char, next, s string

	j.reset()
	j.len = len(data)
	quote, quoted := "", false

	for j.index < j.len {
		oldprev, prev, char, next = j.getSegments(data, prev)

		// If value starts with 0x, parse as hexadecimal
		if j.isNonStringValue(char, "0") && (next == "x" || next == "X") {
			s += j.hexadecimal(data)
			continue
		}

		quote, quoted = j.quoteKey(char, quoted)
		s += quote

		// Trim trailing commas at the end of array or object
		if j.comment == 0 && !j.inStr && ((j.inArr && char == "]") || (j.inObj && char == "}")) {
			s = comma.ReplaceAllString(s, `$1`)
		}

		j.checkArrayObject(char)

		// Append char as is (or it's compliment pair) if inside string or outside comment
		if j.inString(prev, char, next, oldprev) || j.outsideComment(char, next) {
			s += j.compliment(prev, char, next)
			continue
		}

		// Wipe out trailing whitespaces around comment
		if j.hasCommentEnded(char, next) && char == "\n" {
			s = strings.TrimRight(s, "\r\n\t ") + char
		}
	}
	return s
}

// Unmarshal strips and parses the json byte array
func (j *Jsonc) Unmarshal(jsonb []byte, v interface{}) error {
	return json.Unmarshal(j.Strip(jsonb), v)
}

// UnmarshalFile strips and parses the json content from file
func (j *Jsonc) UnmarshalFile(file string, v interface{}) error {
	jsonb, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return j.Unmarshal(jsonb, v)
}

// reset resets the Jsonc with proper defaults
func (j *Jsonc) reset() {
	j.index, j.comment = 0, 0
	j.objDepth, j.arrDepth = 0, 0
	j.inStr, j.inArr, j.inObj = false, false, false
	j.last, j.strDelim = "", ""
}

// getSegments gets look-behind, current and look-ahead chars
func (j *Jsonc) getSegments(json, old string) (oldprev, prev, char, next string) {
	oldprev = old
	if j.index > 0 {
		prev = json[j.index-1 : j.index]
	}
	char = json[j.index : j.index+1]
	if j.index < j.len-1 {
		next = json[j.index+1 : j.index+2]
	}
	j.index++
	return
}

// isNonStringValue checks if char is value outside string (or comment) and matches chars
func (j *Jsonc) isNonStringValue(char, chars string) bool {
	return !j.inStr && j.comment == 0 && strings.ContainsAny(char, chars)
}

// hexadecimal consumes hex (0-9a-fA-F) chars and converts to decimal string
func (j *Jsonc) hexadecimal(data string) string {
	j.index++
	hexa := ""
	for j.index < j.len {
		char := data[j.index : j.index+1]
		if !isNumber(char, true) {
			break
		}
		hexa += char
		j.index++
	}
	dec, _ := strconv.ParseInt(hexa, 16, 32)
	return fmt.Sprintf("%d", dec)
}

// quoteKey double quotes the unquoted object keys
func (j *Jsonc) quoteKey(char string, wasQuoted bool) (q string, quoted bool) {
	quoted = wasQuoted
	inKey := j.inObj && j.comment == 0 && (j.last == "{" || j.last == ",")
	// Object key has just started without quote, so quote it
	if !j.inStr && inKey && !strings.ContainsAny(char, "[]{}'\",/*:\r\n\t ") {
		q = dq
		j.inStr, quoted, j.strDelim = true, true, dq
	}
	// Object key has just ended and was quoted before, so quote it again to compliment
	if j.inStr && wasQuoted && inKey && (char == ":" || char == " " || char == sq) {
		q = dq
		j.inStr, quoted, j.strDelim = false, false, ""
	}
	return
}

// checkArrayObject checks and sets the depth and state of array &/or object notation
func (j *Jsonc) checkArrayObject(char string) {
	if j.isNonStringValue(char, char) {
		// Last non whitespace char
		if !strings.ContainsAny(char, "\r\n\t /") {
			j.last = char
		}
		if char == "{" {
			j.objDepth++
			j.inObj, j.inArr = true, false
		} else if j.objDepth > 0 && char == "}" {
			j.objDepth--
			j.inObj, j.inArr = j.objDepth > 0, j.arrDepth > 0
		} else if char == "[" {
			j.arrDepth++
			j.inObj, j.inArr = false, true
		} else if j.arrDepth > 0 && char == "]" {
			j.arrDepth--
			j.inObj, j.inArr = j.objDepth > 0, j.arrDepth > 0
		}
	}
}

func (j *Jsonc) inString(prev, char, next, oldprev string) bool {
	charnext := char + next
	maybeStr := (char == dq || char == sq) && (!j.inStr || j.strDelim == char)

	// Toggle j.inStr if j.strDelim is not escaped
	if j.comment == 0 && maybeStr && prev != esc {
		if !j.inStr {
			j.strDelim = char
		}
		j.inStr = !j.inStr
		return j.inStr
	}
	if j.inStr && (charnext == `":` || charnext == `",` || charnext == `"]` || charnext == `"}`) {
		j.inStr = oldprev+prev != esc+esc
	}
	return j.inStr
}

// outsideComment checks if char is outside comment
// it also sets the state of comment
func (j *Jsonc) outsideComment(char, next string) bool {
	// Set comment state: `//` => 1 | `/*` => 2
	if !j.inStr && j.comment == 0 {
		if char+next == "//" {
			j.comment = 1
		}
		if char+next == "/*" {
			j.comment = 2
		}
	}
	return j.comment == 0
}

// hasCommentEnded checks if the comment has just ended and resets the state
func (j *Jsonc) hasCommentEnded(char, next string) bool {
	// Single line comment ends with `\n` and multiline ends with `*/`
	singleEnded := j.comment == 1 && char == "\n"
	multiEnded := j.comment == 2 && char+next == "*/"
	if singleEnded || multiEnded {
		j.comment = 0
	}
	if multiEnded {
		j.index++
	}
	return j.comment == 0
}

var spacesPair = map[string]string{"\n": `\n`, "\t": `\t`, "\r": `\r`}

// compliment appends char as is (or it's compliment pair)
// (eg: in string boundary the compliment of single quote is double quote)
// it also normalizes whitespaces inside string and signed &/or decimal numbers
func (j *Jsonc) compliment(prev, char, next string) string {
	if j.inStr && char == esc && next == "\n" {
		j.index++
		return ""
	} else if c, ok := spacesPair[char]; ok && j.inStr {
		return c
	}

	// Signed +ve number
	if j.isNonStringValue(char, "+") && isNumber(next, false) {
		return ""
	}

	// Decimal point number
	if j.isNonStringValue(char, ".") {
		prevNum, nextNum := isNumber(prev, false), isNumber(next, false)
		if !prevNum && nextNum {
			char = "0."
		} else if prevNum && !nextNum {
			char = ".0"
		}
		return char
	}

	// Single quoted string
	if j.strDelim == sq {
		if char+next == esc+sq {
			char = sq
			j.index++
		} else if prev != esc && char == sq {
			char = dq
		} else if char == dq {
			char = `\"`
		}
	}
	return char
}

// isNumber checks if a string char is numeric
func isNumber(char string, hex bool) bool {
	if hex {
		return strings.ContainsAny(char, "0123456789abcdefABCDEF")
	}
	return strings.ContainsAny(char, "0123456789")
}
