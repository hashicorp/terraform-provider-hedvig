package hedvig

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

type createDiskResponse struct {
	Result []struct {
		Name    string `json:"name"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"result"`
	RequestID           string `json:"requestId"`
	Type                string `json:"type"`
	Status              string `json:"status"`
	Residence           string `json:"residence"`
	Message             string `json:"message"`
	ReplicationFactor   string `json:"replicationFactor"`
	Deduplication       string `json:"deduplication"`
	Compressed          string `json:"compressed"`
	BlockSize           string `json:"blockSize"`
	ClusteredFileSystem string `json:"clusteredFileSystem"`
	CacheEnabled        string `json:"cacheEnabled"`
	Scsi3pr             string `json:"scsi3pr"`
	Encryption          string `json:"encryption"`
	ReplicationPolicy   string `json:"replicationPolicy"`
	Description         string `json:"description"`
}

type readDiskResponse struct {
	Result struct {
		VDiskName string `json:"vDiskName"`
		Size      struct {
			Units string `json:"units"`
			Value int    `json:"value"`
		} `json:"size"`
		DiskType string `json:"diskType"`
	} `json:"result"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type updateDiskResponse struct {
	RequestID string `json:"requestId"`
	Result    []struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"result"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

type deleteDiskResponse struct {
	Result []struct {
		Name    string `json:"name"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"result"`
	RequestID string `json:"requestId"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

func resourceVdisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceVdiskCreate,
		Read:   resourceVdiskRead,
		Update: resourceVdiskUpdate,
		Delete: resourceVdiskDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"residence": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "HDD",
				ValidateFunc: validation.StringInSlice([]string{
					"Flash",
					"HDD",
				}, true),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NFS",
					"BLOCK",
				}, true),
			},
			"replicationfactor": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      "3",
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 6),
			},
			"deduplication": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"compressed": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, true),
			},
			"blocksize": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "4096",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"512",
					"4096",
					"4k",
					"65536",
					"64k",
				}, true),
			},
			"clusteredfilesystem": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, true),
			},
			"scsi3pr": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, true),
			},
			"cacheenabled": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "false",
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, true),
			},
			"encryption": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "false",
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, true),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"replicationpolicy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Agnostic",
				ValidateFunc: validation.StringInSlice([]string{
					"Agnostic",
					"DataCenterAware",
					"RackAware",
					//		"RackUnaware",
				}, true),
			},
		},
	}
}

func resourceVdiskCreate(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	sessionID, err := GetSessionId(d, meta.(*HedvigClient))
	if err != nil {
		return err
	}

	compress := "false"

	if (d.Get("deduplication") == "true") && (d.Get("compressed") == "false") {
		return fmt.Errorf("Deduplication enabled, compression must also be enabled.")
	} else if d.Get("compressed") == "true" {
		compress = "true"
	}

	if d.Get("deduplication") == "true" && d.Get("type") == "BLOCK" && d.Get("clusteredfilesystem") == "true" {
		return fmt.Errorf("Deduplication cannot be enabled for a block virtual disk with a clustered file system.")
	}

	if d.Get("blocksize").(string) != "512" && d.Get("type") == "NFS" {
		return fmt.Errorf("Block size must be 512 on NFS disks")
	}

	//if !(d.Get("blocksize").(string) == "4k" || d.Get("blocksize").(string) == "4096") {
	//	if d.Get("deduplication") == "true" {
	//		return fmt.Errorf("Deduplication enabled, block size must be 4k (or 4096)")
	//	}
	//}

	blocksize := ""

	if d.Get("blocksize").(string) == "4k" {
		blocksize = "4096"
	} else if d.Get("blocksize").(string) == "64k" {
		blocksize = "65536"
	} else {
		blocksize = d.Get("blocksize").(string)
	}

	if d.Get("residence").(string) != "HDD" && d.Get("deduplication") == "true" {
		return fmt.Errorf("Deduplication enabled, residence must be HDD.")
	}

	if d.Get("clusteredfilesystem") == "false" && d.Get("type") == "NFS" {
		return fmt.Errorf("Disk type is NFS, clustered file system must be enabled.")
	}

	if d.Get("clusteredfilesystem") == "true" && d.Get("blocksize") != "512" {
		return fmt.Errorf("Block Size must be 512 when Clustered File System is enabled.")
	}

	if d.Get("scsi3pr") == "true" && d.Get("type") == "NFS" {
		return fmt.Errorf("Clustered Shared Volumes (scsi3pr) not supported for NFS disks.")
	}

	if d.Get("cacheenabled") == "false" && d.Get("deduplication") == "true" {
		return fmt.Errorf("Client-side caching should be enabled when deduplication is.")
	}

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:AddVirtualDisk, category:VirtualDiskManagement, params:{name:'%s', size:{unit:'GB', value:%d}, diskType:%s, residence:%s, replicationFactor:%d, deduplication:%t, compressed:%s, blockSize:%s, scsi3pr:%s, cacheEnabled:%s, replicationPolicy:%s, clusteredFileSystem:%s, encryption:%s, description:'%s'}, sessionId:'%s'}", d.Get("name").(string), d.Get("size").(int), d.Get("type").(string), d.Get("residence"), d.Get("replicationfactor").(int), d.Get("deduplication"), compress, blocksize, d.Get("scsi3pr"), d.Get("cacheenabled"), d.Get("replicationpolicy").(string), d.Get("clusteredfilesystem"), d.Get("encryption"), d.Get("description"), sessionID))
	u.RawQuery = q.Encode()
	log.Printf("URL: %v", u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	createResp := createDiskResponse{}
	err = json.Unmarshal(body, &createResp)
	if err != nil {
		return err
	}

	//TODO: check for better way of returning results
	if len(createResp.Result) < 1 {
		return errors.New(createResp.Message)
	}

	if createResp.Result[0].Status != "ok" {
		if strings.HasSuffix(createResp.Message, "Run setkmsinfo command") {
			return fmt.Errorf("Cannot enable encryption without setting up KMS. Please refer to the Hedvig Encrypt360 Guide for assistance.")
		}
		return fmt.Errorf("Error creating vdisk %q: %s", d.Get("name").(string), createResp.Result[0].Message)
	}

	d.SetId("vdisk$" + d.Get("name").(string) + "$" + d.Get("type").(string))

	return resourceVdiskRead(d, meta)
}

