package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func trim(src string) string {
	return strings.TrimSpace(src)
}

func file(name string) io.Reader {
	fd, err := os.Open("testdata/" + name)
	if err != nil {
		panic(err)
	}
	return fd
}

func Test_dump2csv(t *testing.T) {
	type args struct {
		sql io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "test basic file",
			args: args{sql: file("simple.sql")},
			wantOut: `1,0,April,1,0,0,0.778582929065,20140312223924,20140312223929,4657771,20236,0
2,0,August,0,0,0,0.123830928525,20140312221818,20140312221822,4360163,11466,0`,
		},
		{
			name:    "test unicode characters",
			args:    args{sql: file("unicode.sql")},
			wantOut: `0,"Unicode test with ""ã"", ""ç"", and ""ü"" characters",1`,
		},
		{
			name:    "test value types with float, null and timestamps",
			args:    args{sql: file("complex_types.sql")},
			wantOut: `1.3,2021-01-01,2021-01-01 00:00:00,,true,false,NULL,FALSE,TRUE`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := dump2csv(tt.args.sql, out); (err != nil) != tt.wantErr {
				t.Errorf("dump2csv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); trim(gotOut) != trim(tt.wantOut) {
				t.Errorf("dump2csv() unexpected result:\nres  [%#v],\nwant [%#v]", gotOut, tt.wantOut)
			}
		})
	}
}
