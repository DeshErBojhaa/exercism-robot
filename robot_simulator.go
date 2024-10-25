package robot

import (
	"fmt"
)

// See defs.go for other definitions

// Step 1
var (
	N Dir = 0
	E Dir = 1
	S Dir = 2
	W Dir = 3
)

func Right() {
	Step1Robot.Dir = (Step1Robot.Dir + 1) % 4
}

func Left() {
	Step1Robot.Dir = (Step1Robot.Dir + 3) % 4
}

func Advance() {
	switch Step1Robot.Dir {
	case N:
		Step1Robot.Y += 1
	case S:
		Step1Robot.Y -= 1
	case E:
		Step1Robot.X += 1
	case W:
		Step1Robot.X -= 1
	}
}

func (d Dir) String() string {
	return fmt.Sprintf("%d", d)
}

type Action struct {
	command Command
}

var robot Step2Robot

func maxRU(a, b RU) RU {
	if a > b {
		return a
	}
	return b
}

func minRU(a, b RU) RU {
	if a < b {
		return a
	}
	return b
}

func (r *Step2Robot) advance(rect Rect) {
	switch r.Dir {
	case N:
		r.Pos.Northing = minRU(r.Pos.Northing+1, rect.Max.Northing)
	case S:
		r.Pos.Northing = maxRU(r.Pos.Northing-1, rect.Min.Northing)
	case E:
		r.Pos.Easting = minRU(r.Pos.Easting+1, rect.Max.Easting)
	case W:
		r.Pos.Easting = maxRU(r.Pos.Easting-1, rect.Min.Easting)
	}
}

func (r *Step2Robot) Left() {
	r.Dir = (r.Dir + 3) % 4
}

func (r *Step2Robot) Right() {
	r.Dir = (r.Dir + 1) % 4
}

func StartRobot(command chan Command, action chan Action) {
	robot.Dir = N
	robot.Pos = Pos{1, 1}

	for c := range command {
		action <- Action{c}
	}
	close(action)
}

func Room(extent Rect, _ Step2Robot, action chan Action, report chan Step2Robot) {
	for act := range action {
		switch act.command {
		case 'A':
			robot.advance(extent)
		case 'L':
			robot.Left()
		case 'R':
			robot.Right()
		case ' ':
			robot.Dir = N
			robot.Pos = Pos{1, 1}
		}
	}
	report <- robot
	close(report)
}

type Action3 struct{}

func StartRobot3(name, script string, action chan Action3, log chan string) {
	panic("Please implement the StartRobot3 function")
}

func Room3(extent Rect, robots []Step3Robot, action chan Action3, rep chan []Step3Robot, log chan string) {
	panic("Please implement the Room3 function")
}
