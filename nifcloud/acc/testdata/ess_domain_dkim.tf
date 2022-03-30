resource "nifcloud_ess_domain_identity" "basic" {
  domain = "%s"
}

resource "nifcloud_ess_domain_dkim" "basic" {
  domain = nifcloud_ess_domain_identity.basic.domain
}
