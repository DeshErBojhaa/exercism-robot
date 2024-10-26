package robot

import (
	"fmt"
	"sync"
	"time"
)

// See defs.go for other definitions

// Step 1
const (
	N Dir = iota
	E
	S
	W
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
	switch d {
	case N:
		return "N"
	case S:
		return "S"
	case E:
		return "E"
	case W:
		return "W"
	default:
		return fmt.Sprintf("Dir(%d)", d)
	}
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

type Action3 struct {
	name string
	cmd  Command
}

var once sync.Once

func StartRobot3(name, script string, action chan Action3, log chan string) {
	for _, ch := range script {
		action <- Action3{
			name: name,
			cmd:  Command(ch),
		}
	}
	time.Sleep(time.Millisecond * 100)
	once.Do(func() {
		close(action)
	})
}

func Room3(extent Rect, robots []Step3Robot, action chan Action3, rep chan []Step3Robot, log chan string) {
	once = sync.Once{}
	rMap := make(map[string]*Step2Robot)
	positions := make(map[Pos]struct{}, len(robots))
	robotPos := make(map[string]Pos, len(robots))
	for i, r := range robots {
		if r.Name == "" {
			log <- "no name"
		}
		if _, ok := rMap[r.Name]; ok {
			log <- "duplicate name"
		}
		if r.Pos.Easting < extent.Min.Easting || r.Pos.Easting > extent.Max.Easting || r.Pos.Northing < extent.Min.Northing || r.Pos.Northing > extent.Max.Northing {
			log <- "out side of room"
		}
		if _, ok := positions[r.Pos]; ok {
			log <- "duplicate position"
		}
		positions[r.Pos] = struct{}{}
		rMap[r.Name] = &robots[i].Step2Robot
		robotPos[r.Name] = robots[i].Step2Robot.Pos
	}

	for act := range action {
		r, ok := rMap[act.name]
		if !ok {
			log <- "invalid robot"
			continue
		}
		switch act.cmd {
		case 'A':
			oldPos := robotPos[act.name]
			r.advance(extent)
			if r.Pos.Northing == oldPos.Northing && r.Pos.Easting == oldPos.Easting {
				log <- act.name + " bump into wall"
			}
			for k, v := range robotPos {
				if k == act.name {
					continue
				}
				if r.Pos == v {
					log <- act.name + " bump into " + k
					r.Pos = oldPos
					break
				}
			}
			robotPos[act.name] = r.Pos
		case 'L':
			r.Left()
		case 'R':
			r.Right()
		case 'T':
			continue
		default:
			log <- "invalid command"
		}
	}

	ret := make([]Step3Robot, 0, len(robots))
	for _, r := range robots {
		ret = append(ret, Step3Robot{
			Name:       r.Name,
			Step2Robot: *rMap[r.Name],
		})
	}

	rep <- ret
	close(rep)
}
