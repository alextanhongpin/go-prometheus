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

resource "grafana_data_source" "prometheus" {
  type                = "prometheus"
  name                = "prometheus"
  url                 = "http://prometheus:9090"
  basic_auth_enabled  = false
  json_data_encoded  = jsonencode({
    httpMethod = "POST"
  })
}

resource "grafana_data_source" "loki" {
  type                = "loki"
  name                = "loki"
  url                 = "http://loki:3100"
  basic_auth_enabled  = false
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
