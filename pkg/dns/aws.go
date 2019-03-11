package dns

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

// ChangeAWSRecord - Updates record in Route53
func ChangeAWSRecord(Action string, fqdn string, HostedZoneID string, IPAddress string, TTL int64, AccessKeyID string, SecretAccessKey string) (err error) {
	awsConfig, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Printf("unable to load SDK config, %v" + err.Error())
		log.Printf("Will attempt to use the config file")
	}

	// Todo Get creds from file
	/*
		var credsfromjson aws.Credentials
		credsfromjson.AccessKeyID = AccessKeyID
		credsfromjson.SecretAccessKey = SecretAccessKey
		awsConfig.Credentials = credsfromjson
	*/

	r53 := route53.New(awsConfig)

	recordlist := append([]route53.ResourceRecord{}, route53.ResourceRecord{
		Value: &IPAddress,
	})

	var dnsAction route53.ChangeAction
	if Action == "CREATE" {
		dnsAction = route53.ChangeAction("CREATE")
	} else if Action == "DELETE" {
		dnsAction = route53.ChangeAction("DELETE")
	} else {
		return fmt.Errorf("Undefined DNS action")
	}

	dnschangelist := append([]route53.Change{}, route53.Change{
		Action: dnsAction,
		ResourceRecordSet: &route53.ResourceRecordSet{
			Name:            &fqdn,
			Type:            "CNAME",
			TTL:             &TTL,
			ResourceRecords: recordlist,
		},
	})

	comment := "Comment goes here"

	dnschangeBatch := &route53.ChangeBatch{
		Changes: dnschangelist,
		Comment: &comment,
	}

	log.Printf("Change: %s", dnschangeBatch)

	request := r53.ChangeResourceRecordSetsRequest(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &HostedZoneID,
		ChangeBatch:  dnschangeBatch,
	})

	resp, err := request.Send()

	if err != nil {
		return err
	}

	log.Printf("%s", resp)

	return nil
}
