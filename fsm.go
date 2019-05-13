package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	segSize float64 = 30
	gravity         = .49 // estimation
)

type inputs struct{ up, down, left, right, space bool }

type state struct {
	sneamia               bool //= false
	curFrame              int
	updateCounter         int
	updateReset           int //= 0 // ????
	startMS, timeoutMS    float64
	forceMyPositionUpdate bool
	course                *course
	inputs                inputs
}

func (s *state) timeLeft() float64 {
	return math.Round((s.timeoutMS - (getMS() - s.startMS)) / 1000)
}

func (s *state) nextFrame() {
	s.curFrame++
	if s.curFrame == 1 {
		s.loadLevel()
		s.course.me = s.course.guys[0]
		s.course.lasers = []*laser{}
		s.course.placedMines = []*mine{}
		//drawMinimap(mapp);
		//centerCam(me.m)
		//camFollow(me.m)
		s.startMS = getMS()
		s.timeoutMS = 120000
	} else if s.curFrame == 61 {
		for _, guy := range s.course.guys {
			guy.waiting = false
		}
	}
	// some local refs for convenience
	me := s.course.me
	guys := s.course.guys
	course := s.course

	var timeLeft = s.timeLeft()
	//menu_mc.timeLeftBox.text = formatSeconds(timeLeft,"seconds");
	if timeLeft <= 0 && me.sendToSocket != nil {

		me.sendToSocket("race_over`")
	}
	if !me.finished && !me.forfeit {
		me.controlChange("r", s.inputs.right)
		me.controlChange("l", s.inputs.left)
		me.controlChange("u", s.inputs.up)
		me.controlChange("d", s.inputs.down)
		me.controlChange("space", s.inputs.space)

		s.squashTest()
		// if sneamia {
		// 	if win.JustPressed(pixelgl.Key1) {
		// 		sendToSocket("force_item`1")
		// 	} else if win.JustPressed(pixelgl.Key2) {
		// 		sendToSocket("force_item`2")
		// 	} else if win.JustPressed(pixelgl.Key3) {
		// 		sendToSocket("force_item`3")
		// 	} else if win.JustPressed(pixelgl.Key4) {
		// 		sendToSocket("force_item`4")
		// 	} else if win.JustPressed(pixelgl.Key5) {
		// 		sendToSocket("force_item`5")
		// 	} else if win.JustPressed(pixelgl.Key6) {
		// 		sendToSocket("force_item`6")
		// 	} else if win.JustPressed(pixelgl.Key7) {
		// 		sendToSocket("force_item`7")
		// 	}
		// }
	}
	for _, guy := range guys {
		if !guy.finished && !guy.forfeit {
			guy.move()
		}
	}
	s.updateCounter--
	if (s.updateCounter <= 0 || s.forceMyPositionUpdate) && !me.finished && !me.forfeit {
		s.forceMyPositionUpdate = false
		s.updateCounter = s.updateReset
		me.sendPositionUpdate()
	}
	//moveDust() <- this isn't defined anywhere???

	// mines
	tmp1 := course.placedMines[:0]
	for _, m := range course.placedMines {
		if m.placeTimer == 1 {
			course.blocks[int(m.x/30)][int(m.y/30)] = &block{x: m.x, y: m.y, t: blockMine}
		} else {
			m.placeTimer--
			tmp1 = append(tmp1, m)
		}

	}
	for i := len(tmp1); i < len(course.placedMines); i++ {
		course.placedMines[i] = nil
	}
	course.placedMines = tmp1

	// lasers
	tmp2 := course.lasers[:0]
	for _, l := range course.lasers {
		l.x += l.vel
		bx := int(math.Round(l.x / segSize))
		by := int(math.Round(l.y / segSize))
		if bx > -1 && bx < len(course.blocks) && by > -1 && by < len(course.blocks[0]) {
			var _loc5 = course.blocks[bx][by]
			if _loc5 != nil {
				if _loc5.t == blockBrick || _loc5.t == blockMine {
					course.blocks[bx][by] = nil
					continue
				}
				if l.vel > 0 {
					l.x = _loc5.x
				} else {
					l.x = _loc5.x + segSize
				}
				continue //this.onDie()
			}
		}
		for _, guy := range guys {
			if guy != l.shooter && guy.recoveryTimer <= 0 && l.hitTest(guy) {
				if l.vel > 0 {
					guy.xVel += 5
				} else {
					guy.xVel -= 5
				}
				if guy == me && !s.sneamia {
					me.controlChange("bumped", true)
				}
				l.x = guy.x
				continue //this.onDie()
			}
			// note: lasers which don't hit a player or block never expire
		}
		tmp2 = append(tmp2, l)
	}
	for i := len(tmp2); i < len(course.lasers); i++ {
		course.lasers[i] = nil
	}
	course.lasers = tmp2
}

