package main

import (
	_ "embed"
	"image/color"
	"time"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"

	qrcode "github.com/skip2/go-qrcode"
)

const QR_INFO_SIZE = 180

func Info() {

	qr, err := qrcode.New("https://www.paroleostili.it", qrcode.Medium)
	if err != nil {
		println(err, 123)
	}

	qrbytes := qr.Bitmap()
	size := int16(len(qrbytes))

	factor := int16(QR_INFO_SIZE / len(qrbytes))

	bx := (QR_INFO_SIZE - size*factor) / 2
	by := (QR_INFO_SIZE - size*factor) / 2
	display.FillScreen(color.RGBA{255, 255, 0, 255})
	for y := int16(0); y < size; y++ {
		for x := int16(0); x < size; x++ {
			if qrbytes[y][x] {
				display.FillRectangle(bx+x*factor, by+y*factor, factor, factor, colors[0])
			} else {
				display.FillRectangle(bx+x*factor, by+y*factor, factor, factor, colors[1])
			}
		}
	}

	w32, _ := tinyfont.LineWidth(&freesans.Bold18pt7b, "SCAN")
	tinyfont.WriteLine(&display, &freesans.Bold18pt7b, QR_INFO_SIZE+((WIDTH-QR_INFO_SIZE)-int16(w32))/2, 45, "SCAN", colors[DARKBLUE])

	w32, _ = tinyfont.LineWidth(&freesans.Bold18pt7b, "ME")
	tinyfont.WriteLine(&display, &freesans.Bold18pt7b, QR_INFO_SIZE+((WIDTH-QR_INFO_SIZE)-int16(w32))/2, 80, "ME", colors[DARKBLUE])

	w32, _ = tinyfont.LineWidth(&freesans.Bold9pt7b, "Press any")
	tinyfont.WriteLine(&display, &freesans.Bold9pt7b, QR_INFO_SIZE+((WIDTH-QR_INFO_SIZE)-int16(w32))/2, 120, "Press any", colors[DARKBLUE])
	w32, _ = tinyfont.LineWidth(&freesans.Bold9pt7b, "button")
	tinyfont.WriteLine(&display, &freesans.Bold9pt7b, QR_INFO_SIZE+((WIDTH-QR_INFO_SIZE)-int16(w32))/2, 140, "button", colors[DARKBLUE])
	w32, _ = tinyfont.LineWidth(&freesans.Bold9pt7b, "to continue")
	tinyfont.WriteLine(&display, &freesans.Bold9pt7b, QR_INFO_SIZE+((WIDTH-QR_INFO_SIZE)-int16(w32))/2, 160, "to continue", colors[DARKBLUE])

	w32, _ = tinyfont.LineWidth(&freesans.Bold9pt7b, "Visit https://www.paroleostili.it")
	tinyfont.WriteLine(&display, &freesans.Bold9pt7b, (WIDTH-int16(w32))/2, QR_INFO_SIZE+25, "Visit https://www.paroleostili.it", colors[DARKBLUE])
	w32, _ = tinyfont.LineWidth(&freesans.Bold9pt7b, "for more information")
	tinyfont.WriteLine(&display, &freesans.Bold9pt7b, (WIDTH-int16(w32))/2, QR_INFO_SIZE+45, "for more information", colors[DARKBLUE])

	for {
		time.Sleep(100 * time.Millisecond)
		if !btnA.Get() || !btnB.Get() || !btnUp.Get() || !btnLeft.Get() || !btnRight.Get() || !btnDown.Get() {
			break
		}
	}

}
