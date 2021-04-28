job "[[.job-name]]" {
  datacenters = [[.datacenters]]
  type = "service"

  group "webservice-group" {
    count = 1

    network {
      mode = "bridge"

      port "https" {
        static = [[.port]]
        to     = [[.port]]
      }
    }

    task "webservice" {
      driver = "docker"

      config {
        image = [[.image]]
        ports = ["https"]
      }

      env {
        PORT = "${NOMAD_PORT_http}"
      }

      resources {
        memory = [[.memory]]
        cpu = [[.cpu]]
      }
    }
  }
}
