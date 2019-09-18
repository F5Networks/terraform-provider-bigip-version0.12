provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}
resource "bigip_net_vlan" "vlan1" {
	name = "/Common/Internal"
	tag = 101
	interfaces = {
                vlanport = 1.2,
		tagged = false
	}
}