func resourceVdiskRead(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	sessionID, err := GetSessionId(d, meta.(*HedvigClient))
	if err != nil {
		return err
	}

	idSplit := strings.Split(d.Id(), "$")
	log.Printf("idSplit: %v", idSplit)
	if len(idSplit) != 3 {
		return fmt.Errorf("Invalid ID: %s", d.Id())
	}

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:VirtualDiskDetails,category:VirtualDiskManagement,params:{virtualDisk:'%s'},sessionId:'%s'}", idSplit[1], sessionID))

	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return errors.New("Malformed query; aborting")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	readResp := readDiskResponse{}
	err = json.Unmarshal(body, &readResp)
	if err != nil {
		return err
	}

	if readResp.Status == "warning" && strings.HasSuffix(readResp.Message, "t be found") {
		d.SetId("")
		log.Printf("Vdisk not found, clearing from state")
		return nil
	}

	if readResp.Result.DiskType == "NFS_MASTER_DISK" {
		d.Set("type", "NFS")
	} else {
		d.Set("type", readResp.Result.DiskType)
	}
	d.Set("name", readResp.Result.VDiskName)
	d.Set("size", readResp.Result.Size.Value)

	return nil
}

// TODO: Verify and add tests
func resourceVdiskUpdate(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID, err := GetSessionId(d, meta.(*HedvigClient))
	if err != nil {
		return err
	}

	idSplit := strings.Split(d.Id(), "$")
	log.Printf("idSplit: %v", idSplit)
	if len(idSplit) != 3 {
		return fmt.Errorf("Invalid ID : %s", d.Id())
	}

	if d.HasChange("size") {
		q.Set("request", fmt.Sprintf("{type:VirtualDiskDetails,category:VirtualDiskManagement,params:{virtualDisk:'%s'},sessionId:'%s'}", idSplit[1], sessionID))

		u.RawQuery = q.Encode()

		resp, err := http.Get(u.String())

		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		readResp := readDiskResponse{}
		err = json.Unmarshal(body, &readResp)
		if err != nil {
			return err
		}

		if readResp.Result.Size.Value > d.Get("size").(int) {
			return errors.New("Cannot downsize a virtual disk")
		}

		q.Set("request", fmt.Sprintf("{type:ResizeDisks, category:VirtualDiskManagement, params:{virtualDisks:['%s'], size:{unit:'GB', value:%d}}, sessionId:'%s'}", idSplit[1], d.Get("size").(int),
			sessionID))
		u.RawQuery = q.Encode()
		log.Printf("URL: %v", u.String())

		resp, err = http.Get(u.String())

		if err != nil {
			return err
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		updateResp := updateDiskResponse{}
		err = json.Unmarshal(body, &updateResp)
		if err != nil {
			return err
		}

		if updateResp.Status != "ok" {
			return fmt.Errorf("Error updating vdisk: %s", updateResp.Status)
		}

		log.Printf("body: %s", body)
	}

	return resourceVdiskRead(d, meta)
}

func resourceVdiskDelete(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	sessionID, err := GetSessionId(d, meta.(*HedvigClient))
	if err != nil {
		return err
	}

	idSplit := strings.Split(d.Id(), "$")
	log.Printf("idSplit: %v", idSplit)
	if len(idSplit) != 3 {
		return fmt.Errorf("Invalid ID: %s", d.Id())
	}

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:DeleteVDisk, category:VirtualDiskManagement, params:{virtualDisks:['%s']}, sessionId:'%s'}}", idSplit[1], sessionID))

	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	deleteResp := deleteDiskResponse{}
	err = json.Unmarshal(body, &deleteResp)
	if err != nil {
		return err
	}

	if deleteResp.Result[0].Status != "ok" {
		return fmt.Errorf("Error deleting vdisk: %s", deleteResp.Result[0].Message)
	}
	return nil
}
