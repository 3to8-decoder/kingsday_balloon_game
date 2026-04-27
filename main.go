package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	gameTime     = 60 // seconds
)

var (
	orangeLight = color.RGBA{0xFF, 0xA0, 0x00, 0xFF}
	orangeDark  = color.RGBA{0xE8, 0x6C, 0x00, 0xFF}
	red         = color.RGBA{0xE0, 0x30, 0x30, 0xFF}
	white       = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	blue        = color.RGBA{0x20, 0x40, 0xB0, 0xFF}
	yellow      = color.RGBA{0xFF, 0xE0, 0x00, 0xFF}
	bgColor     = color.RGBA{0x10, 0x10, 0x30, 0xFF}

	balloonColors = []color.RGBA{orangeLight, orangeDark, red, white, yellow, blue}
)

type Balloon struct {
	x, y      float64
	speedY    float64
	wobbleOff float64
	radius    float64
	color     color.RGBA
	alive     bool
}

type Particle struct {
	x, y   float64
	vx, vy float64
	life   float64
	color  color.RGBA
	radius float64
}

type Game struct {
	balloons   []Balloon
	particles  []Particle
	score      int
	combo      int
	comboTimer float64
	timeLeft   float64
	state      string // "menu", "playing", "gameover"
	highScore  int
	spawnTimer float64
	rng        *rand.Rand
	popTexts   []PopText
}

type PopText struct {
	x, y  float64
	text  string
	life  float64
	color color.RGBA
}

func NewGame() *Game {
	return &Game{
		state:    "menu",
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		timeLeft: gameTime,
	}
}

func (g *Game) spawnBalloon() {
	r := 18 + g.rng.Float64()*14
	b := Balloon{
		x:         r + g.rng.Float64()*(screenWidth-r*2),
		y:         screenHeight + r,
		speedY:    40 + g.rng.Float64()*60,
		wobbleOff: g.rng.Float64() * math.Pi * 2,
		radius:    r,
		color:     balloonColors[g.rng.Intn(len(balloonColors))],
		alive:     true,
	}
	g.balloons = append(g.balloons, b)
}

func (g *Game) spawnPop(x, y float64, c color.RGBA) {
	for i := 0; i < 12; i++ {
		angle := g.rng.Float64() * math.Pi * 2
		speed := 60 + g.rng.Float64()*120
		g.particles = append(g.particles, Particle{
			x:      x,
			y:      y,
			vx:     math.Cos(angle) * speed,
			vy:     math.Sin(angle) * speed,
			life:   0.4 + g.rng.Float64()*0.3,
			color:  c,
			radius: 2 + g.rng.Float64()*3,
		})
	}
}

func (g *Game) startGame() {
	g.state = "playing"
	g.score = 0
	g.combo = 0
	g.comboTimer = 0
	g.timeLeft = gameTime
	g.balloons = nil
	g.particles = nil
	g.popTexts = nil
	g.spawnTimer = 0
}

func (g *Game) Update() error {
	switch g.state {
	case "menu":
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.startGame()
		}
		return nil

	case "gameover":
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = "menu"
		}
		return nil
	}

	// --- playing ---
	dt := 1.0 / 60.0

	// timer
	g.timeLeft -= dt
	if g.timeLeft <= 0 {
		g.timeLeft = 0
		if g.score > g.highScore {
			g.highScore = g.score
		}
		g.state = "gameover"
		return nil
	}

	// combo decay
	g.comboTimer -= dt
	if g.comboTimer <= 0 {
		g.combo = 0
	}

	// spawn
	g.spawnTimer -= dt
	if g.spawnTimer <= 0 {
		g.spawnBalloon()
		// spawn faster as time goes on
		elapsed := gameTime - g.timeLeft
		interval := 0.8 - elapsed*0.008
		if interval < 0.25 {
			interval = 0.25
		}
		g.spawnTimer = interval + g.rng.Float64()*0.3
	}

	// update balloons
	for i := range g.balloons {
		b := &g.balloons[i]
		b.y -= b.speedY * dt
		b.x += math.Sin(b.wobbleOff+float64(time.Now().UnixNano())/1e9*1.5) * 0.5
		if b.y+b.radius < 0 {
			b.alive = false
		}
	}
	// remove dead
	g.balloons = filterBalloons(g.balloons)

	// click to pop
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		// check from top (last drawn) to bottom
		for i := len(g.balloons) - 1; i >= 0; i-- {
			b := &g.balloons[i]
			dx := float64(mx) - b.x
			dy := float64(my) - b.y
			if dx*dx+dy*dy < b.radius*b.radius {
				// pop!
				g.combo++
				g.comboTimer = 1.0
				points := g.combo
				g.score += points
				g.spawnPop(b.x, b.y, b.color)
				// pop text
				txt := fmt.Sprintf("+%d", points)
				if g.combo >= 5 {
					txt = fmt.Sprintf("+%d COMBO x%d!", points, g.combo)
				} else if g.combo >= 3 {
					txt = fmt.Sprintf("+%d x%d!", points, g.combo)
				}
				g.popTexts = append(g.popTexts, PopText{
					x: b.x, y: b.y, text: txt, life: 1.0, color: b.color,
				})
				b.alive = false
				g.balloons = filterBalloons(g.balloons)
				break
			}
		}
	}

	// update particles
	for i := range g.particles {
		p := &g.particles[i]
		p.x += p.vx * dt
		p.y += p.vy * dt
		p.vy += 200 * dt // gravity
		p.life -= dt
	}
	// remove dead particles
	g.particles = filterParticles(g.particles)

	// update pop texts
	for i := range g.popTexts {
		t := &g.popTexts[i]
		t.y -= 40 * dt
		t.life -= dt
	}
	g.popTexts = filterPopTexts(g.popTexts)

	return nil
}

