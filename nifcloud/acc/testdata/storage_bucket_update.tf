provider "nifcloud" {
  storage_region = "jp-east-1"
}

resource "nifcloud_storage_bucket" "basic" {
  bucket = "%s"
  policy = <<EOT
%s
EOT

  versioning {
    enabled = true
  }
}
