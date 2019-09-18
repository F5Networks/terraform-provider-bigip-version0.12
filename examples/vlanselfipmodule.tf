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

resource "bigip_net_selfip" "selfip" {

        name = "/Common/InternalselfIP"
        ip = "100.1.1.1/24"
        vlan = "/Common/intvlan"
        depends_on = ["module.sjvlan1"]
        }



