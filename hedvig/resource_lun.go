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

type readLunResponse struct {
	Result struct {
		TargetLocations []string `json:"targetLocations"`
	} `json:"result"`
	Message string `json:"message"`
	Status  string `json:"status"`
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
	RequestId string `json:"requestId"`
	Message   string `json:"message"`
	Status    string `json:"status"`
	Type      string `json:"type"`
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

	q.Set("request", fmt.Sprintf("{type:AddLun, category:VirtualDiskManagement, params:{virtualDisks:['%s'], targets:['%s'], readonly:false}, sessionId:'%s'}", d.Get("vdisk").(string), d.Get("controller").(string), sessionID))
	u.RawQuery = q.Encode()

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
		return fmt.Errorf("Error creating export: %s", createResp.Result[0].Targets[0].Message)
	}

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

	idSplit := strings.Split(d.Id(), "$")
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

	readResp := readLunResponse{}
	err = json.Unmarshal(body, &readResp)
	if err != nil {
		return err
	}

	if readResp.Status == "warning" && strings.HasSuffix(readResp.Message, "t be found") {
		d.SetId("")
		log.Print("Lun resource not found in virtual disk, clearing from state")
		return nil
	}

	if readResp.Status != "ok" {
		return fmt.Errorf("Error reading lun details: %s", readResp.Message)
	}

	if len(readResp.Result.TargetLocations) < 1 {
		return errors.New("Not enough results found to define resource")
	}

	for _, target := range readResp.Result.TargetLocations {
		if strings.HasPrefix(target, idSplit[2]) {
			d.Set("controller", idSplit[2]) // cheating
			return nil
		}
	}

	return errors.New("Resource not found")
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

	idSplit := strings.Split(d.Id(), "$")
	if len(idSplit) != 3 {
		return fmt.Errorf("Invalid ID: %s", d.Id())
	}

	q.Set("request", fmt.Sprintf("{type:UnmapLun, category:VirtualDiskManagement, params:{virtualDisk:'%s', target:'%s'}, sessionId: '%s'}", idSplit[1], idSplit[2], sessionID))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	deleteResp := deleteLunResponse{}
	err = json.Unmarshal(body, &deleteResp)
	if err != nil {
		return err
	}

	if deleteResp.Status != "ok" {
		return fmt.Errorf("Error deleting lun: %s", deleteResp.Message)
	}
	return nil
}
