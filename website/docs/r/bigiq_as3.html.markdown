---
layout: "bigip"
page_title: "BIG-IP: bigip_bigiq_as3"
sidebar_current: "docs-bigip-resource-x"
description: |-
    Provides details about bigiq as3 resource
---

# bigip_bigiq_as3

`bigip_bigiq_as3` provides details about bigiq as3 resource

This resource is helpful to configure as3 declarative JSON on BIG-IP through BIG-IQ.

## Example Usage 

```hcl


# Example Usage for json file
resource "bigip_bigiq_as3" "exampletask" {
bigiq_address = "xx.xx.xxx.xx"
bigiq_user = "xxxxx"
bigiq_password = "xxxxxxxxx"
as3_json = "${file("bigiq_example.json")}"
}


```

## Argument Reference


* `bigiq_address` - (Required) Address of the BIG-IQ to which your targer BIG-IP is attached

* `bigiq_user` - (Required) User name  of the BIG-IQ to which your targer BIG-IP is attached

* `bigiq_password` - (Required) Password of the BIG-IQ to which your targer BIG-IP is attached

* `as3_json` - (Required) Path/Filename of Declarative AS3 JSON which is a json file used with builtin ```file``` function

* `bigiq_example.json` - Example  AS3 Declarative JSON file

```hcl

{
    "class": "AS3",
    "action": "deploy",
    "persist": true,
    "declaration": {
        "class": "ADC",
        "schemaVersion": "3.7.0",
        "id": "example-declaration-01",
        "label": "Task1",
        "remark": "Task 1 - HTTP Application Service",
        "target": {
            "address": "xx.xxx.xx.xxx"
        },
        "Task1": {
            "class": "Tenant",
            "MyWebApp1http": {
                "class": "Application",
                "template": "http",


                "serviceMain": {
                    "class": "Service_HTTP",
                    "virtualAddresses": [
                        "10.1.2.10"
                    ],
                    "pool": "web_pool"
                },
                "web_pool": {
                    "class": "Pool",
                    "monitors": [
                        "http"
                    ],
                    "members": [
                        {
                            "servicePort": 80,
                            "serverAddresses": [
                                "192.0.2.33",
                                "192.0.2.13"
                            ],
                            "shareNodes": true
                        }
                    ]
                }
            }
        }
    }
}

```

* `AS3 documentation` - https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/big-iq.html

