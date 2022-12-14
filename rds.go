package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/jessevdk/go-flags"
)

type options struct {
	Region  string `short:"r" long:"region" description:"" required:"false" default:"ap-northeast-1"`
	Profile string `short:"p" long:"profile" description:"" required:"false"`
	DB_NAME string `long:"db-name" description:"" required:"false"`
}

var opts options

func checkContainsList(c string, s []*rds.DBSnapshot) bool {
	for _, s := range s {
		if strings.Index(*s.DBSnapshotIdentifier, c) == 0 {
			return true
		}
	}
	return false
}

func getFormatDate() string {
	t := time.Now().Add(-24 * time.Hour)
	const layout2 = "2006-01-02"
	d := fmt.Sprintf(t.Format(layout2))
	return fmt.Sprintf("rds:%s-%s-", opts.DB_NAME, d)
}

func getSnapShotList() bool {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           opts.Profile,
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
	}))
	svc := rds.New(sess)
	input := &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String(opts.DB_NAME),
		IncludePublic:        aws.Bool(false),
		IncludeShared:        aws.Bool(true),
		SnapshotType:         aws.String("automated"),
	}

	result, err := svc.DescribeDBSnapshots(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBSnapshotNotFoundFault:
				fmt.Println(rds.ErrCodeDBSnapshotNotFoundFault, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	return checkContainsList(getFormatDate(), result.DBSnapshots)
}

func Do() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !getSnapShotList() {
		fmt.Println("fail snapshot")
	}
	fmt.Println("success")
}
