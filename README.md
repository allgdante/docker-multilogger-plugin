# Multilogger Docker Logging Driver

## Overview

This Docker plugins sends container logs to multiple destinations at once. These instructions are for Linux host systems. For other platforms, see the [Docker Engine managed plugin system documentation](https://docs.docker.com/engine/extend/).

If you have any questions or issues using the Docker plugin feel free to open an issue in this [repository](https://github.com/allgdante/docker-multilogger-plugin/issues).

## Plugin Installation

### Install from the Docker store

Pull the plugin from the Docker Store

```
$ docker plugin install --alias multilogger allgdante/docker-multilogger-plugin:0.0.1
```

Enable the plugin

```
$ docker plugin enable multilogger
```

### Install from source

Clone the repository

```
$ cd docker-multilogger-plugin
```

Build the plugin

```
$ make all
```

Enable the plugin

```
$ docker plugin multilogger enable
```

## Plugin Configuration

### Configure the logging driver for a container

If `multilogger` is in use but no logging driver is selected, the `json-file` logging driver will be automatically enabled:

```sh
docker run \
    --log-driver=multilogger \
    nginx/stable-alpine
```

The following command configure the container to start with the `multilogger` driver which will send logs to `gelf` and `syslog` servers.

```sh
docker run \
    --log-driver=multilogger \
    --log-opt gelf-enabled=true \
    --log-opt gelf-address=udp://127.0.0.1:12201 \
    --log-opt syslog-enabled=true \
    --log-opt syslog-address=tcp://127.0.0.1:514 \
    nginx/stable-alpine
```

In the previous example, if we want to have the `json-file` logging driver enabled, it must be explicitly enabled like this:

```sh
docker run \
    --log-driver=multilogger \
    --log-opt json-file-enabled=true \
    --log-opt gelf-enabled=true \
    --log-opt gelf-address=udp://127.0.0.1:12201 \
    --log-opt syslog-enabled=true \
    --log-opt syslog-address=tcp://127.0.0.1:514 \
    nginx/stable-alpine
```

### Configure the default logging driver

To configure the Docker daemon to default to `multilogger` logging driver, you must configure the file `/etc/docker/daemon.json` like this:

```json
{
    "log-driver": "multilogger"
}
```

If we want to include options, it must be configured like this:

```json
{
    "log-driver": "multilogger",
    "log-opts": {
        "json-file-enabled": "true",
        "gelf-enabled": "true",
        "gelf-address": "udp://127.0.0.1:12201",
        "syslog-enabled": "true",
        "syslog-address": "tcp:/127.0.0.1:514"
    }
}
```

### Available options and logging drivers

#### JSON File logging driver

Refer to the official [documentation](https://docs.docker.com/config/containers/logging/json-file/) for more details.
Note that this logging driver is always enabled if `multilogger` driver is selected and no logging driver is enabled

| Option                                    | Description                                       |
|-------------------------------------------|---------------------------------------------------|
| `json-file-enabled`                       | To enable this driver, use `true` here.           |
| `max-size`                                | Same as Docker.                                   |
| `max-file`                                | Same as Docker.                                   |
| `labels`                                  | Same as Docker.                                   |
| `labels-regex`                            | Same as Docker.                                   |
| `env`                                     | Same as Docker.                                   |
| `env-regex`                               | Same as Docker.                                   |
| `compress`                                | Same as Docker.                                   |
| `tag`                                     | Same as Docker.                                   |
| `json-file-log-dir`                       | `/var/log/docker`

#### Amazon CloudWatch Logs logging driver

Refer to the official [documentation](https://docs.docker.com/config/containers/logging/awslogs/) for more details.

| Option                                    | Description                                       |
|-------------------------------------------|---------------------------------------------------|
| `awslogs-enabled`                         | To enable this driver, use `true` here.           |
| `awslogs-region`                          | Same as Docker.                                   |
| `awslogs-endpoint`                        | Same as Docker.                                   |
| `awslogs-group`                           | Same as Docker.                                   |
| `awslogs-create-group`                    | Same as Docker.                                   |
| `awslogs-datetime-format`                 | Same as Docker.                                   |
| `awslogs-multiline-pattern`               | Same as Docker.                                   |
| `awslogs-credentials-endpoint`            | Same as Docker.                                   |
| `awslogs-force-flush-interval-seconds`    | Same as Docker.                                   |
| `awslogs-max-buffered-events`             | Same as Docker.                                   |
| `awslogs-tag`                             | Same as `tag` parameter in Docker docs.           |

#### Fluentd logging driver

Refer to the official [documentation](https://docs.docker.com/config/containers/logging/fluentd/) for more details.

| Option                                    | Description                                       |
|-------------------------------------------|---------------------------------------------------|
| `fluentd-enabled`                         | To enable this driver, use `true` here.           |
| `fluentd-address`                         | Same as Docker.                                   |
| `fluentd-async`                           | Same as Docker.                                   |
| `fluentd-async-connect`                   | Same as Docker.                                   |
| `fluentd-buffer-limit`                    | Same as Docker.                                   |
| `fluentd-max-retries`                     | Same as Docker.                                   |
| `fluentd-request-ack`                     | Same as Docker.                                   |
| `fluentd-retry-wait`                      | Same as Docker.                                   |
| `fluentd-sub-second-precision`            | Same as Docker.                                   |
| `fluentd-labels`                          | Same as `labels` parameter in Docker docs.        |
| `fluentd-labels-regex`                    | Same as `labels-regex` parameter in Docker docs.  |
| `fluentd-env`                             | Same as `env` parameter in Docker docs.           |
| `fluentd-env-regex`                       | Same as `env-regex` parameter in Docker docs.     |
| `fluentd-tag`                             | Same as `tag` parameter in Docker docs.           |

#### Google Cloud Logging driver

Refer to the official [documentation](https://docs.docker.com/config/containers/logging/gcplogs/) for more details.

| Option                                    | Description                                       |
|-------------------------------------------|---------------------------------------------------|
| `gcp-enabled`                             | To enable this driver, use `true` here.           |
| `gcp-project`                             | Same as Docker.                                   |
| `gcp-log-cmd`                             | Same as Docker.                                   |
| `gcp-meta-zone`                           | Same as Docker.                                   |
| `gcp-meta-name`                           | Same as Docker.                                   |
| `gcp-meta-id`                             | Same as Docker.                                   |
| `gcp-labels`                              | Same as `labels` parameter in Docker docs.        |
| `gcp-labels-regex`                        | Same as `labels-regex` parameter in Docker docs.  |
| `gcp-env`                                 | Same as `env` parameter in Docker docs.           |
| `gcp-env-regex`                           | Same as `env-regex` parameter in Docker docs.     |

#### Graylog Extended Format logging driver

Refer to the official [documentation](https://docs.docker.com/config/containers/logging/gelf/) for more details.

| Option                                    | Description                                       |
|-------------------------------------------|---------------------------------------------------|
| `gelf-enabled`                            | To enable this driver, use `true` here.           |
| `gelf-address`                            | Same as Docker.                                   |
| `gelf-compression-level`                  | Same as Docker.                                   |
| `gelf-compresion-type`                    | Same as Docker.                                   |
| `gelf-tcp-max-reconnect`                  | Same as Docker.                                   |
| `gelf-labels`                             | Same as `labels` parameter in Docker docs.        |
| `gelf-labels-regex`                       | Same as `labels-regex` parameter in Docker docs.  |
| `gelf-env`                                | Same as `env` parameter in Docker docs.           |
| `gelf-env-regex`                          | Same as `env-regex` parameter in Docker docs.     |
| `gelf-tag`                                | Same as `tag` parameter in Docker docs.           |

#### Journald logging driver

Refer to the official [documentation](https://docs.docker.com/config/containers/logging/journald/) for more details.

| Option                                    | Description                                       |
|-------------------------------------------|---------------------------------------------------|
| `journald-enabled`                        | To enable this driver, use `true` here.           |

#### Logentries logging driver

Refer to the official [documentation](https://docs.docker.com/config/containers/logging/logentries/) for more details.

| Option                                    | Description                                       |
|-------------------------------------------|---------------------------------------------------|
| `logentries-enabled`                      | To enable this driver, use `true` here.           |
| `logentries-token`                        | Same as Docker.                                   |
| `line-only`                               | Same as Docker.                                   |
| `logentries-labels`                       | Same as `labels` parameter in Docker docs.        | 
| `logentries-labels-regex`                 | Same as `labels-regex` parameter in Docker docs.  | 
| `logentries-env`                          | Same as `env` parameter in Docker docs.           | 
| `logentries-env-regex`                    | Same as `env-regex` parameter in Docker docs.     | 
| `logentries-tag`                          | Same as `tag` parameter in Docker docs.           | 

#### Splunk logging driver

Refer to the official [documentation](https://docs.docker.com/config/containers/logging/splunk/) for more details.

| Option                                    | Description                                       |
|-------------------------------------------|---------------------------------------------------|
| `splunk-enabled`                          | To enable this driver, use `true` here.           |
| `splunk-url`                              | Same as Docker.                                   |
| `splunk-token`                            | Same as Docker.                                   |
| `splunk-source`                           | Same as Docker.                                   |
| `splunk-sourcetype`                       | Same as Docker.                                   |
| `splunk-index`                            | Same as Docker.                                   |
| `splunk-capath`                           | Same as Docker.                                   |
| `splunk-caname`                           | Same as Docker.                                   |
| `splunk-insecureskipverify`               | Same as Docker.                                   |
| `splunk-format`                           | Same as Docker.                                   |
| `splunk-verify-connection`                | Same as Docker.                                   |
| `splunk-gzip`                             | Same as Docker.                                   |
| `splunk-gzip-level`                       | Same as Docker.                                   |
| `splunk-index-acknowledgment`             | Same as Docker.                                   |
| `splunk-labels`                           | Same as `labels` parameter in Docker docs.        |
| `splunk-labels-regex`                     | Same as `labels-regex` parameter in Docker docs.  |
| `splunk-env`                              | Same as `env` parameter in Docker docs.           |
| `splunk-env-regex`                        | Same as `env-regex` parameter in Docker docs.     |
| `splunk-tag`                              | Same as `tag` parameter in Docker docs.           |

#### Syslog logging driver

Refer to the official [documentation](https://docs.docker.com/config/containers/logging/syslog/) for more details.

| Option                                    | Description                                       |
|-------------------------------------------|---------------------------------------------------|
| `syslog-enabled`                          | To enable this driver, use `true` here.           |
| `syslog-address`                          | Same as Docker.                                   |
| `syslog-facility`                         | Same as Docker.                                   |
| `syslog-tls-ca-cert`                      | Same as Docker.                                   |
| `syslog-tls-cert`                         | Same as Docker.                                   |
| `syslog-tls-key`                          | Same as Docker.                                   |
| `syslog-tls-skip-verify`                  | Same as Docker.                                   |
| `syslog-format`                           | Same as Docker.                                   |
| `syslog-labels`                           | Same as `labels` parameter in Docker docs.        |
| `syslog-labels-regex`                     | Same as `labels-regex` parameter in Docker docs.  |
| `syslog-env`                              | Same as `env` parameter in Docker docs.           |
| `syslog-env-regex`                        | Same as `env-regex` parameter in Docker docs.     |
| `syslog-tag`                              | Same as `tag` parameter in Docker docs.           |

#### Syslog5424 logging driver

It's a modified `syslog` driver that puts labels and environment variables as structured data.

| Option                                    | Description                                                                                                               |
|-------------------------------------------|---------------------------------------------------------------------------------------------------------------------------|
| `syslog5424-enabled`                      | To enable this driver, use `true` here.                                                                                   |
| `syslog5424-address`                      | The address of an external syslog server. The URI specifier may be [tcp|udp|tcp+tls]://host:port, unix://path, or unixgram://path. If the transport is tcp, udp, or tcp+tls, the default port is 514. |
| `syslog5424-facility`                     | The syslog facility to use. Can be the number or name for any valid syslog facility. See the [syslog documentation](https://tools.ietf.org/html/rfc5424#section-6.2.1). |
| `syslog5424-time-format`                  | Use `rfc3339` for RFC-5424 compatible format, or `rfc3339micro` for RFC-5424 compatible format with microsecond timestamp resolution. |
| `syslog5424-tls-ca-cert`                  | The absolute path to the trust certificates signed by the CA. Ignored if the address protocol is not `tcp+tls`.           |
| `syslog5424-tls-cert`                     | The absolute path to the TLS certificate file. Ignored if the address protocol is not `tcp+tls`.                          |
| `syslog5424-tls-key`                      | The absolute path to the TLS key file. Ignored if the address protocol is not `tcp+tls`.                                  |
| `syslog5424-tls-skip-verify`              | If set to true, TLS verification is skipped when connecting to the syslog daemon. Defaults to `false`. Ignored if the address protocol is not `tcp+tls`. |
| `syslog5424-hostname`                     | Defaults to `os.Hostname()`, but we could use a literal value or a template using the [info](https://godoc.org/github.com/docker/docker/daemon/logger#Info) struct as reference. | 
| `syslog5424-msgid`                        | Defaults to the `syslog5424-tag` value, but we could use a literal value or a template using the [info](https://godoc.org/github.com/docker/docker/daemon/logger#Info) struct as reference. |
| `syslog5424-disable-framer`               | If `true`, we won't sent the RFC5425 message length framer. Disabled by default.                                          |
| `syslog5424-labels`                       | List of comma-separated labels that will be used as structured data in every message.                                     |
| `syslog5424-labels-regex`                 | Regular expression to match labels that will be used as structured data in every message.                                 |
| `syslog5424-env`                          | List of comma-separated environment variables that will be used as structured data in every message.                      |
| `syslog5424-env-regex`                    | Regular expression to match environment variables that will be used as structured data in every message.                  |
| `syslog5424-tag`                          | Defaults to `{{.ID}}` template, but we could use a literal value or a template using the [info](https://godoc.org/github.com/docker/docker/daemon/logger#Info) struct as reference. |

## Uninstall the plugin

To cleanly disable and remove the plugin, run:

```bash
docker plugin disable multilogger
docker plugin rm multilogger
```
