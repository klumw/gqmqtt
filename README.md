## USB serial to MQTT converter for GQ GMC-500+ Geiger Counter

The aim of this project is to enable an easy integration of the GQ GMC-500+ Geiger Counter into a home automation system (e.g. **openHAB**).  
The software reads out the measured values via the USB serial interface and transfers the values to a MQTT host.  
Since most home automation systems support the MQTT protocol, integration should be easy. 

The software can be compiled for Windows and Linux, further below the installation for Raspberry PI-3 and PI-4 is explained.

For compiling the software the installation of a current version of the [GO](golang.org) programming language is necessary. 


## Installation on Raspberry PI-3 and PI-4
You will need an up to date version of Pi OS.
1. Run *sudo apt update* then *sudo apt upgrade*
2. Follow the [GO installation instructions](https://shores.dev/install-go-language-on-raspberry-pi-3-and-4/)
3. Run *go build* inside the source folder.
 If everything is fine you will get a **gqmqtt** executable.
4. Make sure the device is switched on and is connected to your Raspberry via USB.
5. In the device settings the option **Third party output** must be switched off. Baud rate should be set to the default (115200).
6. For a quick test start your software with the command *./gqmqtt -v*, this will start the converter in verbose mode. Type *--help* to get an overview of all available command line flags.
7. If your mqtt broker runs on another host you need to set the host url with the -h flag
   (e.g. *-h tcp://192.168.178.25:1883*)
8. You can [install the software as a service](https://domoticproject.com/creating-raspberry-pi-service/), so that it is automatically started on startup or after a failure.

