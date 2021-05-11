package log

var _fields = make([]string, 0)

func SetGlobalFields(fields []string) {
	_fields = fields
}

func AddGlobalFields(field string) {
	_fields = append(_fields, field)
}

func GetGlobalFields() []string {
	return _fields
}

func ClearGlobalFields() {
	_fields = make([]string, 0)
}
