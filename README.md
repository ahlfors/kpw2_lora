# kpw2_lora
LORA chipset receive and show data with Kindle PaperWhite 2


## Kindle PaperWhite 2

Hardware Version: Kindle 5.4.3.2 (2317380003)

JailBreak

USBNetwork

```shell
ssh root@192.168.31.150


Welcome to Kindle!

#################################################
#  N O T I C E  *  N O T I C E  *  N O T I C E  # 
#################################################
Rootfs is mounted read-only. Invoke mntroot rw to
switch back to a writable rootfs.
#################################################
[root@kindle root]# 


```

**NOTE** <***mntroot rw***> will make rootfs writable


the kpw2/lora is a extension of KUAL


## LORA


### PINs

there is 3 parts, LORA chipset, Raspberry Pi II and Kindle PaperWhite2

- LORA chipset UART (TX, RX) connects to KPW2 UART (RX, TX)

- LORA chipset UART (MD0, MD1, AUX, VCC, GND) connects to RPI2 UART (GPIO, GPIO, GPIO, VCC 5V, GND)

- KPW2 UART (GND) connects to RPI2 UART (GND)



all componets:

![all](/images/lora_pins_all.jpeg)

the lines:

![all](/images/lora_pins.jpeg)


LORA chipset pins:

![chipset pins](/images/lora_chipset_pins.jpeg)


Kindle pins:

![kindle pins](/images/lora_kindle_pins.jpeg)


![kindle pins](/images/lora_kindle_pins2.jpeg)

Raspberry Pi II pins:

![rpi pins](/images/lora_rpi_pins.jpeg)



### Software

usage:

```shell
./lora -h
Usage of ./lora:
  -b int
    	CMD -b 115200 (default 115200)
  -d string
    	CMD -d  something (default "data")
  -dpi float
    	screen resolution in Dots Per Inch (default 150)
  -fontfile string
    	filename of the ttf font (default "./arialuni.ttf")
  -hinting string
    	none | full (default "none")
  -m string
    	CMD -m r  // r (read) or w (write) (default "r")
  -p string
    	CMD -p /dev/ttymxc0 (default "/dev/ttymxc0")
  -size float
    	font size in points (default 12)
  -spacing float
    	line spacing (e.g. 2 means double spaced) (default 1.5)
  -t int
    	CMD -t 5  // second, <= 0 then stop sleep (default 2)
  -whiteonblack
    	white text on a black background


```



default serial port of KPW2 is /dev/ttymxc0


#### Base Logic


receive part:

- ./lora <br/>read data from LORA chipset, convert data to be gray graphic, save each point to metric.txt file


- ./luajit ./fb.lua <br/> read metric.txt by line, set the color of each pixel with framebuffer, then update screen

install this part to KPW2,

**NOTE:**, kindle uart was set by 
```shell 
getty -L 115200 /dev/ttymxc0
```
so the default UART baund is 115200

send part:

- ./lora -m w -b 9600 -p /dev/ttyUSB0 -d "from send part" <br/> send data via LORA chipset


