package elasticsearch

import (
	"net/url"

	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerOpts the options for the provider
type providerOpts struct {
	url      *url.URL
	username string
	password string
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"elasticsearch_template": resourceTemplate(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ES_URL", nil),
				Description: "The URL of the elasticsearch",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ES_USERNAME", nil),
				Description: "The username if there is a basic auth for elasticsearch",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ES_PASSWORD", nil),
				Description: "The password if there is a basic auth for elasticsearch",
			},
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	u := d.Get("url").(string)
	validURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return &providerOpts{
		url:      validURL,
		username: d.Get("username").(string),
		password: d.Get("password").(string),
	}, nil
}

func getClientES(opts *providerOpts) (*es7.Client, error) {
	cfg := es7.Config{
		Addresses: []string{
			opts.url.String(),
		},
		Username: opts.username,
		Password: opts.password,
	}
	return es7.NewClient(cfg)
}
