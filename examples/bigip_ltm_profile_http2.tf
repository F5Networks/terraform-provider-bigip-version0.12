provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}


resource "bigip_ltm_profile_http2" "nyhttp2"

        {
            name = "/Common/NewYork_http2"
            defaults_from = "/Common/http2"
            concurrent_streams_per_connection = 10
            connection_idle_timeout= 30
            activation_modes = ["alpn","npn"]
        }
