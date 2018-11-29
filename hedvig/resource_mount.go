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
	"strconv"
	"strings"
)

type MountResponse struct {
	Result []string `json:"result"`
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

func resourceMountCreate(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:Mount, category:VirtualDiskManagement, params:{virtualDisk:'%s', targets:['%s']}, sessionId:'%s'}", d.Get("vdisk"), d.Get("controller"),
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

	d.SetId("mount-" + d.Get("vdisk").(string) + "-" + d.Get("controller").(string))

	return resourceMountRead(d, meta)
}

func resourceMountRead(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	dsplit := strings.Split(d.Id(), "-")

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:ListExportedTargets,category:VirtualDiskManagement,params:{virtualDisk:'%s'},sessionId:'%s'}", dsplit[1], sessionID))

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 404 {
		d.SetId("")
		s := strconv.Itoa(resp.StatusCode)
		log.Print("Received " + s + ", removing resource from state")
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	mount := MountResponse{}
	err = json.Unmarshal(body, &mount)

	if err != nil {
		return err
	}

	if len(mount.Result) < 1 {
		return errors.New("Resource not found, returning error.")
	}

	d.Set("controller", mount.Result[0])

	return nil
}

func resourceMountUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("vdisk") || d.HasChange("controller") {
		dOldVDisk, _ := d.GetChange("vdisk")
		dOldController, _ := d.GetChange("controller")

		u := url.URL{}
		u.Host = meta.(*HedvigClient).Node
		u.Path = "/rest/"
		u.Scheme = "http"

		q := url.Values{}

		sessionID := GetSessionId(d, meta.(*HedvigClient))

		q.Set("request", fmt.Sprintf("{type:Unmount, category:VirtualDiskManagement, params:{virtualDisk:'%s', targets:['%s']}, sessionId: '%s'}", dOldVDisk.(string), dOldController.(string), sessionID))

		u.RawQuery = q.Encode()
		log.Printf("URL: %v", u.String())

		_, err := http.Get(u.String())

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

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:Unmount, category:VirtualDiskManagement, params:{virtualDisk:'%s', targets:['%s']}, sessionId: '%s'}", d.Get("vdisk").(string), d.Get("controller").(string),
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
