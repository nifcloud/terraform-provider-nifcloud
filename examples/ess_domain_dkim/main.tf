terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

resource "nifcloud_ess_domain_identity" "example" {
  domain = "example.com"
}

resource "nifcloud_ess_domain_dkim" "example" {
  domain = nifcloud_ess_domain_identity.example.domain
}

resource "nifcloud_dns_record" "example" {
  count   = 3
  zone_id = "ABCDEFGHIJ123"
  name    = "${element(nifcloud_ess_domain_dkim.example.dkim_tokens, count.index)}._domainkey"
  type    = "CNAME"
  ttl     = "600"
  record  = "${element(nifcloud_ess_domain_dkim.example.dkim_tokens, count.index)}.dkim.ess.nifcloud.com"
}
