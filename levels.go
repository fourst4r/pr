package main

// func parseT(t byte) blockType {
// 	switch t {
// 	case '0':
// 		return blockType(0)
// 	case '1':
// 		return blockType(1)
// 	case '2':
// 		return blockType(2)
// 	case '3':
// 		return blockType(3)
// 	case '4':
// 		return blockType(4)
// 	case '5':
// 		return blockType(5)
// 	case '6':
// 		return blockType(6)
// 	case '7':
// 		return blockType(7)
// 	case '8':
// 		return blockType(8)
// 	case '9':
// 		return blockType(9)
// 	case 'A':
// 		return blockType(10)
// 	case 'B':
// 		return blockType(11)
// 	case 'C':
// 		return blockType(12)
// 	case 'D':
// 		return blockType(13)
// 	default:
// 		panic("cant parse t: " + string(t))
// 	}
// }

// func newPlayer(x, y float64) *player {
// 	return &player{
// 		x: x, y: y,
// 		speed: 50, jump: 50, traction: 50,
// 		waiting:       false,
// 		speedMod:      1,
// 		runTimerReset: 7,
// 		recoveryTimer: 75,
// 		xScale:        +1,
// 		xCorrection:   x,
// 		yCorrection:   y,
// 	}
// }

// func loadLevel(name string) {
// 	raw := levels[name]
// 	for i := len(raw)/2 - 1; i >= 0; i-- {
// 		opp := len(raw) - 1 - i
// 		raw[i], raw[opp] = raw[opp], raw[i]
// 	}
// 	width := len(raw[0])                    //+ 1
// 	height := len(raw)                      //+ 5
// 	deathHeight = -float64(height*30 + 300) // might be wrong
// 	guys = make([]*player, 4)

// 	var _loc4 = make([][]*bloc, width)
// 	for x := 0; x < width; x++ {
// 		_loc4[x] = make([]*bloc, height)
// 		for y := 0; y < height; y++ {
// 			byt := raw[y][x]
// 			if byt == ' ' {
// 				continue
// 			}
// 			var _loc1 = &bloc{}
// 			_loc1.t = parseT(byt)
// 			switch _loc1.t {
// 			case blockIce:
// 				_loc1.traction = .2
// 			case blockP1:
// 				guys[0] = newPlayer(float64(x*30+15), float64(y*30+15))
// 				continue
// 			case blockP2:
// 				guys[1] = newPlayer(float64(x*30+15), float64(y*30+15))
// 				continue
// 			case blockP3:
// 				guys[2] = newPlayer(float64(x*30+15), float64(y*30+15))
// 				continue
// 			case blockP4:
// 				guys[3] = newPlayer(float64(x*30+15), float64(y*30+15))
// 				continue
// 			default:
// 				_loc1.traction = 1
// 			}

// 			_loc1.x = float64(x) * segSize
// 			_loc1.y = float64(y) * segSize
// 			_loc4[x][y] = _loc1
// 		}
// 	}
// 	course = _loc4
// }

var levels = map[string][]string{
	"newbieland": {
		"                                                                                                                     000000000000000",
		"                                                                                                                     0             0",
		"                                                                                                                     0          7  0",
		"                                                                                                                     0             0",
		"                                                                                                   000000000000000   0             0",
		"                                                                                                                     066666000000000",
		"                                                                                  000000000000000                                  0",
		"                                                                                                                                   0",
		"                                                                            000000                                                 0",
		"                                                                                                                   00000000000000000",
		"  1234                       0 0                                                                                                    ",
		"                            00 00                                                                                                   ",
		"0000000000000000000000000000000000                                                                                                  ",
		"                                  000                                                                                               ",
		"                                    0                                                                                               ",
		"                                             5                                                                                      ",
		"                                                                                                                                    ",
		"                                                             000000000000000                                                        ",
		"                                     000000000000000000000000                                                                       ",
	},
	"buto":         {},
	"pyramids":     {},
	"robocity":     {},
	"assembly":     {},
	"infernal hop": {},
	"going down":   {},
	"slip":         {},
}

// func defineWalls(m [][]*tile) [][]*tile {
// 	var it1 = 0 // it1 is not defined in pr1
// 	for it1 < len(m) {
// 		var it2 = 0 // it2 is not defined in pr1
// 		for it2 < len(m[it1]) {
// 			if m[it1][it2].wall == true {
// 				if m[it1-1][it2].wall == false {
// 					m[it1][it2].leftWall = true
// 				}
// 				if m[it1+1][it2].wall == false {
// 					m[it1][it2].rightWall = true
// 				}
// 				if m[it1][it2-1].wall == false {
// 					m[it1][it2].topWall = true
// 				}
// 				if m[it1][it2+1].wall == false {
// 					m[it1][it2].bottomWall = true
// 				}
// 			}
// 			it2++
// 		}
// 		it1++
// 	}
// 	return m
// }

