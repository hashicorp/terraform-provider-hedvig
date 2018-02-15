package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func resourceLun() *schema.Resource {
	return &schema.Resource{
		Create: resourceLunCreate,
		Read:   resourceLunRead,
		Update: resourceLunUpdate,
		Delete: resourceLunDelete,

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

func resourceLunCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId("lun-" + d.Get("vdisk").(string))

	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:AddLun, category:VirtualDiskManagement, params:{virtualDisks:['%s'], targets:['%s'], readonly:false}, sessionId:'%s'}", d.Get("vdisk"), d.Get("controller"),
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

func resourceLunRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceLunUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceLunDelete(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:UnmapLun, category:VirtualDiskManagement, params:{virtualDisk:'%s', target:'%s'}, sessionId: '%s'}", d.Get("vdisk"), d.Get("controller"),
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
