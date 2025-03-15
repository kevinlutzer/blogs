# RPi Image Initialization with Cloud Init and The RPi Image

Often, you need to flash multiple Raspberry Pis for use in a cluster or distributed processing application.
After flashing, additional configuration steps may be required to complete the setup.
Flashing and configuring multiple Raspberry Pis can be a tedious process. luckily, there are tools like Cloud Init
that simplify OS initialization! In this tutorial I will show you how to use Cloud Init to initialize a Ubuntu 24.04 Server
with:

1. basic network configuration
1. system (apt) packages
1. files

The OS and any configuration files will be flashed onto the boot device via the [RPi Imager](https://github.com/raspberrypi/rpi-imager)

Lets get started!

## Cloud Init

Cloud Init provides a templating mechanism for operating system initialization, automating the configuration process for
 each device. Cloud Init runs on first boot of the device, and will block normal start up process until
 it completes the initialization of the OS based on the `user-data` and `network-config` yaml files. The outout of
  each command ran can be found in the `/var/log/cloud-init-output.log` file.

### User Data

The `user-data` file is where you can:

- configure new users and groups
- add apt repositories
- write arbitrary files
- execute any commands with sudo permissions

In a directory of your choice, create a yaml file called `user-data`.

### Users

 Lets start by adding the default user a new user
 called `work`. In the follow code snippet I will show how to add a new user with a password, default ssh credentials
 and appended groups.

``` yaml
users:
    - default # This is the default ubuntu user
    - name: work
      groups: users,disk,dialout # dialout and disk give the work user access to storage and any connected serial terminals. 
      plain_text_passwd: newpassword # this will be the password for the user you would login with ssh
      sudo: ALL=(ALL) NOPASSWD:ALL # similar to how the default user is
      ssh_import_id:
        - gh:kevinlutzer # change this to be your Github user, it will add your public key so you can SSH into the Rasberry Pi with the same Github private keys
```

### Apt Packages

Lets add the `apt` source for docker as well as install some packages. Append the following code snippet to the
 `user-data` file.

``` yaml
apt:
    sources:
        docker.list:
            source: deb [arch=amd64] https://download.docker.com/linux/ubuntu $RELEASE stable # Where apt can find docker
            keyid: 9DC858229FC7DD38854AE2D88D81803C0EBFCD88 # The ID of the GPG key docker uses

packages:
    - docker-ce
    - docker-ce-cli
    - build-essential # C toolchain
    - libssl-dev # dev package for Openssl development

```

### Creating Files

Let's create a sample script file that Cloud Init will run. Append the following to your `user-data` file.

``` yaml
write_files:
    - content: |
        echo "Hello World"
    owner: root:root
    permissions: '0755'
    path: /opt/hello_world
```

### Executing Commands

After Cloud Init has setup all of the configuration, it will run any commands in the `runcmds` list with sudo permissions.
 Add the following snippet to the `user-data` file to execute the script
 created in the previous section.

``` yaml
runcmds: 
    - ./opt/hello_world
```

## Network Config

Configure static or dynamic network settings to ensure proper connectivity.

## 