# drone-passwordstate

[![Build Status](https://travis-ci.org/TDabasinskas/drone-passwordstate.svg?branch=master)](https://travis-ci.org/TDabasinskas/drone-passwordstate)

## Description

[drone.io](https://drone.io/) plugin, allowing to export passwords (secrets) from a [Click Studio Passwordstate](https://www.clickstudios.com.au/) password list. The plugin exports secrets to the specified file within the workspace, allowing the file to be used inside further pipeline steps (e.g. deploying them to [Kubernetes](https://kubernetes.io/) via [drone-helm](https://github.com/ipedrazas/drone-helm) plugin.

## Usage

### Simple usage

To simply export all the secrets from the specified Passwordstate to `./secrets.yaml` file, the following pipeline step should be added:

```yaml
pipeline:
  inject_secrets:
    image: tdabasinskas/drone-passwordstate
    api_endpoint: https://passwordstate/api/
    api_key: $PASSWORD_STATE_KEY
    skip_tls_verify: false
    password_list_id: 1231
    output_path: ./secrets.yaml
    secrets: [ PASSWORD_STATE_KEY ]
```

The plugin would connect to the specified Password state instance and extract the passwords as secrets using `UserName` field as the secret key and `Password` field as the secret value. Once finished, the folllowing file would be created within the workspace:

```yaml
secrets:
  some_secret: 'some_secret_value'
  another_secret: 'another_secret_value'
```

### Encoding the secrets

By default, the secrets are exported *as-is*, meaning, they would need to separately encoded with BASE64 if used as Kubernetes secrets. To handle that automatically, `encode_secrets` parameter can be used, e.g.:

```yaml
pipeline:
  inject_secrets:
    image: tdabasinskas/drone-passwordstate
    api_endpoint: https://passwordstate/api/
    api_key: d417b3c2f586b9eaed8b736f95324cd5
    skip_tls_verify: false
    password_list_id: 1231
    output_path: ./secrets.yaml
    encode_secrets: true
```

### Using different Key/Value fields

As mentioned, by default, `UserName` and `Password` fields are used as the Key/Value pair. Anyhow, it's possible, to use different fields, e.g.:

```yaml
pipeline:
  inject_secrets:
    image: tdabasinskas/drone-passwordstate
    api_endpoint: https://passwordstate/api/
    api_key: d417b3c2f586b9eaed8b736f95324cd5
    skip_tls_verify: false
    password_list_id: 1231
    key_field: Title
    value_field: GenericField6
```

### Using the plugin for Kubernetes secrets

One of the most likely use case for the plugin would be combining it with [drone-helm plugin](https://github.com/ipedrazas/drone-helm), allowing you to deploy the secrets as part of the whole [Helm chart](https://github.com/kubernetes/helm). The following example illustrates the pipeline combining these two plugins:

```yaml
pipeline:
  inject_secrets:
    image: tdabasinskas/drone-passwordstate
    api_endpoint: https://passwordstate/api/
    skip_tls_verify: false
    password_list_id: 1231
    output_path: ./secrets.yaml
    encode_secrets: true
    secrets: [ API_KEY ]
  deploy:
    image: quay.io/ipedrazas/drone-helm
    chart: ./helm
    release: app
    values_files: [ ./helm/values.default.yaml, ./secrets.yaml ]
    wait: true
    prefix: DEV
```

Assuming the helm chart under `./helm` contains the following secrets template file, it would be automatically filled with the secrets during the deployment:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: {{ .Release.Name }}-secrets
type: Opaque
data:
  cache__connectionString: {{ .Values.secrets.some_secret | quote }}
  consul__token: {{ .Values.secrets.another_secret | quote }}
```

## Known issues

- The plugin currently supports exporting of all secrets within the password list only, not allowing to specify the exact secrets (passwords) to export.

## Contributing

Feel free to fork the repository and submit changes via a Pull Request.
