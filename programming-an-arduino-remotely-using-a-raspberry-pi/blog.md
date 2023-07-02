# Overview

Why would you even want to remotely program your Arduino? Well for me specifically, my electronic test gear like power supplies, oscillisopes and analyzers are in a different part of my work room from where my computer is. Have isolation between the computer you use for development and the embedded system you are building is also important when working with larger voltages, signals that oscillocate, or many other reasons @KLUTZER TODO. Say if you are working on embedded project where a 120 VAC signal is being manipulated -- One small mistake and you expose your computer's USB bus to this dangerous signal!

I am going to go over how to program an Arduino using a Raspberry Pi without needing to remote desktop into it! This will allow us to run a few simple bash commands for each time we want to compile and upload code to our Arduino. The other benefit is that by using our main development computer to compile the code we can save compilation time. 

# How Will This Work? 

We are doing to use the [https://arduino.github.io/arduino-cli/0.33/](arduino-cli) to compile our code, then we will copy it to the Raspberry Pi. After we copy it, we will then upload it to our Arduino. Here is a flow chart to show this workflow in action. 

![Alt text](./assets/programming-paradigm.jpg "a title")

# Requirements

So to follow along you will need a Raspberry Pi, any model that will have access to wifi will work! You will also need a cable to connect your Arduino board to the Raspbery Pi, and a cable to power your Raspberry Pi. On your development computer you will also need to install the [https://www.raspberrypi.com/software/](Raspberry Pi Imager) to create the bootable SD for the Raspberry Pi, as well as the [https://arduino.github.io/arduino-cli/0.33/](arduino-cli) to be able to compile your commands. Note we will installing the arduino cli on the raspberry pi as well. Note that the commands I am going to give will work for Mac and Linux users, but they can be adapted for Windows as well.

# Setting up a Raspberry Pi

Lets go over how to setup the Raspberry Pi. To prevent the need to enter a password each time we want to either execute a command or copy data to the Raspberry Pi, we are going to setup an SSH Key. After selecting the operating system to be `Raspberry Pi OS Lite (32-bit)` un the `Raspberry Pi OS (Other)` and the SD card with `Choose Storage`, click on the settings icon on the bottom right of the application. From there you might be required to enter your password to allow the tool to get access to the wifi password of the current network you are on. You can either enter your password, or hit `no`. From there, click the settings icon on the bottom right hand side of the window. In this new menu make sure that the OS has SSH enabled, the correct wifi credentails, as well you have selected `Allow for public key authentication only`. This should autopopulate with your computer's correct hostname. Now in that menu you can hit `save` and then `write`. This will take a few minutes so go grab some coffee ☕ while you wait!

Once the OS has been written to the SD card, remove it, put it into the Pi and power it on. Once the little computer boots you should be able to ssh into by running the following command:

`ssh pi@raspberrypi.local`

From there update the Pi's package registry, as well the OS. 

```bash
sudo apt-get update
sudo apt-get -y upgrade
```

Now we can install the Arduino Cli and add it to our path by running. 

```bash
BINDIR=~/local/bin
mkdir -p $BINDIR
curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | BINDIR=$BINDIR sh
echo 'export PATH=$PATH:$BINDIR' >>> ~/.bashrc
source ~/.bashrc
```

This just makes the `~/local/bin` path on the Pi, installs the Arduino CLI in it, and then adds that path to the bash shell's `PATH` environment variable. 

Thats it! We are ready to start compiling and uploading code!

# Using the Arduino CLI to programing the Arduino

The Arduino CLI will allow us to both compile code on our host computer 

So on your development computer create a directory called blink with a file called blink.ino in it. Also in that directory add a folder called `bin` Add the following contents to that file.

```c++
void setup() {
  pinMode(LED_BUILTIN, OUTPUT);
}

void loop() {
  digitalWrite(LED_BUILTIN, HIGH);
  delay(1000);                     
  digitalWrite(LED_BUILTIN, LOW);   
  delay(1000);                    
}
```


- What is a fully qualified board name
- Showing how you can find out what board you want to talk to remotely.

# Uploading a sketch remotely 

- build the sketch on the main computer, 
- copy the files over using scp 
- upload the files using the upload command
- verify

# Automating this with vscode tasks

- Just give them the json file with the necessary tasks they can add to their workspace
