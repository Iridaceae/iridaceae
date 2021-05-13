package rosetta

import (
	"log"
	"strings"
	"sync"
)

func hasPrefix(msg string, prefix string, ignoreCase bool) (string, bool) {
	strToCheck := msg
	if ignoreCase {
		strToCheck = strings.ToLower(strToCheck)
		prefix = strings.ToLower(prefix)
	}
	if strings.HasPrefix(strToCheck, prefix) {
		return msg[len(prefix):], true
	}
	return msg, false
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

func clearMap(m *sync.Map) {
	m.Range(func(key, _ interface{}) bool {
		m.Delete(key)
		return true
	})
}

func getErrorTypeName(e ErrorType) string {
	switch e {
	case ErrTypeCommandExec:
		return "ErrCommandExec"
	case ErrTypeMiddleware:
		return "ErrMiddleware"
	case ErrTypeCommandNotFound:
		return "ErrCommandNotFound"
	case ErrTypeGetChannel:
		return "ErrGetChannel"
	case ErrTypeDeleteCommandMessage:
		return "ErrDeleteCommandMessage"
	case ErrTypeGetGuild:
		return "ErrGetGuild"
	case ErrTypeGuildPrefixGetter:
		return "ErrGuildPrefixGetter"
	case ErrTypeNotExecutableInDM:
		return "ErrNotExecutableInDM"
	}
	return ""
}
