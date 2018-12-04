package hedvig

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type createMountResponse struct {
	Result struct {
		ExportInfo []struct {
			Target string `json:"target"`
			Status string `json:"status"`
		} `json:"exportInfo"`
	} `json:"result"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

type readMountResponse struct {
	Result    []string `json:"result"`
	RequestID string   `json:"requestId"`
	Type      string   `json:"type"`
	Message   string   `json:"message"`
	Status    string   `json:"status"`
}

type deleteMountResponse struct {
	Result []struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"result"`
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

func resourceMount() *schema.Resource {
	return &schema.Resource{
		Create: resourceMountCreate,
		Read:   resourceMountRead,
		Delete: resourceMountDelete,

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

// TODO: Need to test against multiple controllers
func resourceMountCreate(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID, err := GetSessionId(d, meta.(*HedvigClient))

	if err != nil {
		return err
	}

	q.Set("request", fmt.Sprintf("{type:Mount, category:VirtualDiskManagement, params:{virtualDisk:'%s', targets:['%s']}, sessionId:'%s'}", d.Get("vdisk"), d.Get("controller"), sessionID))

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

	log.Printf("body: %s", body)

	createResp := createMountResponse{}
	err = json.Unmarshal(body, &createResp)
	if err != nil {
		return err
	}

	if createResp.Status != "ok" {
		return errors.New("Error creating export: " + createResp.Message)
	}

	if createResp.Result.ExportInfo[0].Status != "ok" {
		return errors.New("Error creating export")
	}
	d.SetId("mount$" + d.Get("vdisk").(string) + "$" + d.Get("controller").(string))

	return resourceMountRead(d, meta)
}

func resourceMountRead(d *schema.ResourceData, meta interface{}) error {
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
		return errors.New("Invalid ID: " + d.Id())
	}

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:ListExportedTargets,category:VirtualDiskManagement,params:{virtualDisk:'%s'},sessionId:'%s'}", idSplit[1], sessionID))

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		d.SetId("")
		s := strconv.Itoa(resp.StatusCode)
		log.Print("Received " + s + ", removing resource from state")
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("body: %s", body)
	readResp := readMountResponse{}
	err = json.Unmarshal(body, &readResp)
	if err != nil {
		return err
	}

	if readResp.Status != "ok" {
		return errors.New("Error: " + readResp.Message)
	}

	// TODO: verify is necessary - should be caught by readResp.Status and readResp.Message
	if len(readResp.Result) < 1 {
		return errors.New("Resource not found: " + idSplit[1])
	}

	d.Set("controller", readResp.Result[0])

	return nil
}

// TODO: remove?
func resourceMountUpdate(d *schema.ResourceData, meta interface{}) error {
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

		q.Set("request", fmt.Sprintf("{type:Unmount, category:VirtualDiskManagement, params:{virtualDisk:'%s', targets:['%s']}, sessionId: '%s'}", dOldVDisk.(string), dOldController.(string), sessionID))

		u.RawQuery = q.Encode()
		log.Printf("URL: %v", u.String())

		_, err = http.Get(u.String())

		if err != nil {
			return err
		}

		resourceMountCreate(d, meta)
	}

	return resourceMountRead(d, meta)
}

func resourceMountDelete(d *schema.ResourceData, meta interface{}) error {
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
		return errors.New("Invalid ID: " + d.Id())
	}

	q.Set("request", fmt.Sprintf("{type:Unmount, category:VirtualDiskManagement, params:{virtualDisk:'%s', targets:['%s']}, sessionId: '%s'}", idSplit[1], idSplit[2], sessionID))

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
	log.Printf("body: %s", body)

	deleteResp := deleteMountResponse{}
	err = json.Unmarshal(body, &deleteResp)
	if err != nil {
		return err
	}
	if deleteResp.Status != "ok" {
		return errors.New("Error deleting mount: " + deleteResp.Message)
	}
	return nil
}
