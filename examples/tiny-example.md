# Tiny Example

This produces a very simple demonstration shard of two VirtualBox VMs,
one master and one worker, plus another VM for the proxy.  The shard
has Mesos installed, and Kubernetes and Swarm playing nicely together
thanks to Mesos.  The networking is Docker bridge networking (which is
good only for a single-worker deployment; see the documentation of
[networking plugins](../docs/ansible.md#networking-plugins) for more
options).  The Swarm master is modified for multi-tenant use.

Install [git](https://git-scm.com/downloads),
[Vagrant](https://www.vagrantup.com/) and
[VirtualBox](https://www.virtualbox.org/wiki/Downloads). You will need
at least Vagrant 1.8.4 and VirtualBox 5.0.24.

Checkout this project:

```bash
git clone --recursive https://github.com/containercafe/containercafe.git
cd openradiant
```

This example proceeds through five steps, as follows.

1. Provision target machines
2. Prepare the installer machine
3. Use the installer machine to deploy an OpenRadiant shard on target machines
4. (Optionally) deploy the OpenRadiant API proxy
5. Exercise the OpenRadiant shard

For OpenRadiant in general, the first two steps can be done in either
order.


## Provision target machines

This example shows just one way to provision machines for use with
OpenRadiant.  In general, you can use OpenRadiant with any
provisioning technology you like.  See
[the inventory contract](../docs/ansible.md#the-inventory-contract) for
the key idea.

Create target machines with Vagrant.  The following creates three, named
`radiant2`, `radiant3` and `proxy`.

```bash
( cd examples/vagrant; vagrant up )
```
In case you face any issues, please follow [vagrant troubleshooting](vagrant/README.md)


### FYI on SSH to Vagrant/VirtualBox VMs

There are two easy ways to open a connection to a shell in a
Vagrant/VirtualBox VM.  One is provided by Vagrant:

```bash
( cd examples/vagrant; vagrant ssh radiant2 )
```

There's not much magic under that hat. You can do the equivalent
directly using `ssh` as follows.

```bash
ssh -i ~/.vagrant.d/insecure_private_key vagrant@192.168.10.2
```


## FYI on using Docker in the deployed environment

In this example's VMs, the `vagrant` user is not authorized to use
docker.  If you `ssh` to one of these VMs and try to use raw Docker
engine commands (i.e., without the settings explained
[below](#exercise-the-shard-directly)) then you will need to `sudo`
them.


## Prepare the installer machine

To create an installer machine, you can either (1) use another Vagrant
VM we have defined to serve this purpose for you or (2) follow
instructions to make your laptop or machine of choice into an
installer machine.  The installer machine has to be able to run
Ansible, which runs only on Linux and MacOS.


### Create installer VM

```bash
( cd examples/vagrant; vagrant up installer-tiny )
```

That creates an installer VM that is specialized to this example.
Connect to it using SSH as described above.  In that VM you will find
most of the contents of the OpenRadiant repository in `~/openradiant`.


### Manually create installer

If you have created the installer VM using [Create Installer VM](../tiny-example.md#create-installer-vm), you can skip this section.   To manually create the installer, see
[the general documentation of the installer machine](../README.md#the-installer-machine)
for the general story.  Following is one concrete realization of that
story for this example.

If you are running Ubuntu on your installer, you may need to install
the following python packages:

```bash
sudo apt-get install python-pip python-dev
```

OTOH, if you are running MacOS 10 on your installer then you probably
need to:

```bash
brew install gnu-tar
```

In any case, install ansible and its `netaddr` module:

```bash
pip install -r requirements.txt
```

Ansible version 2.1.1 or the latest is recommended. See
[our Ansible documentation](../docs/ansible.md#ansible-versions-and-bugs-and-configuration)
for more details.  If you have a different version and run into issues
try the following:

```bash
pip install --upgrade ansible
```


## Deploy an OpenRadiant environment with one shard

Use Ansible on the installer machine to begin the process of deploying
an environment.  This will create the certificates and keys that are
common throughout the environment, and deploy the API proxy.

If you use the installer VM created from [Create Installer VM](../tiny-example.md#create-installer-vm) as the installer machine, use vagrant ssh to ssh into the installer VM and execute the Ansible scripts as the *vagrant* user:
```bash
vagrant ssh installer-tiny
```

```bash
( cd ansible; \
  ansible-playbook -v -i ../examples/envs/dev-vbox/radiant01.hosts env-basics.yml \
      -e "envs=../examples/envs env_name=dev-vbox" )
```


```bash
( cd ansible; \
  ansible-playbook -v -i ../examples/envs/dev-vbox/radiant01.hosts shard.yml \
      -e "envs=../examples/envs cluster_name=dev-vbox-radiant01 network_kind=flannel" )
```

The `cd` makes Ansible 2 find the `ansible.cfg` supplied by OpenRadiant.

The `envs` variable tells the playbook where to find the files that
define the environment and shard.

The `env_name` variable in the first command tells the playbook which
environment to deploy, and the `cluster_name` variable in the second
command tells the playbook which shard to deploy.

The `networking_kind` variable tells the playbook which networking
plugin to deploy (Ansible technicalities make it impossible for the
playbook to use a definition for this variable placed the environment
or shard variables file --- do not put one there, it will just cause
confusion).

See
[the general doc on deployment](../README.md#installing-openradiant)
for the general story about deploying OpenRadiant.


## One-step create and use installer machine


An alternative to creating and then using the installer machine is to
use another Vagrant/VirtualBox VM that we have prepared for you that
is an installer that deploys the shard as the last startup step.

```bash
( cd examples/vagrant; vagrant up active-installer-tiny )
```


## Optionally deploy the OpenRadiant API proxy

Now you have an option to deploy and use the API proxy.  The shard can
be used directly or through the proxy.  The proxy enables
multi-tenancy, multi-sharding and many other
[features](../docs/proxy.md).  Using the shard directly does not
involve these features.

For details on proxy setup and use, please see
[Proxy documentation](../proxy/README.md).


## Exercise the shard directly

The following describes how to exercise the shard without using the
API proxy and the features it provides.

Open an SSH connection to the master node:
```bash
ssh -i ~/.vagrant.d/insecure_private_key vagrant@192.168.10.2
```

On the master you will find both the `kubectl` and `docker`
(currently 1.11) commands on your `$PATH`.

The Swarm master is configured for multi-tenant use.  To prepare to
use it as a tenant, do this on the master:
```bash
cd; mkdir -p radiant/configs/demo; cat > radiant/configs/demo/config.json <<EOF
{
    "HttpHeaders": {
          "X-Auth-TenantId": "demo"
    }
}
EOF
```

Then you will want these commands on the master:
```bash
export DOCKER_TLS_VERIFY=""
export DOCKER_CONFIG=~/radiant/configs/demo
export DOCKER_HOST=localhost:2375
```

To get a listing of this tenant's containers, issue the following command on the master:
```bash
docker ps
```
At first, there will be none.  So create one, like this:
```bash
docker run --name s1 -d -m 128m busybox sleep 864000
```

Then, get a list of containers, with `docker ps`.  You can inspect its network
configuration from inside, like this:
```bash
docker exec s1 ifconfig
```

You can create a Kubernetes "deployment" with a command like this:
```bash
kubectl run k1 --image=busybox sleep 864000
```

The containers in this deployment will be invisible to Swarm because
they lack the label identifying your tenant.  To make containers
visible to Swarm, make a kubernetes pod as follows.  Create a YAML
file prescribing the pod:
```bash
cat > sleepy-pod.yaml <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: sleepy-pod
  annotations:
    containers-annotations.alpha.kubernetes.io: "{ \"com.ibm.radiant.tenant.0\": \"demo\",  \"OriginalName\": \"sleeper\" }"
spec:
  containers:
    - name: sleeper
      image: busybox
      args:
      - sleep
      - "864000"
EOF
```

Then create the pod:
```bash
kubectl create -f sleepy-pod.yaml
```

Then you can watch for it to come up, with
```bash
kubectl get pod
```

Once the pod is created, you can see it with `docker ps`.  You can
inspect its network configuration from inside, like this:

```bash
kubectl exec sleepy-pod ifconfig
```

In a similar vein, you can ping one of these containers from the
other.  For example (in which the Swarm container's IP address
is 172.17.0.5):
```bash
kubectl exec sleepy-pod ping -- -c 2 172.17.0.5
```


## General clues


### To check the HAproxy statistics using the GUI:

On your local browser, enter the following URL:

master_ip:harproxy_GUI_port/haproxy_stats

Example:
http://192.168.10.2:9000/haproxy_stats (port 9000 is statically assigned)
When prompt for the user_namer:password  use  vagrant:radiantHA


### To view the Mesos web UI

On your local browser visit http://192.168.10.2:5050/
