terraform {
  required_version = ">= 1.9.0, < 2.0.0"

  required_providers {
    gitlab {
      version = "~> 18.0"
      source  = "gitlabhq/gitlab"
    }
    null {
      version = "~> 3.0"
      source  = "hashicorp/null"
    }
  }

  backend "http" {
    foo = "bar"
  }
}

variable "test" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))

  description = "this is a test"
  sensitive   = false
  nullable    = false
  ephemeral   = false
  const       = false

  validation {
    condition     = length(var.image_id) > 4 && substr(var.image_id, 0, 4) == "ami-"
    error_message = "some msg"
  }
}

locals {
  k1 = "value"
  k2 = 123
  k3 = false
  k4 = [1, 2.1]
  k5 = { a = "bar", b = "foo" }
  k6 = local.k4[1]
}

resource "label" "tag" {
  k1 = "value"
  k2 = 123
  k3 = false
  k4 = [1, 2.1]
  k5 = { a = "bar", b = "foo" }
}
