package bigip

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
        "strings"
        "reflect"
)

const (
	uriRegkey      = "regkey"
	uriLicenses    = "licenses"
	uriResolver    = "resolver"
	uriDevicegroup = "device-groups"
	uriCmBigip     = "cm-bigip-allBigIpDevices"
	uriDevice      = "device"
	uriMembers     = "members"
	uriTasks       = "tasks"
	uriManagement  = "member-management"
        uriDeclare      = "declare"
)

type BigiqDevice struct {
        Address string `json:"address"`
        Username      string `json:"username"`
        Password      string `json:"password"`
        Port     int    `json:"port,omitempty"`
}

type DeviceRef struct {
	Link string `json:"link"`
}
type ManagedDevice struct {
	DeviceReference DeviceRef `json:"deviceReference"`
}

type UnmanagedDevice struct {
	DeviceAddress string `json:"deviceAddress"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	HTTPSPort     int    `json:"httpsPort,omitempty"`
}

type regKeyPools struct {
	//Items      []struct {
	//	ID       string `json:"id"`
	//	Name     string `json:"name"`
	//	SortName string `json:"sortName"`
	//} `json:"items"`
	RegKeyPoollist []regKeyPool `json:"items"`
}

type regKeyPool struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SortName string `json:"sortName"`
}

type devicesList struct {
	DevicesInfo []deviceInfo `json:"items"`
}
type deviceInfo struct {
	Address           string `json:"address"`
	DeviceURI         string `json:"deviceUri"`
	Hostname          string `json:"hostname"`
	HTTPSPort         int    `json:"httpsPort"`
	IsClustered       bool   `json:"isClustered"`
	MachineID         string `json:"machineId"`
	ManagementAddress string `json:"managementAddress"`
	McpDeviceName     string `json:"mcpDeviceName"`
	Product           string `json:"product"`
	SelfLink          string `json:"selfLink"`
	State             string `json:"state"`
	UUID              string `json:"uuid"`
	Version           string `json:"version"`
}

type MembersList struct {
	Members []memberDetail `json:"items"`
}

type memberDetail struct {
	AssignmentType  string `json:"assignmentType"`
	DeviceAddress   string `json:"deviceAddress"`
	DeviceMachineID string `json:"deviceMachineId"`
	DeviceName      string `json:"deviceName"`
	ID              string `json:"id"`
	Message         string `json:"message"`
	Status          string `json:"status"`
}

type regKeyAssignStatus struct {
	ID             string `json:"id"`
	DeviceAddress  string `json:"deviceAddress"`
	AssignmentType string `json:"assignmentType"`
	DeviceName     string `json:"deviceName"`
	Status         string `json:"status"`
}

type LicenseParam struct {
	Address         string `json:"address,omitempty"`
	Port            int    `json:"port,omitempty"`
	AssignmentType  string `json:"assignmentType,omitempty"`
	Command         string `json:"command,omitempty"`
	Hypervisor      string `json:"hypervisor,omitempty"`
	LicensePoolName string `json:"licensePoolName,omitempty"`
	MacAddress      string `json:"macAddress,omitempty"`
	Password        string `json:"password,omitempty"`
	SkuKeyword1     string `json:"skuKeyword1,omitempty"`
	SkuKeyword2     string `json:"skuKeyword2,omitempty"`
	Tenant          string `json:"tenant,omitempty"`
	UnitOfMeasure   string `json:"unitOfMeasure,omitempty"`
	User            string `json:"user,omitempty"`
}

func (b *BigIP) PostLicense(config *LicenseParam) (string, error) {
	log.Printf("[INFO] %v license to BIGIP device:%v from BIGIQ", config.Command, config.Address)
	resp, err := b.postReq(config, uriMgmt, uriCm, uriDevice, uriTasks, uriLicensing, uriPool, uriManagement)
	if err != nil {
		return "", err
	}
	respRef := make(map[string]interface{})
	json.Unmarshal(resp, &respRef)
	respID := respRef["id"].(string)
	time.Sleep(5 * time.Second)
	return respID, nil
}
func (b *BigIP) GetLicenseStatus(id string) (map[string]interface{}, error) {
	licRes := make(map[string]interface{})
	err, _ := b.getForEntity(&licRes, uriMgmt, uriCm, uriDevice, uriTasks, uriLicensing, uriPool, uriManagement, id)
	if err != nil {
		return nil, err
	}
	licStatus := licRes["status"].(string)
	for licStatus != "FINISHED" {
		//log.Printf(" status response is :%s", licStatus)
		if licStatus == "FAILED" {
			log.Println("[ERROR]License assign/revoke status failed")
			return licRes, nil
		}
		return b.GetLicenseStatus(id)
	}
	log.Printf("License Assignment is :%s", licStatus)
	return licRes, nil
}

func (b *BigIP) GetDeviceLicenseStatus(path ...string) (string, error) {
	licRes := make(map[string]interface{})
	err, _ := b.getForEntity(&licRes, path...)
	if err != nil {
		return "", err
	}
	//log.Printf(" Initial status response is :%s", licRes["status"])
	return licRes["status"].(string), nil
}
func (b *BigIP) GetRegPools() (*regKeyPools, error) {
	var self regKeyPools
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses)
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func (b *BigIP) GetPoolType(poolName string) (*regKeyPool, error) {
	var self regKeyPools
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses)
	if err != nil {
		return nil, err
	}
	for _, pool := range self.RegKeyPoollist {
		if pool.Name == poolName {
			return &pool, nil
		}
	}
	err, _ = b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriUtility, uriLicenses)
	if err != nil {
		return nil, err
	}
	for _, pool := range self.RegKeyPoollist {
		if pool.Name == poolName {
			return &pool, nil
		}
	}
	return nil, nil
}

func (b *BigIP) GetManagedDevices() (*devicesList, error) {
	var self devicesList
	err, _ := b.getForEntity(&self, uriMgmt, uriShared, uriResolver, uriDevicegroup, uriCmBigip, uriDevices)
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func (b *BigIP) GetDeviceId(deviceName string) (string, error) {
	var self devicesList
	err, _ := b.getForEntity(&self, uriMgmt, uriShared, uriResolver, uriDevicegroup, uriCmBigip, uriDevices)
	if err != nil {
		return "", err
	}
	for _, d := range self.DevicesInfo {
		log.Printf("Address=%v,Hostname=%v,UUID=%v", d.Address, d.Hostname, d.UUID)
		if d.Address == deviceName || d.Hostname == deviceName || d.UUID == deviceName {
			log.Printf("SelfLink Type=%T,SelfLink=%v", d.SelfLink, d.SelfLink)
			return d.SelfLink, nil
		}
	}
	return "", nil
}

func (b *BigIP) GetRegkeyPoolId(poolName string) (string, error) {
	var self regKeyPools
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses)
	if err != nil {
		return "", err
	}
	for _, pool := range self.RegKeyPoollist {
		if pool.Name == poolName {
			return pool.ID, nil
		}
	}
	return "", nil
}
func (b *BigIP) RegkeylicenseAssign(config interface{}, poolId string, regKey string) (*memberDetail, error) {
	resp, err := b.postReq(config, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers)
	if err != nil {
		return nil, err
	}
	var resp1 regKeyAssignStatus
	err = json.Unmarshal(resp, &resp1)
	if err != nil {
		return nil, err
	}
	return b.GetMemberStatus(poolId, regKey, resp1.ID)
}

func (b *BigIP) GetMemberStatus(poolId, regKey, memId string) (*memberDetail, error) {
	var self memberDetail
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return nil, err
	}
	for self.Status != "LICENSED" {
		log.Printf("Member status:%+v", self.Status)
		if self.Status == "INSTALLATION_FAILED" {
			return &self, fmt.Errorf("INSTALLATION_FAILED with %s", self.Message)
		}
		return b.GetMemberStatus(poolId, regKey, memId)
	}
	return &self, nil
}
func (b *BigIP) RegkeylicenseRevoke(poolId, regKey, memId string) error {
	log.Printf("Deleting License for Member:%+v", memId)
	_, err := b.deleteReq(uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return err
	}
	r1 := make(map[string]interface{})
	err, _ = b.getForEntity(&r1, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return err
	}
	log.Printf("Response after delete:%+v", r1)
	return nil
}
func (b *BigIP) LicenseRevoke(config interface{}, poolId, regKey, memId string) error {
	log.Printf("Deleting License for Member:%+v from LicenseRevoke", memId)
	_, err := b.deleteReqBody(config, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return err
	}
	r1 := make(map[string]interface{})
	err, _ = b.getForEntity(&r1, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return err
	}
	log.Printf("Response after delete:%+v", r1)
	return nil
}
func (b *BigIP) PostAs3Bigiq(as3NewJson string) (error) {
    return b.post(as3NewJson, uriMgmt, uriShared, uriAppsvcs, uriDeclare )
    
}
func (b *BigIP) GetAs3Bigiq(name string) (string, error) {
as3Json := make(map[string]interface{})
	as3Json["class"] = "AS3"
	as3Json["action"] = "deploy"
	as3Json["persist"] = true
	adcJson := make(map[string]interface{})
        err, ok := b.getForEntityNew(&adcJson, uriMgmt, uriShared, uriAppsvcs, uriDeclare, name)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
       as3Json["declaration"] = adcJson
	out, _ := json.Marshal(as3Json)
	as3String := string(out)
	return as3String, nil
}

func (b *BigIP) DeleteAs3Bigiq(as3NewJson string, tenantName string) (error, string) {
 as3Json, err := tenantTrimToDelete(as3NewJson)
 if err != nil {
        log.Println("[ERROR] Error in trimming the as3 json")
        return err, ""
      }
return b.post(as3Json, uriMgmt, uriShared, uriAppsvcs, uriDeclare ), ""
}

func (b *BigIP) GetTenantList(body interface{}) (string, int) {
	s := make([]string, 0)
	as3json := body.(string)
	resp := []byte(as3json)
	jsonRef := make(map[string]interface{})
	json.Unmarshal(resp, &jsonRef)
	for key, value := range jsonRef {
		if rec, ok := value.(map[string]interface{}); ok && key == "declaration" {
			for k, v := range rec {
				if rec2, ok := v.(map[string]interface{}); ok {
					found := 0
					for k1, v1 := range rec2 {
						if k1 == "class" && v1 == "Tenant" {
							found = 1
						}
					}
					if found == 1 {
						s = append(s, k)
					}
				}
			}
		}
	}
	tenant_list := strings.Join(s[:], ",")
	return tenant_list, len(s)
}
func (b *BigIP) TenantDifference(slice1 []string, slice2 []string) string {
	var diff []string
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}

		}
	}
	diff_tenant_list := strings.Join(diff[:], ",")
	return diff_tenant_list
}
func tenantCompare(t1 string, t2 string) int {
	tenantList1 := strings.Split(t1, ",")
	tenantList2 := strings.Split(t2, ",")
	if len(tenantList1) == len(tenantList2) {
		return 1
	}
	return 0
}

func tenantTrimToDelete(resp string) (string, error) {
jsonRef := make(map[string]interface{})
	json.Unmarshal([]byte(resp), &jsonRef)

	for key, value := range jsonRef {
		if rec, ok := value.(map[string]interface{}); ok && key == "declaration" {
			for k, v := range rec {
                                 if (k == "target" && reflect.ValueOf(v).Kind() == reflect.Map) {
                                                                       continue
                                                                   }
				if rec2, ok := v.(map[string]interface{}); ok {
					for k1, v1 := range rec2 {
						if k1 != "class" && v1 != "Tenant" {
							delete(rec2,k1)
						}
					}       

				}
			}
		}
	}
        
      b, err := json.Marshal(jsonRef)
      if err != nil {
        return "", err
      }
      return string(b), nil
}
