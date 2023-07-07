
Why would you even want to remotely program your Arduino device? Well for me, my electronic test gear like power supplies, oscillisopes and analyzers are in a different part of my work room from where my computer is. So having the ability to remotely program my arduino means I don't have to run excessively long USB cables, or move equipment I am using to develop my project with the Arduino near my computer. Having isolation between the computer you use for development and the embedded system you are building is also important when working with signals that you wouldn't want to expose to your computer's USB bus. Consider this if you are working on embedded project where a 120 VAC signal is being manipulated -- one small mistake and you expose your computer's to this dangerous signal!

I am going to go over how to program an Arduino using a Raspberry Pi without needing to remote desktop into it! This will allow us to run a few simple commands in a linux or mac command line for each time we want to compile and upload code to our Arduino device. 

# How Will This Work? 

We are doing to use the [https://arduino.github.io/arduino-cli/0.33/](arduino-cli) to compile our code on our main computer, then we will copy it to the Raspberry Pi. After we copy it, we will use the Arduino Cli on the Raspberry Pi to upload it to our Arduino. We want to compile on our main computer as it saves time since the Raspberry Pi is a lot slower. Here is a flow chart to show this workflow in action.  

![Alt text](./assets/programming-paradigm.jpg "a title")

# Requirements

So to follow along you will need a Raspberry Pi -- any model that will have access to wifi will work! You will also need a cable to connect your Arduino board to the Raspbery Pi, and a cable to power your Raspberry Pi. 

# Install the Raspberry Pi Image and Setup the Arduino CLI

On your development computer you will also need to install the [https://www.raspberrypi.com/software/](Raspberry Pi Imager) to create the bootable SD for the Pi. Click on the download button for your OS. If you are not using either a OSX, or Ubuntu, you will have to manually download the latest `.iso` for the Raspberry Pi OS, tweak the dhcp config for your wifi network and add the ssh file.

To install the Arduino CLI on OSX, use homebrew:

```
brew update
brew install arduino-cli
```

Otherwise you can use the manually install script for Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | BINDIR=<Your Install Directory> sh
```

Make sure that where ever you install the command, is in your `PATH` variable.

# Setting up a Raspberry Pi

Lets go over how to setup the Raspberry Pi. Start up the `Raspberry Pi Imager` on your main computer. After selecting the operating system to be `Raspberry Pi OS Lite (32-bit)` under the `Raspberry Pi OS (Other)` menu and the SD card with `Choose Storage`, click on the settings icon on the bottom right of the application. From there you might be required to enter your password to allow the tool to get access to the wifi password of the current network your computer is on. You can either enter your password, or hit `no`. If you click no, you will have to enter that information later on. In the settings menu, make sure that SSH enabled, the correct wifi credentails as specified, as well you have selected `Allow for public key authentication only`. We want this option selected so that we don't need to enter a password each time we want to execute a command or copy data to the Pi. Now in that menu you can hit `save` and then `write`. This will take a few minutes so go grab some coffee ☕ while you wait!

Once the OS has been written to the SD card, remove it, put it into the Pi and power it on. Once the Pi boots you should be able to ssh into by running the following command:

`ssh pi@raspberrypi.local`

From there update the Pi's Aptitude package registry, as well as update existing packages on the OS to the newest versions.

```bash
sudo apt-get update
sudo apt-get -y upgrade
```

Now we can install the Arduino Cli to `~/local/bin` and add that bath to the shell's `PATH` environment variable.

```bash
BINDIR=~/local/bin
mkdir -p $BINDIR
curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | BINDIR=$BINDIR sh
echo 'export PATH=$PATH:$(BINDIR)' >> ~/.bashrc
source ~/.bashrc
```

Thats it! We are ready to start compiling and uploading code!

# Using the Arduino CLI to programing the Arduino

The Arduino CLI will allow us to both compile code on our host computer 

So on your computer create a directory called blink with a file called blink.ino in it. Also in that directory add a folder called `bin` Add the following contents to that file.

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

From here we can 

- What is a fully qualified board name
- Showing how you can find out what board you want to talk to remotely.

# Uploading a sketch remotely 

- build the sketch on the main computer, 
- copy the files over using scp 
- upload the files using the upload command
- verify

# Automating this with vscode tasks

- Just give them the json file with the necessary tasks they can add to their workspace
