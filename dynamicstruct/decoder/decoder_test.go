package decoder_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/goldeneggg/structil"
	. "github.com/goldeneggg/structil/dynamicstruct/decoder"
)

const (
	typeJSON int = iota
	typeYAML
	typeHCL
	typeXML // FIXME
)

var (
	singleJSON = []byte(`
{
  "null_field":null,
  "string_field":"かきくけこ",
  "int_field":45678,
  "float32_field":9.876,
  "bool_field":false,
  "struct_ptr_field":{
    "key":"hugakey",
    "value":"hugavalue"
  },
  "array_string_field":[
    "array_str_1",
    "array_str_2"
  ],
  "array_struct_field":[
    {
      "kkk":"kkk1",
      "vvvv":"vvv1"
    },
    {
      "kkk":"kkk2",
      "vvvv":"vvv2"
    },
    {
      "kkk":"kkk3",
      "vvvv":"vvv3"
    }
  ]
}
`)

	arrayJSON = []byte(`
[
  {
    "null_field":null,
    "string_field":"かきくけこ",
    "int_field":45678,
    "float32_field":9.876,
    "bool_field":false,
    "struct_ptr_field":{
      "key":"hugakey",
      "value":"hugavalue"
    },
    "array_string_field":[
      "array_str_1",
      "array_str_2"
    ],
    "array_struct_field":[
      {
        "kkk":"kkk1",
        "vvvv":"vvv1"
      },
      {
        "kkk":"kkk2",
        "vvvv":"vvv2"
      },
      {
        "kkk":"kkk3",
        "vvvv":"vvv3"
      }
    ]
  },
  {
    "null_field":null,
    "string_field":"さしすせそ",
    "int_field":7890,
    "float32_field":4.99,
    "bool_field":true,
    "struct_ptr_field":{
      "key":"hugakeyXXX",
      "value":"hugavalueXXX"
    },
    "array_string_field":[
      "array_str_111",
      "array_str_222"
    ],
    "array_struct_field":[
      {
        "kkk":"kkk99",
        "vvvv":"vvv99"
      },
      {
        "kkk":"kkk999",
        "vvvv":"vvv999"
      },
      {
        "kkk":"kkk9999",
        "vvvv":"vvv9999"
      }
    ]
  }
]
`)

	singleYAML = []byte(`
null_field: null
string_field: かきくけこ
int_field: 45678
float32_field: 9.876
bool_field: false
struct_ptr_field:
  key: hugakey
  value: hugavalue
array_string_field:
  - array_str_1
  - array_str_2
array_struct_field:
  - kkk: kkk1
    vvvv: vvv1
  - kkk: kkk2
    vvvv: vvv2
  - kkk: kkk3
    vvvv: vvv3
`)

	arrayYAML = []byte(`
- null_field: null
  string_field: かきくけこ
  int_field: 45678
  float32_field: 9.876
  bool_field: false
  struct_ptr_field:
    key: hugakey
    value: hugavalue
  array_string_field:
    - array_str_1
    - array_str_2
  array_struct_field:
    - kkk: kkk1
      vvvv: vvv1
    - kkk: kkk2
      vvvv: vvv2
    - kkk: kkk3
      vvvv: vvv3
- null_field: null
  string_field: さしすせそ
  int_field: 7890
  float32_field: 4.99
  bool_field: true
  struct_ptr_field:
    key: hugakeyXXX
    value: hugavalueXXX
  array_string_field:
    - array_str_111
    - array_str_222
  array_struct_field:
    - kkk: kkk99
      vvvv: vvv99
    - kkk: kkk999
      vvvv: vvv999
    - kkk: kkk9999
      vvvv: vvv9999
`)

	//lint:ignore U1000 It's ok because this is for the future.
	singleTOML = []byte(`
string_field = "かきくけこ,"
int_field = 45678
float32_field = "9.876,"
bool_field = false
array_string_field = ["array_str_1", "array_str_2"]

[struct_ptr_field]
  key = "hugakey"
  value = "hugavalue"

[[array_struct_field]]
  kkk = "kkk1"
  vvvv = "vvv1"

[[array_struct_field]]
  kkk = "kkk2"
  vvvv = "vvv2"

[[array_struct_field]]
  kkk = "kkk3"
  vvvv = "vvv3"
`)

	//lint:ignore U1000 It's ok because this is for the future.
	singleXML = []byte(`
<?xml version="1.0" encoding="UTF-8" ?>
<root>
    <null_field/>
    <string_field>かきくけこ</string_field>
    <int_field>45678</int_field>
    <float32_field>9.876</float32_field>
    <bool_field>false</bool_field>
    <struct_ptr_field>
        <key>hugakey</key>
        <value>hugavalue</value>
    </struct_ptr_field>
    <array_string_field>array_str_1</array_string_field>
    <array_string_field>array_str_2</array_string_field>
    <array_struct_field>
        <kkk>kkk1</kkk>
        <vvvv>vvv1</vvvv>
    </array_struct_field>
    <array_struct_field>
        <kkk>kkk2</kkk>
        <vvvv>vvv2</vvvv>
    </array_struct_field>
    <array_struct_field>
        <kkk>kkk3</kkk>
        <vvvv>vvv3</vvvv>
    </array_struct_field>
</root>
`)

	//lint:ignore U1000 It's ok because this is for the future.
	arrayXML = []byte(`
<?xml version="1.0" encoding="UTF-8" ?>
<root>
    <0>
        <null_field/>
        <string_field>かきくけこ</string_field>
        <int_field>45678</int_field>
        <float32_field>9.876</float32_field>
        <bool_field>false</bool_field>
        <struct_ptr_field>
            <key>hugakey</key>
            <value>hugavalue</value>
        </struct_ptr_field>
        <array_string_field>array_str_1</array_string_field>
        <array_string_field>array_str_2</array_string_field>
        <array_struct_field>
            <kkk>kkk1</kkk>
            <vvvv>vvv1</vvvv>
        </array_struct_field>
        <array_struct_field>
            <kkk>kkk2</kkk>
            <vvvv>vvv2</vvvv>
        </array_struct_field>
        <array_struct_field>
            <kkk>kkk3</kkk>
            <vvvv>vvv3</vvvv>
        </array_struct_field>
    </0>
    <1>
        <null_field/>
        <string_field>さしすせそ</string_field>
        <int_field>7890</int_field>
        <float32_field>4.99</float32_field>
        <bool_field>true</bool_field>
        <struct_ptr_field>
            <key>hugakeyXXX</key>
            <value>hugavalueXXX</value>
        </struct_ptr_field>
        <array_string_field>array_str_111</array_string_field>
        <array_string_field>array_str_222</array_string_field>
        <array_struct_field>
            <kkk>kkk99</kkk>
            <vvvv>vvv99</vvvv>
        </array_struct_field>
        <array_struct_field>
            <kkk>kkk999</kkk>
            <vvvv>vvv999</vvvv>
        </array_struct_field>
        <array_struct_field>
            <kkk>kkk9999</kkk>
            <vvvv>vvv9999</vvvv>
        </array_struct_field>
    </1>
</root>
`)

	//lint:ignore U1000 It's ok because this is for the future.
	singleHCL = []byte(`
"null_field" =

"string_field" = "かきくけこ,"

"int_field" = 45678

"float32_field" = "9.876,"

"bool_field" = false

"struct_ptr_field" = {
    "key" = "hugakey"

    "value" = "hugavalue"
}

"array_string_field" = ["array_str_1", "array_str_2"]

"array_struct_field" = {
    "kkk" = "kkk1"

    "vvvv" = "vvv1"
}

"array_struct_field" = {
    "kkk" = "kkk2"

    "vvvv" = "vvv2"
}

"array_struct_field" = {
    "kkk" = "kkk3"

    "vvvv" = "vvv3"
}
`)
)

