package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/viki-org/dnscache"
	"net"
	"net/http"
	"strings"
	"time"
)

func setup() {

}

func main() {

	//dns resolver from dns cache package
	resolver := dnscache.New(time.Minute * 5)

	//creating an http client to use it
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 64,
			Dial: func(network string, address string) (net.Conn, error) {
				fmt.Print("Using dns cache as resolver to look up ", address)
				separator := strings.LastIndex(address, ":")
				ip, _ := resolver.FetchOneString(address[:separator])
				return net.Dial("tcp", ip+address[separator:])
			},
		},
	}

	//create a config for aws service client
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region and httpclient that the service clients should use
	cfg.Region = endpoints.UsEast1RegionID
	cfg.HTTPClient = client

	// Using the Config value, create the s3 client
	svc := s3.New(cfg)

	//input parameter for an s3 ListBuckets request
	input := &s3.ListBucketsInput{}

	//create the request
	req := svc.ListBucketsRequest(input)

	//send the request
	result, err := req.Send()
	if err != nil {
		fmt.Println(err.Error())
	}

	//print the list of buckets in the account/region
	fmt.Println(result)

	return
}
