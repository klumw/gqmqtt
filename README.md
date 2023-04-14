## USB serial to MQTT bridge for GQ GMC-500+ Geiger Counter

The aim of this project is to enable easy integration of the [GQ GMC-500+](https://www.gqelectronicsllc.com/comersus/store/comersus_viewItem.asp?idProduct=5631) Geiger Counter device into home automation systems (e.g. [openHAB](https://www.openhab.org/)).  
The software reads data from the device via USB and transfers the data to a MQTT host.  
Since most home automation systems support the MQTT protocol, sensor integration should be straightforward. 

The bridge software can be compiled for Windows and Linux.
See also the installation section for [Raspberry PI](https://www.raspberrypi.org/) 3 and 4. You can also download a ready to use [executable for Raspberry PI OS](https://github.com/klumw/gqmqtt/releases).

For compiling the software the installation of the latest version of the [GO](golang.org) programming language is necessary. 


## Installation on Raspberry Pi 3 and 4
1. Get the [latest version](https://www.raspberrypi.com/software/operating-systems/) of [Pi OS](https://www.raspberrypi.com/software/).
2. Run *sudo apt update* then *sudo apt upgrade*
3. Follow the [GO installation instructions](https://shores.dev/install-go-language-on-raspberry-pi-3-and-4/)
4. Run *go build* inside the source folder.
 If everything went well you will get a **gqmqtt** executable in the src folder.
5. Make sure the Geiger Counter device is switched on and is connected to your Raspberry Pi via USB. Use only the supplied USB cable if possible.
6. In the device settings the option **Third party output** must be switched off. Baud rate should be set to the default value (115200).
7. Make sure user is member or group "dialout"
8. For a quick test start your software with the command *./gqmqtt -v*, this will start the bridge in verbose mode. Type *--help* to get an overview of all available command line flags.
9. If your mqtt broker runs on another host you will need to set up the host url with the -h flag
   (e.g. *-h tcp://192.168.178.25:1883*)
9. You can also [install the software as a service](https://domoticproject.com/creating-raspberry-pi-service/) for automatic start and restart.
