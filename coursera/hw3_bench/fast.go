package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson3486653aDecodeGolearnCourseraHw3Bench(in *jlexer.Lexer, out *UserRecord) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}


// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserRecord) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3486653aDecodeGolearnCourseraHw3Bench(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserRecord) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3486653aDecodeGolearnCourseraHw3Bench(l, v)
}


//easyjson:json
type UserRecord struct {
	Browsers []string `json:"browsers"`
	Email string      `json:"email"`
	Name  string      `json:"name"`
}

var linesCache []string

func getFileLines() []string {
	if len(linesCache) == 0 {
		file, _ := os.Open(filePath)
		data, _ := ioutil.ReadAll(file)
		file.Close()

		linesCache = strings.Split(string(data), "\n")
	}
	return linesCache
}


func FastSearch(out io.Writer) {
	seenBrowsers := make(map [string]bool, 100)
	const strAndroid = "Android"
	const strMSIE = "MSIE"

	foundUsers := make([]string, 0)

	for i, line := range getFileLines() {
		// only half of strings contains both browsers
		if !(strings.Contains(line, strAndroid) || strings.Contains(line, strMSIE)) {
			continue
		}

		var user UserRecord
		err := user.UnmarshalJSON([]byte(line))
		if err != nil {
			panic(err)
		}

		var seenAndroid, seenMSIE bool

		for _, browser := range user.Browsers {
			if strings.Contains(browser, strAndroid) {
				seenAndroid = true
				seenBrowsers[browser] = true
			}
			if strings.Contains(browser, strMSIE) {
				seenMSIE = true
				seenBrowsers[browser] = true
			}
		}

		if seenAndroid && seenMSIE {
			email := strings.Replace(user.Email, "@", " [at] ", 1)
			foundUsers = append(foundUsers, fmt.Sprintf("[%d] %s <%s>", i, user.Name, email))
		}
	}

	fmt.Fprintln(out, "found users:")
	for _, u := range foundUsers {
		fmt.Fprintln(out, u)
	}
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}

func main() {
	FastSearch(os.Stdout)
}