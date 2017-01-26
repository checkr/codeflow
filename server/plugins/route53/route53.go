package route53

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
)

type Route53 struct {
	events chan agent.Event
}

func init() {
	agent.RegisterPlugin("route53", func() agent.Plugin {
		return &Route53{}
	})
}

func (x *Route53) Description() string {
	return "Set Route53 DNS for Kubernetes services"
}

func (x *Route53) SampleConfig() string {
	return ` `
}

func (x *Route53) Start(e chan agent.Event) error {
	x.events = e
	log.Println("Started Route53")
	return nil
}

func (x *Route53) Stop() {
	log.Println("Stopping Route53")
}

func (x *Route53) Subscribe() []string {
	return []string{
		"plugins.LoadBalancer:info",
	}
}

func (x *Route53) Process(e agent.Event) error {

	switch e.Name {
	case "plugins.LoadBalancer:status":
		log.Printf("Process Route53 event: %s", e.Name)
		x.updateRoute53(e)
	}
	return nil
}

func (x *Route53) updateRoute53(e agent.Event) error {
	payload := e.Payload.(plugins.LoadBalancer)
	if payload.State == plugins.Complete {
		log.Printf("Route53 plugin received LoadBalancer success message for %s, %s.  Processing.", payload.Service.Name, payload.Name)

		// Create the client
		sess, err := session.NewSession()
		if err != nil {
			fmt.Println("failed to create session,", err)
			return err
		}
		client := route53.New(sess)
		// Look for this dns name
		params := &route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String("Z29MODHVVGZ870"), // Required
			//	MaxItems:              aws.String("PageMaxItems"),
			//	StartRecordIdentifier: aws.String("ResourceRecordSetIdentifier"),
			//	StartRecordName:       aws.String("DNSName"),
			//	StartRecordType:       aws.String("RRType"),
		}
		foundRecord := false
		pageNum := 0
		errList := client.ListResourceRecordSetsPages(params,
			func(page *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
				pageNum++
				for _, p := range page.ResourceRecordSets {
					for _, r := range p.ResourceRecords {
						fmt.Printf("RESOURCE_RECORD: %s", *r.Value)
						if *r.Value == "test" {
							foundRecord = true
							// break out of pagination
							return true
						}
					}
				}
				return false
			})
		if errList != nil {
			log.Printf("Error listing ResourceRecordSets for Route53: %s", errList)
			return errList
		}
	}

	return nil
}