func (s *state) loadLevel() {
	raw := levels[s.course.name]
	for i := len(raw)/2 - 1; i >= 0; i-- {
		opp := len(raw) - 1 - i
		raw[i], raw[opp] = raw[opp], raw[i]
	}
	width := len(raw[0])
	height := len(raw)
	s.course.deathHeight = -float64(height*30 + 300)
	s.course.guys = make([]*player, 4)

	parseT := func(t byte) blockType {
		d, err := strconv.ParseInt(string(t), 0, 64)
		if err != nil {
			panic(err)
		}
		return blockType(d)
	}

	newPlayer := func(x, y float64) *player {
		return &player{
			x: x, y: y,
			speed: 50, jump: 50, traction: 50,
			waiting:       true,
			speedMod:      1,
			runTimerReset: 7,
			recoveryTimer: 75,
			xScale:        +1,
			xCorrection:   x,
			yCorrection:   y,
			course:        s.course,
		}
	}

	var _loc4 = make([][]*block, width)
	for x := 0; x < width; x++ {
		_loc4[x] = make([]*block, height)
		for y := 0; y < height; y++ {
			byt := raw[y][x]
			if byt == ' ' {
				continue
			}
			var _loc1 = &block{}
			_loc1.t = parseT(byt)
			switch _loc1.t {
			case blockIce:
				_loc1.traction = .2
			case blockP1:
				s.course.guys[0] = newPlayer(float64(x*30+15), float64(y*30+15))
				continue
			case blockP2:
				s.course.guys[1] = newPlayer(float64(x*30+15), float64(y*30+15))
				continue
			case blockP3:
				s.course.guys[2] = newPlayer(float64(x*30+15), float64(y*30+15))
				continue
			case blockP4:
				s.course.guys[3] = newPlayer(float64(x*30+15), float64(y*30+15))
				continue
			default:
				_loc1.traction = 1
			}

			_loc1.x = float64(x) * segSize
			_loc1.y = float64(y) * segSize
			_loc4[x][y] = _loc1
		}
	}
	s.course.blocks = _loc4
}

func (s *state) squashTest() {
	var _loc7 float64
	var _loc5 float64
	for _, guyA := range s.course.guys {
		if guyA.yVel < -0.5 {
			for _, guyB := range s.course.guys {
				if guyA == guyB {
					continue
				}
				_loc7 = guyA.x - guyB.x
				_loc5 = guyA.y - guyB.y
				// fmt.Println("check1:", _loc5 < (0-pheight)/2)
				// fmt.Println("check2:", _loc5 > 0-pheight+10)
				if math.Abs(_loc7) < 25 && _loc5 < (pheight) && _loc5 > 0-pheight+10 && guyA.recoveryTimer <= 0 && guyB.recoveryTimer <= 0 {
					//_root.startGameSound(guyA.m,75,"headBounce");
					guyA.yVel = 12
					if guyB == s.course.me {
						guyB.controlChange("squashed", true)
						guyB.squashTimer = 80
						guyB.recoveryTimer = guyB.squashTimer + 25
					}
				}
			}
		}
	}
}

type laser struct {
	x, y    float64
	vel     float64
	shooter *player
}

// probably wrong, copied from pr2tas (fix later)
func (l *laser) hitTest(p *player) bool {
	var left, right float64
	if l.vel > 0 {
		right = 44
		left = 0
	} else {
		right = 0
		left = 33
	}
	if l.x-right < p.x+11 && l.x+left > p.x-11 {
		if l.y < p.y && l.y > p.y-57 {
			return true
		}
	}
	return false
}

type mine struct {
	x, y       float64
	placeTimer int
}

type blockType int

const (
	blockBasic blockType = iota
	blockP1
	blockP2
	blockP3
	blockP4
	blockItem   // 5
	blockBrick  // 6
	blockFinish // 7
	blockMine   // 8
	blockIce    // 9
	blockMetal  // A
	blockLeft   // B
	blockRight  // C
	blockUp     // D
)

type itemType int

const (
	itemNone itemType = iota
	itemSuperJump
	itemLightning
	itemSpeed
	itemMine
	itemGun
	itemJetPack
	itemTeleport
)

