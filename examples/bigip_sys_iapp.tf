provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}

resource "bigip_sys_iapp" "waf_asm" {
  name = "policywaf"
  jsonfile = "${file("policywaf.json")}"
}

resource "bigip_sys_iapp" "pool_deployed" {
  name = "sap-dmzpool-rp1-80"
  jsonfile = "${file("sap-dmzpool-rp1-80.json")}"
}

