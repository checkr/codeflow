## Preparing an application for deployment
### Requirements
* Git repository
* __Access to the github repository admin settings__.
    * Navigate to your project -> Settings.  If you cannot see the settings page you do not have access.  Ask your team lead for access.

### Adding to codeflow
* Click __Projects -> Add New__
* Enter the full path to the git repository eg. git@github.com:checkr/codeflow
* On Github, navigate to your project -> Settings -> Webhooks and copy the settings from Codeflow
  * Github webhook, Github secret.  PayloadURL: application/json, Send me everything.
  ![Webhook settings](/images/setup_webhook.png "Webhook settings")


* On Github, navigate to your project -> Settings -> Deploy Keys and create a new deploy key.
  * Key name: codeflow-prod, copy the public key.
  ![Create deploy key](/images/deploy_key.png "Create deploy key")

* All set, now you just have to push a commit and refresh the page.

### Adding Resources
Resources represent command(s) that you run as a daemon inside of a Docker container.  You can have as many resources as you want in your project.
* For each resource:
  * Name: The friendly name for this resource (ie. www if it's the main web service).
  * Command: The command to run in the docker container that will start this service.
  * Container Listeners: This is the default internal port (if any) that your service listens on.  If this is a worker you do not need to specify a port.
![Resources](/images/cf_resources.png "Create a Resource")

### Load Balancers
Load Balancers are what allows your new resource to be discovered by other services or exposed to the internet.

* Service: Select the resource that one you want to expose with a load balancer from the drop-down.
* Access:
  * __Internal__: Internal only connectivity and discovery.  This does not use an ELB (low cost) but provides basic load balancing and service discovery for your Resource.
  * __Office__: Exposes the service to the Office network.  This means services deployed into a private subnet in your VPC. 
  * __External__: Exposes the service to the public internet!  Careful with this as your service will be accessible by the outside world.
* Port: The external port to serve the traffic on.  Eg, for a normal SSL enabled http server this would be 443.
* Service Port: The internal port that your container is listening on (that was specified earlier in the Resource creation).
* Service Protocol:
  * HTTPS: Tell the load balancer to terminate https using the SSL certificate specified in the [settings](settings.md) kubedeploy:ssl_cert_arn wildcard domain.
  * HTTP:  Plain HTTP (no SSL). 
  * SSL:   Terminate SSL with TCP using the certificate from [settings](settings.md) kubedeploy:ssl_cert_arn wildcard domain  (Example: for using with websocket protocol wss://).
  * TCP:   Plain TCP (no SSL).
  * UDP:   Use UDP.  Only available for Internal service type.

![Load Balancers](/images/cf_loadbalancer.png "Create a LoadBalancer")