# https://registry.terraform.io/providers/grafana/grafana/latest/docs
terraform {
    required_providers {
        grafana = {
            source = "grafana/grafana"
            version = ">= 1.28.2"
        }
    }
}


provider "grafana" {
    url = "http://localhost:3000"
    auth = "admin:admin"
}

// Create resources (optional: within the organization)
resource "grafana_folder" "my_folder" {
  title  = "RED Metrics"
}

resource "grafana_dashboard" "test_dashboard" {
  folder = grafana_folder.my_folder.id
  config_json = file("${path.module}/dashboard.json")
}

resource "grafana_dashboard" "another_dashboard" {
  folder = grafana_folder.my_folder.id
  config_json = file("${path.module}/another.json")
}
