package runtime;

import "math"

type Vec2 struct { // using float64 for ease of use
	x float64
	y float64
}

func trns(float angle, float amount) Vec2{
	return set(amount, 0).rotate(angle);
}

func len1(this Vec2) float64{
	return math.Sqrt(this.x * this.x + this.y * this.y);
}


func len2(v Vec2) float64{
	return this.x * this.x + this.y * this.y;
}

/*func set(this Vec2) Vec2{
	x = v.x;
	y = v.y;
	return this;
}*/

public Vec2 set(Position v){
	x = v.getX();
	y = v.getY();
	return this;
}

/**
	* Sets the components of this vector
	* @param x The x-component
	* @param y The y-component
	* @return This vector for chaining
	*/
func set(x, y float64) Vec2{
	this := new Vec2();
	this.x := x;
	this.y := y;
	return this;
}
func rotate(this Vec2, degrees float64) Vec2{
	return rotateRad(this, degrees * Mathf.degreesToRadians);
}
/**
	* Rotates the Vec2 by the given angle, counter-clockwise assuming the y-axis points up.
	* @param radians the angle in radians
	*/
func rotateRad(this Vec2, float radians) Vec2{
	cos := math.Cos(radians);
	sin := math.Sin(radians);

	newX = this.x * cos - this.y * sin;
	newY = this.x * sin + this.y * cos;

	this.x = newX;
	this.y = newY;

	return this;
}
func (this Vec2) isZero() bool{
	return this.x == 0 && this.y == 0;
}

func (this Vec2) setZero{
	this.x = 0;
	this.y = 0;
	return this;
}
