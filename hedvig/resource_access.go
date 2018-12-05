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

type readAccessResponse struct {
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

type deleteAccessResponse struct {
	RequestId string `json:"requestId"`
	Status    string `json:"status"`
	Type      string `json:"type"`
}

type createAccessResponse struct {
	Result []struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"result"`
	RequestId string `json:"requestId"`
	Status    string `json:"status"`
	Type      string `json:"type"`
}

func resourceAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessCreate,
		Read:   resourceAccessRead,
		Delete: resourceAccessDelete,

		Schema: map[string]*schema.Schema{
			"vdisk": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

	sessionID, err := GetSessionId(d, meta.(*HedvigClient))

	if err != nil {
		return err
	}

	q.Set("request", fmt.Sprintf("{type:PersistACLAccess, category:VirtualDiskManagement, params:{virtualDisks:['%s'], host:'%s', address:'%s', type:'%s'}, sessionId:'%s'}", d.Get("vdisk").(string), d.Get("host").(string), d.Get("address").(string),
		d.Get("type").(string), sessionID))
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

	createResp := createAccessResponse{}

	err = json.Unmarshal(body, &createResp)

	log.Printf("body: %s", body)

	d.SetId("access$" + d.Get("vdisk").(string) + "$" + d.Get("host").(string) + "$" + d.Get("address").(string))

	return resourceAccessRead(d, meta)
}

func resourceAccessRead(d *schema.ResourceData, meta interface{}) error {
	u := url.URL{}
	u.Host = meta.(*HedvigClient).Node
	u.Path = "/rest/"
	u.Scheme = "http"

	sessionID, err := GetSessionId(d, meta.(*HedvigClient))

	if err != nil {
		return err
	}

	idSplit := strings.Split(d.Id(), "$")

	if len(idSplit) != 4 {
		return errors.New("Invalid ID: " + d.Id())
	}

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:GetACLInformation,category:VirtualDiskManagement,params:{virtualDisk:'%s'},sessionId:'%s'}", idSplit[1], sessionID))

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		d.SetId("")
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	access := readAccessResponse{}
	err = json.Unmarshal(body, &access)

	if err != nil {
		return err
	}

	if len(access.Result) < 1 {
		return errors.New("Not enough results to find host in")
	}

	d.Set("host", access.Result[0].Host)

	return nil
}

func resourceAccessDelete(d *schema.ResourceData, meta interface{}) error {
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

	if len(idSplit) != 4 {
		return errors.New("Invalid ID: " + d.Id())
	}

	q.Set("request", fmt.Sprintf("{type:RemoveACLAccess, category:VirtualDiskManagement, params:{virtualDisk:'%s', host:'%s', address:['%s']}, sessionId: '%s'}", idSplit[1], idSplit[2], idSplit[3], sessionID))
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

	deleteResp := deleteAccessResponse{}

	err = json.Unmarshal(body, &deleteResp)

	if err != nil {
		return err
	}
	// TODO: Verify
	d.SetId("")

	log.Printf("body: %s", body)

	return nil
}
