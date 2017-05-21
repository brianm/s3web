package fix

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"

	"log"
	"strings"
	"io/ioutil"
	"net/http"
)

type Fix struct {
	Bucket                         string
	Verbose                        bool
	Simulate                       bool
	svc                            *s3.S3
	DetectMimeTypesForUnknownNames bool
}

func (f Fix) CacheControlForType(contentType string) string {
	return "public, max-age=600"
}

func (f Fix) fixicate(obj *s3.Object, contentType string) error {
	hoo, err := f.svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(f.Bucket),
		Key: obj.Key,
	})
	if err != nil {
		return errors.Wrapf(err, "unable to fetch metadata for %s", *obj.Key)
	}

	update := false

	cacheControl := f.CacheControlForType(contentType)

	update = update || *hoo.CacheControl != cacheControl
	update = update || *hoo.ContentType != contentType

	if update {
		_, err = f.svc.CopyObject(&s3.CopyObjectInput{
			ACL:               aws.String("public-read"),
			Bucket:            aws.String(f.Bucket),
			Key:               obj.Key,
			CopySource:        aws.String(strings.Join([]string{f.Bucket, *obj.Key}, "/")),
			MetadataDirective: aws.String("REPLACE"),
			CacheControl:      aws.String(cacheControl),
			ContentType:       aws.String(contentType),
		})
		if err != nil {
			return errors.Wrapf(err, "unable to copy %s", *obj.Key)
		}
		log.Printf("FIXED\t%s", *obj.Key)
	} else {
		if f.Verbose {
			log.Printf("NOOP\t%s", *obj.Key)
		}
	}
	return nil
}

func (f Fix) Fix() error {
	if f.Simulate {
		return nil
	}

	if f.svc == nil {
		sess := session.Must(session.NewSession())
		f.svc = s3.New(sess)
	}
	loo, err := f.svc.ListObjects(&s3.ListObjectsInput{
		Bucket: &f.Bucket,
	})

	if err != nil {
		return errors.Wrapf(err, "unable to fetch objects from  %s", f.Bucket)
	}

	for _, c := range loo.Contents {
		if strings.HasSuffix(*c.Key, ".html") {
			if err := f.fixicate(c, "text/html"); err != nil {
				return errors.Wrapf(err, "unable fix copy %s", *c.Key)
			}
		} else if strings.HasSuffix(*c.Key, ".css") {
			if err := f.fixicate(c, "text/css"); err != nil {
				return errors.Wrapf(err, "unable fix copy %s", *c.Key)
			}
		} else if strings.HasSuffix(*c.Key, ".png") {
			if err := f.fixicate(c, "image/png"); err != nil {
				return errors.Wrapf(err, "unable fix copy %s", *c.Key)
			}
		} else if strings.HasSuffix(*c.Key, ".jpg") {
			if err := f.fixicate(c, "image/jpg"); err != nil {
				return errors.Wrapf(err, "unable fix copy %s", *c.Key)
			}
		} else if strings.HasSuffix(*c.Key, ".js") {
			if err := f.fixicate(c, "application/javascript"); err != nil {
				return errors.Wrapf(err, "unable to fix %s", *c.Key)
			}
		} else if strings.HasSuffix(*c.Key, ".gif") {
			if err := f.fixicate(c, "image/gif"); err != nil {
				return errors.Wrapf(err, "unable to fix %s", *c.Key)
			}
		} else if strings.HasSuffix(*c.Key, ".atom") {
			if err := f.fixicate(c, "application/atom+xml"); err != nil {
				return errors.Wrapf(err, "unable to fix %s", *c.Key)
			}
		} else if strings.HasSuffix(*c.Key, ".txt") {
			if err := f.fixicate(c, "text/plain; charset=utf-8"); err != nil {
				return errors.Wrapf(err, "unable to fix %s", *c.Key)
			}
		} else if strings.HasSuffix(*c.Key, ".ico") {
			if err := f.fixicate(c, "image/x-icon"); err != nil {
				return errors.Wrapf(err, "unable to fix %s", *c.Key)
			}
		} else if f.DetectMimeTypesForUnknownNames {
			obj, err := f.svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(f.Bucket),
				Key: c.Key,
				Range: aws.String("bytes=0-512"),
			})
			if err != nil {
				return errors.Wrapf(err, "unable to fetch first 512 bytes of %s", *c.Key)
			}
			buf, err := ioutil.ReadAll(obj.Body)
			if err != nil {
				obj.Body.Close()
				return errors.Wrapf(err, "unable to read bytes from fetch of %s", *c.Key)
			}
			contentType := http.DetectContentType(buf)
			obj.Body.Close()

			if err := f.fixicate(c, contentType); err != nil {
				return errors.Wrapf(err, "unable to fix %s", *c.Key)
			}
		}
	}
	return nil
}

