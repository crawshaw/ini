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
