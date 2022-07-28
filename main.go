package main

import (
	"image/color"
	"log"

	"math"
	"math/rand"
	"time"

	// "github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/v2"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) calcAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius) //upper is upperrightcorner lower is lowerleftcorner
	avgVelocity := Vector2D{0, 0}
	count := 0.0

	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherBoidId].velocity)
				}
			}
		}
	}

	accel := Vector2D{0, 0}
	if count > 0 {
		avgVelocity = avgVelocity.DivisionV(count)
		accel = avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
	}

	return accel
}

func (b *Boid) moveOne() {
	b.velocity = b.velocity.Add(b.calcAcceleration()).limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	next := b.position.Add(b.velocity)
	if next.x >= screenWidth || next.x < 0 {
		b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	}
	if next.y >= screenHeight || next.y < 0 {
		b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	}
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(bid int) {
	b := Boid{
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
		id:       bid,
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	go b.start()
}

type Vector2D struct {
	x float64
	y float64
}

func (v1 Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{v1.x + v2.x, v1.y + v2.y}
}
func (v1 Vector2D) Subtract(v2 Vector2D) Vector2D {
	return Vector2D{v1.x - v2.x, v1.y - v2.y}
}
func (v1 Vector2D) Multiply(v2 Vector2D) Vector2D {
	return Vector2D{v1.x * v2.x, v1.y * v2.y}
}
func (v1 Vector2D) AddV(d float64) Vector2D {
	return Vector2D{v1.x + d, v1.y + d}
}
func (v1 Vector2D) DivisionV(d float64) Vector2D {
	return Vector2D{v1.x / d, v1.y / d}
}
func (v1 Vector2D) MultiplyV(d float64) Vector2D {
	return Vector2D{v1.x * d, v1.y * d}
}

func (v1 Vector2D) limit(lower, upper float64) Vector2D {
	return Vector2D{math.Min(math.Max(v1.x, lower), upper),
		math.Min(math.Max(v1.y, lower), upper)}
}

func (v1 Vector2D) Distance(v2 Vector2D) float64 {
	return math.Sqrt(math.Pow(v1.x-v2.x, 2) + math.Pow(v1.y-v2.y, 2))
}

const (
	screenWidth, screenHeight = 640, 360
	boidCount                 = 500
	viewRadius                = 13
	adjRate                   = 0.015
)

var (
	green   = color.RGBA{10, 255, 50, 255}
	boids   [boidCount]*Boid
	boidMap [screenWidth + 1][screenHeight + 1]int //a 2d array which will keep track of positions of all particles being simulated THis array will be memory shared
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, boid := range boids {
		screen.Set(int(boid.position.x+1), int(boid.position.y), green)
		screen.Set(int(boid.position.x-1), int(boid.position.y), green)
		screen.Set(int(boid.position.x), int(boid.position.y-1), green)
		screen.Set(int(boid.position.x), int(boid.position.y+1), green)
	}
}

func (g *Game) Layout(_, _ int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}
	for i := 0; i < boidCount; i++ {
		createBoid(i)
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Boids in a box")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
