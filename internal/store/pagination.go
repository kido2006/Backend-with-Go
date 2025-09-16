package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginationQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since" `
	Until  string   `json:"until"`
}

func (fq *PaginationQuery) Parse(r *http.Request) (PaginationQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return *fq, nil
		}
		fq.Limit = l
	}

	offset := qs.Get("offset")
	if limit != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return *fq, nil
		}
		fq.Offset = l
	}
	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {

		fq.Since = parseTime(since)
	}

	until := qs.Get("until")
	if until != "" {
		fq.Until = parseTime(until)
	}

	return *fq, nil
}

func parseTime(timestr string) string {
	t, err := time.Parse(time.DateTime, timestr)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
