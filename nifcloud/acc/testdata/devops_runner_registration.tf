provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_devops_runner_registration" "basic" {
  runner_name          = nifcloud_devops_runner.basic.name
  gitlab_url           = var.devops_gitlab_url
  parameter_group_name = nifcloud_devops_runner_parameter_group.basic.name
  token                = var.devops_runner_token
}

resource "nifcloud_devops_runner" "basic" {
  name          = "%s"
  instance_type = "c-small"
}

resource "nifcloud_devops_runner_parameter_group" "basic" {
  name = "%s"
}

resource "nifcloud_devops_runner_parameter_group" "upd" {
  name = "%s-upd"
}

variable "devops_gitlab_url" {
  description = "test devops GitLab URL"
  type        = string
}

variable "devops_runner_token" {
  description = "test devops runner token"
  type        = string
}
