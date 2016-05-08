
# Z-Way to MQTT

A simple Z-Way to MQTT bridge in GO.

[![Build Status](https://drone.io/github.com/cblomart/zwaymqtt/status.png)](https://drone.io/github.com/cblomart/zwaymqtt/latest)

[![Build Status](https://travis-ci.org/cblomart/zwaymqtt.svg?branch=master)](https://travis-ci.org/cblomart/zwaymqtt)

[Latest binaries](https://drone.io/github.com/cblomart/zwaymqtt/files)

The service sends MQTT message for all handled classes to the a MQTT broker.

Classes can be handled either read only (ro), typicaly for sensors, or read only (ro) ,typicaly for switches.

There is a base subject for MQTT (root topic; per default "razberry").

The nodename will be read form the Z-Way node.

In principle nodes will be split between actuators and sensors and between binray and analogic ones.

The instance is always used. (i)

Currently the bridge is limited to certain ZWave classes (see bellow).

# Usage

> Usage of zwaymqtt:
>   -h string
>     	mqtt topic root or ZWAY_HOME environment variable (default "razberry")
>   -m string
>     	MQTT server or MQTT_SERVER environment variable (default "localhost:1883")
>   -mp string
>     	MQTT password or MQTT_PASSWORD environment variable
>   -mu string
>     	MQTT username or MQTT_USERNAME environment variable
>   -p string
>     	Z-Way passsword or ZWAY_PASSWORD environment variable
>   -profile string
>     	Profile execution (cpu/mem/all)
>   -r int
>     	Z-Way refresh rate in seconds or ZWAY_REFRESH environment variable (default 30)
>   -s string
>     	Z-Way server name or ZWAY_SERVER environment variable (default "localhost:8083")
>   -u string
>     	Z-Way username or ZWAY_USERNAME environment variable (default "admin")
>   -v	Show debug messages

notice that main parameters are mapped by an environment variable

## Instalation

### from release

On your pi:

> $ wget https://github.com/cblomart/zwaymqtt/releases/download/0.1/zwaymqtt-linux-arm5.tgz
> 
> sudo tar -zxvf ./zwaymqtt.tgz -C /
> 
> sudo vi /etc/default/zwaymqtt
> 
> sudo systemctl enable zwaymqtt
> 
> sudo systemctl start zwaymqtt

### from sources

I install the software at home with the RaZberry.

To do this:

- be sure to have GO and ability to cross compile. (If you do not, you can alway run it with go on a pc)

- compile for an arm running linux (raspberry) on a pc:

> $ GOOS=linux GOARCH=arm GOARM=5 go get github.com/cblomart/zwaymqtt

- copy the necessary files to your Pi:

> $ scp $GOPATH/src/github.com/cblomart/zwaymqtt/etc/systemd/system/zwaymqtt.server pi@raspberry.local:/tmp/
> 
> $ scp $GOPATH/src/github.com/cblomart/zwaymqtt/etc/default/zwaymqtt pi@raspberry.local:/tmp/
> 
> $ scp $GOPATH/bin/linux_arm/zwaymqtt pi@raspberry.local:/tmp/zwaymqtt.bin

- on your pi, place the files at the right places:

> $ sudo cp /tmp/zwaymqtt.server /etc/systemd/system/
> 
> $ sudo cp /tmp/zwaymqtt /etc/default/
> 
> $ sudo cp /tmp/zwaymqtt.bin /usr/local/bin/zwaymqtt 
> 

- on your pi, edit the /etc/default/zwaymqtt to match your preferences

- on your pi, enable and start the service:

> $ sudo systemctl enable zwaymqtt
> 
> $ sudo systemctl start zwaymqtt





## Zwave Classes

### BATTERY

Encompassed classes:

- COMMAND\_CLASS\_BATTERY (ro)

The events on this class will be mapped to the "\<root\_topic\>/sensors/analogic/\<nodename\>/\<i\>/battery" topic.

i.e.: razberry/sensors/binary/detector_door_basement/0/general_purpose

### SWITCH

Encompassed classes:

- COMMAND\_CLASS\_SWITCH\_BINARY (rw)

The events on this class will be mapped to the "\<root\_topic\>/actuators/binary/\<nodename\>/\<i\>/switch" topic.

i.e.: razberry/actuators/binary/binary_switch_living/1/switch

### MULTILEVEL SWITCH

Encompassed classes:

- COMMAND\_CLASS\_SWITCH\_MULTILEVEL (rw)

The events on this class will be mapped to the "\<root\_topic\>/actuators/analogic/\<nodename\>/\<i\>/dimmer" topic.

**TODO**: put an example (i have no multilevel switch to test)

### BINARY SENSOR

Encompassed classes:

- COMMAND\_CLASS\_SENSOR\_BINARY (ro)

The utility will be determined by the sensor type described on the node. If it is a generic sensor... "generic" will be used.

The events on this class will be mapped to the "\<root\_topic\>/sensors/binary/\<nodename\>/\<i\>/\<utility\>" topic.

i.e.: razberry/sensors/binary/detector_entry/0/motion

### MULTILEVEL SENSOR

Encompassed classes:

- COMMAND\_CLASS\_SENSOR\_MULTILEVEL (ro)
- COMMAND\_CLASS\_METER (ro)

The utility is still used but the scale type has been added.

The events on this class will be mapped to the "\<root\_topic\>/sensors/binary/\<nodename\>/\<i\>/\<utility\>/\<scale\>" topic.

i.e.: razberry/sensors/analogic/detector_entry/0/temperature/c

 License

The MIT License (MIT)

Copyright (c) 2016 cblomart

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
