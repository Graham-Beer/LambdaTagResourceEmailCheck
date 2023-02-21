package resources

import (
	"context"
	"log"

	"net/mail"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	tag "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	// sestypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

/*
`TagMatch` is a custom type to display the AWS resource information.
*/

type TagMatch struct {
	AccountId string
	Service   string
	Resource  string
	ARN       string
}

/*
`GetTagResources` function takes the *`resourcegroupstaggingapi` package generated client and a
unlimited slice of tag values (using a vardic function).

The helper function `GroupValues` is called. See below for function details.

* see https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi@v1.14.3
*/

func GetTagResources(client *tag.Client, key ...string) (map[string][]TagMatch, error) {
	var group = make(map[string][]TagMatch)
	for _, k := range key {
		res, err := client.GetResources(context.TODO(), &tag.GetResourcesInput{
			IncludeComplianceDetails: aws.Bool(true),
			TagFilters: []types.TagFilter{
				{Key: aws.String(k)},
			},
		})
		if err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		out := GroupValues(res, k)
		for k, v := range out {
			group[k] = v
		}
	}
	return group, nil
}

/*
GroupValues is a helper function.
Function is used inside the GetTagResources function.

GroupValues takes the output type of *`GetResourcesOutput` from the AWS SDK Go v2 along
with a tag key.
The return is of the custom type map[string][]TagMatch.

* see https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi@v1.14.3#GetResourcesOutput
*/

func GroupValues(res *tag.GetResourcesOutput, key string) map[string][]TagMatch {
	var group = make(map[string][]TagMatch)
	for _, res := range res.ResourceTagMappingList {
		for _, tag := range res.Tags {
			if TagValue := *tag.Key == *aws.String(key); TagValue && *tag.Value != "" {
				par, err := arn.Parse(*res.ResourceARN)
				if err != nil {
					log.Fatal(err)
				}

				details := TagMatch{
					AccountId: par.AccountID,
					Service:   par.Service,
					Resource:  par.Resource,
					ARN:       *res.ResourceARN,
				}

				group[*tag.Value] = append(group[*tag.Value], details)
			}
		}
	}
	return group
}

/*
TODO:
`SendSns` will send the data collected.
*/

func SendSns(ses *sesv2.Client, user map[string][]TagMatch) {
	keys := make([]string, 0, len(user))
	for k := range user {
		keys = append(keys, k)
	}

	for _, address := range keys {
		// confirm email address
		if valid(address) {
			log.Printf("\nUser: %v\n", address)
			for _, res := range user[address] {
				log.Printf("AccountId: %v\nAWS Service: %v\nAWS Resource: %v\nARN: %v\n\n",
					res.AccountId, res.Service, res.Resource, res.ARN)
			}
			continue
		}
		log.Printf("[%s] is not a valid email address\n", address)
	}
}

/*
Helper function to validate an email address.
*/

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ses.SendEmail(context.TODO(), &sesv2.SendEmailInput{
// 	Content:                        &sestypes.EmailContent{
// 		Raw:      &sestypes.RawMessage{},
// 		Simple:   &sestypes.Message{},
// 		Template: &sestypes.Template{},
// 	},
// 	ConfigurationSetName:           new(string),
// 	Destination:                    &sestypes.Destination{
// 		BccAddresses: []string{},
// 		CcAddresses:  []string{},
// 		ToAddresses:  []string{},
// 	},
// 	EmailTags:                      []sestypes.MessageTag{},
// 	FeedbackForwardingEmailAddress: new(string),
// 	FeedbackForwardingEmailAddressIdentityArn: new(string),
// 	FromEmailAddress:            new(string),
// 	FromEmailAddressIdentityArn: new(string),
// 	ListManagementOptions:       &sestypes.ListManagementOptions{},
// 	ReplyToAddresses:            []string{},
// })
