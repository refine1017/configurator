package table

type Col struct {
	Idx     int         `bson:"-" json:"idx"` // just for edit by
	Name    string      `bson:"name" json:"name"`
	Type    string      `bson:"type" json:"type"`
	Desc    string      `bson:"desc" json:"desc"`
	Feature uint32      `bson:"-" json:"feature"`
	Index   string      `bson:"index" json:"index"`
	Data    interface{} `bson:"data" json:"data"`
}

func FieldCols() []Col {
	return []Col{
		{
			Name:    "name",
			Type:    TypeString,
			Desc:    "Field name",
			Feature: ColFeature(EditAble, DisableChange, SearchAble),
		},
		{
			Name:    "type",
			Type:    TypeFields,
			Desc:    "Field type",
			Feature: ColFeature(EditAble),
		},
		{
			Name:    "desc",
			Type:    TypeText,
			Desc:    "Field desc",
			Feature: ColFeature(EditAble),
		},
		{
			Name:    "index",
			Type:    TypeEnum,
			Desc:    "index field",
			Feature: ColFeature(EditAble),
			Data:    []string{"primary", "index"},
		},
	}
}
