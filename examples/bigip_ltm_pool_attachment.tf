provider "bigip" {
  address = "xxx.xxx.xxx.xxx"
  username = "xxxx"
  password = "xxxx"
}


resource "bigip_ltm_pool_attachment" "attach_node" {
        pool = "/Common/terraform-pool"
	node = "/Common/11.1.1.101:80"
	depends_on = ["bigip_ltm_pool.pool"]

}

