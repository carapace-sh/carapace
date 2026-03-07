package jsonc

import (
	"os"
	"strings"
	"testing"
)

type Case struct {
	json   string
	expect string
}

func TestStrip(t *testing.T) {
	j := New()
	for name, test := range testcases() {
		t.Run("Strip "+name, func(t *testing.T) {
			actual := j.StripS(test.json)
			if actual != test.expect {
				t.Errorf("[%s] expected %s, got %s",
					name,
					strings.ReplaceAll(test.expect, "\t", "."),
					strings.ReplaceAll(actual, "\t", "."),
				)
			} else if actual != string(j.Strip([]byte(test.json))) {
				t.Error("byte str should match")
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	j := New()
	for name, test := range testcases() {
		t.Run("Unmarshal "+name, func(t *testing.T) {
			var ref map[string]interface{}
			if err := j.Unmarshal([]byte(test.json), &ref); err != nil {
				t.Errorf("[%s] unmarshal should not error, got %#v", name, err)
			}
			if name == "nested subjson" {
				jo := ref["jo"].(string)
				if err := j.Unmarshal([]byte(jo), &ref); err != nil {
					t.Errorf("[%v] unmarshal should not error, got %#v", jo, err)
				}
			}
		})
	}

	t.Run("UnmarshalFile", func(t *testing.T) {
		files := []string{"chromium.json5", "test1.json5", "test.json5"}
		for _, file := range files {
			t.Run(file, func(t *testing.T) {
				var ref map[string]interface{}
				if err := j.UnmarshalFile("./examples/"+file, &ref); err != nil {
					t.Errorf("UnmarshalFile should not error, got %#v", err)
				}
				if file != "test.json5" {
					return
				}
				s := ref["g"]
				if err := j.Unmarshal([]byte(s.(string)), &ref); err != nil {
					t.Errorf("[%v] unmarshal should not error, got %#v", s, err)
				}
				if err := j.UnmarshalFile("./examples/invalid.json5", &ref); err == nil {
					t.Error("invalid file should error, got none")
				}
			})
		}
	})
}

func TestCachedDecoder(t *testing.T) {
	file := "./examples/test1.json5"
	cd, val := NewCachedDecoder(), make(map[string]interface{})
	t.Run("before cache", func(t *testing.T) {
		os.Remove("./examples/test1.cached.json")
		if err := cd.Decode(file, &val); err != nil {
			t.Errorf("[%v] decode should not error, got %#v", file, err)
		}
		t.Run("after cache", func(t *testing.T) {
			if err := cd.Decode(file, &val); err != nil {
				t.Errorf("[%v] decode should not error, got %#v", file, err)
			}
		})
	})
}

func testcases() map[string]Case {
	return map[string]Case{
		"without comment": {
			json:   `{"a":1,"b":2}`,
			expect: `{"a":1,"b":2}`,
		},
		"with trail only": {
			json:   `{"a":1,"b":2,,}`,
			expect: `{"a":1,"b":2}`,
		},
		"single line comment": {
			json: `{"a":1,
			// comment
				"b":2,
			// comment
				"c":3,,}`,
			expect: `{"a":1,
				"b":2,
				"c":3}`,
		},
		"single line comment at end": {
			json: `{"a":1,
				"b":2,// comment
				"c":[1,2,,]}`,
			expect: `{"a":1,
				"b":2,
				"c":[1,2]}`,
		},
		"real multiline comment": {
			json: `{"a":1,
			/*
			 * comment
			 */
			"b":2, "c":3,}`,
			expect: `{"a":1,
			` + `
			"b":2, "c":3}`,
		},
		"inline multiline comment": {
			json: `{"a":1,
				/* comment */"b":2, "c":3}`,
			expect: `{"a":1,
				"b":2, "c":3}`,
		},
		"inline multiline comment at end": {
			json:   `{"a":1, "b":2, "c":3/* comment */,}`,
			expect: `{"a":1, "b":2, "c":3}`,
		},
		"comment inside string": {
			json:   `{"a": "a//b", "b":"a/* not really comment */b"}`,
			expect: `{"a": "a//b", "b":"a/* not really comment */b"}`,
		},
		"escaped string": {
			json:   `{"a": "a//b", "b":"a/* \"not really comment\" */b"}`,
			expect: `{"a": "a//b", "b":"a/* \"not really comment\" */b"}`,
		},
		"string inside comment": {
			json:   `{"a": "ab", /* also comment */ "b":"a/* not a comment */b" /* "comment string" */ }`,
			expect: `{"a": "ab",  "b":"a/* not a comment */b"  }`,
		},
		"literal lf": {
			json:   `{"a":/*literal linefeed*/"apple` + "\n" + `ball","b":"","c\\\\":"",}`,
			expect: `{"a":"apple\nball","b":"","c\\\\":""}`,
		},
		"nested subjson": {
			json: `{
				"jo": "{/* comment */\"url\": \"http://example.com\"//comment
				}",
				"x": {
				/* comment 1
					comment 2 */
					"y": {
						// comment
						"XY\\": "//no comment/*",
					},
				}
			}`,
			expect: `{
				"jo": "{/* comment */\"url\": \"http://example.com\"//comment\n\t\t\t\t}",
				"x": {
				` + `
					"y": {
						"XY\\": "//no comment/*"
					}
				}
			}`,
		},
		"with gap": {
			json: `{/*
				?"\" */
				" a " : 1 ,
				" // " :  " : //" // :, \" ",,
			}`,
			expect: `{
				" a " : 1 ,
				" // " :  " : //"
			}`,
		},
	}
}
