package route53

import (
	"fmt"
	"log"
	"strings"
	"time"

	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
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
		DNS:          lbPayload.DNS,
		Subdomain:    lbPayload.Subdomain,
		FQDN:         viper.GetString("plugins.route53.hosted_zone_name"),
		Environment:  lbPayload.Environment,
		Project:      lbPayload.Project,
	}, nil)
	x.events <- event
}

func (x *Route53) updateRoute53(e agent.Event) error {
	payload := e.Payload.(plugins.LoadBalancer)
	// Sanity checks
	if payload.DNS == "" {
		failMessage := fmt.Sprintf("DNS was blank for %s, skipping Route53.", payload.Project.Slug)
		x.sendRoute53Response(e, plugins.Failed, failMessage, payload)
		return nil
	}
	if payload.Subdomain == "" {
		failMessage := fmt.Sprintf("Subdomain was blank for %s, skipping Route53.", payload.Project.Slug)
		x.sendRoute53Response(e, plugins.Failed, failMessage, payload)
		return nil
	}
	if payload.Type == plugins.Internal {
		fmt.Printf("Internal service type ignored for %s", payload.DNS)
		return nil
	}
	route53Name := fmt.Sprintf("%s.%s", payload.Subdomain, viper.GetString("plugins.route53.hosted_zone_name"))
	if payload.State == plugins.Complete {
		log.Printf("Route53 plugin received LoadBalancer success message for %s, %s, %s.  Processing.\n", payload.Project.Slug, payload.DNS, payload.Name)

		// Wait for DNS from the ELB to settle, abort if it does not resolve in initial_wait
		// Trying to be conservative with these since we don't want to update Route53 before the new ELB dns record is available
		time.Sleep(time.Second * viper.GetDuration("plugins.route53.initial_wait"))

		// Query the DNS until it resolves or timeouts
		dnsTimeout := viper.GetInt("plugins.route53.dns_resolve_timeout_seconds")
		dnsValid := false
		var failMessage string
		var dnsLookup []string
		var dnsLookupErr error
		for dnsValid == false {
			dnsLookup, dnsLookupErr = net.LookupHost(payload.DNS)
			dnsTimeout -= 10
			if dnsLookupErr != nil {
				failMessage = fmt.Sprintf("Error '%s' resolving DNS for: %s", dnsLookupErr, payload.DNS)
			} else if len(dnsLookup) == 0 {
				failMessage = fmt.Sprintf("Error 'found no names associated with ELB record' while resolving DNS for: %s", payload.DNS)
			} else {
				dnsValid = true
			}
			if dnsTimeout <= 0 || dnsValid {
				break
			}
			time.Sleep(time.Second * 10)
			fmt.Println(failMessage + ".. Retrying in 10s")
		}
		if dnsValid == false {
			x.sendRoute53Response(e, plugins.Failed, failMessage, payload)
			return nil
		}
		fmt.Printf("DNS for %s resolved to: %s\n", payload.DNS, strings.Join(dnsLookup, ","))

		// Create the client
		sess := awssession.Must(awssession.NewSessionWithOptions(
			awssession.Options{
				Config: aws.Config{
					Credentials: credentials.NewStaticCredentials(viper.GetString("plugins.route53.aws_access_key_id"), viper.GetString("plugins.route53.aws_secret_key"), ""),
				},
			},
		))

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
			log.Printf("Route53 found existing record for: %s\n", route53Name)
		} else {
			log.Printf("Route53 record not found, creating %s\n", route53Name)
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
									Value: aws.String(payload.DNS),
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
		log.Printf("Route53 record UPSERTed for %s: %s", route53Name, payload.DNS)
		x.sendRoute53Response(e, plugins.Complete, "", payload)
	}

	return nil
}