// function drawMap(num)
// {
//    var _loc2 = addMC(cam.map,"map" + num);
//    var _loc3 = Math.ceil(_loc2._width / segSize) + 1;
//    var _loc4 = Math.ceil(_loc2._height / segSize) + 5;
//    _root.deathHeight = _loc2._height + 300;
//    _root.map = makeBlankMap(_loc3,_loc4);
// }

// func drawMap(name string) {
// 	lvl := levels[name]
// 	width := len(lvl[0]) + 1
// 	height := len(lvl) + 5
// 	deathHeight = float64(height + 300)
// 	mapp = makeBlankMap(width, height)
// }

// func makeBlankMap(xSegs, ySegs int) [][]*tile {
// 	var _loc4 = make([][]*tile, xSegs)
// 	var _loc3 = 0 // loc3 is undefined in pr1
// 	for _loc3 < xSegs {
// 		_loc4[_loc3] = make([]*tile, ySegs)
// 		var _loc2 = 0 // loc2 is undefined in pr1
// 		for _loc2 < ySegs {
// 			var _loc1 = &tile{} //new Object();
// 			//_loc1.things = new Array(); <- never used?
// 			_loc1.wall = false
// 			_loc1.i = float64(_loc3)
// 			_loc1.j = float64(_loc2)
// 			_loc1.x = float64(_loc3) * segSize
// 			_loc1.y = float64(_loc2) * segSize
// 			_loc4[_loc3][_loc2] = _loc1
// 			_loc2++
// 		}
// 		_loc3++
// 	}
// 	return _loc4
// }

// func addStartPos(num int, m mc) {
// 	var _loc1 = mcToTile(m)
// 	var _loc2 = guys[num-1]
// 	_loc2.m._x = _loc1.x + segSize/2
// 	_loc2.m._y = _loc1.y + segSize/2
// }

// func addBG(linkage, ratio){
//    var _loc1 = addMC(cam.bg,linkage);
//    _loc1.ratio = ratio;
//    bgArray.push(_loc1);
// }

// func addBlock(linkage string, m *mTile) *block {
// 	var pos = mcToTile(m)
// 	var block = addMC(cam.mapp, linkage)
// 	block.x = pos.x
// 	block._x = pos.x
// 	block.y = pos.y
// 	block._y = pos.y
// 	block.i = pos.i
// 	block.j = pos.j
// 	block.tile = pos
// 	block.traction = 1
// 	pos.m = block
// 	pos.wall = true
// 	return block
// }
// func addBasicBlock(m *mTile) {
// 	var block = addBlock("basicBlock", m)
// 	//    block.onBump = func(player){
// 	//       _root.startGameSound(this,75,"thump");
// 	//    };
// 	block.onStand = func(player *player) {
// 		player.lastSafeTile = mapp[block.i][block.j-1]
// 	}
// }
// func addMetalBlock(m *mTile) {
// 	var block = addBlock("metalBlock", m)
// 	//    block.onBump = func(player){
// 	//       _root.startGameSound(this,75,"thump");
// 	//    };
// 	block.onStand = func(player *player) {
// 		player.lastSafeTile = mapp[block.i][block.j-1]
// 	}
// }
// func addRightBlock(m *mTile) {
// 	var block = addBlock("rightBlock", m)
// 	//    block.onBump = func(player){
// 	//       _root.startGameSound(this,75,"thump");
// 	//    };
// 	block.onStand = func(player *player) {
// 		player.lastSafeTile = mapp[block.i][block.j-1]
// 		player.xVel = player.xVel + 0.3
// 		//   this.play();
// 	}
// }
// func addLeftBlock(m *mTile) {
// 	var block = addBlock("leftBlock", m)
// 	//    block.onBump = func(player){
// 	//       _root.startGameSound(this,75,"thump");
// 	//    };
// 	block.onStand = func(player *player) {
// 		player.lastSafeTile = mapp[block.i][block.j-1]
// 		player.xVel = player.xVel - 0.3
// 		//this.play();
// 	}
// }
// func addUpBlock(m *mTile) {
// 	var block = addBlock("upBlock", m)
// 	//    block.onBump = func(player){
// 	//       _root.startGameSound(this,75,"thump");
// 	//    };
// 	block.onStand = func(player *player) {
// 		player.lastSafeTile = mapp[block.i][block.j-1]
// 		player.yVel = -10
// 	}
// }
// func addWaffleBlock(m *mTile) {
// 	var block = addBlock("waffleBlock", m)
// 	block.onBump = nil //0
// 	block.onStand = func(player *player) {
// 		player.lastSafeTile = mapp[block.i][block.j-1]
// 	}
// }
// func addFinishBlock(m *mTile) {
// 	var block = addBlock("finishBlock", m)
// 	block.onBump = func(player *player) {
// 		player.finished = true
// 		//   _root.startGameSound(this,75,"victory");
// 		//   if(player == me){
// 		//      _root.menu_mc.gotoAndStop("waiting");
// 		//      _root.camFree();
// 		//      _root.sendToSocket("finish_level`" + math.Round((getMS() - startMS) / 10) / 100);
// 		//   }
// 	}
// 	block.onStand = func(player *player) {
// 		player.lastSafeTile = mapp[block.i][block.j-1]
// 	}
// }
// func addMineBlock(m *mTile) {
// 	var block = addBlock("mineBlock", m)
// 	block.tile.explodable = true
// 	//block.gotoAndStop("rest");
// 	block.onExplode = func() {
// 		//this.gotoAndPlay("explode");
// 		block.tile.wall = false
// 		block.tile.rightWall = false
// 		block.tile.leftWall = false
// 		block.tile.topWall = false
// 		block.tile.bottomWall = false
// 		//_root.startGameSound(this,75,"explosion");
// 	}
// 	temp := func(player *player) {
// 		var _loc4 = player.m._x - block._x
// 		var _loc3 = player.m._y - block._y
// 		player.xVel = _loc4 * 5
// 		player.yVel = 0 - _loc3
// 		if player == me && sneamia == false {
// 			controlChange("bumped", "true")
// 		}
// 		block.onExplode()
// 	}
// 	block.onBump = temp
// 	block.onStand = temp
// 	block.onLeftHit = temp
// 	block.onRightHit = temp
// }
// func addIceBlock(m *mTile) {
// 	var block = addBlock("iceBlock", m)
// 	block.traction = 0.2
// 	// block.onBump = func(player *player) {
// 	// 	_root.startGameSound(this, 75, "thump")
// 	// }
// 	block.onStand = func(player *player) {
// 		player.lastSafeTile = mapp[block.i][block.j-1]
// 	}
// }
// func getMyItem(item string, player *player) {
// 	player.weapon = item
// 	activateWeapon(player)
// 	//_root.menu_mc.weapon.gotoAndStop(item);
// }

