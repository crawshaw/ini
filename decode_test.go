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
					"a": "b#this is a comment",
				},
			},
			// Or maybe this should error.
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
