# Krakend Endpoints Tool

This tool is a utility to work with sparsely defined [krakend.io](https://www.krakend.io/)
endpoints. The intention is to help service authors define their endpoints in
separated directories, and then generate the final krakend.json file.

This aims to also do the inverse: Given a krakend.json file, generate the
endpoints in a directory structure.

Note that this assumes that Krakend.io's [flexible configuration](
https://www.krakend.io/docs/configuration/flexible-config/) is used.

## Usage

### Generate krakend.json

Given a directory structure like this:

    endpoints/
    ├── api
    │   ├── api1
    │   │   ├── endpoint1.json
    │   │   └── endpoint2.json
    │   └── api2
    │       ├── endpoint1.json
    │       └── endpoint2.json
    config/
    └── krakend.tmpl

Where `krakend.tmpl` looks as follows:

    {
        "$schema": "https://www.krakend.io/schema/v3.json",
        "version": 3,
        "name": "my-krakend-instance",
        "port": {{ env "KRAKEND_PORT" }},
        "timeout": "3s",
        "cache_ttl": "3s",
        "output_encoding": "json",
        "plugin": {
            "pattern": ".so",
            "folder": "/opt/krakend/plugins/"
        },
        "endpoints": $ENDPOINTS$,
        "extra_config": {}
    }

Run the following command:

    $ krakend-endpoints-tool generate \
        --endpoints endpoints/ \
        --config config/krakend.tmpl \
        --output krakend.tmpl

This will generate a krakend.tmpl file with the following content:

    {
      "$schema": "https://www.krakend.io/schema/v3.json",
      "version": 3,
      "name": "my-krakend-instance",
      "port": {{ env "KRAKEND_PORT" }},
      "timeout": "3s",
      "cache_ttl": "3s",
      "output_encoding": "json",
      "plugin": {
          "pattern": ".so",
          "folder": "/opt/krakend/plugins/"
      },
      "endpoints": [
        {
          "endpoint": "/api/api1/endpoint1",
          "method": "GET",
          "backend": [
            {
              "url_pattern": "/api/api1/endpoint1",
              "host": [
                "http://localhost:8081"
              ],
              "encoding": "no-op",
              "sd": "static",
              "extra_config": {
                "github_com/devopsfaith/krakend/proxy": {
                  "disable_host_sanitize": true
                }
              }
            }
          ]
        },
        ...
      ]
    }
