package main

import (
	"time"

	"tinygo.org/x/tinydraw"
)

const (
	headRadius     = HEIGHT / 12
	stickmanHeight = HEIGHT / 2
	armLength      = WIDTH / 8
	legLength      = stickmanHeight / 3
	defaultOffset  = 8
)

func stick(day int, hour int) {
	// Persistent position variable, initialized to center
	positionX := WIDTH / 2
	display.FillScreen(colors[GREEN])
	stickperson(positionX)
	display.Display()

	for {
		if !btnA.Get() || !btnB.Get() || !btnUp.Get() || !btnDown.Get() {
			break
		}
		moved := false
		if !btnLeft.Get() {
			positionX -= defaultOffset
			moved = true
		}
		if !btnRight.Get() {
			positionX += defaultOffset
			moved = true
		}
		// Accelerometro
		x, _, _ := ReadAcceleration()
		const threshold = 10
		if x < -threshold {
			positionX += defaultOffset
			moved = true
		} else if x > threshold {
			positionX -= defaultOffset
			moved = true
		}
		if moved {
			display.FillScreen(colors[GREEN])
			stickperson(positionX)
			display.Display()
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// positionX: posizione orizzontale assoluta del corpo
func stickperson(positionX int) {
	// Stickman parameters
	centerX := positionX
	// Calcola la posizione Y in modo che i piedi tocchino il bordo inferiore
	bodyEndY := HEIGHT - legLength
	bodyStartY := bodyEndY - stickmanHeight/2
	headCenterY := bodyStartY - headRadius
	armY := bodyStartY + headRadius

	// Draw head
	tinydraw.Circle(&display, int16(centerX), int16(headCenterY), int16(headRadius), colors[BLACK])
	// Draw body
	tinydraw.Line(&display, int16(centerX), int16(bodyStartY), int16(centerX), int16(bodyEndY), colors[BLACK])
	// Draw arms
	tinydraw.Line(&display, int16(centerX-armLength), int16(armY), int16(centerX+armLength), int16(armY), colors[BLACK])
	// Simula camminata: gambe alternate in base al movimento (usiamo una variabile statica per ricordare l'ultimo movimento)
	var leftLegX, leftLegY, rightLegX, rightLegY int
	// step variable removed (no longer used)
	// Per semplicità, le gambe sono sempre in posizione neutra (come da ultimo requisito)
	leftLegX = centerX - armLength
	leftLegY = bodyEndY + legLength
	rightLegX = centerX + armLength
	rightLegY = bodyEndY + legLength
	// Draw left leg
	tinydraw.Line(&display, int16(centerX), int16(bodyEndY), int16(leftLegX), int16(leftLegY), colors[BLACK])
	// Draw right leg
	tinydraw.Line(&display, int16(centerX), int16(bodyEndY), int16(rightLegX), int16(rightLegY), colors[BLACK])
}
