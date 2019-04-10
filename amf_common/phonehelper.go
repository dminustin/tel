package amf_common

import (
	"regexp"
)

//======================================
//PARSERS

func PH_parseClip(clip string) (string) {
	r,_ := regexp.Compile("\\+CLIP: \"([0-9\\+]+)\"")
	if (r.MatchString(clip)) {
		x := r.FindStringSubmatch(clip)
		return x[1]
	}
	return ""
}
