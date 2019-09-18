provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}

resource "bigip_net_selfip" "selfip1" {
	name = "/Common/internalselfIP"
	ip = "11.1.1.1/24"
	vlan = "/Common/internal"
	depends_on = ["bigip_net_vlan.vlan1"]
	}

