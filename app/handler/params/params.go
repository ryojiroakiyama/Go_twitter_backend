package params

import (
	"net/http"
	"strconv"
)

// Example key name of parameter
// URI/?since_id=10 の部分
const (
	SinceID   = "since_id"
	MaxID     = "max_id"
	Limit     = "limit"
	OnlyMedia = "only_media"
)

func FormValue(r *http.Request, key string, defaut int64, min int64, max int64) int64 {
	if v := r.FormValue(key); v == "" {
		return defaut
	} else {
		if i, err := strconv.ParseInt(v, 10, 64); err != nil {
			return defaut
		} else {
			switch {
			case i < min:
				return min
			case max < i:
				return max
			default:
				return i
			}
		}
	}
}
