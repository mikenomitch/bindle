hook "deploy-prod" {
  cta = "Deploy to Nomad Prod Cluster"
  type = "deploy"
  address = "https://100.36.128.78:4646"
}

hook "deploy-test" {
  cta = "Deploy to Nomad Test Cluster"
  type = "deploy"
  address = "https://100.36.128.78:4646"
}

hook "post-to-custom-flow" {
  cta = "Send to Custom CI/CD Flow"
  type = "webhook"
  address = "https://my-custom-webhook.com/123"
}
