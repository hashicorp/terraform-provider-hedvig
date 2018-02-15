package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

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
	d.SetId("id-" + d.Get("name").(string))

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("body: %s", body)

	return nil
}

func resourceVdiskRead(d *schema.ResourceData, m interface{}) error {
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

	return nil
}

func resourceVdiskDelete(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:DeleteVDisk, category:VirtualDiskManagement, params:{virtualDisks:['%s']}, sessionId:'%s'}, sessionId:'%s'}", d.Get("name").(string),
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
