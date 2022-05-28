# elmasy

[![Go Report Card](https://goreportcard.com/badge/github.com/elmasy-com/elmasy)](https://goreportcard.com/report/github.com/elmasy-com/elmasy)

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/elmasy-com/logo/raw/main/png/logo_white_512.png">
    <img width="25%" height="25%" alt="Elmasy logo" src="https://github.com/elmasy-com/logo/raw/main/png/logo_black_512.png">
  </picture>
</p>


Elmasy *will be* an attack surface analysis tool.Discover and assess as many as possible publicly accessible assets from a domain name.

This project has 3 goals:

- To provide a website for non-professionals/professionals to check the attack surfaces of publicly accessible assets.
- To provide an API for professionals to work with easily.
    - The website will use the same API as your service.
- To provide a wide range of libraries for professionals to work with without relying on external services.
    - The API will use the same library as you.

As we grow, we will know more protocols, more services and assess/evaluate as fast as the programming language allows. The project is starts from lower layers (like TCP and UDP), basic services (like web servers) and gathering basic information about the servers (like TLS ciphers).

## Install

```bash
wget https://raw.githubusercontent.com/elmasy-com/elmasy/main/scripts/install.sh
```

```bash
sudo bash install.sh install
```
```bash
cp /opt/elmasy/elmasy.conf.example /opt/elmasy/elmasy.conf
```

Edit the config file.

## Run

### Systemd

The installer create a service for Elmasy.

To start:

```bash
systemctl start elmasy
```

### Manual

If `-config` not specified, Elmasy look for the config file in the current directory.

```bash
cd /opt/elmasy && ./elmasy
```

## Update

```bash
sudo bash install.sh updateself && sudo bash install.sh install
```

NOTE: If Elmasy was running before update, the `install` command will restart automatically.
