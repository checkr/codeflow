package route53

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/spf13/viper"
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
		"plugins.LoadBalancer",
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

func (x *Route53) sendRoute53Response(e agent.Event, state plugins.State, failureMessage string, lbPayload plugins.LoadBalancer) {
	event := e.NewEvent(plugins.Route53{
		State:        state,
		StateMessage: failureMessage,
		Service:      lbPayload.Service,
		DNSName:      lbPayload.DNSName,
		Route53DNS:   lbPayload.Route53DNS,
	}, nil)
	x.events <- event
}

func (x *Route53) updateRoute53(e agent.Event) error {
	payload := e.Payload.(plugins.LoadBalancer)
	names := strings.Split(payload.Project.Repository, "/")
	route53Name := fmt.Sprintf("%s.%s", names[1], viper.GetString("plugins.route53.hosted_zone_name"))
	if payload.State == plugins.Complete {
		log.Printf("Route53 plugin received LoadBalancer success message for %s, %s.  Processing.", payload.Service.Name, payload.Name)

		// Create the client
		sess := session.New()
		client := route53.New(sess)
		// Look for this dns name
		params := &route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(viper.GetString("plugins.route53.hosted_zone_id")), // Required
		}
		foundRecord := false
		pageNum := 0
		// Route53 has a . on the end of the name
		lookFor := fmt.Sprintf("%s.", route53Name)
		errList := client.ListResourceRecordSetsPages(params,
			func(page *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
				pageNum++
				for _, p := range page.ResourceRecordSets {
					if *p.Name == lookFor {
						foundRecord = true
						// break out of pagination
						return true
					}
				}
				return false
			})
		if errList != nil {
			log.Printf("Error listing ResourceRecordSets for Route53: %s", errList)
			return errList
		}
		if foundRecord {
			log.Printf("Route53 found existing record for: %s", route53Name)
		} else {
			log.Printf("Route53 record not found, creating %s", route53Name)
		}
		updateParams := &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: aws.String(viper.GetString("plugins.route53.hosted_zone_id")),
			ChangeBatch: &route53.ChangeBatch{
				Changes: []*route53.Change{
					{
						Action: aws.String("UPSERT"),
						ResourceRecordSet: &route53.ResourceRecordSet{

							Name: aws.String(route53Name),
							Type: aws.String("CNAME"),
							ResourceRecords: []*route53.ResourceRecord{
								{
									Value: aws.String(payload.DNSName),
								},
							},
							TTL: aws.Int64(60),
						},
					},
				},
			},
		}
		_, err := client.ChangeResourceRecordSets(updateParams)
		if err != nil {
			failMessage := fmt.Sprintf("ERROR '%s' setting Route53 DNS for %s", err, route53Name)
			x.sendRoute53Response(e, plugins.Failed, failMessage, payload)
			return nil
		}
		x.sendRoute53Response(e, plugins.Complete, "", payload)
	}

	return nil
}
