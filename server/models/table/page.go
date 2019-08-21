package table

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PageList struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
	Cols  []Col       `json:"cols"`
}

type PageParams struct {
	Page  int    `binding:"Required;MaxSize(255)"`
	Limit int    `binding:"Required;MaxSize(255)"`
	Sort  string `binding:"Required;MaxSize(255)"`
}

func BuildColList(col *mgo.Collection, condition bson.M, params PageParams, list interface{}, cols []Col) (pageList *PageList, err error) {
	pageList = &PageList{}

	skip := 0
	if params.Page > 0 {
		skip = (params.Page - 1) * params.Limit
	}

	query := col.Find(condition)

	pageList.Total, err = query.Count()
	if err != nil {
		return
	}

	if params.Limit > 0 {
		query.Limit(params.Limit).Skip(skip)
	}

	if params.Sort != "" {
		query.Sort(params.Sort)
	}

	err = query.All(list)
	if err != nil {
		return
	}

	pageList.Cols = cols
	pageList.Items = list

	return
}

func BuildArrayList(list []interface{}, params PageParams, cols []Col) (pageList *PageList, err error) {
	pageList = &PageList{}

	start, end := CalculateLimitStartAndEnd(params.Page, params.Limit, len(list))

	pageList.Total = len(list)
	pageList.Items = list[start:end]
	pageList.Cols = cols

	return
}

func CalculateLimitStartAndEnd(page, limit, total int) (int, int) {
	skip := 0
	if page > 0 {
		skip = (page - 1) * limit
	}

	var start = 0
	if total > skip {
		start = skip
	} else {
		skip = total
	}

	var end = 0
	if total > start+limit {
		end = start + limit
	} else {
		end = total
	}

	return start, end
}
