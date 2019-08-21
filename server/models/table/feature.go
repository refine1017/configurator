package table

// Features
const (
	EditAble uint32 = iota
	SortAble
	SearchAble
	DisableExport
	DisableChange
	DisableShow
)

var featureSetting = map[string]uint32{
	"edit":          EditAble,
	"sort":          SortAble,
	"search":        SearchAble,
	"disableShow":   DisableShow,
	"disableExport": DisableExport,
	"disableChange": DisableChange,
}

func ColFeature(values ...uint32) uint32 {
	var feature uint32

	for _, v := range values {
		feature |= 1 << v
	}

	return feature
}

func AddFeature(oldFeature uint32, values ...uint32) uint32 {
	var feature = oldFeature

	for _, v := range values {
		feature |= 1 << v
	}

	return feature
}

func FeatureSetting() map[string]uint32 {
	return featureSetting
}
