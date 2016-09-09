## OpenRadiant

[![Travis Build Status](https://travis-ci.org/containercafe/containercafe.svg?branch=master)](https://travis-ci.org/containercafe/containercafe)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

OpenRadiant is a modular platform for enterprise container-native
devops. OpenRadiant supports Kubernetes, and will eventually support
and docker APIs on the same infrastructure.  OpenRadiant supports
multi-tenancy and sharding.  The OpenRadiant platform can be extended
by additional componentry, such as: support for using Mesos (including
making Kubernetes act as a Mesos framework), support for using a
CloudFoundry installation as the source of identities, support for
using Neutron tenant networks to isolate OpenRadiant tenants in the
network.

Features of the OpenRadiant platform include:
* Kubernetes
* Eventually Swarm (original, not "swarm mode" introduced in Docker 1.12)
* Multi-tenancy - with or without Bring-Your-Own-IPv4
* Multi-sharding
* HA control plane in each shard
* Ansible-based create/update/destroy tooling
* Openness to a variety of sources of authentication
* Control plane secured by TLS
* Openness to support for a variety of Software-Defined-Network technologies
* Live container crawling

OpenRadiant is a work in progress, as a collaboration between a
central OpenRadiant team and other teams that use it to build more
specific and capable platforms.

* [Architecture Overview](#architecture-overview)
* [Tiny Example](examples/tiny-example.md)
* [Ansible Principles](#ansible-principles)
* [The container images](docs/building-images.md)
* [The installer machine](#the-installer-machine)
* [Installing OpenRadiant](#installing-openradiant)
* [Code of Conduct](#code-of-conduct)
* [Contributing to the project](#contributing-to-the-project)
* [Maintainers](#maintainers)
* [Communication](#communication)
* [Proxy documentation](proxy/README.md)
* [Mesos clues](mesos-clues.md)
* [Live Container crawling](#live-container-crawling)
* [Learn concepts and commands](#learn-concepts-and-commands)
* [License](#license)
* [Issues](#issues)

### Architecture Overview

OpenRadiant is software that you can use, or extend to produce more
capable software that you use, to operate an enterprise
container-native devops service.

One operating instantiation of the full platform is called an
environment, and it contains one or more shards that operate
independently of each other.  Each shard provides an independent
installation of Kubernetes, and with extension additional platforms
such as Mesos and/or Swarm.  There is an outer control plane with a
proxy API server that implements the Kubernetes (and eventually
Docker/SwarmV1) APIs --- with appropriate restrictions and extensions
--- by appropriately transforming each RPC and dispatching it to the
appropriate shard.

You can subset OpenRadiant so that it creates just one shard.  You can
subset OpenRadiant to omit API proxy if you are not interested in
multi-sharding nor the conveniences it supplies for multi-tenancy.

In a shard there are worker nodes and master nodes.  The Kubernetes
(and eventually Swarm) workload is dispatched to the worker nodes.
The master nodes run the Kubernetes, eventually Swarm, and any
extension's master components in an HA configuration.  We are working
on a better solution than Mesos to coordinate resource allocation
decisions between Kubernetes and Swarm.

OpenRadiant includes Ansible-based installation technology to
instantiate an OpenRadiant environment.  An installer machine acts as
Ansible controller to install OpenRadiant in a target environment.
The installation is parameterized by some Ansible variables files that
describe the desired target environment.

OpenRadiant deploys Kubernetes in containers.  You can choose whatever
version you want.  Your configure your choice of image (including
tag).  See the doc about
[the relevant configuration variables](docs/ansible.md#primary-shard-variables-that-have-defaults)
and
[where to put your settings for those variables](docs/ansible.md#additional-files-for-setting-ansible-variable-values).

The Ansible playbooks strive to meet the Ansible ideal of achieving a
prescribed desired state, and can thus be used to update as well as
install.  However, there are limits to the space of initial states
with which these playbooks can cope.

See [the architecture doc](docs/architecture.md) for more details.


### Ansible principles

The structure of the Ansible playbooks and roles, and the related
conventions for how environments and shards are organized and how they
are described by files, are important parts of OpenRadiant's
modularity.  This includes the inventory contracts that establish the
orthogonality between (a) how machines are provisioned and (b) how
software is installed on them.  This also includes the concept of
networking plugins, the convention for where temporary and
not-temporary files are kept on the installer machine(s), and the
considerations for deploying in a context where external repositories
can not be reached.  See [the Ansible doc](docs/ansible.md) for
details.


### The installer machine

Following are instructions on how to create a usable installer
machine.  In the near future we hope to also provide a Docker image of
a usable installer machine.

All the controller machine really needs is to have a copy of
OpenRadiant and be able to run Ansible commands using that.  The
following shows ways to accomplish these using Linux shell commands.
For other operating systems you could do something analogous.

The controller machine must have git installed.

Checkout this project along with all its submodules:

```bash
git clone --recursive https://github.com/containercafe/containercafe.git
cd openradiant
```

The controller machine must have Ansible installed, including its
`netaddr` module.  See
[our Ansible doc](docs/ansible.md#ansible-versions-and-bugs-and-configuration)
for details on Ansible versions.  One way to get these two installed
is to use our `requirements.txt` file, as follows.

```bash
pip install -r requirements.txt
```

In order to use that method, the controller machine must have Python's
`pip` installed.  One way to accomplish that is to install the
`python-pip` package using your operating system's package manager.
Another way is to use Python's `easy_install` to install `pip`.

Another way to get the `netaddr` module is to install `python-netaddr`
using the operating system's package manager.

If the installer machine is running MacOS 10 then it needs to have
done `brew install gnu-tar`, for mysterious reasons.


### Installing OpenRadiant

To create/update an OpenRadiant environment, invoke the
`ansible/env-basics.yml` playbook.  This will create the certificates
and keys that are common throughout the environment, and deploy the
API proxy.  Following is an example invocation.

```bash
cd ansible
ansible-playbook -v env-basics.yml -e env_name==${env_kind}-${env_loc} -e envs=${envs}
```

To create/update an OpenRadiant shard, invoke the `ansible/shard.yml`
playbook.  Following is an example invocation.

```bash
cd ansible
ansible-playbook -v shard.yml -e cluster_name=${shard_name} -e envs=${envs} -e network_kind=flannel
```

The shard name, AKA cluster name, is explained in
[the inventory contract doc](docs/ansible.md#the-inventory-contract)
and this variable must be defined in the Ansible invocation.  Two
other variables must also be defined, as mentioned in that document:
`envs` and `network_kind`.

See also [the deployment doc](docs/deploying.md).


### Live container crawling
It is essentially an agentless system crawler that offers a
native and seamless framework for operational visibility and analytics.
It makes use of virtualization and containerization abstractions together with
introspection techniques to provide complete visibility into running entities like
containers without modifying, instrumenting, or accessing the end user context.

This is an opensource technology. For details visit [agentless-system-crawler](https://github.com/cloudviz/agentless-system-crawler).
For how to use crawler in OpenRadiant environment, see [ link ](crawler/README.md)

### Code of Conduct
Participation in the OpenRadiant community is governed by the OpenRadiant [Code of Conduct](CONDUCT.md)

### Contributing to the project
We welcome contributions to the OpenRadiant Project in many forms. There's always plenty to do! Full details of how to contribute to this project are documented in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

### Maintainers
The project's [maintainers](MAINTAINERS.txt): are responsible for reviewing and merging all pull requests and they guide the over-all technical direction of the project.

### Communication
We use \[TODO] [OpenRadiant Slack](https://OpenRadiant.slack.org/) for communication between developers.

### Learn concepts and commands

Browse TBD to learn more. Here are some topics you may be
interested in:

TBD


### License

Copyright 2015-2016 IBM Corporation

Licensed under the [Apache License, Version 2.0 (the "License")](http://www.apache.org/licenses/LICENSE-2.0.html).

Unless required by applicable law or agreed to in writing, software distributed under the license is distributed on an "as is" basis, without warranties or conditions of any kind, either express or implied. See the license for the specific language governing permissions and limitations under the license.

### Issues

Report bugs, ask questions and request features [here on GitHub](../../issues).
