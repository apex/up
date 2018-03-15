package lambda

import (
	"sort"
	"strconv"

	"github.com/apex/log"
	"github.com/apex/up/internal/colors"
	"github.com/apex/up/internal/table"
	"github.com/apex/up/internal/util"
	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
)

// ShowDeploys implementation.
func (p *Platform) ShowDeploys(region string) error {
	s := session.New(aws.NewConfig().WithRegion(region))
	c := lambda.New(s)

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
		{Text: colors.Bold("Date")},
	})

	t.AddRow(table.Row{
		{
			Span: 4,
		},
	})

	for _, f := range versions {
		if *f.Version != "$LATEST" {
			addFunction(t, f)
		}
	}

	defer util.Pad()()
	t.Println()

	return nil
}

// addFunction adds function to table.
func addFunction(t *table.Table, f *lambda.FunctionConfiguration) {
	commit := f.Environment.Variables["UP_COMMIT"]
	author := f.Environment.Variables["UP_AUTHOR"]
	stage := *f.Environment.Variables["UP_STAGE"]
	created := dateparse.MustParse(*f.LastModified)
	date := util.RelativeDate(created)
	version := *f.Version

	t.AddRow(table.Row{
		{Text: formatStage(stage)},
		{Text: colors.Gray(util.DefaultString(commit, version))},
		{Text: colors.Gray(util.DefaultString(author, "–"))},
		{Text: date},
	})
}

// formatStage returns the stage string format.
func formatStage(s string) string {
	switch s {
	case "production":
		return colors.Purple(s)
	default:
		return colors.Gray(s)
	}
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

		if err != nil {
			return nil, err
		}

		log.Debugf("fetched %d versions", len(res.Versions))
		versions = append(versions, res.Versions...)

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
