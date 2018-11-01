package hedvig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	//"os"
)

type DiskResponse struct {
	Result struct {
		VDiskName string `json:"vDiskName"`
		Size      struct {
			Units string `json:"units"`
			Value int    `json:"value"`
		} `json:"size"`
	} `json:"result"`
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
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"cluster": &schema.Schema{
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

func resourceVdiskCreate(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:AddVirtualDisk, category:VirtualDiskManagement, params:{name:'%s', size:{unit:'GB', value:%d}, diskType:%s, scsi3pr:false}, sessionId:'%s'}", d.Get("name").(string), d.Get("size").(int), d.Get("type").(string),
		sessionID))
	u.RawQuery = q.Encode()
	log.Printf("URL: %v", u.String())

	resp, err := http.Get(u.String())

	if err != nil {
		log.Fatal(err)
	}

        if res.StatusCode == 404 {
                d.SetId("")
                log.Fatal(res.StatusCode)
        }

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("body: %s", body)

        d.SetId("id-" + d.Get("name").(string))

	return resourceVdiskRead(d, meta)
}

func resourceVdiskRead(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:VirtualDiskDetails,category:VirtualDiskManagement,params:{virtualDisk:'%s'},sessionId:'%s'}", d.Get("name").(string), sessionID))

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

	disk := DiskResponse{}
	err = json.Unmarshal(body, &disk)

	if err != nil {
		//os.Stderr.WriteString(string(body))
		log.Fatalf("Error unmarshalling: %s :: %s", err, string(body))
	}

	d.Set("name", disk.Result.VDiskName)
	d.Set("size", disk.Result.Size.Value)

	return nil
}

func resourceVdiskUpdate(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	if d.HasChange("size") {
		q.Set("request", fmt.Sprintf("{type:ResizeDisks, category:VirtualDiskManagement, params:{virtualDisks:['%s'], size:{unit:'GB', value:%d}}, sessionId:'%s'}", d.Get("name").(string), d.Get("size").(int),
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

	}

	return resourceVdiskRead(d, meta)
}

func resourceVdiskDelete(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:DeleteVDisk, category:VirtualDiskManagement, params:{virtualDisks:['%s']}, sessionId:'%s'}, sessionId:'%s'}", d.Get("name").(string), sessionID,
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
