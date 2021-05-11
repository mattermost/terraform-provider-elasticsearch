package elasticsearch

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"
)

func resourceTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateCreate,
		UpdateContext: resourceTemplateUpdate,
		ReadContext:   resourceTemplateRead,
		DeleteContext: resourceTemplateDelete,
		Schema: map[string]*schema.Schema{
			// resource arguments and their specifications go here
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the elasticsearch template resource.",
			},
			"template": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The JSON body for the template you want to use for an elasticsearch index.",
				DiffSuppressFunc: indexDiffSuppressFunc,
				ValidateFunc:     validation.StringIsJSON,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func indexDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	newJson, _ := structure.NormalizeJsonString(new)
	oldJson, _ := structure.NormalizeJsonString(old)
	return newJson == oldJson
}

func resourceTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()

	var diags diag.Diagnostics

	c, err := getClientES(meta.(*providerOpts))
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "Failed to get elasticsearch client"))
	}
	req := esapi.IndicesGetTemplateRequest{
		Name:   []string{id},
		Pretty: true,
	}
	resp, err := req.Do(ctx, c)
	if err != nil {
		return diag.FromErr(errors.Wrapf(err, "Failed to fetch index template with ID: %s", id))
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "Failed to read the response"))
	}
	//nolint
	d.Set("name", d.Id())
	//nolint
	d.Set("template", string(bodyBytes))
	return diags
}

func resourceTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	template := d.Get("template").(string)

	var diags diag.Diagnostics
	c, err := getClientES(meta.(*providerOpts))
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "Failed to get elasticsearch client"))
	}
	req := esapi.IndicesPutTemplateRequest{
		Name:   name,
		Create: esapi.BoolPtr(true),
		Order:  esapi.IntPtr(0),
		Pretty: true,
		Body:   strings.NewReader(template),
	}
	resp, err := req.Do(ctx, c)
	if err != nil {
		return diag.FromErr(errors.Wrapf(err, "Failed to create a template with ID: %s", name))
	}
	defer resp.Body.Close()

	d.SetId(name)
	return diags
}

func resourceTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()

	var diags diag.Diagnostics
	c, err := getClientES(meta.(*providerOpts))
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "Failed to get elasticsearch client"))
	}
	req := esapi.IndicesDeleteTemplateRequest{
		Name:   id,
		Pretty: true,
	}
	resp, err := req.Do(ctx, c)
	if err != nil {
		return diag.FromErr(errors.Wrapf(err, "Failed to delete template with ID: %s", id))
	}
	defer resp.Body.Close()

	d.SetId("")
	return diags
}

func resourceTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceTemplateCreate(ctx, d, meta)
}
