resource "nifcloud_ssl_certificate" "basic" {
  certificate = <<EOT
%sEOT
  key         = <<EOT
%sEOT
  ca          = <<EOT
%sEOT
  description = "memo-upd"
}
