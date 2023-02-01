# connections-looker
Small utility to collect Looker information which is only available with Admin permissions.

## How to use

If you have Golang installed, you can build the binary yourself, otherwise download appropriate binary from the [releases screen](https://github.com/getsynq/connections-looker/releases) (darwin == macOS).

There are two ways to run the binary, using command line arguments and using interactive wizard:

```
❯ ./connections-looker --help
Small utility to collect Looker information which is only available with Admin permissions

Usage:
  connections-looker [flags]

Flags:
      --client_id string       Client ID
      --client_secret string   Client Secret
  -h, --help                   help for connections-looker
      --url string             Full URL of the Looker instance

❯ ./connections-looker
? Full URL of the Looker instance: https://corp.cloud.looker.com
? Client ID: xxxxxxxxxxxxxxxxxxxx
? Client Secret: ************************
Discovered 1 database connections
File connections-2023-02-01T16_19_31Z.json created
```
