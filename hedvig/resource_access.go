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
			IP   string `json:"ip"`
			Name string `json:"name"`
		}
	} `json:"result"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type deleteAccessResponse struct {
	RequestID string `json:"requestId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Type      string `json:"type"`
}

type createAccessResponse struct {
	Result []struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"result"`
	RequestID string `json:"requestId"`
	Status    string `json:"status"`
	Type      string `json:"type"`
}

func resourceAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessCreate,
		Read:   resourceAccessRead,
		Delete: resourceAccessDelete,

		Schema: map[string]*schema.Schema{
			"vdisk": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
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

	q.Set("request", fmt.Sprintf("{type:PersistACLAccess, category:VirtualDiskManagement, params:{virtualDisks:['%s'], host:'%s', address:'%s', type:'%s'}, sessionId:'%s'}", d.Get("vdisk").(string), d.Get("host").(string), d.Get("address").(string), d.Get("type").(string), sessionID))
	u.RawQuery = q.Encode()
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
	if err != nil {
		return err
	}

	if len(createResp.Result) < 1 {
		return errors.New("Insufficient results from search")
	}

	if createResp.Result[0].Status != "ok" {
		return fmt.Errorf("Error creating access: %s", createResp.Result[0].Message)
	}
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
		return fmt.Errorf("Invalid ID: %s", d.Id())
	}

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:GetACLInformation,category:VirtualDiskManagement,params:{virtualDisk:'%s'},sessionId:'%s'}", idSplit[1], sessionID))

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		return errors.New("Malformed query; aborting")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	readAccess := readAccessResponse{}
	err = json.Unmarshal(body, &readAccess)

	if err != nil {
		return err
	}

	if readAccess.Status == "warning" && strings.HasSuffix(readAccess.Message, "t be found") {
		d.SetId("")
		log.Print("Access resource not found for vdisk, clearing from state")
		return nil
	}

	for _, rec := range readAccess.Result {
		if rec.Host == idSplit[2] {
			for _, export := range rec.Initiator {
				if export.IP == idSplit[3] {
					d.Set("host", rec.Host)
					d.Set("address", export.IP)
					return nil
				}
			}
		}
	}

	return errors.New("Could not find address associated with host")
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
		return fmt.Errorf("Invalid ID: %s", d.Id())
	}

	q.Set("request", fmt.Sprintf("{type:RemoveACLAccess, category:VirtualDiskManagement, params:{virtualDisk:'%s', host:'%s', address:['%s']}, sessionId: '%s'}", idSplit[1], idSplit[2], idSplit[3], sessionID))
	u.RawQuery = q.Encode()

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

	if deleteResp.Status != "ok" {
		return fmt.Errorf("Error removing access: %s", deleteResp.Message)
	}
	return nil
}
