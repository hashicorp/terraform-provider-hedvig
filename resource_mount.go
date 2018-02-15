package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func resourceMount() *schema.Resource {
	return &schema.Resource{
		Create: resourceMountCreate,
		Read:   resourceMountRead,
		Update: resourceMountUpdate,
		Delete: resourceMountDelete,

		Schema: map[string]*schema.Schema{
			"vdisk": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"controller": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceMountCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId("lun-" + d.Get("vdisk").(string))

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
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("body: %s", body)

	return nil
}

func resourceMountRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceMountUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMountDelete(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:Unmount, category:VirtualDiskManagement, params:{virtualDisk:'%s', targets:['%s']}, sessionId: '%s'}", d.Get("vdisk"), d.Get("controller"),
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
