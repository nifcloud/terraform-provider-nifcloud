provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_devops_parameter_group" "basic" {
  name        = "%s"
  description = "tfacc-memo"

  sensitive_parameter {
    name  = "smtp_password"
    value = "mystrongpassword"
  }

  parameter {
    name  = "smtp_user_name"
    value = "user1"
  }

  parameter {
    name  = "gitlab_email_from"
    value = "from@mail.com"
  }

  parameter {
    name  = "gitlab_email_reply_to"
    value = "reply-to@mail.com"
  }
}
