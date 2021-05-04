package jog

import (
	"log"
	"strings"
)

func hasPrefix(s string, prefixes []string, ignoreCase bool) (bool, string) {
	for _, prefix := range prefixes {
		strToCheck := s
		if ignoreCase {
			strToCheck = strings.ToLower(strToCheck)
			prefix = strings.ToLower(prefix)
		}

		if strings.HasPrefix(strToCheck, prefix) {
			return true, s[len(prefix):]
		}
	}
	return false, s
}

func trimPreSuffix(s string, preSuffix string) string {
	if !(strings.HasPrefix(s, preSuffix) && strings.HasSuffix(s, log.Prefix())) {
		return s
	}
	return strings.TrimPrefix(strings.TrimSuffix(s, preSuffix), preSuffix)
}

func arrayContains(arr []string, s string, ignoreCase bool) bool {
	if ignoreCase {
		s = strings.ToLower(s)
	}
	for _, v := range arr {
		if ignoreCase {
			v = strings.ToLower(v)
		}
		if v == s {
			return true
		}
	}
	return false
}
