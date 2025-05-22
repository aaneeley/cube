package geometry

import (
	"fmt"
	"math"

	"github.com/aaneeley/cube/model"
	win "github.com/aaneeley/cube/window"
)

var lightDir = model.NewVec3(1, -1, 1).Normalize()

const (
	ambientIntensity = 0.3
)

var faceVertIndices = [6][4]int{
	{0, 1, 2, 3}, // front
	{4, 5, 6, 7}, // back
	{4, 7, 1, 0}, // left
	{3, 2, 6, 5}, // right
	{4, 0, 3, 5}, // top
	{1, 7, 6, 2}, // bottom
}

var texture = []string{"C", "U", "B", "E"}

type Cube struct {
	sideLength float64
	origin     *model.Vec3
	rotation   *model.Vec3
}

type ProjectedCube struct {
	Points []model.Vec2
}

func NewCube(sideLength float64) *Cube {
	return &Cube{
		sideLength: sideLength,
		origin:     model.EmptyVec3(),
		rotation:   model.EmptyVec3(),
	}
}

func (c *Cube) SetOrigin(origin *model.Vec3) {
	c.origin = origin
}

func (c *Cube) SetRotation(rotation *model.Vec3) {
	c.rotation = rotation
}

func (c *Cube) getPoints() [8]*model.Vec3 {
	d := c.sideLength / 2
	o := c.origin

	points := [8]*model.Vec3{
		{X: o.X - d, Y: o.Y - d, Z: o.Z - d},
		{X: o.X - d, Y: o.Y + d, Z: o.Z - d},
		{X: o.X + d, Y: o.Y + d, Z: o.Z - d},
		{X: o.X + d, Y: o.Y - d, Z: o.Z - d},
		{X: o.X - d, Y: o.Y - d, Z: o.Z + d},
		{X: o.X + d, Y: o.Y - d, Z: o.Z + d},
		{X: o.X + d, Y: o.Y + d, Z: o.Z + d},
		{X: o.X - d, Y: o.Y + d, Z: o.Z + d},
	}

	for i, p := range points {
		points[i] = model.RotateEuler(p, c.origin, c.rotation)
	}

	return points
}

func (c *Cube) visibleFaces() [6]bool {
	points := c.getPoints()

	var result [6]bool
	for i, faceVertInd := range faceVertIndices {
		edge1 := points[faceVertInd[1]].Sub(points[faceVertInd[0]])
		edge2 := points[faceVertInd[2]].Sub(points[faceVertInd[0]])
		normal := edge1.Cross(edge2)
		view := model.NewVec3(0, 0, -1)
		if normal.Dot(view) < 0 {
			result[i] = true
		}
	}
	return result
}

func (c *Cube) getBrightness(faceIndex int) float64 {
	points := c.getPoints()
	faceVertInd := faceVertIndices[faceIndex]
	edge1 := points[faceVertInd[1]].Sub(points[faceVertInd[0]])
	edge2 := points[faceVertInd[2]].Sub(points[faceVertInd[0]])
	normal := edge1.Cross(edge2).Normalize()

	dot := normal.Dot(lightDir)
	diffuse := math.Max(dot, 0)
	brightness := ambientIntensity + diffuse*(1.0-ambientIntensity)

	return math.Min(1.0, math.Max(0.0, brightness))
}

func brightnessToChar(brightness float64, textureIndex int) string {
	targetChar := texture[textureIndex%len(texture)]
	r := int(brightness*255) - 50
	g := int(brightness*255) - 10
	b := int(brightness * 255)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, targetChar)
}

func (c *Cube) DrawToBuf(window *win.Window, yScaleOrigin float64) {
	// Project
	var projPts []*model.Vec2
	for _, p := range c.getPoints() {
		projPts = append(projPts, model.NewVec2(p.X, p.Y))
	}

	minX := math.MaxInt
	maxX := math.MinInt
	minY := math.MaxInt
	maxY := math.MinInt

	// Scale for pixel height and get bounds.
	for _, pp := range projPts {
		pp.Y = yScaleOrigin + (pp.Y-yScaleOrigin)*0.45
		if thisX := int(math.Floor(pp.X)); thisX < minX {
			minX = thisX
		}
		if thisX := int(math.Ceil(pp.X)); thisX > maxX {
			maxX = thisX
		}
		if thisY := int(math.Floor(pp.Y)); thisY < minY {
			minY = thisY
		}
		if thisY := int(math.Ceil(pp.Y)); thisY > maxY {
			maxY = thisY
		}
	}

	visibleFaces := c.visibleFaces()

	// Iterate through all pixels in bounds.
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			window.Buffer.WriteString(model.PosCode(x, y))
			p := model.Vec2{X: float64(x), Y: float64(y)}

			pointLandsOnCube := false
			for i, face := range faceVertIndices {
				if !visibleFaces[i] {
					continue
				}
				if model.PointInTriangle(p, *projPts[face[0]], *projPts[face[2]], *projPts[face[3]]) ||
					model.PointInTriangle(p, *projPts[face[0]], *projPts[face[1]], *projPts[face[2]]) {
					brightness := c.getBrightness(i)
					window.Buffer.WriteString(brightnessToChar(brightness, x))
					pointLandsOnCube = true
				}
			}
			if !pointLandsOnCube {
				window.Buffer.WriteString(" ")
			}
		}
	}
}
