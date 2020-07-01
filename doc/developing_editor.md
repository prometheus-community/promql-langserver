Development of new client editor
=============================

This documentation will give you some tips to help when you are going to support the promql-server in a new editor.

## Config
There is two different config used by the promql-server. 

### YAML configuration (cold configuration)
The first one is a classic yaml file which can be used to customize the server when it is starting. It has the following structure:

```yaml
activate_rpc_log: false # It's a boolean in order to activate or deactivate the rpc log. It's deactivated by default
log_format: "text" # The format of the log printed. Possible value: json, text. Default value: "text"
prometheus_url: "http://localhost:9090" # the HTTP URL of the prometheus server.
rest_api_port: 8080 # When set, the server will be started as an HTTP server. Default value: 0
```

### JSON configuration (hot configuration)
There is a second configuration which is used only at the runtime.

```json
{
  "promql": {
    "url": "http://localhost:9090" # the HTTP URL of the prometheus server.
  }
}
```

This configuration will be watched thanks to the method implemented `DidChangeConfiguration`
