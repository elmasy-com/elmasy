# elmasy

[![Go Report Card](https://goreportcard.com/badge/github.com/elmasy-com/elmasy)](https://goreportcard.com/report/github.com/elmasy-com/elmasy)

Elmasy *will be* an attack surface analysis tool.

## Install

```bash
wget https://raw.githubusercontent.com/elmasy-com/elmasy/main/scripts/install.sh
```

```bash
sudo bash install.sh
```

**Dont forget to edit the config file!**

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