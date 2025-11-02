package main

import (
	_ "embed"
	"image/color"
	"math/rand"
	"time"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
)

func Fosdem() {
	var (
		colorArray = [9]color.RGBA{
			{R: 255, G: 0, B: 0, A: 255},     // RED
			{R: 0, G: 255, B: 0, A: 255},     // GREEN
			{R: 0, G: 0, B: 255, A: 255},     // BLUE
			{R: 255, G: 255, B: 0, A: 255},   // YELLOW
			{R: 0, G: 255, B: 255, A: 255},   // CYAN
			{R: 255, G: 0, B: 255, A: 255},   // MAGENTA
			{R: 255, G: 165, B: 0, A: 255},   // ORANGE
			{R: 255, G: 255, B: 255, A: 255}, // WHITE
			{R: 128, G: 0, B: 128, A: 255},   // PURPLE
		}
	)
	backgroundColor := colorArray[len(colorArray)-1]
	display.FillScreen(backgroundColor)
	var beatWord string = "FOSDEM"
	var oldBeatWord string = ""
	w32, _ := tinyfont.LineWidth(&freesans.Bold18pt7b, beatWord)
	tinyfont.WriteLine(&display, &freesans.Bold18pt7b, (WIDTH-int16(w32))/2, HEIGHT/2-20, beatWord, colors[WHITE])

	ledColors := make([]color.RGBA, 2)
	ledColors[0] = backgroundColor
	ledColors[1] = backgroundColor
	leds.WriteColors(ledColors)

	// Draw a 4px circle at the pos of the screen
	var (
		posX          = int16(WIDTH / 2) // Current position in pixels
		posY          = int16(HEIGHT / 2)
		deltaTime     = int16(1) // Higher value means slower movement
		radius        = int16(8)
		colorIndex    = int(0)
		oldColorIndex = int(0)
		tColor        = int(0)
		stepTime      = int(1000) // milliseconds
	)

	for {

		if colorIndex > len(colorArray)-2 { // don't use the last color (it's the background)
			colorIndex = 0
		}

		if tColor >= stepTime {
			oldColorIndex = colorIndex
			colorIndex++
			tColor = 0
		}
		tColor += 50

		if !btnA.Get() || !btnB.Get() || !btnUp.Get() || !btnLeft.Get() || !btnRight.Get() || !btnDown.Get() {
			break
		}

		// Leds!
		ledColors[0] = colorArray[colorIndex]
		ledColors[1] = colorArray[oldColorIndex]
		leds.WriteColors(ledColors)

		// Update beat word Should be a function
		oldBeatWord = beatWord
		beatWord = BeatWord()
		updateBeatWord(w32, oldBeatWord, backgroundColor, beatWord)

		accX, accY, accZ := ReadAcceleration()
		oldPosX := posX
		oldPosY := posY
		oldRadius := radius
		// Update position based on velocity
		posX -= accX / deltaTime
		posY -= accY / deltaTime
		posX = PosX(posX, radius)
		posY = PosY(posY, radius)
		radius = PosZ(radius, accZ)

		updateCircle(oldPosX, oldPosY, oldRadius, backgroundColor, posX, posY, radius, colorArray, colorIndex)

		time.Sleep(50 * time.Millisecond)

	}

	// Turn off leds shoudl be a function
	shutdownLeds(ledColors)
}

func updateCircle(oldPosX int16,
	oldPosY int16,
	oldRadius int16,
	backgroundColor color.RGBA,
	posX int16,
	posY int16,
	radius int16,
	colorArray [9]color.RGBA,
	colorIndex int) {
	tinydraw.FilledCircle(&display, oldPosX, oldPosY, oldRadius, backgroundColor)
	tinydraw.FilledCircle(&display, posX, posY, radius, colorArray[colorIndex])
}

func updateBeatWord(w32 uint32, oldBeatWord string, backgroundColor color.RGBA, beatWord string) {
	w32, _ = tinyfont.LineWidth(&freesans.Bold18pt7b, oldBeatWord)
	tinyfont.WriteLine(&display, &freesans.Bold18pt7b, (WIDTH-int16(w32))/2, HEIGHT/2-20, oldBeatWord, backgroundColor)
	w32, _ = tinyfont.LineWidth(&freesans.Bold18pt7b, beatWord)
	tinyfont.WriteLine(&display, &freesans.Bold18pt7b, (WIDTH-int16(w32))/2, HEIGHT/2-20, beatWord, colors[WHITE])
}

func shutdownLeds(ledColors []color.RGBA) {
	ledColors[0] = color.RGBA{0, 0, 0, 255}
	ledColors[1] = color.RGBA{0, 0, 0, 255}
	leds.WriteColors(ledColors)
	time.Sleep(50 * time.Millisecond)
	ledColors[0] = color.RGBA{0, 0, 0, 255}
	ledColors[1] = color.RGBA{0, 0, 0, 255}
	leds.WriteColors(ledColors)
	time.Sleep(50 * time.Millisecond)
}

func PosX(x int16, r int16) int16 {
	if x <= 0 {
		return 0 + r
	}
	if x >= WIDTH {
		return WIDTH - r
	}
	return x
}

func PosY(y int16, r int16) int16 {
	if y <= 0 {
		return 0 + r
	}
	if y >= HEIGHT {
		return HEIGHT - r
	}
	return y
}

// This function doen't really work
func PosZ(r int16, a int16) int16 {
	min := int16(1)
	max := int16(12)
	step := float32((max - min) / 128)
	r += int16(float32(a*10) * step)
	if r <= min {
		return min
	}
	if r >= max {
		return max
	}
	return r
}

func ReadAcceleration() (int16, int16, int16) {
	x, y, z := accel.ReadRawAcceleration()
	x = x / 250
	y = y / 250
	z = z / 250
	if x > 128 {
		x = 128
	}
	if y > 128 {
		y = 128
	}
	if z > 128 {
		z = 128
	}
	if x < -128 {
		x = -128
	}
	if y < -128 {
		y = -128
	}
	if z < -128 {
		z = -128
	}
	return x, y, z
}

func BeatWord() string {
	var retroTechWords = [52]string{
		"FOSDEM",
		"Exterminate",
		"Synthesizer",
		"Cyberpunk",
		"Dystopia",
		"BPM",
		"Mainframe",
		"Neon",
		"Analog",
		"LaserDisc",
		"Trance",
		"Hologram",
		"Technotron",
		"Cyberspace",
		"Beat",
		"Retrowave",
		"Matrix",
		"Circuit",
		"Warp",
		"Pulse",
		"Cyborg",
		"Groove",
		"Android",
		"Bass",
		"Quantum",
		"Rhythm",
		"Neuromancer",
		"Beats",
		"Synthwave",
		"Electro",
		"Datastream",
		"Vocoder",
		"Hyperspace",
		"Acid",
		"Replicant",
		"Vector",
		"Plasma",
		"Cyberdeck",
		"Techno",
		"Dystopian",
		"Sync",
		"Chrome",
		"Console",
		"Flux",
		"Arcade",
		"Binary",
		"Darkwave",
		"Amplitude",
		"Raygun",
		"Modem",
		"Vaporwave",
		"Zion",
	}
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	randomIndex := r.Intn(len(retroTechWords))
	return retroTechWords[randomIndex]
}
