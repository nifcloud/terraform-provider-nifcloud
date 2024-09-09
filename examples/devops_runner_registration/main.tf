terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_devops_runner_registration" "example" {
  runner_name          = nifcloud_devops_runner.example.name
  gitlab_url           = "https://gitlab.com/"
  parameter_group_name = nifcloud_devops_runner_parameter_group.example.name
  token                = "glrt-thegitlabrunnertoken"
}

resource "nifcloud_devops_runner" "example" {
  name          = "example"
  instance_type = "c-small"
}

resource "nifcloud_devops_runner_parameter_group" "example" {
  name = "example"
}
