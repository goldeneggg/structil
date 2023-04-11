package decoder_test

import (
	"testing"

	. "github.com/goldeneggg/structil/dynamicstruct/decoder"
)

var (
	singleJSON = []byte(`
{
  "null_field":null,
  "string_field":"かきくけこ",
  "int_field":45678,
  "float64_field":9.876,
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
    "float64_field":9.876,
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
    "float64_field":4.99,
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
float64_field: 9.876
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
  float64_field: 9.876
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
  float64_field: 4.99
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
float64_field = "9.876,"
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
    <float64_field>9.876</float64_field>
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
        <float64_field>9.876</float64_field>
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
        <float64_field>4.99</float64_field>
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
null_field =

string_field = "かきくけこ,"

int_field = 45678

float64_field = "9.876,"

bool_field = false

struct_ptr_field = {
  key = "hugakey"
	value = "hugavalue"
}

array_string_field = ["array_str_1", "array_str_2"]

array_struct_field = {
  kkk = "kkk1"
  vvvv = "vvv1"
}

"array_struct_field" = {
  kkk = "kkk2"
  vvvv = "vvv2"
}

array_struct_field = {
  kkk = "kkk3"
  vvvv = "vvv3"
}
`)
)

func BenchmarkFromJSON_singleJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromJSON(singleJSON)
	}
}

func BenchmarkDynamicStruct_singleJSON_nonNest_nonUseTag(b *testing.B) {
	d, _ := FromJSON(singleJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStruct_singleJSON_nest_nonUseTag(b *testing.B) {
	d, _ := FromJSON(singleJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, false)
	}
}

func BenchmarkDynamicStruct_singleJSON_nest_useTag(b *testing.B) {
	d, _ := FromJSON(singleJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, true)
	}
}

func BenchmarkFromJSON_arrayJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromJSON(arrayJSON)
	}
}

func BenchmarkDynamicStruct_arrayJSON_nonNest_nonUseTag(b *testing.B) {
	d, _ := FromJSON(arrayJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStruct_arrayJSON_nest_nonUseTag(b *testing.B) {
	d, _ := FromJSON(arrayJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, false)
	}
}

func BenchmarkDynamicStruct_arrayJSON_nest_useTag(b *testing.B) {
	d, _ := FromJSON(arrayJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, true)
	}
}

func BenchmarkJSONToGetter_singleJSON_nonNest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JSONToGetter(singleJSON, false)
	}
}

func BenchmarkJSONToGetter_singleJSON_nest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JSONToGetter(singleJSON, true)
	}
}

func BenchmarkJSONToGetter_arrayJSON_nonNest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JSONToGetter(arrayJSON, false)
	}
}

func BenchmarkJSONToGetter_arrayJSON_nest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JSONToGetter(arrayJSON, true)
	}
}

func BenchmarkFromYAML_singleYAML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromYAML(singleYAML)
	}
}

func BenchmarkDynamicStruct_singleYAML_nonNest_nonUseTag(b *testing.B) {
	d, _ := FromYAML(singleYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStruct_singleYAML_nest_nonUseTag(b *testing.B) {
	d, _ := FromYAML(singleYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, false)
	}
}

func BenchmarkDynamicStruct_singleYAML_nest_useTag(b *testing.B) {
	d, _ := FromYAML(singleYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, false)
	}
}

func BenchmarkFromYAML_arrayYAML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromYAML(arrayYAML)
	}
}

func BenchmarkDynamicStruct_arrayYAML_nonNest_nonUseTag(b *testing.B) {
	d, _ := FromYAML(arrayYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStruct_arrayYAML_nest_nonUseTag(b *testing.B) {
	d, _ := FromYAML(arrayYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, false)
	}
}

func BenchmarkDynamicStruct_arrayYAML_nest_useTag(b *testing.B) {
	d, _ := FromYAML(arrayYAML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.DynamicStruct(true, true)
	}
}

func BenchmarkYAMLToGetter_singleYAML_nonNest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = YAMLToGetter(singleYAML, false)
	}
}

func BenchmarkYAMLToGetter_singleYAML_nest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = YAMLToGetter(singleYAML, true)
	}
}

func BenchmarkYAMLToGetter_arrayYAML_nonNest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = YAMLToGetter(arrayYAML, false)
	}
}

func BenchmarkYAMLToGetter_arrayYAML_nest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = YAMLToGetter(arrayYAML, true)
	}
}
