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

		// TODO Property names - can have underscores?
		{
			conf: "user_name=user1234",
			want: map[string]map[string]string{
				Default: map[string]string{
					"user_name": "user1234",
				},
			},
		},
		// TODO Property names - case sensitive?
		{
			conf: `username=user1234
USERNAME=user`,
			want: map[string]map[string]string{
				Default: map[string]string{
					"username": "user1234",
					"USERNAME": "user",
				},
			},
		},
		// TODO Property names - can they have spaces in them?
		{
			conf: "user name=foo",
			want: map[string]map[string]string{
				Default: map[string]string{
					"user name": "foo",
				},
			},
		},

		// Subsections
		{
			conf: `
global=1
user=user2
[database]
user=user1
password=1234
[foo]
user=user2`,
			want: map[string]map[string]string{
				Default: map[string]string{
					"global": "1",
					"user":   "user2",
				},
				"database": map[string]string{
					"user":     "user1",
					"password": "1234",
				},
				"foo": map[string]string{
					"user": "user2",
				},
			},
		},
		// Subsection - with a comment
		{
			conf: `[database] # this is the db section
	user=foobar`,
			want: map[string]map[string]string{
				"database": map[string]string{
					"user": "foobar",
				},
			},
		},

		// Blank subsection - go to global?
		{
			conf: `user=1
[foo]
user=2
[]
user=3`,
			want: map[string]map[string]string{
				Default: map[string]string{
					"user": "3",
				},
				"foo": map[string]string{
					"user": "2",
				},
			},
		},
		// Multiple subsections with identical name - keep overwriting that section as if they were contiguous
		// in the file
		{
			conf: `
[foo]
bar=1
[baz]
bar=2
[foo]
bar=3`,
			want: map[string]map[string]string{
				"foo": map[string]string{
					"bar": "3",
				},
				"baz": map[string]string{
					"bar": "2",
				},
			},
		},

		// TODO Invalid section header - what should happen?
		{
			conf: `
[foo
bar=1`,
			wantErr: "Invalid section header",
		},
		// TODO Two section headers on same line - what should happen?
		{
			conf: `
[foo][bar]
baz=1`,
			wantErr: "Invalid section header",
		},
	}
	for _, test := range tests {
		c, err := Decode(test.conf)
		if !reflect.DeepEqual(c, test.want) {
			t.Errorf("Decode(%s): got %v want %v", test.conf, c, test.want)
		} else if test.wantErr == "" && err != nil {
			t.Errorf("Decode(%s): got error %v wanted no error", test.conf, err)
		} else if test.wantErr != "" && err == nil {
			t.Errorf("Decode(%s): got nil error wanted one containing %q", test.conf, test.wantErr)
		} else if test.wantErr != "" && !strings.Contains(err.Error(), test.wantErr) {
			t.Errorf("Decode(%s): got error %v wanted one containing %q", test.conf, err, test.wantErr)
		}
	}
}
