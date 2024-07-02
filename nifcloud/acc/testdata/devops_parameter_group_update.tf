provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_devops_parameter_group" "basic" {
  name        = "%s-upd"
  description = "tfacc-memo-upd"

  sensitive_parameter {
    name  = "smtp_password"
    value = "mynewstrongpassword"
  }

  parameter {
    name  = "smtp_user_name"
    value = "user101"
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
