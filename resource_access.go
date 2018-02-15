package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

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
	d.SetId("access-" + d.Get("vdisk").(string))

	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:PersistACLAccess, category:VirtualDiskManagement, params:{virtualDisks:['%s'], host:'%s', address:'%s', type:'%s'}, sessionId:'%s'}", d.Get("vdisk"), d.Get("host"), d.Get("address"),
		d.Get("type"), sessionID))
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

func resourceAccessRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAccessUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAccessDelete(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}

	sessionID := GetSessionId(d, meta.(*HedvigClient))

	q.Set("request", fmt.Sprintf("{type:RemoveACLAccess, category:VirtualDiskManagement, params:{virtualDisk:'%s', host:'%s', address:['%s']}, sessionId: '%s'}", d.Get("vdisk"), d.Get("host"), d.Get("address"),
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
