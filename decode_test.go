package ini

import (
	"reflect"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		conf string
		want Config
		// if non-empty, a substring that is expected to occur within the error's Error() string.
		wantErr string
	}{
		{
			conf: "",
		},
		{
			conf: "a=b",
			want: map[string]map[string]string{
				Default: map[string]string{
					"a": "b",
				},
			},
		},
		{
			// Both a=b and c:d are allowed - see http://en.wikipedia.org/wiki/INI_file section "Name/value delimiter"
			conf: "a:b",
			want: map[string]map[string]string{
				Default: map[string]string{
					"a": "b",
				},
			},
		},
		// Comments at start of lines are ignored
		{
			conf: `#this is a comment a=b
; and so is this`,
		},
		// TODO What should we do with comments in the middle of a line? See http://en.wikipedia.org/wiki/INI_file#Comments
		{
			conf: "a=b#this is a comment",
			want: map[string]map[string]string{
				Default: map[string]string{
					"a": "b",
				},
			},
			// Or maybe this should error.
		},

		// ESCAPE SEQUENCES
		{
			conf: `a=b\#this is a comment\;`,
			want: map[string]map[string]string{
				Default: map[string]string{
					"a": "b#this is a comment;",
				},
			},
		},

		// Blank lines ignored
		{
			conf: `foo=bar
			
bar=foo`,
			want: map[string]map[string]string{
				Default: map[string]string{
					"foo": "bar",
					"bar": "foo",
				},
			},
		},
		// Duplicate k/v pairs - last one is used. TODO is this right? See http://en.wikipedia.org/wiki/INI_file#Duplicate_names
		{
			conf: `foo=bar
foo=baz`,
			want: map[string]map[string]string{
				Default: map[string]string{
					"foo": "baz",
				},
			},
		},

		// WHITESPACE
		// http://en.wikipedia.org/wiki/INI_file#Whitespace - there are a few options here.
		{
			conf: "  foo=    bar flag",
			want: map[string]map[string]string{
				Default: map[string]string{
					"foo": "    bar flag",
				},
			},
		},

		// Quoted values - http://en.wikipedia.org/wiki/INI_file#Quoted_values
		// TODO do we agree that quotes should just be included literally?
		{
			conf: `
bar="foo"
baz='bar'
bat=baseball`,
			want: map[string]map[string]string{
				Default: map[string]string{
					"bar": `"foo"`,
					"baz": `'bar'`,
					"bat": "baseball",
				},
			},
		},
	}
	for _, test := range tests {
		c, err := Decode(test.conf)
		if !reflect.DeepEqual(c, test.want) {
			t.Errorf("Decode(%s): got %v want %v", test.conf, c, test.want)
		} else if test.wantErr == "" && err != nil {
			t.Errorf("Decode(%s): got error %v wanted no error", test.conf, err)
		} else if test.wantErr != "" && !strings.Contains(err.Error(), test.wantErr) {
			t.Errorf("Decode(%s): got error %v wanted one containing %q", test.conf, err, test.wantErr)
		}
	}
}
