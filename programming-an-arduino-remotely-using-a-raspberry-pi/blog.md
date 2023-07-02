# Overview

Why would you even want to remotely program your Arduino? Well for me specifically, my electronic test gear like power supplies, oscillisopes and analyzers are in a different part of my work room from where my computer is. Have isolation between the computer you use for development and the embedded system you are building is also important when working with larger voltages, signals that oscillocate, or many other reasons @KLUTZER TODO. Say if you are working on embedded project where a 120 VAC signal is being manipulated -- One small mistake and you expose your computer's USB bus to this dangerous signal!

I am going to go over how to program an Arduino using a Raspberry Pi without needing to remote desktop into it! This will allow us to run a few simple bash commands for each time we want to compile and upload code to our Arduino. The other benefit is that by using our main development computer to compile the code we can save compilation time. 

# How Will This Work? 

We are doing to use the [https://arduino.github.io/arduino-cli/0.33/](arduino-cli) to compile our code, then we will copy it to the Raspberry Pi. After we copy it, we will then upload it to our Arduino. Here is a flow chart to show this workflow in action. 

![Alt text](./assets/programming-paradigm.jpg "a title")

# Requirements

So to follow along you will need a Raspberry Pi, any model that will have access to wifi will work! You will also need a cable to connect your Arduino board to the Raspbery Pi, and a cable to power your Raspberry Pi. On your development computer you will also need to install the [https://www.raspberrypi.com/software/](Raspberry Pi Imager) to create the bootable SD for the Raspberry Pi, as well as the [https://arduino.github.io/arduino-cli/0.33/](arduino-cli) to be able to compile and upload your commands. Note that the commands I am going to give will work for Mac and Linux users, but they can be adapted for Windows as well.

# Setting up a Raspberry

Lets go over how to setup the Raspberry Pi. To prevent the need to enter a password each time we want to either execute a command or copy data to the Raspberry Pi, we are going to setup an SSH Key. After selecting the operating system to be `Raspberry Pi OS Lite (32-bit)` and the  

- Make sure to emphazise the need for ssh keys, that will make the commands easier. 

# Using the Arduino CLI to programing the Arduino

- Installing the arduino cli, 
- Intro in the commands you will be using. 
- What is a fully qualified board name
- Showing how you can find out what board you want to talk to remotely.

# Uploading a sketch remotely 

- build the sketch on the main computer, 
- copy the files over using scp 
- upload the files using the upload command
- verify

# Automating this with vscode tasks

- Just give them the json file with the necessary tasks they can add to their workspace
