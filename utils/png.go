package utils

import (
	"fmt"
        _"image"
        "image/jpeg"
        "image/png"
        "io"
        "os"
)

// Converts image in reader to png format if previously in jpg format
func convertJPG2PNG(w io.Writer, r io.Reader) error{

        m, err := jpeg.Decode(r)
        if err != nil {
                return err
        }
        err = png.Encode(w, m)
        if err != nil {
                return err
        }

        return err
}

// Converts image in reader to jpg format if previously in png format
func convertPNG2JPG(w io.Writer, r io.Reader) error{

        m, err := png.Decode(r)
        if err != nil {
                return err
        }
        err = jpeg.Encode(w, m, nil)
        if err != nil {
                return err
        }

        return err
}

// Prints out histogram of image at provided path
func histogram(path string) {

        rdr, err := os.Open(path) //("../img/home.png")
        if err != nil {
                fmt.Println(err)
        }
        defer rdr.Close()
        im, err := png.Decode(rdr)
        if err != nil {
                fmt.Println(err)
        }
        bounds := im.Bounds()

        var histogram [16][4]int
        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
                for x := bounds.Min.X; x < bounds.Max.X; x++ {
                        r,g,b,a := im.At(x,y).RGBA()
                        histogram[r>>12][0]++
                        histogram[g>>12][1]++
                        histogram[b>>12][2]++
                        histogram[a>>12][3]++
                }
        }
        fmt.Printf("%-14s %6s %6s %6s %6s \n", "bin", "red", "green", "blue", "alpha")
        for i,x := range histogram {
                fmt.Printf("0x%04x-0x%04x: %6d %6d %6d %6d \n", i<<12, (i+1)<<12-1, x[0],x[1],x[2],x[3])
        }

}

