package anon

import (
	"encoding/json"
	"reflect"
	"testing"
)

type test struct {
	ToStar      string   `json:"to_star" anon:"stars"`
	ToEmpty     string   `json:"to_empty" anon:"empty"`
	ToLen       string   `json:"to_len" anon:"stars_with_len"`
	ToInfo      string   `json:"to_info" anon:"with_info"`
	ToInfoRune  string   `json:"to_info_rune" anon:"with_info"`
	ToIgnore    string   `json:"to_ignore"`
	ToSHA       string   `json:"to_sha" anon:"sha512"`
	Inner       Inner    `json:"inner"`
	Slice       []string `json:"slice" anon:"stars_with_len"`
	SliceIgnore []string `json:"slice_ignore"`
	StructSlice []Inner  `json:"struct_slice"`
}

type Inner struct {
	InnerString string `json:"inner_string" anon:"stars"`
}

type ShouldErr struct {
	Field string `json:"field" anon:"unknown!"`
}

type Simple struct {
	Field string `json:"field" anon:"stars"`
}

type SliceErr struct {
	Fields []string `json:"fields" anon:"unknown!"`
}

var four = 4

func Test(t *testing.T) {
	tcs := []struct {
		name      string
		toMarshal any
		want      []byte
		wantErr   bool
	}{
		{
			name: "all good",
			toMarshal: &test{
				ToStar:      "hello, world",
				ToEmpty:     "erase me",
				ToLen:       "swear",
				ToInfo:      "Through the fence, between the curling flower spaces, I could see them hitting.",
				ToInfoRune:  "Para a aventura indefinida, para o Mar Absoluto, para realizar o Impossível!",
				ToIgnore:    "keep me as-is",
				ToSHA:       "hash me please",
				Inner:       Inner{InnerString: "aa"},
				Slice:       []string{"123", "1234", "á2"},
				SliceIgnore: []string{"as-is"},
				StructSlice: []Inner{{InnerString: "aa"}, {InnerString: "bb"}},
			},
			want:    []byte(`{"to_star":"****","to_empty":"","to_len":"*****","to_info":"len:79,is_ascii:true","to_info_rune":"len:77,is_ascii:false","to_ignore":"keep me as-is","to_sha":"\ufffd\ufffd|\ufffd-\r\\\u00186*\ufffd\ufffd\ufffd%\ufffd/~\ufffdt\ufffd]\u001fP\ufffd\ufffd\u003e\ufffd0\ufffd\ufffdK\ufffd\u0018\ufffdd\"5\ufffd\ufffd[\ufffd2\ufffd)\ufffd\ufffd\u000f\"\u0014\ufffd\ufffd:8\ufffd{H\ufffd\u001c#Cc\ufffd\ufffdV","inner":{"inner_string":"****"},"slice":["***","****","***"],"slice_ignore":["as-is"],"struct_slice":[{"inner_string":"****"},{"inner_string":"****"}]}`),
			wantErr: false,
		},
		{
			name:      "unknown anon func",
			toMarshal: &ShouldErr{Field: "a"},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "not a ptr",
			toMarshal: Simple{Field: "a"},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "unknown anon func in slice",
			toMarshal: &SliceErr{Fields: []string{"a", "b"}},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "random slice, does nothing because it's not tagged",
			toMarshal: &[]string{"a", "b"},
			want:      []byte(`["a","b"]`),
			wantErr:   false,
		},
		{
			name:      "random int, does nothing because it's not tagged",
			toMarshal: &four,
			want:      []byte("4"),
			wantErr:   false,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Marshal(tc.toMarshal, json.Marshal)
			if err != nil && !tc.wantErr {
				t.Fatalf("no err expected, got %s", err.Error())
			} else if err == nil && tc.wantErr {
				t.Fatalf("should return error")
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Not equal!\ngot: %s\nwant: %s", string(got), string(tc.want))
			}
		})
	}
}
