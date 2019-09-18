provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}

resource "bigip_sys_snmp_traps" "snmp_traps" {
name = "snmptraps"
community = "f5community"
host = "195.10.10.1"
description = "Setup snmp traps"
port = 111
}

