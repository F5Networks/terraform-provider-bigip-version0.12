provider "bigip" {
 address = "xxx.xxx.xxx.xxx"
 username = "xxxxx"
 password = "xxxxx"
}


module  "sjvlan1" {
  source = "./vlanmodule"
  name = "/Common/intvlan"
  tag = 101
  vlanport = "1.1"
  tagged = true
 }

module "sjvlan2"  {
  source = "./vlanmodule"
  name = "/Common/extvlan"
  tag = 102
  vlanport = "1.2"
  tagged = true
 }

