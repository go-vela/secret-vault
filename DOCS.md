## Description

This plugin enables the ability pull secrets from [Vault](https://www.vaultproject.io/) into the secret mount within a Vela pipeline.

Source Code: https://github.com/go-vela/secret-vault

Registry: https://hub.docker.com/r/target/secret-vault

## Usage

Sample of retrieving a secret using token authentication:

```yaml
secrets:
  - origin:
      name: vault
      image: target/secret-vault:latest
      parameters:
        addr: vault.company.com
        token: superSecretVaultToken
        auth_method: token
        items:
          # Written to path: "/vela/secrets/docker/<key>"
          - source: secret/vela/username
            path: docker
```

Sample of retrieving a secret using ldap authentication:

```diff
secrets:
  - origin:
      name: vault
      image: target/secret-vault:latest
      parameters:
        addr: vault.company.com
+       username: octocat
+       password: superSecretPassword
-       token: superSecretVaultToken
+       auth_method: ldap
        items:
          # Written to path: "/vela/secrets/docker/<key>"
          - source: secret/vela/username
            path: docker
```

Sample of reading a secret using ldap authentication with verbose logging:

```diff
secrets:
  - origin:
      name: vault
      image: target/secret-vault:latest
      parameters:
        addr: vault.company.com
        username: octocat
        password: superSecretPassword
        token: superSecretVaultToken
        auth_method: ldap
+       log_level: trace        
        items:
          # Written to path: "/vela/secrets/docker/<key>"
          - source: secret/vela/username
            path: docker
```

Sample of retrieving a secret and customizing environment targets for the value
```yaml
secrets:
  - origin:
      name: vault
      image: target/secret-vault:latest
      secrets:
        - source: superSecretToken
          target: vault_token
      parameters:
        addr: vault.company.com
        auth_method: token
        items:
          # assuming user_A has two keys: `username` and `password`
          - source: secret/vela/user_A
            keys:
              - name: username
                target: [ KANIKO_USERNAME, ARTIFACTORY_USERNAME ]
              - name: password
                target: [ KANIKO_PASSWORD, ARTIFACTORY_PASSWORD ]
```

## Secrets

**NOTE: Users should refrain from configuring sensitive information in your pipeline in plain text.**

**NOTE: Secrets used within the secret plugin must exist as Vela secrets.**

You can use Vela secrets to substitute sensitive values at runtime:

```diff
secrets:
  # Repo secret created within Vela
  - name: vault_token
  
  # Example using token authentication method
  - origin:
      name: vault
      image: target/secret-vault:latest
      secret: [ vault_token ]
      parameters:
        addr: vault.company.com
-       token: superSecretVaultToken
        auth_method: token
        items:
          # Written to path: "/vela/secrets/docker/<key>"
          - source: secret/vela/username
            path: docker
```

## Parameters

The following parameters are used to configure the image:

| Name          | Description                                              | Required  | Default |
| ------------- | -------------------------------------------------------- | --------- | ------- |
| `addr`        | address to the instance                                  | `true`    | `N/A`   |
| `auth_method` | authentication method for interfacing (i.e. token, ldap) | `true`    | `N/A`   |
| `log_level`   | set the log level for the plugin                         | `true`    | `info`  |
| `password`    | password for server authentication with ldap             | `false`   | `N/A`   |
| `token`       | token for server authentication                          | `false`   | `N/A`   |
| `username`    | set the log level for the plugin                         | `false`   | `N/A`   |
| `items`       | set of secrets to retrieve and write to workspace        | `true`    | `N/A`   |

### Items

| Name          | Description                                              | Required                  | Default      |
| ------------- | -------------------------------------------------------- | ------------------------- | ------------ |
| `source`      | path to secret                                           | `true`                    | `N/A`        |
| `path`        | desired file path under `vela/secrets/` directory        | `path` or `keys` required | `N/A`        |
| `keys`        | custom environment variable or file path targets for key | `path` or `keys` required | `N/A`        |

### Keys

| Name          | Description                                                        | Required                    | Default      |
| ------------- | ------------------------------------------------------------------ | --------------------------- | ------------ |
| `name`        | name of key in a standard K-V vault                                | `true`                      | `N/A`        |
| `target`      | desired environment variable(s) for key value                      | `target` or `path` required | `N/A`        |
| `path`        | custom file path for key value (auto prefixed by `/vela/secrets/`) | `target` or `path` required | `N/A`        |



## Template

COMING SOON!

## Troubleshooting

Below are a list of common problems and how to solve them:
