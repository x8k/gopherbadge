package main

import (
	"time"

	"tinygo.org/x/drivers/st7789"
	"tinygo.org/x/tinydraw"
)

type MovementDirection int

const (
	Stopped MovementDirection = 0
	Left    MovementDirection = -1
	Right   MovementDirection = 1
)

// MovementStrategy interface for flexible movement behaviors
type MovementStrategy interface {
	GetDirection() MovementDirection
}

// ButtonMovementStrategy implements MovementStrategy for button input
type ButtonMovementStrategy struct{}

func (b *ButtonMovementStrategy) GetDirection() MovementDirection {
	if !btnLeft.Get() {
		return Left
	}
	if !btnRight.Get() {
		return Right
	}
	return Stopped
}

// AccelerometerMovementStrategy implements MovementStrategy for accelerometer input
type AccelerometerMovementStrategy struct{}

func (a *AccelerometerMovementStrategy) GetDirection() MovementDirection {
	x, _, _ := ReadAcceleration()
	const threshold = 10
	if x < -threshold {
		return Right
	} else if x > threshold {
		return Left
	}
	return Stopped
}

// CompositeMovementStrategy combines multiple strategies
type CompositeMovementStrategy struct {
	strategies []MovementStrategy
}

func (c *CompositeMovementStrategy) GetDirection() MovementDirection {
	for _, s := range c.strategies {
		dir := s.GetDirection()
		if dir != Stopped {
			return dir
		}
	}
	return Stopped
}

// Stickperson struct encapsulates state and rendering
type Stickperson struct {
	PositionX int
	Direction MovementDirection
	Display   *st7789.Device // Dependency injection for display
	Waving    bool
}

func NewStickperson(centerX int, display *st7789.Device) *Stickperson {
	return &Stickperson{
		PositionX: centerX,
		Direction: Stopped,
		Display:   display,
	}
}

func (s *Stickperson) Move(dir MovementDirection, step int) {
	s.Direction = dir
	switch dir {
	case Left:
		s.PositionX -= step
	case Right:
		s.PositionX += step
	case Stopped:
		// No movement
	}
	// Clamp position so arms never fully leave screen
	minX := armLength
	maxX := WIDTH - 1 - armLength
	if s.PositionX < minX {
		s.PositionX = minX
	}
	if s.PositionX > maxX {
		s.PositionX = maxX
	}
}

func (s *Stickperson) Draw() {
	display := s.Display
	centerX := s.PositionX
	bodyEndY := HEIGHT - legLength
	bodyStartY := bodyEndY - stickmanHeight/2
	headCenterY := bodyStartY - headRadius
	armY := bodyStartY + headRadius

	// Draw head
	tinydraw.Circle(display, int16(centerX), int16(headCenterY), int16(headRadius), colors[BLACK])
	// Draw body
	tinydraw.Line(display, int16(centerX), int16(bodyStartY), int16(centerX), int16(bodyEndY), colors[BLACK])

	// Draw arms
	if s.Waving {
		// Left arm waving (raised)
		tinydraw.Line(display, int16(centerX-armLength), int16(armY-headRadius*2), int16(centerX), int16(armY), colors[BLACK])
	} else {
		// Left arm normal
		tinydraw.Line(display, int16(centerX-armLength), int16(armY), int16(centerX), int16(armY), colors[BLACK])
	}
	// Right arm always normal
	tinydraw.Line(display, int16(centerX), int16(armY), int16(centerX+armLength), int16(armY), colors[BLACK])

	// Gambe animate
	var leftLegX, leftLegY, rightLegX, rightLegY int
	step := armLength / 2
	switch s.Direction {
	case Right:
		leftLegX = centerX - armLength
		leftLegY = bodyEndY + legLength
		rightLegX = centerX + armLength + step
		rightLegY = bodyEndY + legLength - step
	case Left:
		leftLegX = centerX - armLength - step
		leftLegY = bodyEndY + legLength - step
		rightLegX = centerX + armLength
		rightLegY = bodyEndY + legLength
	default:
		leftLegX = centerX - armLength
		leftLegY = bodyEndY + legLength
		rightLegX = centerX + armLength
		rightLegY = bodyEndY + legLength
	}
	tinydraw.Line(display, int16(centerX), int16(bodyEndY), int16(leftLegX), int16(leftLegY), colors[BLACK])
	tinydraw.Line(display, int16(centerX), int16(bodyEndY), int16(rightLegX), int16(rightLegY), colors[BLACK])
}

const (
	headRadius     = HEIGHT / 12
	stickmanHeight = HEIGHT / 2
	armLength      = WIDTH / 8
	legLength      = stickmanHeight / 3
	defaultOffset  = 8
)

func stick(day int, hour int) {
	// Use Stickperson abstraction
	sp := NewStickperson(WIDTH/2, &display)
	movement := &CompositeMovementStrategy{
		strategies: []MovementStrategy{
			&ButtonMovementStrategy{},
			&AccelerometerMovementStrategy{},
		},
	}
	display.FillScreen(colors[GREEN])
	sp.Draw()
	display.Display()

	wavingState := false
	for {
		if !btnA.Get() || !btnUp.Get() || !btnDown.Get() {
			break
		}
		dir := movement.GetDirection()
		if dir != Stopped {
			sp.Move(dir, defaultOffset)
		} else {
			sp.Direction = Stopped
		}
		// Waving logic
		if !btnB.Get() {
			wavingState = !wavingState // alternate up/down each frame while pressed
			sp.Waving = true
		} else {
			sp.Waving = false
		}
		display.FillScreen(colors[GREEN])
		sp.Draw()
		display.Display()
		time.Sleep(150 * time.Millisecond)
	}
}

// stick function uses Stickperson abstraction
