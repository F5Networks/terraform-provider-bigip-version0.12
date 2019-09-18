provider "bigip" {
  address = "xxx.xxx.xxx.xxx" // bigip ip address //
  username = "xxxxx" 
  password = "xxxxx"
}


resource "bigip_cm_device" "my_new_device"

        {
            name = "bigip300.f5.com"
            configsync_ip = "2.2.2.2"
            mirror_ip = "10.10.10.10"
            mirror_secondary_ip = "11.11.11.11"
        }
