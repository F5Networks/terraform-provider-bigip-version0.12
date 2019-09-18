provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  alias = "east"
  username = "xxxxx"
  password = "xxxxx"
}

provider "bigip" {
   alias = "west"
   address = "xxx.xxx.xxx.xxx"
   username = "xxxxx"
   password = "xxxxx"
}


resource "bigip_ltm_node" "node_west" {
  name = "/Common/terraform_node1"
  provider = "bigip.west"
  address = "1.1.1.1"
  state = "user-up"
}

resource "bigip_ltm_node" "node_east" {
  name = "/Common/terraform_node1"
  provider = "bigip.east"
  address = "1.1.1.1"
  state = "user-down"
}



