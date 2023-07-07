

# Programming An Arduino Remotely Using a Raspberry Pi


Why would you even want to remotely program your Arduino device? Well for me, my electronic test gear like power supplies, oscilloscopes and analyzers are in a different part of my work room from where my computer is. So having the ability to remotely program my Arduino means I don't have to run excessively long USB cables, or move equipment I am using to develop my project with the Arduino near my computer. When working with signals that are dangerous to expose to your computer's USB bus, it's important to have some layer of isolation. Consider this if you are working on an embedded project where a 120 VAC signal is being manipulated -- one small wiring mistake and you expose your computer's to this dangerous signal!


I am going to go over how to program an Arduino using a Raspberry Pi without needing to remote desktop into it! This will allow us to run a few simple commands on a Linux or Mac command line for each time we want to compile and upload code to our Arduino device.


## How Will This Work?


We are doing to use the [https://arduino.github.io/arduino-cli/0.33/](arduino-cli) to compile our code on our main computer, then we will copy it to the Raspberry Pi. After we copy it, we will use the Arduino Cli on the Raspberry Pi to upload it to our Arduino. We want to compile on our main computer as it saves time since the Raspberry Pi is a lot slower. Here is a flow chart to show this workflow in action. 


![Alt text](./assets/programming-paradigm.jpg "a title")


## Requirements


So to follow along you will need a Raspberry Pi -- any model that will have access to wifi will work! You will also need a cable to connect your Arduino board to the Raspbery Pi, a cable to power your Raspberry Pi, and an SD card to boot your Pi from.


## Install the Raspberry Pi Image and Setup the Arduino Cli


On your development computer you will also need to install the [https://www.raspberrypi.com/software/](Raspberry Pi Imager) to create the bootable SD for the Pi. Click on the download button for your OS. If you are not using either a Mac, or a computer running Ubuntu, you will have to manually download the latest `.img` file for the Raspberry Pi OS, tweak the dhcp config for your wifi network and add the SSH file.


To install the Arduino Cli on OSX, use homebrew:


```
brew update
brew install arduino-cli
```


Otherwise you can use the manual install script for Linux:


```bash
curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | BINDIR=<Your Install Directory> sh
```


Make sure that the directory you installed the Arduino Cli to is in your computer cli's  `PATH` variable. Now we need to install the board files for the specific Arduino Device you are using. For me I am using an Arduino Nano, so I will have to run:


`arduino-cli core install arduino:avr`


You can find the name of the package you have to install by running:


`arduino-cli core search <keywords>`


## Setting up a Raspberry Pi


Start up the `Raspberry Pi Imager` on your main computer. After selecting the operating system to be `Raspberry Pi OS Lite (32-bit)` under the `Raspberry Pi OS (Other)` menu and the SD card with `Choose Storage`, click on the settings icon on the bottom right of the application. From there you might be required to enter your password to allow the tool to get access to the wifi password of the current network your computer is on. You can either enter your password, or hit `no`. If you click no, you will have to enter that information later on. In the settings menu, make sure that SSH enabled, the correct wifi credentials as specified, as well you have selected `Allow for public key authentication only`. We want this option selected so that we don't need to enter a password each time we want to execute a command or copy data to the Pi. Now in that menu you can hit `save` and then `write`. This will take a few minutes so go grab some coffee ☕ while you wait!


Once the OS has been written to the SD card, remove it, put it into the Pi and power it on. Once the Pi boots you should be able to SSH into by running the following command:


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


That's it! We are ready to start compiling and uploading code!


## Using the Arduino CLI to programing the Arduino


On your computer create a directory called blink with a file called blink.ino in it.  Also in that directory add a folder called `bin`.  Add the following contents to the blink.ino file.


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


In the blink directory, we can run the compile command to build our sketch and place the created binaries in the `bin` folder. Since I am going am using the Arduino Nano under the AVR family, the command I will have to run is:


`arduino-cli compile -b arduino:avr:nano --output-dir=./bin`


You will have to replace the board name with whatever board you are specifically using.


We need to make sure that the Raspberry Pi has a folder on it to transfer our binaries too. We are going to use SSH to accomplish this. Note that you only need to do this once. To create this directory from your computer, run:


`ssh pi@raspberrypi.local 'mkdir -p ~/blink'`


We are going to use secure copy to upload our binary files to the Raspberry Pi. To upload the binaries to the `~/blink` directory on the Pi from your computer, run:


`scp bin/* pi@raspberrypi.local:~/blink`


Now to upload the code from the Raspberry Pi we can use the Arduino Cli upload command. From your computer you can just run:


`ssh pi@rasperrypi.local './home/pi/local/bin/arduino-cli upload -b arduino:avr:nano -P /dev/USB0'`.


In my case, my Arduino Nano is mounted to `/dev/USB0` on my Pi. You will have to replace this with the port your programmer is mounted too. The Arduino Cli should detect what programmer is attracted to that port and run the appropriate upload commands. If this command is not working, add a `-P` option to the command with the name of the programmer you are using.


If you want to run only one command instead of three each time, you can concatenate the command with the `&&` bash operator:


`arduino-cli compile -b arduino:avr:nano --output-dir=./bin && scp bin/* pi@raspberrypi.local:~/blink && ssh pi@rasperrypi.local './home/pi/local/bin/arduino-cli upload -b arduino:avr:nano -P /dev/USB0'`.


This will run all three commands synchronously. Be careful though, if one of those commands fails, the next won't run.


There you go! You now have a blink sketch running on your Arduino that you uploaded remotely from your main computer!
