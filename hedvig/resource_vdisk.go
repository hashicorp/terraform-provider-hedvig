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
	RequestID string `json:"requestId"`
	Type      string `json:"type"`
	Status    string `json:"status"`
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

type diskDeleteResponse struct {
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NFS",
					"BLOCK",
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

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:AddVirtualDisk, category:VirtualDiskManagement, params:{name:'%s', size:{unit:'GB', value:%d}, diskType:%s, scsi3pr:false}, sessionId:'%s'}", d.Get("name").(string), d.Get("size").(int), d.Get("type").(string), sessionID))

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
		return errors.New("Unknown error creating Vdisk")
	}

	if createResp.Result[0].Status != "ok" {
		return fmt.Errorf("Error creating vdisk %q: %s", d.Get("name").(string), createResp.Result[0].Message)
	}

	// if resp.StatusCode != 200 {
	// 	d.SetId("")
	// 	// strresp := strconv.Itoa(resp.StatusCode)
	// 	// log.Print("Received " + strresp + " error, removing resource from state.")
	// 	log.Printf("Received %q error, removing resrouce from state.", resp.StatusCode)
	// 	return nil
	// }

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
		return errors.New("Invalid ID: " + d.Id())
	}

	q := url.Values{}
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

	//TODO: check status == ok

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
		return errors.New("Invalid ID")
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

		q.Set("request", fmt.Sprintf("{type:ResizeDisks, category:VirtualDiskManagement, params:{virtualDisks:['%s'], size:{unit:'GB', value:%d}}, sessionId:'%s'}", idSplit[2], d.Get("size").(int),
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
			return errors.New("Error in update response")
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
		return errors.New("Invalid ID")
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

	deleteResp := diskDeleteResponse{}
	err = json.Unmarshal(body, &deleteResp)
	if err != nil {
		return err
	}

	// if resp.StatusCode != 200 {
	// 	}

	if deleteResp.Result[0].Status != "ok" {
		return fmt.Errorf("Error deleting vdisk: %s", deleteResp.Result[0].Message)
	}

	d.SetId("")

	return nil
}
