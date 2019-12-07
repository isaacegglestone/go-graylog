# graylog_event_definition

* [Example](https://github.com/suzuki-shunsuke/go-graylog/blob/master/terraform/example/v0.12/event_definition.tf)
* [Source Code](https://github.com/suzuki-shunsuke/go-graylog/blob/master/terraform/graylog/resource_event_definition.go)

## Argument Reference

### Required Argument

name | type | description
--- | --- | ---
title | string |
config | string | JSON string
notifications[].notification_id | string |

`config` is a JSON string.
The format of `config` depends on the Event Notification type.
Please see the [example](https://github.com/suzuki-shunsuke/go-graylog/blob/master/terraform/example/v0.12/event_definition.tf).
Using the [Graylog's API browser](https://docs.graylog.org/en/3.1/pages/configuration/rest_api.html) you can check the format of `config`.

### Optional Argument

name | default | type | description
--- | --- | --- | ---
description | "" | string |
priority | 0 | int |
alert | false | bool |
field_spec | "" | string | JSON string
notification_settings.grace_period_ms | 0 | int |
notification_settings.backlog_size | 0 | int |
notifications | [] | []object |

## Attrs Reference

None.
