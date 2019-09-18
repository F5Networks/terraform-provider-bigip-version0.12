provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}

resource "bigip_ltm_policy" "test-policy" {
 name = "my_policy"
 strategy = "first-match"
  requires = ["http"]
 published_copy = "Drafts/my_policy"
  controls = ["forwarding"]
  rule  {
  name = "rule6"

   action = {
     tm_name = "20"
     forward = true
      pool = "/Common/mypool"
   }
  }
depends_on = ["bigip_ltm_pool.mypool"]
}

resource "bigip_ltm_pool" "mypool" {
    name = "/Common/mypool"
    monitors = ["/Common/http"]
    allow_nat = "yes"
    allow_snat = "yes"
    load_balancing_mode = "round-robin"
}

