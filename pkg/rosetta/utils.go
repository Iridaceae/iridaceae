package rosetta

import (
	"log"
	"strings"
	"sync"
)

// hasPrefix will returns striped message without prefix and true if message contains prefix.
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

// TODO: better way to handle error.
func getErrorTypeName(e ErrorType) string {
	switch e {
	case ErrTypeCommandExec:
		return ErrCommandExec.Error()
	case ErrTypeMiddleware:
		return ErrMiddleware.Error()
	case ErrTypeCommandNotFound:
		return ErrCommandNotFound.Error()
	case ErrTypeGetChannel:
		return ErrGetChannel.Error()
	case ErrTypeDeleteCommandMessage:
		return ErrDeleteCommandMessage.Error()
	case ErrTypeGetGuild:
		return ErrGetGuild.Error()
	case ErrTypeGuildPrefixGetter:
		return ErrGuildPrefixGetter.Error()
	case ErrTypeNotExecutableInDM:
		return ErrNotExecutableInDMs.Error()
	}
	return ""
}
