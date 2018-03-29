provider "bigip" {
  address = "${var.bigip_management_ip}"
  username = "admin"
  password = "admin"
}




resource "bigip_ltm_monitor" "monitor" {
        name = "/Common/terraform_monitor"
        parent = "/Common/http"
        send = "GET /some/path\r\n"
        timeout = "999"
        interval = "999"
}

resource "bigip_ltm_pool"  "pool" {
        name = "/Common/terraform-pool"
        load_balancing_mode = "round-robin"
        nodes = ["10.0.2.167:80", "10.0.2.168:80"]
        monitors = ["/Common/terraform_monitor"]
        allow_snat = "yes"
        allow_nat = "yes"
}

resource "bigip_ltm_virtual_server" "http" {
	pool = "/Common/terraform-pool"
        name = "/Common/terraform_vs_http"
	destination = "10.0.1.100"
	port = 80
	source_address_translation = "automap"
        depends_on = ["bigip_ltm_pool.pool"]
}


