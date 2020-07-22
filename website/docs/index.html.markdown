---
layout: "bigip"
page_title: "BIG-IP Provider : Index"
sidebar_current: "docs-bigip-index"
description: |-
    Provides details about provider bigip
---

# F5 BIG-IP Provider

A [Terraform](https://terraform.io) provider for F5 BIG-IP. Resources are currently available for LTM.

### Requirements

This provider uses the iControlREST API. All the resources are validated with BigIP v12.1.1

## Example

```
provider "bigip" {
  address = "${var.url}"
  username = "${var.username}"
  password = "${var.password}"
}
```

## Reference

- `address` - (Required) Address of the device
- `username` - (Required) Username for authentication
- `password` - (Required) Password for authentication
- `token_auth` - (Optional, Default=false) Enable to use an external authentication source (LDAP, TACACS, etc)
- `login_ref` - (Optional, Default="tmos") Login reference for token authentication (see BIG-IP REST docs for details)


### Note

```
By default our provider will send telemetry data for this resource to TEEM production server.
If you don't want to send telemetry data to TEEM Server you can do with by setting below environment flag.

export TEEM_DISABLE=true

```