func (i itemType) String() string {
	switch i {
	case itemNone:
		return "none"
	case itemSuperJump:
		return "superJump"
	case itemLightning:
		return "lightning"
	case itemSpeed:
		return "speed"
	case itemMine:
		return "mine"
	case itemGun:
		return "gun"
	case itemJetPack:
		return "jetPack"
	case itemTeleport:
		return "teleport"
	default:
		return "???"
	}
}

var items = [...]itemType{itemNone, itemSuperJump, itemLightning, itemSpeed, itemMine, itemGun, itemJetPack, itemTeleport}

type playerMode int

const (
	modeStand playerMode = iota
	modeRun
	modeJump
	modeSuperJump
	modeHurt
	modeSquashed
	modeBumped
	modeFinish
)

func (m playerMode) String() string {
	switch m {
	case modeStand:
		return "stand"
	case modeRun:
		return "run"
	case modeJump:
		return "jump"
	case modeSuperJump:
		return "superJump"
	case modeHurt:
		return "hurt"
	case modeSquashed:
		return "squashed"
	case modeBumped:
		return "bumped"
	case modeFinish:
		return "finish"
	default:
		return "???"
	}
}

type player struct {
	username                 string
	u, r, l, d, space        bool
	squashed, bumped         bool
	finished, forfeit        bool
	waiting                  bool
	ground                   bool
	killJump                 bool
	weapon                   itemType
	x, y                     float64
	bullets, jetFuel         int
	speed, jump, traction    int
	xVel, yVel, xVelTarget   float64
	superJump                int
	attackTimer              int
	jumpVel                  float64
	speedMod, speedModTimer  float64
	runTimer, runTimerReset  int
	recoveryTimer            int
	correctionTimer          int
	xScale                   int
	xCorrection, yCorrection float64
	stand, safe              *block
	mode                     playerMode

	course       *course
	sneamia      bool
	sendToSocket func(string)
	squashTimer  int
	bumpedTimer  int
}

