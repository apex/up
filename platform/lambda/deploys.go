package lambda

import (
	"sort"
	"strconv"

	"github.com/apex/log"
	"github.com/apex/up"
	"github.com/apex/up/internal/colors"
	"github.com/apex/up/internal/table"
	"github.com/apex/up/internal/util"
	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	humanize "github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

// ShowDeploys implementation.
func (p *Platform) ShowDeploys(region string) error {
	s := session.New(aws.NewConfig().WithRegion(region))
	c := lambda.New(s)

	stages, err := getCurrentVersions(c, p.config)
	if err != nil {
		return errors.Wrap(err, "fetching current versions")
	}

	versions, err := getVersions(c, p.config.Name)
	if err != nil {
		return errors.Wrap(err, "fetching versions")
	}

	versions = filterLatest(versions)
	sortVersionsDesc(versions)
	versions = filterN(versions, 25)
	t := table.New()

	t.AddRow(table.Row{
		{Text: colors.Bold("Stage")},
		{Text: colors.Bold("Version")},
		{Text: colors.Bold("Author")},
		{Text: colors.Bold("Size")},
		{Text: colors.Bold("Date")},
	})

	t.AddRow(table.Row{
		{
			Span: 5,
		},
	})

	for _, f := range versions {
		if *f.Version != "$LATEST" {
			addDeployment(t, f, stages)
		}
	}

	defer util.Pad()()
	t.Println()

	return nil
}

// addDeployment adds the release to table.
func addDeployment(t *table.Table, f *lambda.FunctionConfiguration, stages map[string]string) {
	commit := f.Environment.Variables["UP_COMMIT"]
	author := f.Environment.Variables["UP_AUTHOR"]
	stage := *f.Environment.Variables["UP_STAGE"]
	created := dateparse.MustParse(*f.LastModified)
	date := util.RelativeDate(created)
	version := *f.Version
	current := stages[stage] == version

	t.AddRow(table.Row{
		{Text: formatStage(stage, current)},
		{Text: colors.Gray(util.DefaultString(commit, version))},
		{Text: colors.Gray(util.DefaultString(author, "–"))},
		{Text: humanize.Bytes(uint64(*f.CodeSize))},
		{Text: date},
	})
}

// formatStage returns the stage string format.
func formatStage(s string, current bool) string {
	var c colors.Func

	switch s {
	case "production":
		c = colors.Purple
	default:
		c = colors.Gray
	}

	s = c(s)

	if current {
		s = c("*") + " " + s
	} else {
		s = colors.Gray("\u0020") + " " + s
	}

	return s
}

// getCurrentVersions returns the current stage versions.
func getCurrentVersions(c *lambda.Lambda, config *up.Config) (map[string]string, error) {
	m := make(map[string]string)

	for _, s := range config.Stages.List() {
		if s.IsLocal() {
			continue
		}

		res, err := c.GetAlias(&lambda.GetAliasInput{
			FunctionName: &config.Name,
			Name:         aws.String(s.Name),
		})

		if util.IsNotFound(err) {
			continue
		}

		if err != nil {
			return nil, errors.Wrapf(err, "fetching %s alias", s.Name)
		}

		m[s.Name] = *res.FunctionVersion
	}

	return m, nil
}

// getVersions returns all function versions.
func getVersions(c *lambda.Lambda, name string) (versions []*lambda.FunctionConfiguration, err error) {
	var marker *string

	log.Debug("fetching versions")
	for {
		res, err := c.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
			FunctionName: &name,
			MaxItems:     aws.Int64(10000),
			Marker:       marker,
		})

		if util.IsNotFound(err) {
			goto skip
		}

		if err != nil {
			return nil, err
		}

		log.Debugf("fetched %d versions", len(res.Versions))
		versions = append(versions, res.Versions...)

	skip:
		marker = res.NextMarker
		if marker == nil {
			break
		}
	}
	log.Debug("fetched versions")

	return
}

// filterN returns a slice of the first n versions.
func filterN(in []*lambda.FunctionConfiguration, n int) (out []*lambda.FunctionConfiguration) {
	for i, v := range in {
		if i-1 == n {
			break
		}
		out = append(out, v)
	}
	return
}

// filterLatest returns a slice without $LATEST.
func filterLatest(in []*lambda.FunctionConfiguration) (out []*lambda.FunctionConfiguration) {
	for _, v := range in {
		if *v.Version != "$LATEST" {
			out = append(out, v)
		}
	}
	return
}

// sortVersionsDesc sorts versions descending.
func sortVersionsDesc(versions []*lambda.FunctionConfiguration) {
	sort.Slice(versions, func(i int, j int) bool {
		a := mustParseInt(*versions[i].Version)
		b := mustParseInt(*versions[j].Version)

		return a > b
	})
}

// mustParseInt returns an int from string.
func mustParseInt(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(errors.Wrapf(err, "parsing integer string %v", s))
	}
	return n
}
