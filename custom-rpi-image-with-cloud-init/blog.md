# RPi Image Initialization with Cloud Init and The RPi Imager Tool

Flashing and configuring multiple Raspberry Pis (RPi) can be a tedious process, especially when preparing them for a
 cluster or distributed processing setup. Fortunately,
 tools like Cloud Init make operating system (OS) initialization much easier.

In this tutorial, I’ll show you how to use Cloud Init to initialize a RPi running Ubuntu 24.04 Server with:

1. Basic network configuration
1. Common system (APT) packages
1. Custom files

All needed Cloud Init files and the OS will be flashed directly to the boot device using the RPi Imager.

Let's get started!

## Cloud Init

Cloud Init provides a templating mechanism for automating OS initialization and configuration.
 Cloud Init runs on first boot and blocks the normal startup process until the OS has been fully initialized as described
  by the provided `user-data` and `network-config` YAML files. When the RPi has booted the image, you
   can find the output of each command Cloud Init executed during initialization in
    the `/var/log/cloud-init-output.log` file.

### User Data

The `user-data` file allows you to do things like:

1. Configure users and groups
1. Add `ppa` repositories and install packages with `apt`
1. Write arbitrary files to the system
1. Execute commands with sudo permissions

To get started, create a YAML file named `user-data` in a directory of your choice.

### Users

Let's start by adding the default user and a new user
  called `work`. In the follow code snippet I will show how to add a new user with a password, default SSH credentials
  and appended groups. We will also include the `default` ubuntu user which has a username and password of `ubuntu`.

``` yaml
#cloud-config
users:
  - default # This is the default ubuntu user
  - name: work
    groups: users,disk,dialout # dialout and disk give the work user access to storage and any connected serial devices. 
    plain_text_passwd: newpassword # this will be the password for the user you would login with ssh
    lock_passwd: false # unlock password to be able to lock into user via SSH
    sudo: ["ALL=(ALL) NOPASSWD:ALL"] # all unrestricted sudo access
    ssh_import_id:
      - gh:kevinlutzer # change this to be your Github user, it will add your public Github key so you can SSH into the Raspberry Pi with the same Github private key

ssh_pwauth: true # needed to remotely access the RPi with SSH by password, this is an sshd config
```

Note that the tag `#cloud-init` is needed at the top of the file for it to be a valid Cloud Init `user-data` file.

### Apt Packages

Lets add package repository source for docker as well as install some packages. Append the following code snippet to the
 `user-data` file.

``` yaml
apt:
  sources:
    docker.list:
      source: deb [arch=amd64] https://download.docker.com/linux/ubuntu $RELEASE stable # Where apt can find docker
      keyid: 9DC858229FC7DD38854AE2D88D81803C0EBFCD88 # The ID of the GPG key docker uses

package_update: true
package_upgrade: true

packages:
  - docker-ce
  - docker-ce-cli
  - build-essential # C toolchain
  - libssl-dev # dev package for Openssl development

```

### Creating Files

Let's create a sample script file that we will configure Cloud Init to run in the next step. Append the following to
 your `user-data` file.

``` yaml
write_files:
  - content: |
    #!/bin/bash
    echo "Hello World"
  owner: root:root
  permissions: '0755'
  path: /opt/hello_world
```

### Executing Commands

After Cloud Init has setup all of the configuration, it will run any commands in the `runcmd` list with sudo permissions.
 Add the following snippet to the `user-data` file to execute the script
 created in the previous section.

``` yaml
runcmd: 
  - ./opt/hello_world
```

## Network Config

Lets add some simple network configuration for the RPi so that it will get an IP from a DHCP server on the same
 network, and be accessible with a direct link local connection. Create a new file called `network-config` with the
 following content.

``` yaml
network:
  version: 2
  ethernets:
    eth0:
      dhcp4: true
      optional: true
      addresses: [169.254.0.5/16]
```

The `169.254.0.5/16` is a special IP that is only valid when connected directly to the RPi. This is handy when
 you don't have access to a DHCP server. With a link local connection, you can plug your RPi into your computer and get
  access via SSH using the IP `169.254.0.5`.

## Validating The Schema Of The Cloud Init Configuration Files (Optional)

If you are following this tutorial on a Linux based system, you could install the `cloud-init` tool, and validate your `user-data`
 and `network-config` files. To do this, first install the `cloud-init` tool with one of the following commands:

``` bash
sudo apt install cloud-init # Debian and Ubuntu
sudo dnf install cloud-init # REHL, Rocky, Fedora 
sudo pacman -S cloud-init # Arch 
```

Next run the following commands to validate schema of the `user-data` and `network-config` files.

``` bash
sudo cloud-init schema --config-file=user-data
sudo cloud-init schema -t network-config --config-file=network-config
```

Note that this does not validate that Cloud Init will be able to execute all the operations defined by the `user-data`
 and `network-config` files.

## Flashing Raspberry Pi Boot Media

Grab your SD card or USD based storage and plug it into your computer. We will be using the Raspberry Pi Imager tool's CLI
 to flash the device as the Raspberry Pi Imager's GUI does not have functionality we need. Get the path to the
 storage device on your Mac or Linux computer. For a Mac the path will have the format `/dev/diskN` and for Linux the path
 will be `/dev/sdX`. Run the following command to flash the device.

```bash
/Applications/Raspberry\ Pi\ Imager.app/Contents/MacOS/rpi-imager --cli https://cdimage.ubuntu.com/releases/24.04/release/ubuntu-24.04.2-preinstalled-server-arm64+raspi.img.xz /dev/disk4 --cloudinit-userdata user-data --cloudinit-networkconfig network-config # mac
rpi-imager --cli https://cdimage.ubuntu.com/releases/24.04/release/ubuntu-24.04.2-preinstalled-server-arm64+raspi.img.xz /dev/sdb --cloudinit-userdata user-data --cloudinit-networkconfig network-config # linux
```

Once the Raspberry Pi Imager has flashed the storage device, plug the storage device, network cable, and power supply
 into the Raspberry Pi. After a couple of minutes the initialized Raspberry Pi will show up on your network! To view the
  operations Cloud Init completed, SSH into the RPi and read the
   `/var/log/cloud-init-output.log` file.
