package dotfile

import "embed"

//go:embed all:template
var TemplateFS embed.FS

//go:embed VERSION
var Version string

//go:embed internal/hook/pre-push internal/hook/post-merge
var HookFS embed.FS
