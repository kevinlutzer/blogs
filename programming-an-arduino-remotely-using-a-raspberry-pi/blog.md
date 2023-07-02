# Overview

There are a lot of reasons as too why you would want to remotely program your Arduino boards! For me specifically my test bench with my power supplies, oscillscope, and spectrum analyzer is in a seperate part of my work room from where my computer is. You could also just want the electrical isolation from your computer, say if you are working on embedded project where a 120 VAC signal is being manipulated -- One small mistake and you expose your computer's USB bus to this dangerous signal!

I am going to go over how to program an arduino using a raspberry pi without needing to remote desktop into it! This will allow us to run a few simple bash commands for each time we want to compile and upload code to our arduino. 


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
