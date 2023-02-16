# Commercelayer Plugin for Mach Composer

This repository contains the Commercelayer plugin for Mach Composer. It requires Mach Composer > 2.5


## Usage

```yaml
mach_composer:
  version: 1
  plugins:
    commercelayer:
      source: mach-composer/commercelayer
      version: 0.0.3

global:
  # ...

sites:
  - identifier: my-site

    comemrcelayer:
      client_id: client-id
      client_secret: client-secret
      domain: https://<your-domain>.commercelayer.io

```
