provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_separate_instance_rule" "basic" {
  instance_id        = [nifcloud_instance.basic1.instance_id, nifcloud_instance.basic2.instance_id]
  availability_zone  = "east-21"
  description        = "memo"
  name               = "%s"
}

resource "nifcloud_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}

resource "nifcloud_key_pair" "basic" {
  key_name   = "%s"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCZ1FEVjJpcjBTWjUvWTBCRm9DK1pRMVU4SUpISWZTWkc2QUljbHFCclhqaTNYZ2h3eG9PYzgxUkZmTW55aVB3OGRsakVodlFTcnl0eXpZNkhkVDZZZVR1OWhYWE9sckw3SlExbDVWbEZmT3VsZGlWQi92YTVzL2ZNQlR2SG50aHh4a3hiTm9BYkphQ1lxQVJucStHemU2clNGOEFHOC9DckUwckxuK2tlK1Jkb0d6Mk9uRlc0MDZId01uZVBkRm1QSzFKYjhUZVZMNzUyN3pUaUs0anV2SXU2TlQ2MU96aDh4OHZzRkhzNm52NWRRR0FCdm8rMjUycDJMdUlwczlnNDIydmg1VGhpQ0FPTmRXdjQvZHZrVWg4NDN6a1VRL0tISGNhWkpjcG1zdXNPNUhnbzdKLzk4VVVBU0NPVGgwSVZxZjFtQXdxRkZLVjFkTEw2YnJES2lTTFMwQVkwWUdkMHMvN3lGMTdIK2o1VDVPNjd2Z0RqbTR3K041MFhvUVIwbU5BY0t3UVM0NHhkWkRxallXTzVuc0ZVOWZZY3RsejQ2Qk5xTk51My9GOWJVbFhBM0dkY2FHRmw5elZZQjVwWTdqOW9jbFQ1VWNXdkY1UXByYWFRZGhxVEkxZjFRclRLRkN6Vm1Dc1ROWkZBZU1VMVcwTWFUU1QreVljK0NNc2xSa009IFNDSjAwMDg3QHVidW50dQo="
}

resource "nifcloud_instance" "basic1" {
  instance_id       = "%s1"
  availability_zone = "east-21"
  image_id          = data.nifcloud_image.ubuntu.id
  key_name          = nifcloud_key_pair.basic.key_name
  security_group    = nifcloud_security_group.basic.group_name
  instance_type     = "mini"
  accounting_type   = "2"

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }

  network_interface {
    network_id = "net-COMMON_PRIVATE"
  }

}

resource "nifcloud_instance" "basic2" {
  instance_id       = "%s2"
  availability_zone = "east-21"
  image_id          = data.nifcloud_image.ubuntu.id
  key_name          = nifcloud_key_pair.basic.key_name
  security_group    = nifcloud_security_group.basic.group_name
  instance_type     = "mini"
  accounting_type   = "2"

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }

  network_interface {
    network_id = "net-COMMON_PRIVATE"
  }

}

data "nifcloud_image" "ubuntu" {
  image_name = "Ubuntu Server 22.04 LTS"
}
