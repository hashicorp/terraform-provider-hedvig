package hedvig

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type LunResponse struct {
	Result struct {
		TargetLocations []string `json:"targetLocations"`
	} `json:"result"`
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

	log.Printf("body: %s", body)

	d.SetId("lun-" + d.Get("vdisk").(string) + "-" + d.Get("controller").(string))

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

	dsplit := strings.Split(d.Id(), "-")

	if len(dsplit) < 2 {
		return errors.New("Too few fields in ID")
	}

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:VirtualDiskDetails,category:VirtualDiskManagement,params:{virtualDisk:'%s'},sessionId:'%s'}", dsplit[1], sessionID))

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
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

		_, err := http.Get(u.String())

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

	log.Printf("body: %s", body)

	return nil
}
