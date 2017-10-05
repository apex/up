package secret

import (
	"fmt"
	"strings"

	"github.com/apex/up/platform"
)

// delimiter.
const delim = "/"

// Reserved env variables.
var Reserved = []string{
	"UP_STAGE",
	"PATH",
	"LANG",
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"AWS_SESSION_TOKEN",
}

// Format returns a normalized env var name.
func Format(app, stage, name string) string {
	if stage == "" {
		stage = "all"
	}

	return "/" + strings.Join([]string{
		"up",
		app,
		stage,
		name,
	}, delim)
}

// Parse returns the name stripped of prefix.
func Parse(s string) (app, stage, name string) {
	p := strings.Split(s, delim)
	app, stage, name = p[2], p[3], p[4]
	return
}

// GroupByStage returns secrets mapped by stage.
func GroupByStage(secrets []*platform.Secret) map[string][]*platform.Secret {
	m := make(map[string][]*platform.Secret)

	for _, s := range secrets {
		m[s.Stage] = append(m[s.Stage], s)
	}

	return m
}

// FilterByApp returns secrets filtered by name.
func FilterByApp(secrets []*platform.Secret, name string) (v []*platform.Secret) {
	for _, s := range secrets {
		if s.App == name {
			v = append(v, s)
		}
	}
	return
}

// Env returns a slice of env variables.
func Env(secrets []*platform.Secret) (env []string) {
	for _, s := range secrets {
		env = append(env, fmt.Sprintf("%s=%s", s.Name, s.Value))
	}
	return
}
