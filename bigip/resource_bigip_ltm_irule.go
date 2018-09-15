package bigip

import (
	"fmt"
	"log"
	"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmIRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmIRuleCreate,
		Read:   resourceBigipLtmIRuleRead,
		Update: resourceBigipLtmIRuleUpdate,
		Delete: resourceBigipLtmIRuleDelete,
		Exists: resourceBigipLtmIRuleExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the iRule",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},

			"irule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The iRule body",
				StateFunc: func(s interface{}) string {
					return strings.TrimSpace(s.(string))
				},
			},
		},
	}
}

func resourceBigipLtmIRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Println("[INFO] Creating iRule " + name)

	err := client.CreateIRule(name, d.Get("irule").(string))
	if err != nil {
		log.Printf("[ERROR] Unable to Create Irule %s %v ", name, err)
		return err
	}

	d.SetId(name)

	return resourceBigipLtmIRuleRead(d, meta)
}

func resourceBigipLtmIRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	irule, err := client.IRule(name)
	if err != nil {
		log.Printf("[ERROR] Unbale to retrieve iRule %s: %v", name, err)
		return err
	}
	if irule == nil {
		log.Printf("[WARN] iRule (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}
	if err := d.Set("irule", irule.Rule); err != nil {
		return fmt.Errorf("Error saving iRule (%s) to state: %s", name, err)
	}

	return nil
}

func resourceBigipLtmIRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching iRule " + name)

	irule, err := client.IRule(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrive iRule (%s) (%v) ", name, err)
		return false, err
	}
	if irule == nil {
		log.Printf("[WARN] irule (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, nil
	}
	return irule != nil, nil
}

func resourceBigipLtmIRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	r := &bigip.IRule{
		FullPath: name,
		Rule:     d.Get("irule").(string),
	}

	err := client.ModifyIRule(name, r)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify iRule (%s) (%v) ", name, err)
		return err
	}
	return resourceBigipLtmIRuleRead(d, meta)
}

func resourceBigipLtmIRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	err := client.DeleteIRule(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete iRule (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}
