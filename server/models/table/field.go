package table

const (
	TypeString = "string"
	TypeText   = "text"
	TypeBool   = "bool"
	TypeFloat  = "float"
	TypeDouble = "double"
	TypeInt32  = "int32"
	TypeInt64  = "int64"
	TypeUint32 = "uint32"
	TypeUint64 = "uint64"

	TypeStringArray = "string_array"
	TypeBoolArray   = "bool_array"
	TypeUint32Array = "uint32_array"
	TypeUint64Array = "uint64_array"
	TypeFloatArray  = "float_array"
	TypeDoubleArray = "double_array"
	TypeInt32Array  = "int32_array"
	TypeInt64Array  = "int64_array"

	TypeEnum = "enum"

	TypeFields = "fields"
)

var FieldTypeFeature = map[string]uint32{
	TypeString:      ColFeature(EditAble, SortAble),
	TypeText:        ColFeature(EditAble, SortAble),
	TypeBool:        ColFeature(EditAble, SortAble),
	TypeFloat:       ColFeature(EditAble, SortAble),
	TypeDouble:      ColFeature(EditAble, SortAble),
	TypeInt32:       ColFeature(EditAble, SortAble),
	TypeInt64:       ColFeature(EditAble, SortAble),
	TypeUint32:      ColFeature(EditAble, SortAble),
	TypeUint64:      ColFeature(EditAble, SortAble),
	TypeStringArray: ColFeature(EditAble, SortAble),
	TypeBoolArray:   ColFeature(EditAble, SortAble),
	TypeUint32Array: ColFeature(EditAble, SortAble),
	TypeUint64Array: ColFeature(EditAble, SortAble),
	TypeFloatArray:  ColFeature(EditAble, SortAble),
	TypeDoubleArray: ColFeature(EditAble, SortAble),
	TypeInt32Array:  ColFeature(EditAble, SortAble),
	TypeInt64Array:  ColFeature(EditAble, SortAble),
}

var fieldTypeList = []string{
	TypeString,
	TypeText,
	TypeBool,
	TypeFloat,
	TypeDouble,
	TypeUint32,
	TypeUint64,
	TypeInt32,
	TypeInt64,
	TypeStringArray,
	TypeBoolArray,
	TypeUint32Array,
	TypeUint64Array,
	TypeFloatArray,
	TypeDoubleArray,
	TypeInt32Array,
	TypeInt64Array,
}

func FieldTypeList() []string {
	return fieldTypeList
}