func (p *player) move() {
	if p.correctionTimer > 0 {
		p.x += p.xCorrection
		p.y += p.yCorrection
		p.correctionTimer--
	}

	p.yVel -= gravity
	p.xVelTarget = 0
	p.attackTimer--
	if !p.bumped && !p.squashed && !p.waiting {
		if p.r {
			p.xVelTarget = 3 + float64(p.speed)/40
			p.xScale = +1
		} else if p.l {
			p.xVelTarget = 0 - (3 + float64(p.speed)/40)
			p.xScale = -1
		}
		if p.u {
			if p.ground && p.superJump <= 25 {
				p.killJump = false
				p.jumpVel = 2.4 + float64(p.jump)/60
				//_root.startGameSound(_loc3,35,"jump");
			}
			if !p.killJump {
				p.yVel += p.jumpVel
				p.jumpVel *= 0.75
			}
		} else {
			p.killJump = true
		}
		if p.d {
			if !p.ground {
				if !p.sneamia {
					p.yVel -= 0.5
				} else {
					p.yVel -= 5
				}
				p.superJump = 0
			} else {
				if p.superJump < 100 {
					p.superJump++
				}
				if p.superJump > 25 {
					p.xVelTarget = 0
				}
			}
		} else {
			if p.superJump > 25 {
				p.yVel = float64(p.superJump) / 5
				//startGameSound(c.m,75,"superJump");
			}
			p.superJump = 0
		}
		if p.space && p.attackTimer <= 0 {
			p.useItem()
		}
	}
	if p.speedModTimer > 0 {
		p.speedModTimer--
		p.xVelTarget *= p.speedMod
		if p.speedModTimer == 1 {
			//_root.startGameSound(c.m,75,"slowDown");
		}
	}
	if p.bumpedTimer > 0 {
		p.bumpedTimer--
		if p.bumpedTimer == 1 {
			p.bumped = false
		}
	}
	if p.squashTimer > 0 {
		p.squashTimer--
		if p.squashTimer == 1 {
			p.squashed = false
			p.yVel = 7
		}
	}
	var _loc5 = (float64(p.traction) + 20) / 1000
	if p.stand != nil {
		_loc5 *= p.stand.traction
	}
	p.xVel = p.xVel - (p.xVel-p.xVelTarget)*_loc5
	if p.xVel > 29 {
		p.xVel = 29
	} else if p.xVel < -29 {
		p.xVel = -29
	}
	if p.yVel < -29 {
		p.yVel = -29
	}
	if p.ground == false {
		p.runTimer = 0
	}
	p.x += p.xVel
	p.y += p.yVel
	if p.recoveryTimer > 0 {
		p.recoveryTimer--
		// var _loc12 = p.recoveryTimer % 8
		// if _loc12 >= 4 {
		// 	_loc3._alpha = 50
		// } else {
		// 	_loc3._alpha = 75
		// }
	} else {
		//_loc3._alpha = 100
	}
	if p.finished != false || p.forfeit != false {
		p.mode = modeFinish
	} else if p.bumped {
		p.mode = modeBumped
	} else if p.squashed {
		p.mode = modeSquashed
	} else if !p.ground {
		p.mode = modeJump
	} else if p.superJump > 25 {
		p.mode = modeSuperJump
	} else if p.l || p.r {
		p.mode = modeRun
		p.runTimer--
		if p.runTimer <= 0 {
			//_root.startGameSound(_loc3,50,"run" + Math.ceil(Math.random() * 4));
			p.runTimer = p.runTimerReset
		}
	} else {
		p.mode = modeStand
	}
	if p.y < p.course.deathHeight {
		p.x = p.safe.x + (segSize / 2)
		p.y = p.safe.y + segSize + segSize/2
		p.xVel = 0
		p.xVelTarget = 0
		p.yVel = 0
		p.recoveryTimer = 50
	}
	var _loc4 *block
	var _loc10 int
	var _loc11 int
	if p.yVel > 0 {
		var _loc7 = 50.0 // loc7 is undefined in pr1
		_loc10 = int(math.Floor(p.x / segSize))
		_loc11 = int(math.Floor((p.y + _loc7) / segSize))
		if _loc10 > -1 && _loc10 < len(p.course.blocks) && _loc11 > -1 && _loc11 < len(p.course.blocks[0]) {
			_loc4 = p.course.blocks[(_loc10)][(_loc11)]
			if _loc4 != nil {
				_loc4.onBump(p)
				//bumpAnimation(_loc4.m)
				for _, guy := range p.course.guys { //bumpPlayersTest(_loc4)
					if guy.stand != nil && guy.stand.x == _loc4.x && guy.stand.y == _loc4.y {
						guy.onBump()
					}
				}
				p.y = float64(_loc11)*segSize - /*segSize - segSize*/ _loc7
				p.yVel *= -0.5
			}
		}
		p.ground = false
		p.stand = nil /*"none"*/
	} else if p.yVel < 0 {
		_loc10 = int(math.Floor(p.x / segSize))
		_loc11 = int(math.Floor(p.y / segSize))
		if _loc10 > -1 && _loc10 < len(p.course.blocks) && _loc11 > -1 && _loc11 < len(p.course.blocks[0]) {
			_loc4 = p.course.blocks[(_loc10)][(_loc11)]
		}
		if _loc4 != nil {
			p.y = (float64(_loc11) * segSize) + segSize
			p.yVel = 0
			p.ground = true
			p.stand = _loc4
			_loc4.onStand(p)
		} else {
			p.ground = false
			p.stand = nil
		}
	}
	var _loc9 = 10.0 // loc9 is undefined in pr1
	if p.xVel > 0 {
		_loc10 = int(math.Floor((p.x + _loc9) / segSize))
		_loc11 = int(math.Floor((p.y /*- 10*/) / segSize))
		if _loc10 > -1 && _loc10 < len(p.course.blocks) && _loc11 > -1 && _loc11 < len(p.course.blocks[0]) {
			_loc4 = p.course.blocks[(_loc10)][(_loc11)]
			if _loc4 != nil {
				p.x = float64(_loc10)*segSize - _loc9
				p.xVel = (0 - p.xVel) * 0.25
				p.xVelTarget = 0
				_loc4.onLeftHit(p)
			}
		}
	} else {
		_loc9 = 10
		_loc10 = int(math.Floor((p.x - _loc9) / segSize))
		_loc11 = int(math.Floor((p.y /*- 10*/) / segSize))
		if _loc10 > -1 && _loc10 < len(p.course.blocks) && _loc11 > -1 && _loc11 < len(p.course.blocks[0]) {
			_loc4 = p.course.blocks[(_loc10)][(_loc11)]
			if _loc4 != nil {
				p.x = float64(_loc10)*segSize + segSize + _loc9
				p.xVel = (0 - p.xVel) * 0.25
				p.xVelTarget = 0
				_loc4.onRightHit(p)
			}
		}
	}
	//    p.nameBox._x = _loc3._x;
	//    p.nameBox._y = _loc3._y - 75;
}

