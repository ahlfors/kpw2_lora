package main

import (
	_ "bufio"
	"flag"
	"fmt"
	"github.com/golang/freetype"
	"github.com/jacobsa/go-serial/serial"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"os/exec"

)

var wait sync.WaitGroup

var mode string // read or write
var port string // serial port
var baud int64
var sleep int
var data string // write the data to port

var receivedStr []string
var stringToPNGChan chan int

var (
	dpi      = flag.Float64("dpi", 150, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "./arialuni.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 12, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)

func init() {
	flag.StringVar(&mode, "m", "r", "CMD -m r  // r (read) or w (write)")
	flag.StringVar(&port, "p", "/dev/ttymxc0", "CMD -p /dev/ttymxc0")
	flag.Int64Var(&baud, "b", 115200, "CMD -b 115200")
	flag.IntVar(&sleep, "t", 2, "CMD -t 5  // second, <= 0 then stop sleep")
	flag.StringVar(&data, "d", "data", "CMD -d  something")
	//flag.StringVar(&fontfile, "f", "./xx.ttf", "CMD -f /font/path/font.ttf")

	receivedStr = []string{"English Name", "中文名字", "Number123456()"}

	stringToPNGChan = make(chan int, 1)
}

func main() {

	// parse flag
	flag.Parse()

	for {
		wait.Add(1)

		if mode == "w" {
			go writePort(data)
		} else {
			go readPort()
		}

		wait.Wait()

		if sleep > 0 {
			// wait for close the serial port
			time.Sleep(time.Duration(sleep) * time.Second)
		}
	}

}

func stringToPNG(text []string, ch chan int) {

	fmt.Println("creating text image")
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context.
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	if *wonb {
		fg, bg = image.White, image.Black
		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 758, 100))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(*size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(*size * *spacing)
	}

	outputMetric(rgba)

	fmt.Println("created text image")
	// release channel
	<-ch

	/*
		// Save that RGBA image to disk.
		outFile, err := os.Create("out.png")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		defer outFile.Close()
		b := bufio.NewWriter(outFile)
		err = png.Encode(b, rgba)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		err = b.Flush()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		fmt.Println("Wrote out.png OK.")
	*/
}

func outputMetric(src *image.RGBA) {

	metricFile, err := os.OpenFile("./metric_orig.txt", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {

	}



	// Create a new grayscale image
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(bounds)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := src.At(x, y)

			grayColor := color.GrayModel.Convert(oldColor)

			gray.Set(x, y, grayColor)
			//r,g,b,a := grayColor.RGBA()
			//log.Println("x:",x, " y:", y, " color:", grayColor)

			grayVal := strings.Replace(fmt.Sprint(grayColor), "{", "", 1)
			grayVal = strings.Replace(fmt.Sprint(grayVal), "}", "", 1)

			//fmt.Println(fmt.Sprint(grayVal))

			line := fmt.Sprintf("%d,%d,%s\n", x, y, grayVal)
			metricFile.Write([]byte(line))
		}
	}

	metricFile.Close()


	cmd := exec.Command("cp", "./metric_orig.txt", "./metric.txt")
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println("failed to copay metric file")
		fmt.Println(err.Error())
	}
}

func writePort(data string) {

	defer wait.Done()

	options := serial.OpenOptions{
		PortName:        port,
		BaudRate:        uint(baud),
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 100,
	}

	port, err := serial.Open(options)
	if err != nil {
		fmt.Printf("serial.Open:%v\n", err)
		return
	}

	defer port.Close()

	fmt.Println(data)

	b := []byte(data + "")
	n, err := port.Write(b)
	if err != nil {
		fmt.Printf("port.Write: %v\n", err)
	}

	fmt.Println("Wrote", n, "bytes.")
	time.Sleep(2 * time.Second)
}

func readPort() {

	defer wait.Done()

	options := serial.OpenOptions{
		PortName:        port,
		BaudRate:        uint(baud),
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 100,
	}

	port, err := serial.Open(options)
	if err != nil {
		fmt.Printf("serial.Open:%v\n", err)
		return
	}

	defer port.Close()

	receive := make([]byte, 512)
	counter := 0

	for {

		buf := make([]byte, 10)

		n, err := port.Read(buf)
		if err != nil {
			fmt.Println("Error reading from serial port:", err)
			break

		}

		if len(buf) > 0 && string(buf) != "" {
			fmt.Println("RECEIVED PART: ", string(buf))
			receivedStr = []string{time.Now().String(),string(buf)}//append(receivedStr, string(buf))
			receive = append(receive[:], buf[:n]...)
			counter += n

			/*
			stringToPNGChan<-1
			stringToPNG(receivedStr, stringToPNGChan)
			*/

			select {
			case stringToPNGChan <- 1:
				go stringToPNG(receivedStr, stringToPNGChan)
			default:

			}

		}

	}

	fmt.Println("READ COUNT: ", counter)
	fmt.Println("READ DATA: ", string(receive))

	//receivedStr = []string{}
}
