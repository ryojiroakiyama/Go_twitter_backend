package timelines_test

type params struct {
	only_media string
	max_id     string
	since_id   string
	limit      string
}

//?only_media=1&max_id=1
func (p params) asURI() string {
	var s string
	switch {
	case p.only_media != "":
		s = addParam(s, "only_media="+p.only_media)
		fallthrough
	case p.max_id != "":
		s = addParam(s, "max_id="+p.max_id)
		fallthrough
	case p.since_id != "":
		s = addParam(s, "since_id="+p.since_id)
		fallthrough
	case p.limit != "":
		s = addParam(s, "limit="+p.limit)
	}
	return s
}

func addParam(prev string, new string) string {
	if prev == "" {
		return "?" + new
	} else {
		return prev + "&" + new
	}
}
