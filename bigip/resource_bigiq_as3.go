package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
	"sync"
	"time"
)

var x = 0
var m sync.Mutex

func resourceBigiqAs3() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigiqAs3Create,
		Read:   resourceBigiqAs3Read,
		Update: resourceBigiqAs3Update,
		Delete: resourceBigiqAs3Delete,
		Exists: resourceBigiqAs3Exists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bigiq_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The registration key pool to use",
			},
			"bigiq_user": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_token_auth": {
				Type:      schema.TypeBool,
				Optional:  true,
				Sensitive: true,
				Default:   false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					//log.Printf("Value of k=%v,old=%v,new%v", k, old, new)
					if old != new {
						return true
					}
					return false
				},
				Description: "Enable to use an external authentication source (LDAP, TACACS, etc)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_TOKEN_AUTH", nil),
			},
			"bigiq_login_ref": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "tmos",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					//log.Printf("Value of k=%v,old=%v,new%v", k, old, new)
					if old != new {
						return true
					}
					return false
				},
				Description: "Login reference for token authentication (see BIG-IQ REST docs for details)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_LOGIN_REF", nil),
			},
			"as3_json": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "AS3 json",
			},
			"tenant_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of Tenant",
			},
		},
	}
}

func resourceBigiqAs3Create(d *schema.ResourceData, meta interface{}) error {
	bigipRef := meta.(*bigip.BigIP)
	log.Println(bigipRef)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	m.Lock()
	defer m.Unlock()
	as3_json := d.Get("as3_json").(string)
	tenantList, _ := bigiqRef.GetTenantList(as3_json)
	log.Println(tenantList)
	err = bigiqRef.PostAs3Bigiq(as3_json)
	if err != nil {
		return fmt.Errorf("Error posting as3 from bigiq :%v", err)
	}
	_ = d.Set("tenant_list", tenantList)
	d.SetId("tenantList")
	x = x + 1
	return resourceBigiqAs3Read(d, meta)
}

func resourceBigiqAs3Read(d *schema.ResourceData, meta interface{}) error {
	time.Sleep(20 * time.Second)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	log.Printf("[INFO] Reading As3 config")
	name := d.Get("tenant_list").(string)
	as3Resp, err := bigiqRef.GetAs3Bigiq(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		return err
	}
	if as3Resp == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

func resourceBigiqAs3Exists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}

func resourceBigiqAs3Update(d *schema.ResourceData, meta interface{}) error {
	time.Sleep(20 * time.Second)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	as3Json := d.Get("as3_json").(string)
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Updating As3 Config :%s", as3Json)
	name := d.Get("tenant_list").(string)
	tenantList, _ := bigiqRef.GetTenantList(as3Json)
	if tenantList != name {
		d.Set("tenant_list", tenantList)
		new_list := strings.Split(tenantList, ",")
		old_list := strings.Split(name, ",")
		deleted_tenants := bigiqRef.TenantDifference(old_list, new_list)
		if deleted_tenants != "" {
			//err, _ := bigiqRef.DeleteAs3Bigip(deleted_tenants)
			//if err != nil {
			//	log.Printf("[ERROR] Unable to Delete removed tenants: %v :", err)
			//	return err
			//}
		}
	}
	err = bigiqRef.PostAs3Bigiq(as3Json)
	if err != nil {
		return fmt.Errorf("Error updating json  %s: %v", tenantList, err)
	}
	x = x + 1
	return resourceBigiqAs3Read(d, meta)
}

func resourceBigiqAs3Delete(d *schema.ResourceData, meta interface{}) error {
	time.Sleep(20 * time.Second)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Deleting As3 config")
	name := d.Get("tenant_list").(string)
	as3Json := d.Get("as3_json").(string)
	err, failedTenants := bigiqRef.DeleteAs3Bigiq(as3Json, name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete: %v :", err)
		return err
	}
	if failedTenants != "" {
		_ = d.Set("tenant_list", name)
		return resourceBigipAs3Read(d, meta)
	}
	x = x + 1
	//m.Unlock()
	d.SetId("")
	return nil
}
