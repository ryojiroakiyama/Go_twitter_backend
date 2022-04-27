package params

import (
	"net/http"
	"strconv"
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
