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


### KPW2 Disk Partition

```shell
[root@kindle root]# fdisk -l

Disk /dev/mmcblk0: 1958 MB, 1958739968 bytes
4 heads, 16 sectors/track, 59776 cylinders
Units = cylinders of 64 * 512 = 32768 bytes

        Device Boot      Start         End      Blocks  Id System
/dev/mmcblk0p1   *        1025       12224      358400  83 Linux
/dev/mmcblk0p2           12225       14272       65536  83 Linux
/dev/mmcblk0p3           14273       16320       65536  83 Linux
/dev/mmcblk0p4           16321       59776     1390592   b Win95 FAT32

Disk /dev/mmcblk0boot0: 1 MB, 1048576 bytes
4 heads, 16 sectors/track, 32 cylinders
Units = cylinders of 64 * 512 = 32768 bytes

Disk /dev/mmcblk0boot0 doesn't contain a valid partition table
[root@kindle root]# df -h
Filesystem                Size      Used Available Use% Mounted on
/dev/root               340.2M    292.8M     30.3M  91% /
tmpfs                   124.8M    104.0K    124.7M   0% /dev
tmpfs                   124.8M         0    124.8M   0% /dev/shm
tmpfs                    32.0M    368.0K     31.6M   1% /var
/dev/loop/2               2.5M      2.5M         0 100% /usr/share/X11/xkb
/dev/loop/3             100.1M    100.1M         0 100% /usr/java/lib/fonts
/dev/loop/4               1.3M      1.3M         0 100% /etc/kdb.src
/dev/loop/5               8.0M      8.0M         0 100% /usr/lib/locale
/dev/loop/6              16.0M     16.0M         0 100% /usr/share/keyboard
/dev/mmcblk0p3           62.0M     21.1M     37.6M  36% /var/local
/dev/loop/0               1.3G    623.5M    728.5M  46% /mnt/base-us
fsp                       1.3G    623.5M    728.5M  46% /mnt/us
/dev/loop/7              23.6M     23.6M         0 100% /var/local/font/mnt/ja_font
/dev/loop/8              46.1M     46.1M         0 100% /var/local/font/mnt/zh-Hans_font
[root@kindle root]# 

```

create /lib/ld-linux-armhf.so.3 file for some file, renamed from /lib/ld-linux.so.3 

```shell
[root@kindle root]# ls /lib/ -al
drwxrwxr-x    5 root     root          2048 Jan 17 10:30 .
drwxr-xr-x   13 root     root          1024 Dec  4 16:08 ..
lrwxrwxrwx    1 root     root             8 Apr 23  2014 cpp -> /bin/cpp
-rwxr-xr-x    1 root     root          7656 Apr 23  2014 e2initrd_helper
drwxrwxr-x    4 root     root          1024 Apr 23  2014 firmware
-rwxr-xr-x    1 root     root         76308 Apr 23  2014 klibc-KukkXgqjfSmsQ1r5iaQqrMdfe3M.so
-rwxr-xr-x    1 root     root        101792 Sep 11  2010 ld-2.12.1.so
-rwxr-xr-x    1 root     root        101792 Jan 17 10:30 ld-linux-armhf.so.3
lrwxrwxrwx    1 root     root            12 Apr 23  2014 ld-linux.so.3 -> ld-2.12.1.so
``

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

```shell
[root@kindle sbin]# cat /usr/sbin/usbserial 
#!/bin/sh
# usbserial needs to be called backgrounded in order to avoid trouble !!!!!

_GETTY_CMD="/sbin/getty -L 115200 ttygserial -l /bin/login"

# if -r then rmmod and go back to regular stuff
if [ "x$1" == "x-r" ]; then
	#First kill the other usbserial
	_OTHER=`ps ax | grep "$0" | grep -v grep | awk '{print $1}' | xargs echo`
	for pid in $_OTHER; do
		if [ "$$" != "$pid" ]; then
			kill $pid
		fi
	done
	sleep 1
	kill `ps ax | grep "$_GETTY_CMD" | grep -v grep | awk '{print $1}'` | xargs echo 
	lipc-set-prop com.lab126.volumd useUsbForSerial 0
else
	if [ "`lipc-get-prop com.lab126.volumd useUsbForSerial`" == "1" ]; then
		echo "Already in SerialOverUSB mode"
	else
		lipc-set-prop com.lab126.volumd useUsbForSerial 1
		while [ 1 ]; do
			$_GETTY_CMD
		done;
	fi
fi
```

send part:

- ./lora -m w -b 9600 -p /dev/ttyUSB0 -d "from send part" <br/> send data via LORA chipset


