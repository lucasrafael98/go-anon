## go-anon

![Tests](https://github.com/lucasrafael98/go-anon/actions/workflows/ci.yml/badge.svg)

This is a fairly simple anonymisation library for Go. If you're logging some requests with sensitive fields,
you should obfuscate them. This library offers several ways to obfuscate data, some completely erasing any
traceability, others keeping some metadata (length, presence of non-ascii chars) to help with any debugging
you'd like to do.

The following functions are available: 
 - `Stars/"stars"`: replaces with `"****"`
 - `Empty/"empty"`: replaces with `""`
 - `StarsWithLen/"stars_with_len"`: replaces with a number of asterisks equal to length
 - `WithInfo/"with_info"`: replaces with info (original string length and whether it is ASCII)
 - `SHA512/"sha512"`: replaces with a SHA-512 hash.

There are three ways to use this library:
 - `Marshal`: Call this with a pointer to a struct and the marshal function you'd like to use, and it will
   anonymise the fields you've setup with tags:
   ```go 
	type Thing struct {
		ID 		  string `json:"id"`
		Sensitive string `json:"sensitive" anon:"stars"`
	}

	...

	thing := Thing{
		ID: "111",
		Sensitive: "password123",
	}
	str, err := anon.Marshal(&thing, json.Marshal)
	// Result:
	// {
	// 	"id": "111",
	// 	"sensitive": "****",
	// }
   ```
 - `Anonymise`: same as above, but you're just changing the struct and not outputting marshalled data.
 	```go 
	str, err := anon.Anonymise(&thing)
	// Result:
	// Thing {
	// 	ID: "111",
	// 	Sensitive: "****",
	// }
	```
 - Calling the functions themselves to replace a single string.
	```go 
	thing.Sensitive := anon.WithInfo(thing.Sensitive) 
	// Result: 
	// Sensitive = "len:11,is_ascii:true"
	```

Here's a more complex example: 
```go 
type test struct {
	ToStar      string   `json:"to_star" anon:"stars"`
	ToEmpty     string   `json:"to_empty" anon:"empty"`
	ToLen       string   `json:"to_len" anon:"stars_with_len"`
	ToInfo      string   `json:"to_info" anon:"with_info"`
	ToInfoRune  string   `json:"to_info_rune" anon:"with_info"`
	ToIgnore    string   `json:"to_ignore"`
	ToSHA       string   `json:"to_sha" anon:"sha512"`
	Inner       inner    `json:"inner"`
	Slice       []string `json:"slice" anon:"stars_with_len"`
	SliceIgnore []string `json:"slice_ignore"`
	StructSlice []inner  `json:"struct_slice"`
}

type inner struct {
	innerString string `json:"inner_string" anon:"stars"`
}

anon.Marshal(&test{
	ToStar:      "hello, world",
	ToEmpty:     "erase me",
	ToLen:       "swear",
	ToInfo:      "Through the fence, between the curling flower spaces, I could see them hitting.",
	ToInfoRune:  "Para a aventura indefinida, para o Mar Absoluto, para realizar o Impossível!",
	ToIgnore:    "keep me as-is",
	ToSHA:       "hash me please",
	Inner:       inner{innerString: "aa"},
	Slice:       []string{"123", "1234", "á2"},
	SliceIgnore: []string{"as-is"},
	StructSlice: []inner{{innerString: "aa"}, {innerString: "bb"}},
}, json.Marshal)
```
With the result below:
```json 
{
  "to_star": "****",
  "to_empty": "",
  "to_len": "*****",
  "to_info": "len:79,is_ascii:true",
  "to_info_rune": "len:77,is_ascii:false",
  "to_ignore": "keep me as-is",
  "to_sha": "��|�-\r\\\u00186*���%�/~�t�]\u001fP��>�0��K�\u0018�d\"5��[�2�)��\u000f\"\u0014��:8�{H�\u001c#Cc��V",
  "inner": {
    "inner_string": "****"
  },
  "slice": [
    "***",
    "****",
    "***"
  ],
  "slice_ignore": [
    "as-is"
  ],
  "struct_slice": [
    {
      "inner_string": "****"
    },
    {
      "inner_string": "****"
    }
  ]
}
```
