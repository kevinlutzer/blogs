# Overview

- Why would we want to remotely program an arduino device
- We want to do the compilation on your main computer for performance reasons. 
- Mention that this setup is completely headless! No remote desktop needed. 
- Talk about how we can automate all of this with vscode tasks.

# Requirements

- Base requirements for this tutorial, this is just the arduino, a usb cable to power your raspberry pi, 
and a usb cable to power the arduino. 
- Talk about how the host computer and how we can use it to program the primary device

# Setting up a Raspberry

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
