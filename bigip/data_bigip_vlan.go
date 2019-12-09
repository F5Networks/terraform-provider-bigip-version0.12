/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceBigipNetVlan() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceBigipNetVlanRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the VLAN",
			},

			"tag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VLAN ID (tag)",
			},

			"interfaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Interface(s) attached to the VLAN",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vlanport": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Vlan name",
						},

						"tagged": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Interface tagged",
						},
					},
				},
			},
		},
	}

}

func dataSourceBigipNetVlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Reading VLAN %s", name)

	vlan, err := client.Vlan(name)
	if err != nil {
		return fmt.Errorf("Error retrieving VLAN %s: %v", name, err)
	}
	if vlan == nil {
		log.Printf("[DEBUG] VLAN %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	d.Set("name", vlan.FullPath)
	d.Set("tag", vlan.Tag)

	log.Printf("[DEBUG] Reading VLAN %s Interfaces", name)

	vlanInterfaces, err := client.GetVlanInterfaces(name)
	if err != nil {
		return fmt.Errorf("Error retrieving VLAN %s Interfaces: %v", name, err)
	}

	var interfaces []map[string]interface{}
	var ifaceTagged bool
	for _, iface := range vlanInterfaces.VlanInterfaces {
		if iface.Tagged {
			ifaceTagged = true
		} else {
			ifaceTagged = false
		}
		log.Printf("[DEBUG] Retrieved VLAN Interface %s, tagging is set to %t", iface.Name, ifaceTagged)

		vlanIface := map[string]interface{}{
			"vlanport": iface.Name,
			"tagged":   ifaceTagged,
		}

		interfaces = append(interfaces, vlanIface)
	}

	if err := d.Set("interfaces", interfaces); err != nil {
		return fmt.Errorf("Error updating Interfaces in state for VLAN %s: %v", name, err)
	}

	return nil
}
