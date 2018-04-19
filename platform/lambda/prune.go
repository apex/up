package lambda

import (
	"sort"
	"time"

	"github.com/apex/log"
	"github.com/apex/up/platform/event"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

// Prune implementation.
func (p *Platform) Prune(region, stage string, versions int) error {
	p.events.Emit("prune", nil)

	if err := p.createRole(); err != nil {
		return errors.Wrap(err, "creating iam role")
	}

	s := s3.New(session.New(aws.NewConfig().WithRegion(region)))
	b := aws.String(p.getS3BucketName(region))
	prefix := p.config.Name + "/" + stage + "/"

	params := &s3.ListObjectsInput{
		Bucket: b,
		Prefix: &prefix,
	}

	start := time.Now()
	var objects []*s3.Object
	var count int
	var size int64

	// fetch objects
	err := s.ListObjectsPages(params, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		for _, o := range page.Contents {
			objects = append(objects, o)
		}
		return *page.IsTruncated
	})

	if err != nil {
		return errors.Wrap(err, "listing s3 objects")
	}

	// sort by time descending
	sort.Slice(objects, func(i int, j int) bool {
		a := objects[i]
		b := objects[j]
		return (*b).LastModified.Before(*a.LastModified)
	})

	// remove old versions
	for i, o := range objects {
		ctx := log.WithFields(log.Fields{
			"index":         i,
			"key":           *o.Key,
			"size":          *o.Size,
			"last_modified": *o.LastModified,
		})

		if i < versions {
			ctx.Debug("retain")
			continue
		}

		ctx.Debug("remove")
		size += *o.Size
		count++

		_, err := s.DeleteObject(&s3.DeleteObjectInput{
			Bucket: b,
			Key:    o.Key,
		})

		if err != nil {
			return errors.Wrap(err, "removing object")
		}
	}

	p.events.Emit("prune.complete", event.Fields{
		"duration": time.Since(start),
		"size":     size,
		"count":    count,
	})

	return nil
}
