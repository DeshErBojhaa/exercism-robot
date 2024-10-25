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
	robot *Step2Robot
	ok    bool
}

var robot Step2Robot

func (r *Step2Robot) copy() Step2Robot {
	return Step2Robot{
		Dir: r.Dir,
		Pos: r.Pos,
	}
}

func (r *Step2Robot) advance() {
	switch r.Dir {
	case N:
		r.Pos.Northing += 1
	case S:
		r.Pos.Northing -= 1
	case E:
		r.Pos.Easting += 1
	case W:
		r.Pos.Easting -= 1
	}
}

func (r *Step2Robot) Left() {
	r.Dir += 3
	r.Dir %= 4
}

func (r *Step2Robot) Right() {
	r.Dir += 1
	r.Dir %= 4
}

func StartRobot(command chan Command, action chan Action) {
	for c := range command {
		switch c {
		case 'A':
			action <- Action{&robot, false}
			if r := <-action; r.ok {
				r.robot.advance()
			}
		case 'L':
			(&robot).Left()
		case 'R':
			(&robot).Right()
		case ' ':
			robot.Dir = N
			robot.Pos = Pos{1, 1}
		}
	}
}

func Room(extent Rect, robot Step2Robot, action chan Action, report chan Step2Robot) {
	for act := range action {
		robot := act.robot.copy()
		robot.advance()

		if robot.Pos.Northing > extent.Max.Northing || robot.Pos.Northing < extent.Min.Northing {
			continue
		}

		if robot.Pos.Easting > extent.Max.Easting || robot.Pos.Easting < extent.Min.Easting {
			continue
		}

		act.robot.advance()
	}
	report <- robot
	close(action)
	close(report)
}

type Action3 struct{}

func StartRobot3(name, script string, action chan Action3, log chan string) {
	panic("Please implement the StartRobot3 function")
}

func Room3(extent Rect, robots []Step3Robot, action chan Action3, rep chan []Step3Robot, log chan string) {
	panic("Please implement the Room3 function")
}
