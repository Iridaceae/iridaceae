package startupmsg

import (
	// embed will be used for processing template files.
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/Iridaceae/iridaceae/internal/components"
)

//go:embed template.txt
var templateTxt string

type information struct {
	Appname   string
	Copyright string
	Version   string
	Commit    string
	Release   bool
	Repo      string
}

func getInfo() information {
	return information{
		Appname:   "iridaceae",
		Copyright: fmt.Sprintf("© %d Aaron Pham (@aarnphm)", time.Now().Year()),
		Version:   components.AppVersion,
		Commit:    components.AppCommit,
		Release:   components.IsRelease(),
		Repo:      components.Repo,
	}
}

func pass(err error) {
	if err != nil {
		panic(err)
	}
}

func Output(w io.Writer) {
	t, err := template.New("startupmsg").Parse(templateTxt)
	pass(err)
	pass(t.Execute(w, getInfo()))
}
