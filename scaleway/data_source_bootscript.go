package scaleway

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	api "github.com/nicolai86/scaleway-sdk"
)

func dataSourceScalewayBootscript() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is deprecated and will be removed in the next major version",

		Read: dataSourceScalewayBootscriptRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "exact name of the desired bootscript",
			},
			"name_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "partial name of the desired bootscript to filter with",
			},
			"architecture": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "architecture of the desired bootscript",
			},
			// Computed values.
			"organization": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "organization owning the bootscript",
			},
			"public": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "indication if the bootscript is public",
			},
			"boot_cmd_args": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "commandline boot options used",
			},
			"dtb": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "path to Device Tree Blob detailing hardware information",
			},
			"initrd": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to initial ramdisk content",
			},
			"kernel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to used kernel",
			},
		},
	}
}

func bootscriptDescriptionAttributes(d *schema.ResourceData, script api.Bootscript) error {
	_ = d.Set("architecture", script.Arch)
	_ = d.Set("organization", script.Organization)
	_ = d.Set("public", script.Public)
	_ = d.Set("boot_cmd_args", script.Bootcmdargs)
	_ = d.Set("dtb", script.Dtb)
	_ = d.Set("initrd", script.Initrd)
	_ = d.Set("kernel", script.Kernel)
	d.SetId(script.Identifier)

	return nil
}

func dataSourceScalewayBootscriptRead(d *schema.ResourceData, meta interface{}) error {
	scaleway := meta.(*Meta).deprecatedClient

	scripts, err := scaleway.GetBootscripts()
	if err != nil {
		return err
	}

	isMatch := func(s api.Bootscript) bool { return true }

	architecture := d.Get("architecture")
	if name, ok := d.GetOk("name"); ok {
		isMatch = func(s api.Bootscript) bool {
			architectureMatch := true
			if architecture != "" {
				architectureMatch = architecture == s.Arch
			}
			return s.Title == name.(string) && architectureMatch
		}
	} else if nameFilter, ok := d.GetOk("name_filter"); ok {
		exp, err := regexp.Compile(nameFilter.(string))
		if err != nil {
			return err
		}

		isMatch = func(s api.Bootscript) bool {
			nameMatch := exp.MatchString(s.Title)
			architectureMatch := true
			if architecture != "" {
				architectureMatch = architecture == s.Arch
			}
			return nameMatch && architectureMatch
		}
	}

	var matches []api.Bootscript
	for _, script := range scripts {
		if isMatch(script) {
			matches = append(matches, script)
		}
	}

	if len(matches) > 1 {
		return fmt.Errorf("The query returned more than one result. Please refine your query.")
	}
	if len(matches) == 0 {
		return fmt.Errorf("The query returned no result. Please refine your query.")
	}

	return bootscriptDescriptionAttributes(d, matches[0])
}
