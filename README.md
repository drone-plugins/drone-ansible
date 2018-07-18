# drone-ansible

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-ansible/status.svg)](http://beta.drone.io/drone-plugins/drone-ansible)
[![Join the discussion at https://www.reddit.com/r/droneci/](https://img.shields.io/badge/reddit-forum-orange.svg)](https://www.reddit.com/r/droneci/)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-ansible?status.svg)](http://godoc.org/github.com/drone-plugins/drone-ansible)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-ansible)](https://goreportcard.com/report/github.com/drone-plugins/drone-ansible)
[![](https://images.microbadger.com/badges/image/plugins/ansible.svg)](https://microbadger.com/images/plugins/ansible "Get your own image badge on microbadger.com")

Drone plugin to provision infrastructure with [Ansible](https://www.ansible.com/). For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-ansible/).

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-ansible
docker build --rm -t plugins/ansible .
```

### Usage

```
docker run --rm \
  -e PLUGIN_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" \
  -e PLUGIN_PLAYBOOK="deployment/playbook.yml" \
  -e PLUGIN_INVENTORY="deployment/hosts.yml" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/ansible
```
