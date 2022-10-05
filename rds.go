package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func ExampleRDS_DescribeDBSnapshots_shared00() {
	svc := rds.New(session.New())
	input := &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String(""),
		IncludePublic:        aws.Bool(false),
		IncludeShared:        aws.Bool(true),
		SnapshotType:         aws.String("manual"),
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
		return
	}

	fmt.Println(result)
}

func Do() {
	ExampleRDS_DescribeDBSnapshots_shared00()
}