// func addStarBlock(m *mTile) {
// 	var block = addBlock("starBlock", m)
// 	//    block.onBump = func(player){
// 	//       if(player == me){
// 	//          _root.startGameSound(this,75,"star");
// 	//          _root.startGameSound(this,75,"thump");
// 	//          this.gotoAndStop("off");
// 	//          _root.sendToSocket("request_item`");
// 	//          this.onBump = func(){
// 	//             _root.startGameSound(this,75,"thump");
// 	//          };
// 	//       }
// 	//    };
// 	block.onStand = func(player) {
// 		player.lastSafeTile = mapp[block.i][block.j-1]
// 	}
// }
// func addBrickBlock(m *mTile) {
// 	var block = addBlock("brickBlock", m)
// 	block.tile.explodable = true
// 	block.onBump = func(player) {
// 		block.onExplode()
// 	}
// 	block.onExplode = func() {
// 		//   _root.startGameSound(this,75,"brickBreak");
// 		//   for(it < 10){
// 		//      var brickPiece = addMC(_root.cam.effects,"brickPiece");
// 		//      brickPiece.gotoAndStop(Math.ceil(Math.random() * _loc3._totalframes));
// 		//      brickPiece._x = Math.random() * 30 + this._x;
// 		//      brickPiece._y = Math.random() * 30 + this._y;
// 		//      brickPiece._rotation = Math.random() * 360;
// 		//      brickPiece.rotVel = Math.random() * 10 - 5;
// 		//      brickPiece.xVel = Math.random() * 10 - 5;
// 		//      brickPiece.yVel = Math.random() * 10 - 10;
// 		//      brickPiece.onEnterFrame = func(){
// 		//         this.yVel = this.yVel + _root.gravity;
// 		//         this.xVel = this.xVel * _root.friction;
// 		//         this.yVel = this.yVel * _root.friction;
// 		//         this._x = this._x + this.xVel;
// 		//         this._y = this._y + this.yVel;
// 		//         this._rotation = this._rotation + this.rotVel;
// 		//         this._alpha = this._alpha - 5;
// 		//         if(this._alpha <= 0){
// 		//            removeMovieClip(this);
// 		//         }
// 		//      };
// 		//      it++;
// 		//   }
// 		block.tile.wall = false
// 		block.tile.rightWall = false
// 		block.tile.leftWall = false
// 		block.tile.topWall = false
// 		block.tile.bottomWall = false
// 		removeMovieClip(block)
// 	}
// }

// func mcToTile(m mc) *tile {
// 	var _loc1 = mcToSeg(m)
// 	var tile = mapp[_loc1.xSeg][_loc1.ySeg]
// 	return tile
// }

// func mcToSeg(m mc) *block {
// 	var _loc1 = posToSeg(m._x, m._y)
// 	return _loc1
// }

// func posToSeg(x, y float64) *block {
// 	var _loc1 = &block{} //new Object();
// 	_loc1.xSeg = math.Round(x / segSize)
// 	_loc1.ySeg = math.Round(y / segSize)
// 	_loc1.x = _loc1.xSeg * segSize
// 	_loc1.y = _loc1.ySeg * segSize
// 	return _loc1
// }

// func posToTile(x, y float64) *tile {
// 	var _loc1 = posToSeg(x, y)
// 	var _loc2 = mapp[_loc1.xSeg][_loc1.ySeg]
// 	return _loc2
// }