type decoderTest struct {
	name               string
	data               []byte
	dt                 int
	nest               bool
	useTag             bool
	wantNumF           int
	wantDefinition     string
	fieldAndNestFields map[string][]string
	wantErrorNew       bool
	wantErrorDs        bool
}

func TestDynamicStructJSON(t *testing.T) {
	t.Parallel()

	tests := []decoderTest{
		{
			name: "HasOnlyPrimitive",
			data: []byte(`
{
	"null_field":null,
	"string_field":"かきくけこ",
	"int_field":45678,
	"float32_field":9.876,
	"bool_field":false
}
`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 5,
			wantDefinition: `type DynamicStruct struct {
	BoolField bool
	Float32Field float64
	IntField float64
	NullField interface {}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"BoolField":    nil,
				"Float32Field": nil,
				"IntField":     nil,
				"NullField":    nil,
				"StringField":  nil,
			},
		},
		{
			name: "HasObj",
			data: []byte(`
{
	"string_field":"あああ",
	"obj_field":{
		"id":123,
		"name":"Test Tarou"
	}
}
`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField map[string]interface {}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"StringField": nil,
				"ObjField":    nil,
			},
		},
		{
			name: "HasObjWithTag",
			data: []byte(`
{
	"string_field":"あああ",
	"obj_field":{
		"id":123,
		"name":"Test Tarou"
	}
}
`),
			dt:       typeJSON,
			nest:     false,
			useTag:   true,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField map[string]interface {} ` + "`json:\"obj_field\"`" + `
	StringField string ` + "`json:\"string_field\"`" + `
}`,
			fieldAndNestFields: map[string][]string{
				"StringField": nil,
				"ObjField":    nil,
			},
		},
		{
			name: "HasObjWithNest",
			data: []byte(`
{
	"string_field":"あああ",
	"obj_field":{
		"id":123,
		"name":"Test Tarou"
	}
}
`),
			dt:       typeJSON,
			nest:     true,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Id float64
		Name string
	}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"StringField": nil,
				"ObjField":    {"Id", "Name"},
			},
		},
		{
			name: "HasObjWithTagNest",
			data: []byte(`
{
	"string_field":"あああ",
	"obj_field":{
		"id":123,
		"name":"Test Tarou",
		"obj_array_field": [
			{
				"k1":"v1",
				"k2":"v2"
			},
			{
				"k1":"v111",
				"k2":"v222"
			}
		]
	}
}
`),
			dt:       typeJSON,
			nest:     true,
			useTag:   true,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Id float64 ` + "`json:\"id\"`" + `
		Name string ` + "`json:\"name\"`" + `
		ObjArrayField []struct {
			K1 string ` + "`json:\"k1\"`" + `
			K2 string ` + "`json:\"k2\"`" + `
		} ` + "`json:\"obj_array_field\"`" + `
	} ` + "`json:\"obj_field\"`" + `
	StringField string ` + "`json:\"string_field\"`" + `
}`,
			fieldAndNestFields: map[string][]string{
				"StringField": nil,
				"ObjField":    {"Id", "Name", "ObjArrayField"},
			},
		},
		{
			name: "HasObjTwoNest",
			data: []byte(`
{
	"string_field":"あああ",
	"obj_field":{
		"id":45,
		"name":"Test Jiou",
		"boss":true,
		"objobj_field":{
			"user_id":678,
			"status":"progress"
		}
	}
}
`),
			dt:       typeJSON,
			nest:     true,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Boss bool
		Id float64
		Name string
		ObjobjField struct {
			Status string
			UserId float64
		}
	}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"StringField": nil,
				"ObjField":    {"Boss", "Id", "Name", "ObjobjField"},
			},
		},
		{
			name: "HasArrayString1",
			data: []byte(`
{
	"string_field":"あああ",
	"string_array_field":[
		"id1"
	]
}
`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	StringArrayField []string
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"StringArrayField": nil,
				"StringField":      nil,
			},
		},
		{
			name: "HasArrayString2",
			data: []byte(`
{
	"string_field":"あああ",
	"string_array_field":[
		"id1",
		"id2"
	]
}
`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	StringArrayField []string
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"StringArrayField": nil,
				"StringField":      nil,
			},
		},
		{
			name: "HasArrayStringInArray",
			data: []byte(`
{
	"string_field":"あああ",
	"string_array_field":[
		["id11","id12"],
		["id21","id22"]
	]
}
`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			// FIXME: [][]interface ではなく [][]string になるように修正したい
			wantDefinition: `type DynamicStruct struct {
	StringArrayField [][]interface {}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"StringArrayField": nil,
				"StringField":      nil,
			},
		},
		{
			name: "HasObjectInArray",
			data: []byte(`
{
	"string_field":"あああ",
	"object_array_field":[
		{"nest_str":"aaa","nest_num":23},
		{"nest_str":"bbb","nest_num":34}
	]
}
`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjectArrayField []map[string]interface {}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"ObjectArrayField": nil,
				"StringField":      nil,
			},
		},
		// Note: トップレベルが配列のJSONは、配列の1番目の要素を使って処理する
		{
			name: "TopLevelIsArrayNestFalse",
			data: []byte(`
[
	{
		"string_field":"あああ",
		"objobj_field":{
			"user_id":678,
			"status":"progress"
		},
		"string_array_field":[
			"id1",
			"id2"
		]
	},
	{
		"string_field":"いいいい",
		"string_array_field":[
			"id4",
			"id5",
			"id6"
		]
	}	
]
`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 3,
			wantDefinition: `type DynamicStruct struct {
	ObjobjField map[string]interface {}
	StringArrayField []string
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"ObjobjField":      nil,
				"StringArrayField": nil,
				"StringField":      nil,
			},
		},
		// Note: トップレベルが配列のJSONは、配列の1番目の要素を使って処理する（nest=trueでも同様）
		{
			name: "TopLevelIsArrayNestTrue",
			data: []byte(`
[
	{
		"string_field":"あああ",
		"objobj_field":{
			"user_id":678,
			"status":"progress"
		},
		"string_array_field":[
			"id1",
			"id2"
		]
	},
	{
		"string_field":"いいいい",
		"string_array_field":[
			"id4",
			"id5",
			"id6"
		]
	}	
]
`),
			dt:       typeJSON,
			nest:     true,
			useTag:   false,
			wantNumF: 3,
			wantDefinition: `type DynamicStruct struct {
	ObjobjField struct {
		Status string
		UserId float64
	}
	StringArrayField []string
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"ObjobjField":      {"Status", "UserId"},
				"StringArrayField": nil,
				"StringField":      nil,
			},
		},
		{
			name:     "BracketOnly",
			data:     []byte(`{}`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 0,
			wantDefinition: `type DynamicStruct struct {
}`,
		},
		{
			name:     "ArrayBracketOnly",
			data:     []byte(`[]`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 0,
			wantDefinition: `type DynamicStruct struct {
}`,
		},
		{
			name:           "OnlyLiteral",
			data:           []byte(`aiueo`),
			dt:             typeJSON,
			nest:           false,
			useTag:         false,
			wantNumF:       0,
			wantDefinition: ``,
			wantErrorNew:   true,
		},
		{
			name:           "Empty",
			data:           []byte(``),
			dt:             typeJSON,
			nest:           false,
			useTag:         false,
			wantNumF:       0,
			wantDefinition: ``,
			wantErrorNew:   true,
		},
		{
			name:         "NullData",
			data:         nil,
			dt:           typeJSON,
			nest:         false,
			useTag:       false,
			wantErrorNew: true,
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dec, err := FromJSON(tt.data)
			if err != nil {
				if !tt.wantErrorNew {
					t.Fatalf("unexpected error is returned from FromJSON: %v", err)
				}
				return
			} else if tt.wantErrorNew {
				t.Fatalf("error is expected but it does not occur from FromJSON. data: %q", string(tt.data))
			}

			testCorrectCase(t, tt, dec)
		})
	}
}

func TestDynamicStructYAML(t *testing.T) {
	t.Parallel()

	tests := []decoderTest{
		{
			name: "HasOnlyPrimitive",
			data: []byte(`
null_field: null
string_field: かきくけこ
int_field: 45678
float32_field: 9.876
bool_field: false
`),
			dt:       typeYAML,
			nest:     false,
			useTag:   false,
			wantNumF: 5,
			wantDefinition: `type DynamicStruct struct {
	BoolField bool
	Float32Field float64
	IntField int
	NullField interface {}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"BoolField":    nil,
				"Float32Field": nil,
				"IntField":     nil,
				"NullField":    nil,
				"StringField":  nil,
			},
		},
		{
			name: "HasObj",
			data: []byte(`
string_field: かきくけこ
obj_field:
  id: 123
  name: Test Tarou
`),
			dt:       typeYAML,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField map[string]interface {}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"ObjField":    nil,
				"StringField": nil,
			},
		},
		{
			name: "HasObjWithTag",
			data: []byte(`
string_field: かきくけこ
obj_field:
  id: 123
  name: Test Tarou
`),
			dt:       typeYAML,
			nest:     false,
			useTag:   true,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField map[string]interface {} ` + "`yaml:\"obj_field\"`" + `
	StringField string ` + "`yaml:\"string_field\"`" + `
}`,
			fieldAndNestFields: map[string][]string{
				"ObjField":    nil,
				"StringField": nil,
			},
		},
		{
			name: "HasObjWithNest",
			data: []byte(`
string_field: かきくけこ
obj_field:
  id: 123
  name: Test Tarou
`),
			dt:       typeYAML,
			nest:     true,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Id int
		Name string
	}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"ObjField":    {"Id", "Name"},
				"StringField": nil,
			},
		},
		{
			name: "HasObjWithTagNest",
			data: []byte(`
string_field: かきくけこ
obj_field:
  id: 123
  name: Test Tarou
`),
			dt:       typeYAML,
			nest:     true,
			useTag:   true,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Id int ` + "`yaml:\"id\"`" + `
		Name string ` + "`yaml:\"name\"`" + `
	} ` + "`yaml:\"obj_field\"`" + `
	StringField string ` + "`yaml:\"string_field\"`" + `
}`,
			fieldAndNestFields: map[string][]string{
				"ObjField":    {"Id", "Name"},
				"StringField": nil,
			},
		},
		{
			name: "HasObjTwoNest",
			data: []byte(`
string_field: あああ
obj_field:
  id: 45
  name: Test Jiou
  boss: true
  objobj_field:
    user_id: 678
    status: progress
`),
			dt:       typeYAML,
			nest:     true,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Boss bool
		Id int
		Name string
		ObjobjField struct {
			Status string
			UserId int
		}
	}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"ObjField":    {"Boss", "Id", "Name", "ObjobjField"},
				"StringField": nil,
			},
		},
		{
			name: "HasArrayString",
			data: []byte(`
string_field: あああ
string_array_field:
  - id1
  - id2
`),
			dt:       typeYAML,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	StringArrayField []string
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"StringArrayField": nil,
				"StringField":      nil,
			},
		},
		{
			name: "HasArrayObject",
			data: []byte(`
string_field: いいい
obj_field:
  user_id: 678
  status: progress
arr_obj_field:
  - aid: 45
    aname: Test Mike
  - aid: 678
    aname: Test Davis
`),
			dt:       typeYAML,
			nest:     true,
			useTag:   false,
			wantNumF: 3,
			wantDefinition: `type DynamicStruct struct {
	ArrObjField []struct {
		Aid int
		Aname string
	}
	ObjField struct {
		Status string
		UserId int
	}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"ArrObjField": nil,
				"ObjField":    {"Status", "UserId"},
				"StringField": nil,
			},
		},
		{
			name:         "OnlyLiteral",
			data:         []byte(`aiueo`),
			dt:           typeYAML,
			nest:         false,
			useTag:       false,
			wantNumF:     0,
			wantErrorNew: true,
		},
		{
			name: "HasArrayStringInArray",
			data: []byte(`
string_field: あああ
string_array_field:
 - - id11
   - id12
 - - id21
   - id22
`),
			dt:       typeYAML,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			// FIXME: [][]interface ではなく [][]string になるように修正したい
			wantDefinition: `type DynamicStruct struct {
	StringArrayField [][]interface {}
	StringField string
}`,
			fieldAndNestFields: map[string][]string{
				"StringArrayField": nil,
				"StringField":      nil,
			},
		},
		{
			name: "TopLevelIsArrayObject",
			data: []byte(`
- fieldA: aaa
  fieldB: 23
- fieldA: bbb
  fieldB: 34
`),
			dt:       typeYAML,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	FieldA string
	FieldB int
}`,
			fieldAndNestFields: map[string][]string{
				"FieldA": nil,
				"FieldB": nil,
			},
		},
		{
			name:         "Empty",
			data:         []byte(``),
			dt:           typeYAML,
			nest:         false,
			useTag:       false,
			wantNumF:     0,
			wantErrorNew: true,
		},
		{
			name:         "NullData",
			data:         nil,
			dt:           typeYAML,
			nest:         false,
			useTag:       false,
			wantErrorNew: true,
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dec, err := FromYAML(tt.data)
			if err != nil {
				if !tt.wantErrorNew {
					t.Fatalf("unexpected error is returned from FromYAML: %v", err)
				}
				return
			} else if tt.wantErrorNew {
				t.Fatalf("error is expected but it does not occur from FromYAML. data: %q", string(tt.data))
			}

			// FIXME: エラーケース対応
			testCorrectCase(t, tt, dec)
		})
	}
}

