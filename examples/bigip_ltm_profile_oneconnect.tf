provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}


resource "bigip_ltm_profile_oneconnect" "oneconnect-sanjose"

        {  
            name = "sanjose"
            partition = "Common"
            defaults_from = "/Common/oneconnect"
            idle_timeout_override = "disabled"
            max_age = 3600
            max_reuse = 1000
            max_size = 1000
            share_pools = "disabled"
            source_mask = "255.255.255.255"
        }


