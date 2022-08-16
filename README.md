# drone-ansible

[![go-ci](https://github.com/hay-kot/drone-ansible/actions/workflows/go.yaml/badge.svg)](https://github.com/hay-kot/drone-ansible/actions/workflows/go.yaml)

Drone plugin to provision infrastructure with [Ansible](https://www.ansible.com/). For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-ansible/).


## This is a Fork of the [drone-ansible](http://plugins.drone.io/drone-plugins/drone-ansible/)

The following changes have been made (or are in progress):

- [x] Updated Ansible to the latest version
- [x] Build and Publish x86 Images to Github Docker Registry
- [x] Use GoReleaser to manage builds
- [x] Update dependencies
- [ ] Additional Image Builds for ARM
- [ ] Improve SSH Key Reading

## Documentation

*Documentation Copied from Plugin Page*

### Workflow Examples

#### Drone

```yaml
kind: pipeline
name: default

steps:
- name: check ansible syntax
  image: ghcr.io/hay-kot/drone-ansible:latest
  settings:
    playbook: ansible/playbook.yml
    galaxy: ansible/requirements.yml
    inventory: ansible/inventory
    syntax_check: true
```

#### Harness

```yaml
pipeline:
  stages:
    - identifier: default
      name: default
      steps:
        - identifier: check ansible syntax
          name: check ansible syntax
          spec:
            connectorRef: account.docker
            image: ghcr.io/hay-kot/drone-ansible:latest
            type: Plugin
            settings:
              playbook: ansible/playbook.yml
              galaxy: ansible/requirements.yml
              inventory: ansible/inventory
              syntax_check: true
```

### Configuration Properties
**become** (boolean): run operations with become
 - Default: False
 - Required: False
 - Secret: False


**become_method** (string): privilege escalation method to use
 - Default:
 - Required: False
 - Secret: False


**become_user** (string): run operations as this user
 - Default:
 - Required: False
 - Secret: False


**check** (boolean): run a check, do not apply any changes
 - Default: False
 - Required: False
 - Secret: False


**connection** (string): connection type to use
 - Default:
 - Required: False
 - Secret: False


**diff** (boolean): show the differences, may print secrets
 - Default: False
 - Required: False
 - Secret: False


**extra_vars** (string): set additional variables as key=value e.g. 'key1=value1,[key2=value2]'
 - Default:
 - Required: False
 - Secret: False


**flush_cache** (boolean): clear the fact cache for every host in inventory
 - Default: False
 - Required: False
 - Secret: False


**force_handlers** (boolean): run handlers even if a task fails
 - Default: False
 - Required: False
 - Secret: False


**forks** (number): specify number of parallel processes to use
 - Default: 5
 - Required: False
 - Secret: False


**galaxy** (string): path to galaxy requirements
 - Default:
 - Required: False
 - Secret: False


**galaxy_force** (boolean): force overwriting an existing role or collection
 - Default: True
 - Required: False
 - Secret: False


**inventory** (string): specify (multiple) inventory host path(s) e.g. 'path1,[path2]'
 - Default:
 - Required: False
 - Secret: False


**limit** (string): further limit selected hosts to an additional pattern
 - Default:
 - Required: False
 - Secret: False


**list_hosts** (boolean): outputs a list of matching hosts
 - Default: False
 - Required: False
 - Secret: False


**list_tags** (boolean): list all available tags
 - Default: False
 - Required: False
 - Secret: False


**list_tasks** (boolean): list all tasks that would be executed
 - Default: False
 - Required: False
 - Secret: False


**module_path** (string): prepend paths to module library e.g. 'path1,[path2]'
 - Default:
 - Required: False
 - Secret: False


**playbook** (string): list of playbooks to apply e.g. 'playbook1,[playbook2]'
 - Default:
 - Required: False
 - Secret: False


**private_key** (string): use this key to authenticate the ssh connection
 - Default:
 - Required: False
 - Secret: True


**requirements** (string): path to python requirements
 - Default:
 - Required: False
 - Secret: False


**scp_extra_args** (string): specify extra arguments to pass to scp only
 - Default:
 - Required: False
 - Secret: False


**sftp_extra_args** (string): specify extra arguments to pass to sftp only
 - Default:
 - Required: False
 - Secret: False


**ssh_common_args** (string): specify common arguments to pass to sftp/scp/ssh
 - Default:
 - Required: False
 - Secret: False


**ssh_extra_args** (string): specify extra arguments to pass to ssh only
 - Default:
 - Required: False
 - Secret: False


**skip_tags** (array): only run plays and tasks whose tags do not match
 - Default:
 - Required: False
 - Secret: False


**start_at_task** (string): start the playbook at the task matching this name
 - Default:
 - Required: False
 - Secret: False


**syntax_check** (boolean): perform a syntax check on the playbook
 - Default: False
 - Required: False
 - Secret: False


**tags** (array): only run plays and tasks tagged with these values
 - Default:
 - Required: False
 - Secret: False


**timeout** (number): override the connection timeout in seconds
 - Default: 0
 - Required: False
 - Secret: False


**user** (string): connect as this user
 - Default:
 - Required: False
 - Secret: False


**vault_id** (string): the vault identity to use
 - Default:
 - Required: False
 - Secret: False


**vault_password** (string): the vault password to use
 - Default:
 - Required: False
 - Secret: True


**verbose** (number): level of verbosity, 0 up to 4
 - Default: 0
 - Required: False
 - Secret: False