func TestDynamicStructHCL(t *testing.T) {
	t.Parallel()

	tests := []decoderTest{
		{
			name: "HasOnlyPrimitive",
			data: []byte(`
string_field = "かきくけこ"
int_field = 45678
float32_field = 9.876
bool_field = false
`),
			dt:       typeHCL,
			nest:     false,
			useTag:   false,
			wantNumF: 4,
			wantDefinition: `type DynamicStruct struct {
	BoolField bool
	Float32Field float64
	IntField float64
	StringField string
}`,
		},
		{
			name:     "BracketOnly",
			data:     []byte(`{}`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 0,
			wantDefinition: `type DynamicStruct struct {
}`,
		},
		{
			name:     "ArrayBracketOnly",
			data:     []byte(`[]`),
			dt:       typeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 0,
			wantDefinition: `type DynamicStruct struct {
}`,
		},
		{
			name:           "OnlyLiteral",
			data:           []byte(`aiueo`),
			dt:             typeJSON,
			nest:           false,
			useTag:         false,
			wantNumF:       0,
			wantDefinition: ``,
			wantErrorNew:   true,
		},
		{
			name:           "Empty",
			data:           []byte(``),
			dt:             typeJSON,
			nest:           false,
			useTag:         false,
			wantNumF:       0,
			wantDefinition: ``,
			wantErrorNew:   true,
		},
		{
			name:         "NullData",
			data:         nil,
			dt:           typeJSON,
			nest:         false,
			useTag:       false,
			wantErrorNew: true,
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dec, err := FromHCL(tt.data)
			if err != nil {
				if !tt.wantErrorNew {
					t.Fatalf("unexpected error is returned from FromHCL: %v", err)
				}
				return
			} else if tt.wantErrorNew {
				t.Fatalf("error is expected but it does not occur from FromHCL. data: %q", string(tt.data))
			}

			testCorrectCase(t, tt, dec)
		})
	}
}

