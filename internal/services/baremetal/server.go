package baremetal

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	ipamAPI "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	sdkValidation "github.com/scaleway/scaleway-sdk-go/validation"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/cdf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/dsf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

func ResourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceServerCreate,
		ReadContext:   ResourceServerRead,
		UpdateContext: ResourceServerUpdate,
		DeleteContext: ResourceServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultServerTimeout),
			Create:  schema.DefaultTimeout(defaultServerTimeout),
			Update:  schema.DefaultTimeout(defaultServerTimeout),
			Delete:  schema.DefaultTimeout(defaultServerTimeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the server",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hostname of the server",
			},
			"offer": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID or name of the server offer",
				DiffSuppressFunc: func(_, oldValue, newValue string, d *schema.ResourceData) bool {
					// remove the locality from the IDs when checking diff
					if locality.ExpandID(newValue) == locality.ExpandID(oldValue) {
						return true
					}
					// if the offer was provided by name
					offerName, ok := d.GetOk("offer_name")
					return ok && newValue == offerName
				},
			},
			"offer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the server offer",
			},
			"offer_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the server offer",
			},
			"os": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The base image of the server",
				DiffSuppressFunc: dsf.Locality,
				ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
			},
			"os_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The base image name of the server",
			},
			"ssh_key_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: verify.IsUUID(),
				},
				Optional: true,
				Description: `Array of SSH key IDs allowed to SSH to the server

**NOTE** : If you are attempting to update your SSH key IDs, it will induce the reinstall of your server. 
If this behaviour is wanted, please set 'reinstall_on_ssh_key_changes' argument to true.`,
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User used for the installation.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password used for the installation.",
			},
			"service_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User used for the service to install.",
			},
			"service_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password used for the service to install.",
			},
			"reinstall_on_config_changes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If True, this boolean allows to reinstall the server on SSH key IDs, user or password changes",
			},
			"install_config_afterward": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If True, this boolean allows to create a server without the install config if you want to provide it later",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
				Description:  "Some description to associate to the server, max 255 characters",
			},
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				Description: "Array of tags to associate with the server",
			},
			"zone":            zonal.Schema(),
			"organization_id": account.OrganizationIDSchema(),
			"project_id":      account.ProjectIDSchema(),
			"ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IP addresses attached to the server.",
				Elem:        ResourceServerIP(),
			},
			"ipv4": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IPv4 addresses attached to the server",
				Elem:        ResourceServerIP(),
			},
			"ipv6": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IPv6 addresses attached to the server",
				Elem:        ResourceServerIP(),
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"options": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The options to enable on server",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "IDs of the options",
							Required:    true,
						},
						"expires_at": {
							Type:             schema.TypeString,
							Description:      "Auto expire the option after this date",
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: verify.IsDate(),
							DiffSuppressFunc: dsf.TimeRFC3339,
						},
						// computed
						"name": {
							Type:        schema.TypeString,
							Description: "name of the option",
							Computed:    true,
						},
					},
				},
			},
			"private_network": {
				Type:        schema.TypeSet,
				Optional:    true,
				Set:         privateNetworkSetHash,
				Description: "The private networks to attach to the server",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:             schema.TypeString,
							Description:      "The private network ID",
							Required:         true,
							ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
							StateFunc: func(i interface{}) string {
								return locality.ExpandID(i.(string))
							},
						},
						// computed
						"vlan": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The VLAN ID associated to the private network",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private network status",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time of the creation of the private network",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time of the last update of the private network",
						},
					},
				},
			},
			"private_ip": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of private IP addresses associated with the resource",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the IP address resource",
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private IP address",
						},
					},
				},
			},
		},
		CustomizeDiff: customdiff.Sequence(
			cdf.LocalityCheck("private_network.#.id"),
			customDiffPrivateNetworkOption(),
		),
	}
}

func ResourceServerIP() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the IPv6",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the IPv6",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv6 address",
			},
			"reverse": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Reverse of the IPv6",
			},
		},
	}
}

func ResourceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zone, err := newAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	privateNetworkAPI, _, err := newPrivateNetworkAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	offerID := zonal.ExpandID(d.Get("offer"))
	if !sdkValidation.IsUUID(offerID.ID) {
		o, err := api.GetOfferByName(&baremetal.GetOfferByNameRequest{
			OfferName: offerID.ID,
			Zone:      zone,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		offerID = zonal.NewID(zone, o.ID)
	}

	if !d.Get("install_config_afterward").(bool) {
		if diags := validateInstallConfig(ctx, d, m); len(diags) > 0 {
			return diags
		}
	}

	server, err := api.CreateServer(&baremetal.CreateServerRequest{
		Zone:        zone,
		Name:        types.ExpandOrGenerateString(d.Get("name"), "bm"),
		ProjectID:   types.ExpandStringPtr(d.Get("project_id")),
		Description: d.Get("description").(string),
		OfferID:     offerID.ID,
		Tags:        types.ExpandStrings(d.Get("tags")),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zonal.NewID(server.Zone, server.ID).String())

	_, err = waitForServer(ctx, api, zone, server.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.Get("install_config_afterward").(bool) {
		_, err = api.InstallServer(&baremetal.InstallServerRequest{
			Zone:            server.Zone,
			ServerID:        server.ID,
			OsID:            zonal.ExpandID(d.Get("os")).ID,
			Hostname:        types.ExpandStringWithDefault(d.Get("hostname"), server.Name),
			SSHKeyIDs:       types.ExpandStrings(d.Get("ssh_key_ids")),
			User:            types.ExpandStringPtr(d.Get("user")),
			Password:        types.ExpandStringPtr(d.Get("password")),
			ServiceUser:     types.ExpandStringPtr(d.Get("service_user")),
			ServicePassword: types.ExpandStringPtr(d.Get("service_password")),
		}, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = waitForServerInstall(ctx, api, zone, server.ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	options, optionsExist := d.GetOk("options")
	if optionsExist {
		opSpecs, err := expandOptions(options)
		if err != nil {
			return diag.FromErr(err)
		}
		for i := range opSpecs {
			_, err = api.AddOptionServer(&baremetal.AddOptionServerRequest{
				Zone:      server.Zone,
				ServerID:  server.ID,
				OptionID:  opSpecs[i].ID,
				ExpiresAt: opSpecs[i].ExpiresAt,
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	privateNetworkIDs, pnExist := d.GetOk("private_network")
	if pnExist {
		createBaremetalPrivateNetworkRequest := &baremetal.PrivateNetworkAPISetServerPrivateNetworksRequest{
			Zone:              zone,
			ServerID:          server.ID,
			PrivateNetworkIDs: expandPrivateNetworks(privateNetworkIDs),
		}

		baremetalPrivateNetwork, err := privateNetworkAPI.SetServerPrivateNetworks(
			createBaremetalPrivateNetworkRequest,
			scw.WithContext(ctx),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = waitForServerPrivateNetwork(ctx, privateNetworkAPI, zone, baremetalPrivateNetwork.ServerPrivateNetworks[0].ServerID, d.Timeout(schema.TimeoutCreate))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}
	}

	return ResourceServerRead(ctx, d, m)
}

func ResourceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zonedID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	privateNetworkAPI, _, err := newPrivateNetworkAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	server, err := api.GetServer(&baremetal.GetServerRequest{
		Zone:     zonedID.Zone,
		ServerID: zonedID.ID,
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	offer, err := api.GetOffer(&baremetal.GetOfferRequest{
		Zone:    server.Zone,
		OfferID: server.OfferID,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	var os *baremetal.OS
	if server.Install != nil {
		os, err = api.GetOS(&baremetal.GetOSRequest{
			Zone: server.Zone,
			OsID: server.Install.OsID,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	_ = d.Set("name", server.Name)
	_ = d.Set("zone", server.Zone.String())
	_ = d.Set("organization_id", server.OrganizationID)
	_ = d.Set("project_id", server.ProjectID)
	_ = d.Set("offer_id", zonal.NewIDString(server.Zone, offer.ID))
	_ = d.Set("offer_name", offer.Name)
	_ = d.Set("offer", zonal.NewIDString(server.Zone, offer.ID))
	_ = d.Set("tags", server.Tags)
	_ = d.Set("domain", server.Domain)
	_ = d.Set("ips", flattenIPs(server.IPs))
	_ = d.Set("ipv4", flattenIPv4s(server.IPs))
	_ = d.Set("ipv6", flattenIPv6s(server.IPs))
	if server.Install != nil {
		_ = d.Set("os", zonal.NewIDString(server.Zone, os.ID))
		_ = d.Set("os_name", os.Name)
		_ = d.Set("ssh_key_ids", server.Install.SSHKeyIDs)
		_ = d.Set("user", server.Install.User)
		_ = d.Set("service_user", server.Install.ServiceUser)
	}
	_ = d.Set("description", server.Description)
	_ = d.Set("options", flattenOptions(server.Zone, server.Options))

	listPrivateNetworks, err := privateNetworkAPI.ListServerPrivateNetworks(&baremetal.PrivateNetworkAPIListServerPrivateNetworksRequest{
		Zone:     server.Zone,
		ServerID: &server.ID,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to list server's private networks: %w", err))
	}

	pnRegion, err := server.Zone.Region()
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("private_network", flattenPrivateNetworks(pnRegion, listPrivateNetworks.ServerPrivateNetworks))

	privateNetworkIDs := make([]string, 0, len(listPrivateNetworks.ServerPrivateNetworks))
	for _, pn := range listPrivateNetworks.ServerPrivateNetworks {
		privateNetworkIDs = append(privateNetworkIDs, pn.PrivateNetworkID)
	}

	var allPrivateIPs []map[string]interface{}
	for _, privateNetworkID := range privateNetworkIDs {
		resourceType := ipamAPI.ResourceTypeBaremetalPrivateNic
		opts := &ipam.GetResourcePrivateIPsOptions{
			ResourceType:     &resourceType,
			PrivateNetworkID: &privateNetworkID,
		}
		privateIPs, err := ipam.GetResourcePrivateIPs(ctx, m, pnRegion, opts)
		if err != nil {
			return diag.FromErr(err)
		}
		if privateIPs != nil {
			allPrivateIPs = append(allPrivateIPs, privateIPs...)
		}
	}
	_ = d.Set("private_ip", allPrivateIPs)

	return nil
}

//gocyclo:ignore
func ResourceServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zonedID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	privateNetworkAPI, zone, err := newPrivateNetworkAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	server, err := api.GetServer(&baremetal.GetServerRequest{
		Zone:     zonedID.Zone,
		ServerID: zonedID.ID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var serverGetOptionIDs []*baremetal.ServerOption
	serverGetOptionIDs = append(serverGetOptionIDs, server.Options...)

	if d.HasChange("options") {
		options, err := expandOptions(d.Get("options"))
		if err != nil {
			return diag.FromErr(err)
		}
		optionsToDelete := compareOptions(options, serverGetOptionIDs)
		for i := range optionsToDelete {
			_, err = api.DeleteOptionServer(&baremetal.DeleteOptionServerRequest{
				Zone:     server.Zone,
				ServerID: server.ID,
				OptionID: optionsToDelete[i].ID,
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}

		_, err = waitForServerOptions(ctx, api, zonedID.Zone, zonedID.ID, d.Timeout(schema.TimeoutDelete))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}

		optionsToAdd := compareOptions(serverGetOptionIDs, options)
		for i := range optionsToAdd {
			_, err = api.AddOptionServer(&baremetal.AddOptionServerRequest{
				Zone:      server.Zone,
				ServerID:  server.ID,
				OptionID:  optionsToAdd[i].ID,
				ExpiresAt: optionsToAdd[i].ExpiresAt,
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("private_network") {
		privateNetworkIDs := d.Get("private_network")

		updateBaremetalPrivateNetworkRequest := &baremetal.PrivateNetworkAPISetServerPrivateNetworksRequest{
			Zone:              zone,
			ServerID:          server.ID,
			PrivateNetworkIDs: expandPrivateNetworks(privateNetworkIDs),
		}

		baremetalPrivateNetwork, err := privateNetworkAPI.SetServerPrivateNetworks(
			updateBaremetalPrivateNetworkRequest,
			scw.WithContext(ctx),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = waitForServerPrivateNetwork(ctx, privateNetworkAPI, zone, baremetalPrivateNetwork.ServerPrivateNetworks[0].ServerID, d.Timeout(schema.TimeoutUpdate))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}
	}

	req := &baremetal.UpdateServerRequest{
		Zone:     zonedID.Zone,
		ServerID: zonedID.ID,
	}

	hasChanged := false

	if d.HasChange("name") {
		req.Name = types.ExpandUpdatedStringPtr(d.Get("name"))
		hasChanged = true
	}

	if d.HasChange("description") {
		req.Description = types.ExpandUpdatedStringPtr(d.Get("description"))
		hasChanged = true
	}

	if d.HasChange("tags") {
		req.Tags = types.ExpandUpdatedStringsPtr(d.Get("tags"))
		hasChanged = true
	}

	if hasChanged {
		_, err = api.UpdateServer(req, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	installReq := &baremetal.InstallServerRequest{
		Zone:            zonedID.Zone,
		ServerID:        zonedID.ID,
		Hostname:        types.ExpandStringWithDefault(d.Get("hostname"), d.Get("name").(string)),
		SSHKeyIDs:       types.ExpandStrings(d.Get("ssh_key_ids")),
		User:            types.ExpandStringPtr(d.Get("user")),
		Password:        types.ExpandStringPtr(d.Get("password")),
		ServiceUser:     types.ExpandStringPtr(d.Get("service_user")),
		ServicePassword: types.ExpandStringPtr(d.Get("service_password")),
	}

	if d.HasChange("os") {
		if diags := validateInstallConfig(ctx, d, m); len(diags) > 0 {
			return diags
		}
		err = installServer(ctx, d, api, installReq)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = waitForServerInstall(ctx, api, zonedID.Zone, zonedID.ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	var diags diag.Diagnostics

	if d.HasChanges("ssh_key_ids", "user", "password", "reinstall_on_config_changes") {
		if !d.Get("reinstall_on_config_changes").(bool) && !d.HasChange("os") {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Changes have been made on your config",
				Detail: "[WARN] This change induce the reinstall of your server. " +
					"If this behaviour is wanted, please set 'reinstall_on_config_changes' argument to true",
			})
		} else {
			if diags := validateInstallConfig(ctx, d, m); len(diags) > 0 {
				return diags
			}
			err = installServer(ctx, d, api, installReq)
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = waitForServerInstall(ctx, api, zonedID.Zone, zonedID.ID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return append(diags, ResourceServerRead(ctx, d, m)...)
}

func ResourceServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zonedID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = detachAllPrivateNetworkFromServer(ctx, d, m, zonedID.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = api.DeleteServer(&baremetal.DeleteServerRequest{
		Zone:     zonedID.Zone,
		ServerID: zonedID.ID,
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	_, err = waitForServer(ctx, api, zonedID.Zone, zonedID.ID, d.Timeout(schema.TimeoutDelete))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	return nil
}

func installAttributeMissing(field *baremetal.OSOSField, d *schema.ResourceData, attribute string) bool {
	if field != nil && field.Required && field.DefaultValue == nil {
		if _, attributeExists := d.GetOk(attribute); !attributeExists {
			return true
		}
	}
	return false
}

// validateInstallConfig validates that schema contains attribute required for OS install
func validateInstallConfig(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	baremetalAPI, zone, err := newAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	os, err := baremetalAPI.GetOS(&baremetal.GetOSRequest{
		Zone: zone,
		OsID: locality.ExpandID(d.Get("os")),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	diags := diag.Diagnostics(nil)
	installAttributes := []struct {
		Attribute string
		Field     *baremetal.OSOSField
	}{
		{
			"user",
			os.User,
		},
		{
			"password",
			os.Password,
		},
		{
			"service_user",
			os.ServiceUser,
		},
		{
			"service_password",
			os.ServicePassword,
		},
	}
	for _, installAttr := range installAttributes {
		if installAttributeMissing(installAttr.Field, d, installAttr.Attribute) {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       installAttr.Attribute + " attribute is required",
				Detail:        installAttr.Attribute + " is required for this os",
				AttributePath: cty.GetAttrPath(installAttr.Attribute),
			})
		}
	}
	return diags
}
