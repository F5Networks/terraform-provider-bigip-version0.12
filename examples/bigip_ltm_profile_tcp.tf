provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}


resource "bigip_ltm_profile_tcp" "sanjose-tcp-lan-profile"

        {  
            name = "sanjose-tcp-lan-profile"
            idle_timeout = 200
            close_wait_timeout = 5
            finwait_2timeout = 5
            finwait_timeout = 300
            keepalive_interval = 1700
            deferred_accept = "enabled"
            fast_open = "enabled"
        }


