package hedvig

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type AccessResponse struct {
	RequestID string `json:"requestId"`
	Result    []struct {
		Host      string `json:"host"`
		Initiator []struct {
			Ip   string `json:"ip"`
			Name string `json:"name"`
		}
	} `json:"result"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

func resourceAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessCreate,
		Read:   resourceAccessRead,
		Update: resourceAccessUpdate,
		Delete: resourceAccessDelete,

		Schema: map[string]*schema.Schema{
			"cluster": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vdisk": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAccessCreate(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:PersistACLAccess, category:VirtualDiskManagement, params:{virtualDisks:['%s'], host:'%s', address:'%s', type:'%s'}, sessionId:'%s'}", d.Get("vdisk").(string), d.Get("host").(string), d.Get("address").(string),
		d.Get("type").(string), sessionID))
	u.RawQuery = q.Encode()
	log.Printf("URL: %v", u.String())

	resp, err := http.Get(u.String())

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("body: %s", body)

        d.SetId("access-" + d.Get("vdisk").(string) + "-" + d.Get("address").(string))

	return resourceAccessRead(d, meta)
}

func resourceAccessRead(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:GetACLInformation,category:VirtualDiskManagement,params:{virtualDisk:'%s'},sessionId:'%s'}", d.Get("vdisk").(string), sessionID))

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
        if resp.StatusCode == 404 {
                d.SetId("")
                log.Fatal(err)
        }
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	access := AccessResponse{}
	err = json.Unmarshal(body, &access)

	if err != nil {
		log.Fatalf("Error unmarshalling: %s", err)
	}

        if len(access.Result) < 1 {
                log.Fatal("Incorrect Array Size")
        }

	d.Set("host", access.Result[0].Host)

	return nil
}

func resourceAccessUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("address") || d.HasChange("host") || d.HasChange("cluster") || d.HasChange("vdisk") || d.HasChange("type") {
		dOldDisk, _ := d.GetChange("vdisk")
		dOldHost, _ := d.GetChange("host")
		dOldAddress, _ := d.GetChange("address")

		u := url.URL{}
		u.Host = meta.(*HedvigClient).Node
		u.Path = "/rest/"
		u.Scheme = "http"

		q := url.Values{}

		sessionID := GetSessionId(d, meta.(*HedvigClient))

		log.Printf("dOldDisk: %s", dOldDisk.(string))
		log.Printf("dOldHost: %s", dOldHost.(string))
		log.Printf("dOldAddress: %s", dOldAddress.(string))

		q.Set("request", fmt.Sprintf("{type:RemoveACLAccess, category:VirtualDiskManagement, params:{virtualDisk:'%s', host:'%s', address:['%s']}, sessionId: '%s'}", dOldDisk.(string), dOldHost.(string), dOldAddress.(string), sessionID))
		u.RawQuery = q.Encode()
		log.Printf("URL: %v", u.String())

		resp, err := http.Get(u.String())

		if err != nil {
			log.Fatal(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("body: %s", body)

		resourceAccessCreate(d, meta)
	}

	return resourceAccessRead(d, meta)
}

func resourceAccessDelete(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:RemoveACLAccess, category:VirtualDiskManagement, params:{virtualDisk:'%s', host:'%s', address:['%s']}, sessionId: '%s'}", d.Get("vdisk").(string), d.Get("host").(string), d.Get("address").(string),
		sessionID))
	u.RawQuery = q.Encode()
	log.Printf("URL: %v", u.String())

	resp, err := http.Get(u.String())

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("body: %s", body)

	return nil
}
