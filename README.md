![master](https://github.com/github/docs/actions/workflows/main.yml/badge.svg?branch=master)
[![report card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/mattermost/terraform-provider-elasticsearch)

# Terraform Provider Elasticsearch

Run the following command to build the provider and place in the terraform plugins in your local machine:

```shell
make
```

## Test Example configuration

You can find an examples sample [here](examples):

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
cd examples
terraform init && terraform apply
```

## Run acceptance tests

```
make run-elasticsearch 
ES_URL=http://localhost:9200 make testacc
```