func TestDynamicStructFixmeXml(t *testing.T) {
	t.Parallel()

	tests := []decoderTest{
		{
			name: "ValidYamlButTypeIsInvalid",
			data: []byte(`
null_field: null
string_field: かきくけこ
int_field: 45678
float32_field: 9.876
bool_field: false
`),
			dt:           typeXML,
			nest:         false,
			useTag:       false,
			wantNumF:     0,
			wantErrorNew: true,
		},
		{
			name:         "Empty",
			data:         []byte(``),
			dt:           typeXML,
			nest:         false,
			useTag:       false,
			wantNumF:     0,
			wantErrorNew: true,
		},
		{
			name:         "NullData",
			data:         nil,
			dt:           typeXML,
			nest:         false,
			useTag:       false,
			wantErrorNew: true,
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := FromXML(tt.data)
			if err != nil {
				if !tt.wantErrorNew {
					t.Fatalf("unexpected error is returned from FromXML: %v", err)
				}
				return
			} else if tt.wantErrorNew {
				t.Fatalf("error is expected but it does not occur from FromXML	. data: %q", string(tt.data))
			}
		})
	}
}

func testCorrectCase(t *testing.T, tt decoderTest, dec *Decoder) {
	t.Helper()

	if d := cmp.Diff(dec.OrgData(), tt.data); d != "" {
		t.Fatalf("mismatch OrgData: (-got +want)\n%s", d)
	}

	ds, err := dec.DynamicStruct(tt.nest, tt.useTag)
	if err != nil {
		if !tt.wantErrorDs {
			t.Fatalf("unexpected error is returned from DynamicStruct: %v", err)
		}
		return
	} else if tt.wantErrorDs {
		t.Fatalf("error is expected but it does not occur from DynamicStruct. data: %q", string(tt.data))
	}

	if ds == nil {
		t.Fatalf("unexpected DynamicStruct is null. got: is null, want: is not null")
	}

	if ds.NumField() != tt.wantNumF {
		t.Fatalf("unmatch numfield. got: %d, want: %d, ds.Definition:\n%s", ds.NumField(), tt.wantNumF, ds.Definition())
	}

	if d := cmp.Diff(ds.Definition(), tt.wantDefinition); d != "" {
		t.Fatalf("mismatch Definition: (-got +want)\n%s", d)
	}

	if len(tt.fieldAndNestFields) > 0 {
		var g *structil.Getter

		switch tt.dt {
		case typeJSON:
			g, err = JSONToGetter(tt.data, tt.nest)
			if err != nil {
				t.Fatalf("unexpected error is returned from JSONToGetter: %v", err)
			}
		case typeYAML:
			g, err = YAMLToGetter(tt.data, tt.nest)
			if err != nil {
				t.Fatalf("unexpected error is returned from YAMLToGetter: %v", err)
			}
		}

		for n, nests := range tt.fieldAndNestFields {
			testToGetter(t, n, nests, g)
		}
	}
}

func testToGetter(t *testing.T, n string, nests []string, g *structil.Getter) {
	t.Helper()

	if !g.Has(n) {
		t.Fatalf("Getter does not have %s field", n)
	}
	if len(nests) > 0 {
		ni, _ := g.Get(n)
		gg, err := structil.NewGetter(ni)
		if err != nil {
			t.Fatalf("unexpected error is returned from structil.NewGetter(ni). ni = %+v: %v", ni, err)
		}

		for _, nn := range nests {
			testToGetter(t, nn, nil, gg)
		}
	}
}