func (p *player) getRandomItem() {
	p.weapon = items[rand.Intn(len(items)-1)+1]
	p.activateWeapon()
}

func (p *player) activateWeapon() {
	switch p.weapon {
	case itemGun:
		p.bullets = 3
	case itemJetPack:
		p.jetFuel = 200
	case itemLightning:
	case itemMine:
	case itemTeleport:
	case itemSuperJump:
	case itemSpeed:
	}
}

func (p *player) useItem() {
	//var _loc6 = p.weapon
	//player.m.anim.weapon.gotoAndStop(_loc6);

	switch p.weapon {
	case itemGun:
		if !p.sneamia {
			p.attackTimer = 50
		}
		// var _loc5 = new Object();
		// _loc5.x = player.m.anim.weapon._x;
		// _loc5.y = player.m.anim.weapon._y;
		// player.m.anim.localToGlobal(_loc5);
		// _root.cam.effects.globalToLocal(_loc5);
		l := &laser{
			y:       p.y + (pheight / 2),
			shooter: p,
		}
		p.course.lasers = append(p.course.lasers, l)
		if p.xScale > 0 {
			p.xVel -= 5
			l.vel = 20
			l.x = p.x + 20
		} else {
			p.xVel += 5
			l.vel = -20
			l.x = p.x - 20
		}
		p.bullets--
		if p.bullets <= 0 && !p.sneamia {
			p.weapon = itemNone
		}
	case itemLightning:
		for _, guy := range p.course.guys {
			if guy != p {
				//    var _loc4 = addMC(_root.cam.effects,"lightningStrike");
				//    _loc4._x = guy.m._x;
				//    _loc4._y = guy.m._y;
				//    _loc4.onEnterFrame = func(){
				//       this._alpha = this._alpha - 15;
				//       if(this._alpha <= 0){
				//          removeMovieClip(this);
				//       }
				//    };
				if !p.sneamia {
					guy.bumped = true
				}
			}
		}
		if !p.sneamia {
			p.weapon = itemNone
		}
	case itemMine:
		placeX := int(math.Floor(p.x / segSize))
		placeY := int(math.Floor((p.y + 20) / segSize))
		var _loc4 = p.course.blocks[placeX][placeY]
		if _loc4 == nil {
			//_root.startGameSound(this.m,75,"mineAppear");
			p.course.placedMines = append(p.course.placedMines, &mine{x: float64(placeX) * segSize, y: float64(placeY) * segSize, placeTimer: 41})
			if !p.sneamia {
				p.weapon = itemNone
			}
		}
	case itemTeleport:
		toX := int(math.Floor((p.x + (115 * float64(p.xScale))) / segSize))
		toY := int(math.Floor((p.y - pheight/2) / segSize))
		var toTile = p.course.blocks[toX][toY]
		// TODO: prevent teleporting to deleted mines and bricks -- broken flag?
		if toTile != nil /*&& toTile.explodable == nil*/ /*(toTile.explodable == null || toTile.explodable == undefined)*/ {
			// _root.startGameSound(this.m,75,"teleport");
			// var _loc4 = addMC(_root.cam.effects,"teleportEffect");
			// var _loc3 = addMC(_root.cam.effects,"teleportEffect");
			// _loc4._y = player.m._y - player.m._height/2
			// _loc3._y = player.m._y - player.m._height/2
			// _loc4._x = player.m._x
			if p.xScale > 0 {
				p.x += 150
			} else {
				p.x -= 150
			}
			//_loc3._x = player.m._x
			if !p.sneamia {
				p.weapon = itemNone
			}
		}
	case itemSuperJump:
		//_root.startGameSound(this.m,75,"superJump");
		p.yVel += 20
		if !p.sneamia {
			p.weapon = itemNone
		}
	case itemJetPack:
		if !p.sneamia && p.yVel < 5 {
			p.yVel += 0.65
		} else if p.sneamia && p.yVel < 30 {
			p.yVel += 5
		}
		p.jetFuel--
		//this.m.anim.weapon.jetPack.gotoAndStop("on");
		if p.jetFuel <= 0 && !p.sneamia {
			p.weapon = itemNone
		}
	case itemSpeed:
		//_root.startGameSound(this.m,75,"speedUp");
		p.speedMod = 3
		p.speedModTimer = 100
		p.weapon = itemNone
	}
}

