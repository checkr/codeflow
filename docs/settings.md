## Codeflow settings explained

### Redis
```
---
redis:
  username:
  password:
  server: "redis:6379"
  database: "0"
  pool: "30"
  process: "1"
```

### Plugins

Codeflow is a fully pluggable architecture.  ALL Codeflow components are a plugin.  They can be configured/enabled as shown below.  To disable a plugin you can remove it's section from the config file.

```
plugins:
```

#### Webhooks Plugin

The webhooks plugin is what listens for and consumes Github webhook postback for all projects.

```
  webhooks:
    workers: 1
    service_address: ":3002"
    github:
      path: "/github"
```

#### Codeflow Plugin (API and Dashboard specific settings)

The Codeflow plugin encompasses the API and also acts as a workflow engine.

```
  codeflow:
    workers: 1
    # Dashboard URL is the URL used for accessing the codeflow dashboard.
    dashboard_url: "http://localhost:3000"

    # Logs URL is for linking to Kibana logs from a projects page.
    logs_url: "https://example.com/##PROJECT-NAMESPACE##"

    # JWT Secret Key is the key used to validate with Okta.
    jwt_secret_key: somevalue

    # Allowed Origins is the allowed list of origins for Okta login.
    allowed_origins:
      - "http://localhost:3000"
      - "http://localhost:3001/oauth2/callback/okta"	

    # Default service spec is the default resource settings for all kubernetes deployed projects.
    default_service_spec:
      cpu: "500m"
      cpu_burst: "1000m"
      memory: "512Mi"
      memory_burst: "1Gi"
      termination_grace_period_seconds: 600

    # MongoDB Connection settings
    mongodb:
      database: "codeflow"
      uri: "mongodb://mongo:27017"
      ssl: false

    #  Api service interface/port
    service_address: ":3001"

    # Codeflow routes
    builds:
      path: "/builds"
    projects:
      path: "/projects"
    # Authentication can be either "okta" (SSO), or "demo" (no auth)
    auth:
      path: "/auth"
      # handler: "okta"
      handler: "demo"
      # This is important to set for Okta SSO
      okta_org: "myorg"
    users:
      path: "/users"
    features:
      path: "/features"
    websockets:
      path: "/ws"
    bookmarks:
      path: "/bookmarks"
    stats:
      path: "/stats"
```

#### Docker Build Plugin

The Docker Build plugin will automatically build and push Docker Images to your docker repository for every Feature/PR that lands on Master.

```
  docker_build:
    workers: 1
    registry_host: "docker.io"
    registry_username: ""
    registry_password: ""
    registry_user_email: ""
    build_path: "/tmp"
    docker_host: "unix:///var/run/docker.sock"
```

#### Kubernetes Deployment Plugin

The Kubernetes Deployment plugin handles scheduling deployments and creating load balancers on a Kubernetes cluster.

```
  kubedeploy:
    # Optional: Environment overrides the default setting which is used to create kubernetes Namespaces
    environment: staging 
    # Optional: SSL Cert ARN can be used to create Load Balancers of type HTTPS/SSL using Kubernetes annotations.
    ssl_cert_arn: arn:aws:acm:us-east-1:xxx:certificate/xxxx-xxxx-xxxx-xxxx-xxxx
    # Optional: Node selector can be used for targeting specific types of nodes for deployment
    node_selector: "kubernetes.io/hostname=ip-10-10-10-10.ec2.internal"

```

#### Websockets Plugin

The Websockets plugin provides websocket service for the Dashboard single page app.

```
  websockets:
    workers: 1
    service_address: ":3003"
```

#### Slack Plugin

The Slack Plugin notifies via slack when deployments are complete.

```
  slack:
    workers: 1
    webhook_url: ""
```