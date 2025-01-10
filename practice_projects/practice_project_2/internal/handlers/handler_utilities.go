package handlers

import (
	"example/practice_project_2/internal/templates"

	"github.com/a-h/templ"
)

func wrapComponentWithLayout(c templ.Component) templ.Component {
	return templates.Base(c)
}
