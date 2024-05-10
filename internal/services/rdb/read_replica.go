package rdb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/cdf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/dsf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

func ResourceReadReplica() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceRdbReadReplicaCreate,
		ReadContext:   ResourceRdbReadReplicaRead,
		UpdateContext: ResourceRdbReadReplicaUpdate,
		DeleteContext: ResourceRdbReadReplicaDelete,
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultInstanceTimeout),
			Read:    schema.DefaultTimeout(defaultInstanceTimeout),
			Update:  schema.DefaultTimeout(defaultInstanceTimeout),
			Delete:  schema.DefaultTimeout(defaultInstanceTimeout),
			Default: schema.DefaultTimeout(defaultInstanceTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Id of the rdb instance to replicate",
				DiffSuppressFunc: dsf.Locality,
			},
			"same_zone": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "Defines whether to create the replica in the same availability zone as the main instance nodes or not.",
			},
			"direct_access": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Direct access endpoint, it gives you an IP and a port to access your read-replica",
				MaxItems:    1,
				Elem:        &schema.Resource{},
			},
			"endpoints_spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Endpoints specs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Endpoints common specs
						"endpoint_id": {
							Type:        schema.TypeString,
							Description: "UUID of the endpoint (UUID format).",
							Computed:    true,
						},
						"ip": {
							Type:        schema.TypeString,
							Description: "IPv4 address of the endpoint (IP address). Only one of ip and hostname may be set.",
							Computed:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "TCP port of the endpoint.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the endpoint.",
							Computed:    true,
						},
						"hostname": {
							Type:        schema.TypeString,
							Description: "Hostname of the endpoint. Only one of ip and hostname may be set.",
							Computed:    true,
						},
					},
				},
			},
			"private_network": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Private network endpoints",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Private network specific
						"private_network_id": {
							Type:             schema.TypeString,
							Description:      "UUID of the private network to be connected to the read replica (UUID format)",
							ValidateFunc:     verify.IsUUIDorUUIDWithLocality(),
							DiffSuppressFunc: dsf.Locality,
							Required:         true,
						},
						"service_ip": {
							Type:         schema.TypeString,
							Description:  "The IP network address within the private subnet",
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"enable_ipam": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether or not the private network endpoint should be configured with IPAM",
						},
						"zone": {
							Type:        schema.TypeString,
							Description: "Private network zone",
							Computed:    true,
						},
					},
				},
			},
			// Common
			"region": regional.Schema(),
		},
		CustomizeDiff: cdf.LocalityCheck("instance_id", "private_network.#.private_network_id"),
	}
}

func ResourceRdbReadReplicaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rdbAPI, region, err := newAPIWithRegion(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	endpointSpecs := []*rdb.ReadReplicaEndpointSpec(nil)
	if directAccess := expandReadReplicaEndpointsSpecDirectAccess(d.Get("direct_access")); directAccess != nil {
		endpointSpecs = append(endpointSpecs, directAccess)
	}

	ipamConfig, staticConfig := getIPConfigCreate(d, "service_ip")
	if pn, diags := expandReadReplicaEndpointsSpecPrivateNetwork(d.Get("private_network"), ipamConfig, staticConfig); err != nil || pn != nil {
		if diags.HasError() {
			return diags
		}
		for _, warning := range diags {
			tflog.Warn(ctx, warning.Detail)
		}
		endpointSpecs = append(endpointSpecs, pn)
	}

	rr, err := rdbAPI.CreateReadReplica(&rdb.CreateReadReplicaRequest{
		Region:       region,
		InstanceID:   locality.ExpandID(d.Get("instance_id")),
		EndpointSpec: endpointSpecs,
		SameZone:     types.ExpandBoolPtr(d.Get("same_zone")),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create read-replica: %w", err))
	}

	d.SetId(regional.NewIDString(region, rr.ID))

	_, err = waitForRDBReadReplica(ctx, rdbAPI, region, rr.ID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceRdbReadReplicaRead(ctx, d, m)
}

func ResourceRdbReadReplicaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rdbAPI, region, ID, err := NewAPIWithRegionAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rr, err := waitForRDBReadReplica(ctx, rdbAPI, region, ID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	enableIpam, err := getIPAMConfigRead(rr, m)
	if err != nil {
		return diag.FromErr(err)
	}
	directAccess, privateNetwork := flattenReadReplicaEndpoints(rr.Endpoints, enableIpam)
	_ = d.Set("direct_access", directAccess)
	_ = d.Set("private_network", privateNetwork)

	regionStr := region.String()
	_ = d.Set("same_zone", rr.SameZone)
	_ = d.Set("region", regionStr)
	_ = d.Set("instance_id", regional.NewIDString(region, rr.InstanceID))

	return nil
}

//gocyclo:ignore
func ResourceRdbReadReplicaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rdbAPI, region, ID, err := NewAPIWithRegionAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// verify resource is ready
	rr, err := waitForRDBReadReplica(ctx, rdbAPI, region, ID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	newEndpoints := []*rdb.ReadReplicaEndpointSpec(nil)

	if d.HasChange("direct_access") {
		// delete old endpoint
		for _, e := range rr.Endpoints {
			if e.DirectAccess != nil {
				err := rdbAPI.DeleteEndpoint(&rdb.DeleteEndpointRequest{
					Region:     region,
					EndpointID: e.ID,
				}, scw.WithContext(ctx))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		// retrieve state
		rr, err = waitForRDBReadReplica(ctx, rdbAPI, region, ID, d.Timeout(schema.TimeoutRead))
		if err != nil {
			return diag.FromErr(err)
		}
		// create a new one if defined
		if directAccess, directAccessExists := d.GetOk("direct_access"); directAccessExists {
			newEndpoints = append(newEndpoints, expandReadReplicaEndpointsSpecDirectAccess(directAccess))
		}
	}

	if d.HasChange("private_network") {
		// delete old endpoint
		for _, e := range rr.Endpoints {
			if e.PrivateNetwork != nil {
				err := rdbAPI.DeleteEndpoint(&rdb.DeleteEndpointRequest{
					Region:     region,
					EndpointID: e.ID,
				}, scw.WithContext(ctx))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		// retrieve state
		_, err = waitForRDBReadReplica(ctx, rdbAPI, region, ID, d.Timeout(schema.TimeoutRead))
		if err != nil {
			return diag.FromErr(err)
		}
		// create a new one if defined
		if pn, pnExists := d.GetOk("private_network"); pnExists {
			ipamConfig, staticConfig := getIPConfigUpdate(d, "service_ip")
			pnEndpoint, diags := expandReadReplicaEndpointsSpecPrivateNetwork(pn, ipamConfig, staticConfig)
			if diags.HasError() {
				return diags
			}
			for _, warning := range diags {
				tflog.Warn(ctx, warning.Detail)
			}
			newEndpoints = append(newEndpoints, pnEndpoint)
		}
	}

	if len(newEndpoints) > 0 {
		_, err := waitForRDBReadReplica(ctx, rdbAPI, region, ID, d.Timeout(schema.TimeoutRead))
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = rdbAPI.CreateReadReplicaEndpoint(&rdb.CreateReadReplicaEndpointRequest{
			Region:        region,
			ReadReplicaID: ID,
			EndpointSpec:  newEndpoints,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	_, err = waitForRDBReadReplica(ctx, rdbAPI, region, ID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceRdbReadReplicaRead(ctx, d, m)
}

func ResourceRdbReadReplicaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rdbAPI, region, ID, err := NewAPIWithRegionAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// We first wait in case the instance is in a transient state
	_, err = waitForRDBReadReplica(ctx, rdbAPI, region, ID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = rdbAPI.DeleteReadReplica(&rdb.DeleteReadReplicaRequest{
		Region:        region,
		ReadReplicaID: ID,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	// Lastly wait in case the instance is in a transient state
	_, err = waitForRDBReadReplica(ctx, rdbAPI, region, ID, d.Timeout(schema.TimeoutDelete))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	return nil
}
