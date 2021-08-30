package runtime

import "math"

type Vec2 struct { // using float64 for ease of use
	x float64
	y float64
}

func (this Vec2) trns(angle, amount float64) Vec2 {
	return this.set(amount, 0).rotate(angle)
}

func (this Vec2) len1() float64 {
	return math.Sqrt(this.x*this.x + this.y*this.y)
}
func (this Vec2) sub(v Vec2) Vec2 {
	this.x -= v.x
	this.y -= v.y
	return this
}

func (this Vec2) angle() float64 {
	angle := FastAtan2(this.x, this.y) * 3.1415927 / 180
	if angle < 0 {
		angle += 360
	}
	return angle
}

func (this Vec2) len2() float64 {
	return this.x*this.x + this.y*this.y
}

func (this Vec2) limit2(l float64) Vec2 {
	tl := this.len2()
	if tl > l {
		return this.scl(math.Sqrt(l / tl))
	}
	return this
}

func (this Vec2) limit(l float64) Vec2 {
	return this.limit2(l * l)
}

func (this Vec2) scl(scalar float64) Vec2 {
	this.x *= scalar
	this.y *= scalar
	return this
}

func (this Vec2) add(v Vec2) Vec2 {
	this.x += v.x
	this.y += v.y
	return this
}

/*func set(this Vec2) Vec2{
	x = v.x;
	y = v.y;
	return this;
}

func (this Vec2) set(Position v) Vec2{
	this.x = v.getX();
	this.y = v.getY();
	return this;
}*/

/**
* Sets the components of this vector
* @param x The x-component
* @param y The y-component
* @return This vector for chaining
 */
func (this Vec2) set(x, y float64) Vec2 {
	this.x = x
	this.y = y
	return this
}
func (this Vec2) rotate(degrees float64) Vec2 {
	return this.rotateRad(degrees * 3.1415927 / 180)
}

/**
* Rotates the Vec2 by the given angle, counter-clockwise assuming the y-axis points up.
* @param radians the angle in radians
 */
func (this Vec2) rotateRad(radians float64) Vec2 {
	cos := math.Cos(radians)
	sin := math.Sin(radians)

	newX := this.x*cos - this.y*sin
	newY := this.x*sin + this.y*cos

	this.x = newX
	this.y = newY

	return this
}

func (this Vec2) isZero() bool {
	return this.x == 0 && this.y == 0
}
func (this Vec2) approachDelta(target Vec2, alpha float64) Vec2 {
	return this.approach(target /*Time.delta * */, alpha)
}

func (this Vec2) approach(target Vec2, alpha float64) Vec2 {
	dx := this.x - target.x
	dy := this.y - target.y
	alpha2 := alpha * alpha
	len2 := Vec2{x: dx, y: dy}.len2()

	if len2 > alpha2 {
		scl := math.Sqrt(alpha2 / len2)
		dx *= scl
		dy *= scl

		return this.sub(Vec2{dx, dy})
	} else {
		return this.set(target.x, target.y)
	}
}

func (this Vec2) setZero() Vec2 {
	this.x = 0
	this.y = 0
	return this
}
