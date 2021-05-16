package accessors

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dandersoncouncil/covid_help/datamodels/ddb"
)

type S3Accessor struct {
	SVC *s3.S3
}

func (this *S3Accessor) CreatePreSignedGetUrl(bucket, key string) string {
	req, _ := this.SVC.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
	}

	log.Println("The URL is", urlStr)
	return urlStr
}

func (this *S3Accessor) CreatePreSignedPutUrl(bucket, key string) string {

	req, _ := this.SVC.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	str, err := req.Presign(15 * time.Minute)

	log.Println("The URL is:", str, " err:", err)
	return str
}

func (this *S3Accessor) Downloadfile(bucketname string, path string) ([]*ddb.CovidResource, error) {

	filename := "/tmp/temp_file"

	// The session the S3 Downloader will use
	sess := session.Must(session.NewSession())

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(filename)
	if err != nil {
		return []*ddb.CovidResource{}, fmt.Errorf("failed to create file %q, %v", filename, err)
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(path),
	})
	if err != nil {
		return []*ddb.CovidResource{}, fmt.Errorf("failed to download file, %v", err)
	}
	log.Printf("file downloaded, %d bytes\n", n)

	defer f.Close()

	scanner := bufio.NewScanner(f)
	covidResource := []*ddb.CovidResource{}

	for scanner.Scan() {
		text := scanner.Text()
		log.Printf("Scanned line %v", text)
		lineSpilt := strings.Split(text, "\t")
		name := lineSpilt[0]
		contact := lineSpilt[1]
		locationA := lineSpilt[2]
		comment := lineSpilt[3]
		verified := lineSpilt[4]
		var isVerfified = false
		if verified == "true" {
			isVerfified = true
		}
		category := lineSpilt[5]
		id := fmt.Sprintf("%v", md5.Sum([]byte(name)))

		resource := ddb.CovidResource{ResourceId: id, Name: name,
			Category: category, AddrLine: locationA, PhoenNo: contact, IsVerfied: isVerfified, Remarks: comment}

		covidResource = append(covidResource, &resource)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return covidResource, nil
}
