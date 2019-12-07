# terraform-provider-graylog

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/suzuki-shunsuke/go-graylog/terraform)
[![Build Status](https://cloud.drone.io/api/badges/suzuki-shunsuke/go-graylog/status.svg)](https://cloud.drone.io/suzuki-shunsuke/go-graylog)
[![codecov](https://codecov.io/gh/suzuki-shunsuke/go-graylog/branch/master/graph/badge.svg)](https://codecov.io/gh/suzuki-shunsuke/go-graylog)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/go-graylog)](https://goreportcard.com/report/github.com/suzuki-shunsuke/go-graylog)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/go-graylog.svg)](https://github.com/suzuki-shunsuke/go-graylog)
[![GitHub tag](https://img.shields.io/github/tag/suzuki-shunsuke/go-graylog.svg)](https://github.com/suzuki-shunsuke/go-graylog/releases)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/go-graylog/master/LICENSE)

terraform provider for [Graylog](https://www.graylog.org/).

This is sub project of [go-graylog](https://github.com/suzuki-shunsuke/go-graylog).

## Motivation

https://docs.graylog.org/en/latest/pages/users_and_roles/permission_system.html

The Graylog permission system is extremely flexible but you can't utilize this flexibility from Web UI.
By using this provider, you can utilize this flexibility and manage the infrastructure as code.

## Install

[Download binary](https://github.com/suzuki-shunsuke/go-graylog/releases) and install it.

https://www.terraform.io/docs/configuration/providers.html#third-party-plugins

## Docker Image

https://quay.io/repository/suzuki_shunsuke/terraform-graylog

Docker image which is installed terraform and terraform-provider-graylog on Alpine.

## Example

```hcl
provider "graylog" {
  web_endpoint_uri = "${var.web_endpoint_uri}"
  auth_name = "${var.auth_name}"
  auth_password = "${var.auth_password}"
}

// Role my-role-2
resource "graylog_role" "my-role-2" {
  name = "my-role-2"
  permissions = ["users:edit"]
  description = "Created by terraform"
}
```

And please see [example v0.11](example/v0.11) and [example v0.12](example/v0.12) also.

## Variables

### Required

name | Environment variable | description
--- | --- | ---
web_endpoint_uri | GRAYLOG_WEB_ENDPOINT_URI | API endpoint, for example https://graylog.example.com/api
auth_name | GRAYLOG_AUTH_NAME | Username or API token or Session Token
auth_password | GRAYLOG_AUTH_PASSWORD | Password or the literal `"token"` or `"session"`

About `auth_name` and `auth_password`, please see the [Graylog's Documentation](https://docs.graylog.org/en/latest/pages/configuration/rest_api.html).

You can authenticate with either password or access token or session token.

password

```
auth_name = "<user name>"
auth_password = "<password>"
```

access token

```
auth_name = "<access token>"
auth_password = "token"
```

session token

```
auth_name = "<session token>"
auth_password = "session"
```

### Optional

name | Environment variable | default | description
--- | --- | --- | ---
x_requested_by | GRAYLOG_X_REQUESTED_BY | terraform-go-graylog | [X-Requested-By Header](https://github.com/Graylog2/graylog2-server/blob/370dd700bc8ada5448bf66459dec9a85fcd22d58/UPGRADING.rst#protecting-against-csrf-http-header-required)
api_version | GRAYLOG_API_VERSION | "v2" | Graylog's API version. The default value is "v2" for compatibility. If you use Graylog v3, please set "v3".

## Resources

* [alarm_callback](docs/alarm_callback.md)
* [alert_condition](docs/alert_condition.md)
* [dashboard](docs/dashboard.md)
* [dashboard_widget](docs/dashboard_widget.md)
* [dashboard_widget_positions](docs/dashboard_widget_positions.md)
* [extractor](docs/extractor.md)
* [event_definition](docs/event_definition.md)
* [event_notification](docs/event_notification.md)
* [grok_pattern](docs/grok_pattern.md)
* [index_set](docs/index_set.md)
* [input](docs/input.md)
* [input_static_fields](docs/input_static_fields.md)
* [ldap_setting](docs/ldap_setting.md)
* [output](docs/output.md)
* [pipeline](docs/pipeline.md)
* [pipeline_rule](docs/pipeline_rule.md)
* [pipeline_connection](docs/pipeline_connection.md)
* [role](docs/role.md)
* [stream](docs/stream.md)
* [stream_output](docs/stream_output.md)
* [stream_rule](docs/stream_rule.md)
* [user](docs/user.md)

## Data sources

* [dashboard](docs/data_source_dashboard.md)
* [index_set](docs/data_source_index_set.md)
* [stream](docs/data_source_stream.md)

## Unsupported resources

We can't support these resources for some reasons.

### CollectorConfiguration (includes input, output snippet)

We can't support these resources because graylog API doesn't return the created resource id (response body: no content).

The following APIs doesn't return the created resource id (response body: no content).

* POST /plugins/org.graylog.plugins.collector/configurations/{id}/inputs Create a configuration input
* POST /plugins/org.graylog.plugins.collector/configurations/{id}/outputs Create a configuration output
* POST /plugins/org.graylog.plugins.collector/configurations/{id}/snippets Create a configuration snippet
