---
layout: "bigip"
page_title: "BIG-IP: bigip_sys_ntp"
sidebar_current: "docs-bigip-resource-ntp-x"
description: |-
    Provides details about bigip ntp
---

# bigip\_ntp

`bigip_sys_ntp` provides details about a specific bigip

This resource is helpful when configuring NTP server on the BIG-IP.
## Example Usage


```hcl
provider "bigip" {
  address = "10.192.74.73"
  username = "admin"
  password = "admin"
}


resource "bigip_sys_ntp" "ntp1" {

description = "/Common/NTP1"
  servers = ["time.facebook.com"]
  timezone = "America/Los_Angeles"
}

```      

## Argument Reference

* `bigip_sys_ntp` - Is the resource is used to configure ntp server on the BIG-IP.

* `/Common/NTP1` - Is the description of the NTP server in the main or common partition of BIG-IP.

* `time.facebook.com` - Is the  NTP server configured on the BIG-IP.
