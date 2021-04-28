name = "Simple Docker Application"
description = "Deploy a simple docker webapp to Nomad"

nomad_job "webservice-api" {
  description = "A nomad job that runs a docker image and exposes it via an https port"
  template_file = "service.hcl"
}









