terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}


resource "nifcloud_separate_instance_rule" "example" {
  instance_id        = ["testrun001", "testrun002"]     
  availability_zone  = "east-11"               
  description        = "成果報告会用test"            
  name               = "test001"               
}

