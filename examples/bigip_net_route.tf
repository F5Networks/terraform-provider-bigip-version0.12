provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}

resource "bigip_net_route" "route2" {
  name = "sanjay-route2"
  network = "10.10.10.0/24"
  gw      = "1.1.1.2"
}

