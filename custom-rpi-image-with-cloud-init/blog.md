# RPi Image Initialization with Cloud Init and The RPi Image

Often, you need to flash multiple Raspberry Pis for use in a cluster or distributed processing application.
After flashing, additional configuration steps may be required to complete the setup.
Flashing and configuring multiple Raspberry Pis can be a tedious process, but luckily, there are tools like Cloud Init
that simplify OS initialization! In this tutorial I will show you how to use Cloud Init to initialize a Ubuntu 24.04 Server
with:

1. basic network configuration
1. common system (apt) packages
1. files

The OS and any configuration files will be flashed onto the boot device via the [RPi Imager](https://github.com/raspberrypi/rpi-imager)

Lets get started!

## Cloud Init

Cloud Init provides a templating mechanism for operating system initialization, automating the configuration process for
 each device. Cloud Init runs on first boot of the device, and will block normal start up process until
 it completes the initialization of the OS based on the `user-data` and `network-config` yaml files. The output of
  each command ran can be found in the `/var/log/cloud-init-output.log` file of the booted system.

### User Data

The `user-data` file is where you can:

- configure new users and groups
- add apt repositories and packages
- write arbitrary files
- execute any commands with sudo permissions

In a directory of your choice, create a yaml file called `user-data`.

### Users

 Lets start by adding the default user and a new user
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
        - gh:kevinlutzer # change this to be your Github user, it will add your public Github key so you can SSH into the Raspberry Pi with the same Github private key
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

Lets add some simple configuration for the Raspberry Pi so that it will get allocated a IP from a DHCP server on the same
 network, and be accessible with a direct "Link Local" connection. Create a new file called network with the following content.

``` yaml
network:
  version: 2
  ethernets:
    - eth0:
      dhcp4: true
      optional: true
      addresses: [ 169.254.0.5/16 ]
```

The `169.254.0.5/16` is a special IP that is only valid when connected directly to the Raspberry Pi. This is handy when
 you don't have access to a DHCP server as you can plug your Raspberry Pi into your computer and get access via SSH.

## Flashing Raspberry Pi Boot Media

Grab you SD card or USD based storage and plug it into your computer. We will be using the Raspberry Pi Imager tool's CLI
 to flash the device. Thats because the Raspberry Pi Imager's GUI does not have functionality we need. Get the path to the
 storage device on your Mac or Linux computer. For a Mac the path will have the format `/dev/diskN` and for Linux the path
 will be `/dev/sdX`. Run the following command to flash the device.

```bash 
/Applications/Raspberry\ Pi\ Imager.app/Contents/MacOS/rpi-imager --cli https://cdimage.ubuntu.com/releases/24.04/release/ubuntu-24.04.2-preinstalled-server-arm64+raspi.img.xz /dev/disk4 --cloudinit-userdata user-data --cloudinit-networkconfig network-config # mac
rpi-imager --cli https://cdimage.ubuntu.com/releases/24.04/release/ubuntu-24.04.2-preinstalled-server-arm64+raspi.img.xz /dev/sdb --cloudinit-userdata user-data --cloudinit-networkconfig network-config # linux
```

Now plug the storage device, network cable, and power supply into the Raspberry Pi. After a couple of minutes the
 initialized Raspberry Pi will show up on your network.
