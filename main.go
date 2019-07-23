package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	gwidth  = 550
	gheight = 400

	pwidth  = 20
	pheight = 50

	lwidth  = 43.25
	lheight = 21

	camFollowSpeed = .5
)

var bcolors = []color.Color{
	hex2col(0x9E6E59),
	color.Transparent,
	color.Transparent,
	color.Transparent,
	color.Transparent,
	hex2col(0xE9EB13),
	hex2col(0xA87057),
	hex2col(0xFFFFFF),
	hex2col(0x61614F),
	hex2col(0xCDCDCD),
	hex2col(0xCDCDCD),
	hex2col(0xCDCDCD),
	hex2col(0xCDCDCD),
	hex2col(0xCDCDCD),
}

func hex2col(h uint) pixel.RGBA {
	r := h & 0xFF0000 >> 16
	g := h & 0x00FF00 >> 8
	b := h & 0x0000FF
	return pixel.RGB((float64(r) / 0xFF), (float64(g) / 0xFF), (float64(b) / 0xFF))
}

const deg2rad = math.Pi / 180

func run() {
	rand.Seed(time.Now().Unix())
	cfg := pixelgl.WindowConfig{
		Title:  "Platform Racing | fbf=false",
		Bounds: pixel.R(0, 0, gwidth, gheight),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	txt := text.New(pixel.V(1, 390), text.NewAtlas(basicfont.Face7x13, text.ASCII))
	imd := imdraw.New(nil)
	camPos := pixel.ZV
	fbf := false

	st := &state{
		course: &course{
			// name: "newbieland",
			name: "robocity",
		},
	}

	fps := time.Tick(time.Second / 30)
	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyA) {
			fbf = !fbf
			win.SetTitle(fmt.Sprint("Platform Racing | fbf=", fbf))
		}

		var me *player
		if !fbf || (fbf && (win.JustPressed(pixelgl.KeyS) || win.Repeated(pixelgl.KeyS))) {
			st.inputs = inputs{
				up:    win.Pressed(pixelgl.KeyUp),
				down:  win.Pressed(pixelgl.KeyDown),
				left:  win.Pressed(pixelgl.KeyLeft),
				right: win.Pressed(pixelgl.KeyRight),
				space: win.Pressed(pixelgl.KeySpace),
			}
			st.nextFrame()
			if st.curFrame == 1 {
				st.course.guys = []*player{st.course.me}
			}

			me = st.course.me
			var _loc8 = me.x - gwidth/2
			var _loc7 = me.y - gheight/2 + 25
			var _loc10 = (_loc8) - camPos.X
			var _loc9 = (_loc7) - camPos.Y
			var _loc6 = _loc10 * camFollowSpeed
			var _loc5 = _loc9 * camFollowSpeed
			camPos.X += _loc6
			camPos.Y += _loc5
			cam := pixel.IM.Moved(camPos.Scaled(-1))
			imd.SetMatrix(cam)

			imd.Clear()

			// draw map
			brect := pixel.R(0, 0, segSize, segSize)
			for _, row := range st.course.blocks {
				for _, b := range row {
					if b == nil {
						continue
					}
					br := brect.Moved(pixel.V(b.x, b.y))
					imd.Color = bcolors[int(b.t)]
					imd.Push(br.Min, br.Max)
					imd.Rectangle(0)
					isArrow := b.t == blockUp || b.t == blockLeft || b.t == blockRight
					if isArrow {
						imd.Color = colornames.Black
						line1 := pixel.L(pixel.V(-5, 0), pixel.V(0, +10))  // /
						line2 := pixel.L(pixel.V(0, -10), pixel.V(0, +10)) //  |
						line3 := pixel.L(pixel.V(+5, 0), pixel.V(0, +10))  //   \
						mat := pixel.IM.Moved(br.Center())
						switch b.t {
						case blockLeft:
							mat = mat.Rotated(br.Center(), deg2rad*90)
						case blockRight:
							mat = mat.Rotated(br.Center(), deg2rad*-90)
						}
						imd.Push(mat.Project(line1.A), mat.Project(line1.B))
						imd.Line(1)
						imd.Push(mat.Project(line2.A), mat.Project(line2.B))
						imd.Line(1)
						imd.Push(mat.Project(line3.A), mat.Project(line3.B))
						imd.Line(1)
					}
				}
			}

			// draw player
			prect := pixel.R(0, 0, pwidth, pheight)
			for _, guy := range st.course.guys {
				pr := prect.Moved(pixel.V(guy.x-(prect.W()/2), guy.y))
				imd.Color = colornames.Red
				imd.Push(pr.Min, pr.Max)
				imd.Rectangle(1)
				imd.Color = colornames.Aqua
				imd.Push(pixel.V(me.x, me.y+5), pixel.V(me.x, me.y-5))
				imd.Line(1)
				imd.Push(pixel.V(me.x+5, me.y), pixel.V(me.x-5, me.y))
				imd.Line(1)
			}

			// draw lasers
			lrect := pixel.R(0, 0, lwidth, lheight)
			for _, laser := range st.course.lasers {
				lr := lrect.Moved(pixel.V(laser.x-lrect.W(), laser.y-lrect.H()/2))
				imd.Color = colornames.Yellow
				imd.Push(lr.Min, lr.Max)
				imd.Rectangle(1)
				imd.Color = colornames.Aqua
				imd.Push(pixel.V(laser.x, laser.y+5), pixel.V(laser.x, laser.y-5))
				imd.Line(1)
				imd.Push(pixel.V(laser.x+5, laser.y), pixel.V(laser.x-5, laser.y))
				imd.Line(1)
			}

			win.Clear(colornames.Black)
			imd.Draw(win)

			txt.Clear()
			fmt.Fprintf(txt,
				"tme=%d\n"+
					"pos=%f,%f\n"+
					"vel=%f,%f\n"+
					"vtg=%f\n"+
					"rec=%d\n"+
					"sjv=%d\n"+
					"wpn=%s %d %d\n"+
					"mod=%s\n"+
					"lsr=%v\n", int(st.timeLeft()), me.x, me.y, me.xVel, me.yVel, me.xVelTarget, me.recoveryTimer, me.superJump, me.weapon, me.bullets, me.jetFuel, me.mode, st.course.lasers)
			txt.Draw(win, pixel.IM)
		}

		win.Update()
		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
