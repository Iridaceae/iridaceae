package rosetta

import (
	"testing"

	helpers2 "github.com/Iridaceae/iridaceae/internal/helpers"
)

// we need a general testRunner to avoid duplication like below.
// I may work this out later.
func TestDefaultHelpCommand_GetterSetter(t *testing.T) {
	cmd := &DefaultHelpCommand{}

	tests := []struct {
		name   string
		actual func() interface{}
	}{
		{"get help invoker", func() interface{} { return cmd.GetInvokers() }},
		{"get help description", func() interface{} { return cmd.GetDescription() }},
		{"get usage", func() interface{} { return cmd.GetUsage() }},
		{"get group", func() interface{} { return cmd.GetGroup() }},
		{"get domain", func() interface{} { return cmd.GetDomain() }},
		{"get sub perms", func() interface{} { return cmd.GetSubPermissionRules() }},
		{"get executable", func() interface{} { return cmd.IsExecutableInDM() }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.actual()
		})
	}
}

func TestDefaultHelpCommand_Exec(t *testing.T) {
	r := NewRouter(makeTestConfig())
	ctx := makeTestCtx(false, false)
	ctx.session = helpers2.MakeTestSession()
	msg := makeTestMsg(t, "!help")

	help := &DefaultHelpCommand{}
	r.Register(help)
	r.(*routerImpl).trigger(ctx.session, msg)
}
