Development of new client editor
=============================

This documentation will give you some tips to help when you are going to support the promql-server in a new editor.

## Config
There is two different config used by the promql-server. 

### Cold configuration (YAML or by environment)
The first one is a classic yaml file which can be used to customize the server when it is starting. It's specified by adding the `--config-file` command line flag when starting the language server.

It has the following structure:

```yaml
activate_rpc_log: false # It's a boolean in order to activate or deactivate the rpc log. It's deactivated by default and mainly useful for debugging the language server, by inspecting the communication with the language client.
log_format: "text" # The format of the log printed. Possible value: json, text. Default value: "text"
prometheus_url: "http://localhost:9090" # the HTTP URL of the prometheus server.
rest_api_port: 8080 # When set, the server will be started as an HTTP server that provides a REST API instead of the language server protocol. Default value: 0
```

In case the file is not provided, it will read the configuration from the environment variables with the following structure:

```bash
export LANGSERVER_ACTIVATERPCLOG="true"
export LANGSERVER_PROMETHEUSURL="http://localhost:9090"
export LANGSERVER_RESTAPIPORT"="8080"
export LANGSERVER_LOGFORMAT"="json"
```

Note: documentation and default value are the same for both configuration (yaml and environment)

### JSON configuration (hot configuration)
There is a second configuration which is used only at the runtime and can be sent by the language client over the `DidChangeConfiguration` API. It's used to sync configuration from the text editor to the language server.

It has the following structure:

```json
{
  "promql": {
    "url": "http://localhost:9090" # the HTTP URL of the prometheus server.
  }
}
```

Using this way of changing configuration requires both providing those config options in the Text editor and sending them to the language server whenever they are changed.
