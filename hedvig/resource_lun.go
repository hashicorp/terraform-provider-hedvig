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
)

type LunResponse struct {
	Result struct {
		TargetLocations []string `json:"targetLocations"`
	} `json:"result"`
}

type createLunResponse struct {
	Result []struct {
		Name    string `json:"name"`
		Targets []struct {
			Name    string `json:"name"`
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"targets"`
		Status string `json:"status"`
	} `json:"result"`
	RequestID string `json:"requestId"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

type deleteLunResponse struct {
	RequestId string `json: "requestId"`
	Status    string `json: "status"`
	Type      string `json: "type"`
}

func resourceLun() *schema.Resource {
	return &schema.Resource{
		Create: resourceLunCreate,
		Read:   resourceLunRead,
		Delete: resourceLunDelete,

		Schema: map[string]*schema.Schema{
			"vdisk": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"controller": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLunCreate(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID, err := GetSessionId(d, meta.(*HedvigClient))

	if err != nil {
		return err
	}

	q.Set("request", fmt.Sprintf("{type:AddLun, category:VirtualDiskManagement, params:{virtualDisks:['%s'], targets:['%s'], readonly:false}, sessionId:'%s'}", d.Get("vdisk").(string), d.Get("controller").(string),
		sessionID))
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

	createResp := createLunResponse{}
	err = json.Unmarshal(body, &createResp)

	if err != nil {
		return err
	}

	if createResp.Result[0].Targets[0].Status != "ok" {
		return errors.New("Error creating export: " + createResp.Result[0].Targets[0].Message)
	}

	log.Printf("body: %s", body)

	d.SetId("lun$" + d.Get("vdisk").(string) + "$" + d.Get("controller").(string))

	return resourceLunRead(d, meta)
}

func resourceLunRead(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	sessionID, err := GetSessionId(d, meta.(*HedvigClient))

	if err != nil {
		return err
	}

	IdSplit := strings.Split(d.Id(), "$")

	if len(IdSplit) != 3 {
		return errors.New("Incorrect number of fields in ID")
	}

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:VirtualDiskDetails,category:VirtualDiskManagement,params:{virtualDisk:'%s'},sessionId:'%s'}", IdSplit[1], sessionID))

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		d.SetId("")
		log.Print("Lun resource not found in virtual disk, clearing from state")
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	lun := LunResponse{}
	err = json.Unmarshal(body, &lun)

	if err != nil {
		return err
	}

	if len(lun.Result.TargetLocations) < 1 {
		return errors.New("Not enough results found to define resource")
	}

	controllerparts := strings.Split(lun.Result.TargetLocations[0], ":")[0]

	if len(controllerparts) < 1 {
		return errors.New("Insufficient data in lun.Result")
	}

	d.Set("controller", controllerparts)

	return nil
}

func resourceLunUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("vdisk") || d.HasChange("controller") {
		dOldVDisk, _ := d.GetChange("vdisk")
		dOldController, _ := d.GetChange("controller")

		u := url.URL{}
		u.Host = meta.(*HedvigClient).Node
		u.Path = "/rest/"
		u.Scheme = "http"

		q := url.Values{}

		sessionID, err := GetSessionId(d, meta.(*HedvigClient))

		if err != nil {
			return err
		}

		q.Set("request", fmt.Sprintf("{type:UnmapLun, category:VirtualDiskManagement, params:{virtualDisk:'%s', target:'%s'}, sessionId: '%s'}", dOldVDisk.(string), dOldController.(string), sessionID))

		u.RawQuery = q.Encode()
		log.Printf("URL: %v", u.String())

		_, err = http.Get(u.String())

		if err != nil {
			return err
		}

		resourceLunCreate(d, meta)
	}

	return resourceLunRead(d, meta)
}

func resourceLunDelete(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID, err := GetSessionId(d, meta.(*HedvigClient))

	if err != nil {
		return err
	}

	q.Set("request", fmt.Sprintf("{type:UnmapLun, category:VirtualDiskManagement, params:{virtualDisk:'%s', target:'%s'}, sessionId: '%s'}", d.Get("vdisk"), d.Get("controller"),
		sessionID))
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

	deleteLunResp := deleteLunResponse{}

	err = json.Unmarshal(body, &deleteLunResp)

	log.Printf("body: %s", body)

	return nil
}
