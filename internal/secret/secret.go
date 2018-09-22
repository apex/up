package secret

import (
	"fmt"
	"strings"

	"github.com/apex/up"
)

// delimiter.
const delim = "/"

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
func GroupByStage(secrets []*up.Secret) map[string][]*up.Secret {
	m := make(map[string][]*up.Secret)

	for _, s := range secrets {
		m[s.Stage] = append(m[s.Stage], s)
	}

	return m
}

// FilterByApp returns secrets filtered by app name.
func FilterByApp(secrets []*up.Secret, name string) (v []*up.Secret) {
	for _, s := range secrets {
		if s.App == name {
			v = append(v, s)
		}
	}
	return
}

// Env returns a slice of env variables.
func Env(secrets []*up.Secret) (env []string) {
	for _, s := range secrets {
		env = append(env, fmt.Sprintf("%s=%s", s.Name, s.Value))
	}
	return
}

// String returns a string representation of the secret,
// using * as the wildcard when the secret is sensitive.
func String(s *up.Secret) string {
	switch s.Type {
	case "String":
		return s.Value
	case "SecureString":
		return strings.Repeat("*", len(s.Value))
	default:
		return ""
	}
}
