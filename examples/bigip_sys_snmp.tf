provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}

resource "bigip_sys_snmp" "snmp" {
  sys_contact = " NetOPsAdmin s.shitole@f5.com" 
  sys_location = "SeattleHQ"
  allowedaddresses = ["202.10.10.2"]
}

