package model

import "math"

type Vec3 struct {
	X, Y, Z float64
}

type Vec2 struct {
	X, Y float64
}

func PointInTriangle(p, a, b, c Vec2) bool {
	// Use area-based method to determine if point is in triangle
	// This method works correctly regardless of coordinate system orientation
	area := 0.5 * ((b.X-a.X)*(c.Y-a.Y) - (c.X-a.X)*(b.Y-a.Y))

	// Calculate barycentric coordinates
	s := (a.Y*c.X - a.X*c.Y + (c.Y-a.Y)*p.X + (a.X-c.X)*p.Y) / (2 * area)
	t := (a.X*b.Y - a.Y*b.X + (a.Y-b.Y)*p.X + (b.X-a.X)*p.Y) / (2 * area)
	u := 1 - s - t

	// Point is in triangle if barycentric coordinates are all positive (or zero)
	return s >= 0 && t >= 0 && u >= 0
}

func EmptyVec3() *Vec3 {
	return &Vec3{0, 0, 0}
}

func EmptyVec2() *Vec2 {
	return &Vec2{0, 0}
}

func NewVec3(x, y, z float64) *Vec3 {
	return &Vec3{x, y, z}
}

func NewVec2(x, y float64) *Vec2 {
	return &Vec2{x, y}
}

func rotateX(p Vec3, angle float64) Vec3 {
	s, c := math.Sin(angle), math.Cos(angle)
	return Vec3{
		X: p.X,
		Y: c*p.Y - s*p.Z,
		Z: s*p.Y + c*p.Z,
	}
}

func rotateY(p Vec3, angle float64) Vec3 {
	s, c := math.Sin(angle), math.Cos(angle)
	return Vec3{
		X: c*p.X + s*p.Z,
		Y: p.Y,
		Z: -s*p.X + c*p.Z,
	}
}

func rotateZ(p Vec3, angle float64) Vec3 {
	s, c := math.Sin(angle), math.Cos(angle)
	return Vec3{
		X: c*p.X - s*p.Y,
		Y: s*p.X + c*p.Y,
		Z: p.Z,
	}
}

func RotateEuler(p, origin, rotation *Vec3) *Vec3 {
	// Step 1: Translate to origin
	pTranslated := Vec3{
		X: p.X - origin.X,
		Y: p.Y - origin.Y,
		Z: p.Z - origin.Z,
	}

	// Step 2: Apply rotations in X-Y-Z order (or choose your order)
	pRot := rotateX(pTranslated, rotation.X)
	pRot = rotateY(pRot, rotation.Y)
	pRot = rotateZ(pRot, rotation.Z)

	// Step 3: Translate back
	return &Vec3{
		X: pRot.X + origin.X,
		Y: pRot.Y + origin.Y,
		Z: pRot.Z + origin.Z,
	}
}

func (v *Vec3) Sub(other *Vec3) *Vec3 {
	return &Vec3{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

func (v *Vec3) Cross(other *Vec3) *Vec3 {
	return &Vec3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

func (v *Vec3) Normalize() *Vec3 {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return &Vec3{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

func (v *Vec3) Dot(other *Vec3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}
