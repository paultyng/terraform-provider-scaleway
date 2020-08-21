package scaleway

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

func resourceScalewayLbSubscriberBeta() *schema.Resource {
	return &schema.Resource{
		Create:        resourceScalewayLbSubscriberBetaCreate,
		Read:          resourceScalewayLbSubscriberBetaRead,
		Update:        resourceScalewayLbSubscriberBetaUpdate,
		Delete:        resourceScalewayLbSubscriberBetaDelete,
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Subscriber name.",
				Optional:    true,
				Computed:    true,
			},
			"email_config": {
				ConflictsWith: []string{"webhook_config"},
				MaxItems:      1,
				Description:   "Email address configuration.",
				Type:          schema.TypeList,
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Email receiving the alert.",
						},
					},
				},
			},
			"webhook_config": {
				ConflictsWith: []string{"email_config"},
				MaxItems:      1,
				Description:   "WebHook URI configuration.",
				Type:          schema.TypeList,
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "URI receiving the POST request.",
						},
					},
				},
			},
			"organization_id": organizationIDSchema(),
		},
	}
}
func resourceScalewayLbSubscriberBetaCreate(d *schema.ResourceData, m interface{}) error {
	lbAPI, region, err := lbAPIWithRegion(d, m)
	if err != nil {
		return err
	}

	createReq := &lb.CreateSubscriberRequest{
		Region:         region,
		OrganizationID: d.Get("organization_id").(string),
		Name:           expandOrGenerateString(d.Get("name"), "lb-subscriber"),
		EmailConfig:    expandLbSubscriberEmailConfig(d.Get("email_config")),
		WebhookConfig:  expandLbSubscriberWebhookConfig(d.Get("webhook_config")),
	}

	res, err := lbAPI.CreateSubscriber(createReq)
	if err != nil {
		return err
	}

	d.SetId(newRegionalId(region, res.ID))

	return resourceScalewayLbSubscriberBetaRead(d, m)
}

func resourceScalewayLbSubscriberBetaRead(d *schema.ResourceData, m interface{}) error {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(m, d.Id())
	if err != nil {
		return err
	}

	res, err := lbAPI.GetSubscriber(&lb.GetSubscriberRequest{
		SubscriberID: ID,
		Region:       region,
	})

	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	_ = d.Set("name", res.Name)
	_ = d.Set("email_config", res.EmailConfig)
	_ = d.Set("webhook_config", res.WebhookConfig)

	return nil
}

func resourceScalewayLbSubscriberBetaUpdate(d *schema.ResourceData, m interface{}) error {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(m, d.Id())
	if err != nil {
		return err
	}

	req := &lb.UpdateSubscriberRequest{
		Region:        region,
		SubscriberID:  ID,
		Name:          d.Get("name").(string),
		EmailConfig:   expandLbSubscriberEmailConfig(d.Get("email_config")),
		WebhookConfig: expandLbSubscriberWebhookConfig(d.Get("webhook_config")),
	}

	_, err = lbAPI.UpdateSubscriber(req)
	if err != nil {
		return err
	}

	return resourceScalewayLbSubscriberBetaRead(d, m)
}

func resourceScalewayLbSubscriberBetaDelete(d *schema.ResourceData, m interface{}) error {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(m, d.Id())
	if err != nil {
		return err
	}

	err = lbAPI.DeleteSubscriber(&lb.DeleteSubscriberRequest{
		Region:       region,
		SubscriberID: ID,
	})

	if err != nil && !is404Error(err) {
		return err
	}

	return nil
}
