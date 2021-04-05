provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_load_balancer" "basic" {
  load_balancer_name = "%s"
  instance_port      = 8080
  load_balancer_port = 8080
}

resource "nifcloud_load_balancer_listener" "basic" {
  load_balancer_name                          = nifcloud_load_balancer.basic.load_balancer_name
  instance_port                               = 8083
  load_balancer_port                          = 80
  balancing_type                              = "2"
  instances                                   = [nifcloud_instance.basic.instance_id]
  policy_type                                 = "standard"
  ssl_certificate_id                          = nifcloud_ssl_certificate.basic.id
  ssl_policy_id                               = 2
  filter_type                                 = 2
  filter                                      = ["192.168.1.2"]
  session_stickiness_policy_enable            = true
  session_stickiness_policy_expiration_period = 5
  sorry_page_enable                           = true
  sorry_page_status_code                      = 200
  health_check_interval                       = 11
  health_check_target                         = "ICMP"
  healthy_threshold                           = 1
  unhealthy_threshold                         = 3
  depends_on                                  = [nifcloud_instance.basic, nifcloud_ssl_certificate.basic]
}

resource "nifcloud_instance" "basic" {
  instance_id       = "%s"
  availability_zone = "east-21"
  image_id          = "221"
  key_name          = nifcloud_key_pair.basic.key_name
  depends_on        = [nifcloud_key_pair.basic]

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }

  network_interface {
    network_id = "net-COMMON_PRIVATE"
  }
}

resource "nifcloud_key_pair" "basic" {
  key_name   = "%s"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
}

resource "nifcloud_ssl_certificate" "basic" {
  certificate = <<EOT
%s
EOT
  key         = <<EOT
%s
EOT
  ca          = <<EOT
%s
EOT
  description = "memo"
}