func filterBalloons(bs []Balloon) []Balloon {
	n := 0
	for _, b := range bs {
		if b.alive {
			bs[n] = b
			n++
		}
	}
	return bs[:n]
}

func filterParticles(ps []Particle) []Particle {
	n := 0
	for _, p := range ps {
		if p.life > 0 {
			ps[n] = p
			n++
		}
	}
	return ps[:n]
}

func filterPopTexts(ts []PopText) []PopText {
	n := 0
	for _, t := range ts {
		if t.life > 0 {
			ts[n] = t
			n++
		}
	}
	return ts[:n]
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)

	switch g.state {
	case "menu":
		g.drawMenu(screen)
	case "playing":
		g.drawGame(screen)
	case "gameover":
		g.drawGame(screen)
		g.drawGameOver(screen)
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	// title
	ebitenutil.DebugPrintAt(screen, "KING'S DAY", screenWidth/2-60, 80)
	ebitenutil.DebugPrintAt(screen, "BALLOON POP", screenWidth/2-60, 100)

	// decorative balloons
	for i := 0; i < 5; i++ {
		x := 100 + float64(i)*110
		y := 220 + math.Sin(float64(float64(time.Now().UnixNano())/1e9)*2+float64(i))*15
		drawBalloon(screen, x, y, 22, balloonColors[i%len(balloonColors)])
	}

	ebitenutil.DebugPrintAt(screen, "Click or press SPACE to start!", screenWidth/2-120, 340)
	if g.highScore > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("High Score: %d", g.highScore), screenWidth/2-60, 380)
	}
}

func (g *Game) drawGame(screen *ebiten.Image) {
	// draw balloons
	for _, b := range g.balloons {
		drawBalloon(screen, b.x, b.y, b.radius, b.color)
	}

	// particles
	for _, p := range g.particles {
		alpha := uint8(255 * (p.life / 0.7))
		if alpha > 255 {
			alpha = 255
		}
		c := color.RGBA{p.color.R, p.color.G, p.color.B, alpha}
		ebitenutil.DrawCircle(screen, p.x, p.y, p.radius, c)
	}

	// pop texts
	for _, t := range g.popTexts {
		ebitenutil.DebugPrintAt(screen, t.text, int(t.x)-20, int(t.y))
	}

	// HUD
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Time: %0.1f", g.timeLeft), screenWidth-100, 10)
	if g.combo >= 2 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Combo x%d", g.combo), 10, 28)
	}
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	overlay := ebiten.NewImage(screenWidth, screenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	ebitenutil.DebugPrintAt(screen, "TIME'S UP!", screenWidth/2-40, 140)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score), screenWidth/2-40, 200)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("High Score: %d", g.highScore), screenWidth/2-55, 230)
	ebitenutil.DebugPrintAt(screen, "Click or SPACE to continue", screenWidth/2-100, 300)
}

func drawBalloon(screen *ebiten.Image, x, y, r float64, c color.RGBA) {
	// balloon body
	ebitenutil.DrawCircle(screen, x, y, r, c)
	// highlight
	hx := x - r*0.3
	hy := y - r*0.3
	highlight := color.RGBA{
		min(c.R+80, 255),
		min(c.G+80, 255),
		min(c.B+80, 255),
		255,
	}
	ebitenutil.DrawCircle(screen, hx, hy, r*0.25, highlight)
	// string
	ebitenutil.DrawCircle(screen, x, y+r+3, 1.5, white)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	restore := suppressMetalStartupNoise()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("King's Day Balloon Pop 🧡")
	// Metal init happens inside RunGame on first frame;
	// restore stderr after a brief delay so the startup noise is swallowed
	time.AfterFunc(500*time.Millisecond, restore)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
