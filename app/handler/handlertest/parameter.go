package handlertest

// generate string like '?only_media=1&max_id=1'
func ParamAsURI(p map[string]string) string {
	var s string
	for k, v := range p {
		if s == "" {
			return "?" + k + "=" + v
		} else {
			return s + "&" + k + "=" + v
		}
	}
	return s
}
