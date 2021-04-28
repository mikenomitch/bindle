variable "job-name" {
  display-name = "Job Name"
  description = "The job name is a unique identifier for your Nomad Job"
  type        = string
  default     = "simple-docker-webservice"
}

variable "datacenters" {
  display-name = "Datacenters"
  description = "The Nomad datacenters this job will be deployed on."
  type        = list(string)
  default     = ["dc1"]

  meta {
    options = ["dc1", "dc2", "dc3"]
  }
}

variable "image" {
  display-name = "Image"
  description = "The Docker image used for the service. Pulled from Docker Hub by default."
  type        = string
  default     = "hashicorp/hello-world-webserver"
}

variable "port" {
  display-name = "HTTP Port"
  description = "The static port your server will be accessed on."
  type        = number
  default     = 4000

  meta = {
    min = 3000
    max = 9000
  }
}