func (p *player) onBump() {
	p.yVel = 7
	if p == p.course.me {
		p.controlChange("bumped", true)
		p.bumpedTimer = 97
	}
}

func bumpAnimation(b *block) {
	const bumpVel = -20

}

func (p *player) sendPositionUpdate() {
	if len(p.course.guys) >= 2 && p.sendToSocket != nil {
		p.sendToSocket("update_pos`" + p.getPositionString())
	}
}

func (p *player) getPositionString() string {
	return fmt.Sprintf("%d,%d,%d,%d,%d",
		int(math.Round(p.x)),
		int(math.Round(p.y)),
		int(math.Round(p.xVel*10)/10),
		int(math.Round(p.yVel*10)/10),
		int(math.Round(p.xVelTarget*10)/10))
}

func getMS() float64 { return float64(time.Now().UnixNano() / 1e6) }

func (p *player) controlChange(variable string, val bool) {
	var meVar *bool
	// switch here in place of the reflection
	switch variable {
	case "r":
		meVar = &p.r
	case "l":
		meVar = &p.l
	case "u":
		meVar = &p.u
	case "d":
		meVar = &p.d
	case "space":
		meVar = &p.space
	case "bumped":
		meVar = &p.bumped
	case "squashed":
		meVar = &p.squashed
	default:
		panic("unknown controlChange(): " + variable)
	}
	if *meVar != val {
		*meVar = val
		if len(p.course.guys) > 1 && p.sendToSocket != nil {
			p.sendToSocket("control_change`" + variable + "`" + strconv.FormatBool(val) + "`" + p.getPositionString())
		}
	}
}

type course struct {
	name        string
	deathHeight float64
	blocks      [][]*block
	me          *player
	guys        []*player
	placedMines []*mine
	lasers      []*laser
}

type block struct {
	x, y     float64
	t        blockType
	traction float64
}

func (b *block) onStand(p *player) {
	bx := int(math.Floor(b.x / 30))
	by := int(math.Floor(b.y / 30))
	switch b.t {
	case blockLeft:
		p.safe = p.course.blocks[bx][by]
		p.xVel -= 0.3
	case blockRight:
		p.safe = p.course.blocks[bx][by]
		p.xVel += 0.3
	case blockUp:
		p.safe = p.course.blocks[bx][by]
		p.yVel = -10
	case blockBasic, blockMetal, blockItem, blockFinish, blockIce:
		p.safe = p.course.blocks[bx][by]
	case blockMine:
		p.xVel = (p.x - b.x) * 5
		p.yVel = 0 - (p.y - b.y)
		if p == p.course.me && !p.sneamia {
			p.controlChange("bumped", true)
			p.bumpedTimer = 97
		}
		p.course.blocks[bx][by] = nil
	}
}

func (b *block) onBump(p *player) {
	bx := int(math.Floor(b.x / 30))
	by := int(math.Floor(b.y / 30))
	switch b.t {
	case blockBrick:
		p.course.blocks[bx][by] = nil
	case blockItem:
		if p == p.course.me {
			p.getRandomItem()
			if p.sendToSocket != nil {
				p.sendToSocket("request_item`")
			}
		}
	case blockFinish:
		p.finished = true
	case blockMine:
		p.xVel = (p.x - b.x) * 5
		p.yVel = 0 - (p.y - b.y)
		if p == p.course.me && !p.sneamia {
			p.controlChange("bumped", true)
			p.bumpedTimer = 97
		}
		p.course.blocks[bx][by] = nil
	}
}

func (b *block) onLeftHit(p *player) {
	bx := int(math.Floor(b.x / 30))
	by := int(math.Floor(b.y / 30))
	switch b.t {
	case blockMine:
		p.xVel = (p.x - b.x) * 5
		p.yVel = 0 - (p.y - b.y)
		if p == p.course.me && !p.sneamia {
			p.controlChange("bumped", true)
			p.bumpedTimer = 97
		}
		p.course.blocks[bx][by] = nil
	}
}

func (b *block) onRightHit(p *player) {
	bx := int(math.Floor(b.x / 30))
	by := int(math.Floor(b.y / 30))
	switch b.t {
	case blockMine:
		p.xVel = (p.x - b.x) * 5
		p.yVel = /*0 -*/ (p.y - b.y)
		if p == p.course.me && !p.sneamia {
			p.controlChange("bumped", true)
			p.bumpedTimer = 97
		}
		p.course.blocks[bx][by] = nil
	}
}
