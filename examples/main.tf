terraform {
  required_providers {
    elasticsearch = {
      version = "0.1"
      source  = "mattermost.com/terraform/elasticsearch"
    }
  }
}

provider "elasticsearch" {
  url = "http://localhost:9200"
}

resource "elasticsearch_template" "new" {
  name = "logstash"

  template = <<EOF
  {
    "index_patterns" :["logstash-*"],
    "order" : 0,
    "settings": {
      "index.mapping.total_fields.limit" : "7000"
    }
  }
  EOF
}

