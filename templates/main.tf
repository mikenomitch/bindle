terraform {
  required_providers {
    nomad = {
      source = "hashicorp/nomad"
      version = "1.4.13"
    }
  }
}

provider "nomad" {
  address = "{{.Address}}"
}