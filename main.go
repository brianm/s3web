package main

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var REPLACE string = "REPLACE"
var CC_VAL string = "public, max-age=600"
var CONTENT_TYPE_HTML string = "text/html"
var CONTENT_TYPE_CSS string = "text/css"
var CONTENT_TYPE_PNG string = "image/png"
var CONTENT_TYPE_JPG string = "image/jpg"
var CONTENT_TYPE_GIF string = "image/gif"
var CONTENT_TYPE_JS string = "application/javascript"

func main() {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	bucket := "skife.org"
	loo, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: &bucket,
	})
	if err != nil {
		log.Fatalf("unable to list objects: %+v", err)
	}

	for _, c := range loo.Contents {
		//fmt.Println(*c.Key)
		//if 1+1 == 2{
		//	continue
		//}
		//log.Printf("examining %s", *c.Key)
		src := strings.Join([]string{bucket, *c.Key}, "/")
		if strings.HasSuffix(*c.Key, ".html") {
			_, err := svc.CopyObject(&s3.CopyObjectInput{
				ACL:               aws.String("public-read"),
				Bucket:            aws.String(bucket),
				Key:               c.Key,
				CopySource:        aws.String(src),
				MetadataDirective: aws.String("REPLACE"),
				CacheControl:      aws.String("public, max-age=600"),
				ContentType:       aws.String("text/html"),
			})
			if err != nil {
				log.Fatalf("unable to copy %s: %+v", *c.Key, err)
			}
			log.Printf("fixed %s", *c.Key)
		} else if strings.HasSuffix(*c.Key, ".css") {
			_, err := svc.CopyObject(&s3.CopyObjectInput{
				ACL:               aws.String("public-read"),
				Bucket:            &bucket,
				Key:               c.Key,
				CopySource:        &src,
				MetadataDirective: &REPLACE,
				CacheControl:      &CC_VAL,
				ContentType:       &CONTENT_TYPE_CSS,
			})
			if err != nil {
				log.Fatalf("unable to copy %s: %+v", *c.Key, err)
			}
			log.Printf("fixed %s", *c.Key)
		} else if strings.HasSuffix(*c.Key, ".png") {
			_, err := svc.CopyObject(&s3.CopyObjectInput{
				ACL:               aws.String("public-read"),
				Bucket:            &bucket,
				Key:               c.Key,
				CopySource:        &src,
				MetadataDirective: &REPLACE,
				CacheControl:      &CC_VAL,
				ContentType:       &CONTENT_TYPE_PNG,
			})
			if err != nil {
				log.Fatalf("unable to copy %s: %+v", *c.Key, err)
			}
			log.Printf("fixed %s", *c.Key)
		} else if strings.HasSuffix(*c.Key, ".jpg") {
			_, err := svc.CopyObject(&s3.CopyObjectInput{
				ACL:               aws.String("public-read"),
				Bucket:            &bucket,
				Key:               c.Key,
				CopySource:        &src,
				MetadataDirective: &REPLACE,
				CacheControl:      &CC_VAL,
				ContentType:       &CONTENT_TYPE_JPG,
			})
			if err != nil {
				log.Fatalf("unable to copy %s: %+v", *c.Key, err)
			}
			log.Printf("fixed %s", *c.Key)
		} else if strings.HasSuffix(*c.Key, ".js") {
			_, err := svc.CopyObject(&s3.CopyObjectInput{
				ACL:               aws.String("public-read"),
				Bucket:            &bucket,
				Key:               c.Key,
				CopySource:        &src,
				MetadataDirective: &REPLACE,
				CacheControl:      &CC_VAL,
				ContentType:       &CONTENT_TYPE_JS,
			})
			if err != nil {
				log.Fatalf("unable to copy %s: %+v", *c.Key, err)
			}
			log.Printf("fixed %s", *c.Key)
		} else if strings.HasSuffix(*c.Key, ".gif") {
			_, err := svc.CopyObject(&s3.CopyObjectInput{
				ACL:               aws.String("public-read"),
				Bucket:            &bucket,
				Key:               c.Key,
				CopySource:        &src,
				MetadataDirective: &REPLACE,
				CacheControl:      &CC_VAL,
				ContentType:       &CONTENT_TYPE_GIF,
			})
			if err != nil {
				log.Fatalf("unable to copy %s: %+v", *c.Key, err)
			}
			log.Printf("fixed %s", *c.Key)
		} else if strings.HasSuffix(*c.Key, ".atom") {
			_, err := svc.CopyObject(&s3.CopyObjectInput{
				ACL:               aws.String("public-read"),
				Bucket:            &bucket,
				Key:               c.Key,
				CopySource:        &src,
				MetadataDirective: &REPLACE,
				CacheControl:      &CC_VAL,
				ContentType:       aws.String("application/atom+xml"),
			})
			if err != nil {
				log.Fatalf("unable to copy %s: %+v", *c.Key, err)
			}
			log.Printf("fixed %s", *c.Key)
		}else {
			log.Printf("unknown type for %s", *c.Key)
		}
	}
}
