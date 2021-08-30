package runtime

import (
	"math"
	"reflect"
)

func NewUnitType() UnitType {
	return UnitType{
		name:                  "unit",
		speed:                 1.1,
		boostMultiplier:       1,
		rotateSpeed:           5,
		baseRotateSpeed:       5,
		drag:                  0.3,
		accel:                 0.5,
		landShake:             0,
		rippleScale:           1,
		riseSpeed:             0.08,
		fallSpeed:             0.018,
		health:                200,
		urange:                -1,
		miningRange:           70,
		armor:                 0,
		maxRange:              -1,
		crashDamageMultiplier: 1,
		targetAir:             true,
		targetGround:          true,
		faceTarget:            true,
		rotateShooting:        true,
		isCounted:             true,
		lowAltitude:           false,
		circleTarget:          false,
		canBoost:              false,
		destructibleWreck:     true,
		hovering:              false,
		omniMovement:          true,
		commandLimit:          8,
		commandRadius:         150,
		//groundLayer: Layer.groundUnit,
		payloadCapacity: 8,
		aimDst:          -1,
		allowLegStep:    false,
		itemCapacity:    -1,
		ammoCapacity:    -1,
		clipSize:        -1,
		canDrown:        true,
		ammoType:        0,
		mineTier:        -1,
		buildSpeed:      -1,
		mineSpeed:       1,
		strafePenalty:   0.5,
		hitSize:         6,
		//abilities: new Seq<>(),
		//weapons: new Seq<>()
	}
}

func NewUnit() Unit {
	return Unit{
		speedMultiplier: 1,
		hitSize:         6,
	}
}

func (this Unit) moveAt(vector Vec2) {
	this.moveAtA(vector, this.utype.accel) // imported from *where?*
}
func (this Unit) moveAtA(vector Vec2, acceleration float64) {
	switch this.elevation > 0 {
	case true:
		t := Vec2{}.set(vector.x, vector.y)                                                                 //target vector
		t.sub(this.vel).limit(acceleration * vector.len1() /* * Time.delta*/ * this.floorSpeedMultiplier()) //delta vector
		this.vel.add(t)
	case false:
		//mark walking state when moving in a controlled manner
		if !vector.isZero() {
			this.walked = true
		}
	}
}
func (this Unit) floorSpeedMultiplier() float64 {
	/*Floor on = isFlying() || hovering ? Blocks.air.asFloor() : floorOn();
	return on.speedMultiplier * speedMultiplier;*/
	return this.speedMultiplier
}

func (this Unit) approach(vector Vec2) {
	this.vel.approachDelta(vector, this.utype.accel*this.realSpeed())
}

func (this Unit) rotateMove(vec Vec2) {
	this.moveAt(Vec2{}.trns(this.rotation, vec.len1()))

	if !vec.isZero() {
		this.rotation = moveToward(this.rotation, vec.angle(), this.utype.rotateSpeed /* * math.Max(Time.delta, 1)*/)
	}
}

func (this Unit) aimLook(pos Vec2) {
	//this.aim(pos)
	this.lookAtPos(pos.x, pos.y)
}

/*func (this Unit) aimLook(x, y float64) {
	this.aim(x, y)
	this.lookAt(x, y)
}*/

/** @return approx. square size of the physical hitbox for physics */
func (this Unit) physicSize() float64 {
	return this.hitSize * 0.7
}

/** @return whether there is solid, un-occupied ground under this unit. */
func (this Unit) canLand() bool {
	if this.onSolid() {
		return false
	}
	xleft := this.x - this.physicSize()/2
	xright := this.x + this.physicSize()/2
	ydown := this.y - this.physicSize()/2
	yup := this.y + this.physicSize()/2
	for _, unit := range Groups.unit { // replace with your own unit attachment
		if unit == this || !unit.isGrounded() {
			continue
		}
		oxleft := unit.x - unit.physicSize()/2
		oxright := unit.x + unit.physicSize()/2
		oydown := unit.y - unit.physicSize()/2
		oyup := unit.y + unit.physicSize()/2
		if oxleft < xleft < oxright || oxleft < xright < oxright &&
			oydown < ydown < oyup || oydown < yup < oyup {
			return false
		}
	}
	return true
}

func (this Unit) onSolid() bool {
	return false
}

func (this Unit) inRange(other Vec2) bool {
	return other.len1() < this.utype.urange
	//return this.within(other, this.utype.urange) // declared where?
}

func (this Unit) hasWeapons() bool {
	return len(this.utype.weapons) > 0
}
func (this Unit) isGrounded() bool {
	return this.elevation < 0.001
}
func (this Unit) speed() float64 {
	a := this.vel.angle()
	b := 0.0
	if this.isGrounded() || !this.isPlayer() {
		b = 1
	}
	this.strafePenalty = math.Min(b, 1+(this.utype.strafePenalty-1)*math.Min(math.Mod(a-this.rotation+360, 360), math.Mod(this.rotation-a+360, 360))/180) // Angles.angleDist
	//limit speed to minimum formation speed to preserve formation
	ns := this.utype.speed
	if this.isCommanding() {
		ns = this.minFormationSpeed * 0.98
	}
	return ns * this.strafePenalty
}

func (this Unit) isCommanding() bool {
	for _, unit := range this.formation {
		if !unit.dead {
			return true
		}
	}
	return false
}

/** @return speed with boost multipliers factored in. */
func (this Unit) realSpeed() float64 {
	b := 1.0
	if this.utype.canBoost {
		b = this.utype.boostMultiplier
	}
	return (1 + (b-1)*this.elevation) * this.speed() * this.floorSpeedMultiplier()
}

/** Iterates through this unit and everything it is controlling. */
// BROKEN, TODO fix
/*func (this Unit) eachGroup(cons interface{}){ // originally an arc.func.Cons
	for _, unit := range formation {
		cons(unit);
	}
}*/

func (this Unit) angleTo(x, y float64) float64 {
	return Vec2{x - this.x, y - this.y}.angle()
}

/** @return where the unit wants to look at. */
func (this Unit) prefRotation() float64 {
	if this.isBuilding() {
		return this.angleTo(float64(this.buildX), float64(this.buildY))
	} else if this.mining != false {
		return this.angleTo(float64(this.mineX), float64(this.mineY))
	} else if this.vel.len1() > 0 && this.utype.omniMovement {
		return this.vel.angle()
	}
	return this.rotation
}

func (this Unit) urange() float64 {
	return this.utype.maxRange
}

func (this Unit) isBuilding() bool {
	return false // no ucontrol build yet
}

/*
func (this Unit) clipSize(){ // ???
	if(this.isBuilding()){
		return world.rules.infiniteResources ? Float.MAX_VALUE : math.Max(this.type.clipSize, this.type.region.width) + this.buildingRange + tilesize*4f; // tilesize of what?
	}
	if(this.mining()){
		return this.type.clipSize + this.type.miningRange;
	}
	return this.type.clipSize;
}*/

func (this Unit) sense(sensor string) float64 {
	switch sensor {
	case "totalItems":
		return float64(this.stack.amount)
	case "itemCapacity":
		return float64(this.utype.itemCapacity)
	case "rotation":
		return this.rotation
	case "health":
		return this.health
	case "maxHealth":
		return this.maxHealth
	case "ammo":
		if world.rules.unitAmmo {
			return this.ammo
		}
		return float64(this.utype.ammoCapacity)
	case "ammoCapacity":
		return float64(this.utype.ammoCapacity)
	case "x":
		return this.x
	case "y":
		return this.y
	case "dead":
		if this.dead {
			return 1.0
		}
		return 0.0
	case "team":
		return float64(this.team) // Team struct?
	case "shooting":
		if this.shooting {
			return 1.0
		}
		return 0.0
	case "boosting":
		if this.utype.canBoost && this.elevation > 0 {
			return 1.0
		}
		return 0.0
	case "range":
		return this.urange()
	case "shootX":
		return this.shootX
	case "shootY":
		return this.shootY
	case "mining":
		if this.mining {
			return 1.0
		}
		return 0.0
	case "mineX":
		if this.mining {
			return float64(this.mineX)
		}
		return -1
	case "mineY":
		if this.mining {
			return float64(this.mineY)
		}
		return -1
	case "flag":
		return this.flag
	case "controlled":
		return float64(this.getControllerType())
	case "commanded":
		if this.controller.ctype == 3 && !this.dead {
			return 1.0
		}
		return 0.0
	case "payloadCount":
		return float64(len(this.payloads))
	case "size":
		return this.hitSize
	default:
		return math.NaN()
	}
}

func (this Unit) senseObject(sensor string) interface{} {
	switch sensor {
	case "type":
		return this.utype
	case "name":
		return this.getControllerName()
	case "firstItem":
		return nil //stack().amount == 0 ? null : item();
	case "controller":
		return this.getController()
	case "payloadType":
		return reflect.TypeOf(this.payloads[len(this.payloads)-1])
	default:
		return nil
	}
}

/*
func (this Unit) sense(Content content){ // TODO items
	if(content == stack().item) return stack().amount;
	return Float.NaN;
}*/

func (this Unit) canDrown() bool {
	return this.isGrounded() && !this.hovering && this.utype.canDrown
}

func (this Unit) canShoot() bool {
	//cannot shoot while boosting
	return !this.disarmed && !(this.utype.canBoost && this.elevation > 0)
}

func (this Unit) isCounted() bool {
	return this.utype.isCounted
}

func (this Unit) bounds() float64 {
	return this.hitSize * 2
}

func (this Unit) getController() interface{} {
	return this.controller.parent
}

func (this Unit) getControllerType() int {
	if this.dead {
		return 0
	}
	return this.controller.ctype
	/*
		if casted, ok := this.controller.(ExecutionContext); ok {
			casted = casted
			return 1
		} else if casted, ok := this.controller.(Player); ok {
			casted = casted
			return 2
		} else if casted, ok := this.controller.(Unit); ok {
			casted = casted
			return 3
		}
		return 0*/
}

func (this Unit) resetController() {
	this.controller = nil
}

func (this Unit) set(def UnitType, controller UnitController) {
	if this.utype != &def {
		this.setType(def)
	}
	this.controller = &controller
}

/** @return pathfinder path type for calculating costs */
/*func (this Unit) pathType() int{
	return Pathfinder.costGround;
}*/

func (this Unit) lookAt(angle float64) {
	this.rotation = moveToward(this.rotation, angle, this.utype.rotateSpeed* /*Time.delta**/ this.speedMultiplier)
}

func (this Unit) lookAtPos(x, y float64) {
	this.lookAt(this.angleTo(x, y))
}

func (this Unit) isAI() bool {
	if this.controller.ctype == 1 || this.controller.ctype == 3 {
		return true
	}
	return false
}

/*
public int count(){
	return team.data().countType(type);
}

public int cap(){
	return Units.getCap(team);
}*/

func (this Unit) setType(ntype UnitType) {
	this.utype = &ntype
	this.maxHealth = ntype.health
	this.drag = ntype.drag
	this.armor = ntype.armor
	this.hitSize = ntype.hitSize
	this.hovering = ntype.hovering
}

//if(mounts().length != type.weapons.size) setupWeapons(type);
/*if(len(this.abilities) != ntype.abilities.size){
		abilities = ntype.abilities.map(Ability::copy);
	}
}

@Override
func (this Unit) afterSync(){
	//set up type info after reading
	setType(this.type);
	controller.unit(self());
}

@Override
func (this Unit) afterRead(){
	afterSync();
	//reset controller state
	controller(type.createController());
}

@Override
func (this Unit) add(){
	team.data().updateCount(type, 1);

	//check if over unit cap
	if(count() > cap() && !spawnedByCore && !dead && !state.rules.editor){
		Call.unitCapDeath(self());
		team.data().updateCount(type, -1);
	}

}

@Override
func (this Unit) remove(){
	team.data().updateCount(type, -1);
	controller.removed(self());
}

func (this Unit) landed(){
	if utype.landShake > 0 {
		Effect.shake(type.landShake, type.landShake, this);
	}

	type.landed(self());
}*/

/*func (this Unit) heal(amount float64){
	if this.health < this.maxHealth && amount > 0 {
		this.wasHealed = true;
	}
}*/

func (this Unit) advanceTick() {
	//this.utype.update(self());

	/*if wasHealed && healTime <= -1f {
		healTime = 1f;
	}
	healTime -= Time.delta / 20f;
	this.wasHealed = false;

	//check if environment is unsupported
	if(!type.supportsEnv(state.rules.environment) && !dead){
		Call.unitCapDeath(self());
		team.data().updateCount(type, -1);
	}*/

	//
	//if world.rules.unitAmmo && this.ammo < (float64(this.utype.ammoCapacity)-0.0001) {
	//	this.resupplyTime += 1 /*Time.delta*/
	//
	//	//resupply only at a fixed interval to prevent lag
	//	if this.resupplyTime > 10 {
	//		this.utype.ammoType.resupply(self())
	//		this.resupplyTime = 0
	//	}
	//}

	/*if(abilities.size > 0){
		for(Ability a : abilities){
			a.update(self());
		}
	}*/

	g := 0.0
	if this.isGrounded() {
		g = 1
	}
	this.drag = this.utype.drag * (math.Min(g /**this.floorOn().dragMultiplier*/, 1)) * this.dragMultiplier

	//apply knockback based on spawns
	/*
		if(team != world.rules.waveTeam && world.hasSpawns() && (!net.client() || isLocal())){
			float relativeSize = state.rules.dropZoneRadius + hitSize/2f + 1f;
			for(Tile spawn : spawner.getSpawns()){
				if(within(spawn.worldx(), spawn.worldy(), relativeSize)){
					velAddNet(Tmp.v1.set(this).sub(spawn.worldx(), spawn.worldy()).setLength(0.1f + 1f - dst(spawn) / relativeSize).scl(0.45f * Time.delta));
				}
			}
		}*/

	//simulate falling down
	if this.dead || this.health <= 0 {
		//less drag when dead
		this.drag = 0.01

		/*
			//standard fall smoke
			if(Mathf.chanceDelta(0.1)){
				Tmp.v1.rnd(Mathf.range(hitSize));
				type.fallEffect.at(x + Tmp.v1.x, y + Tmp.v1.y);
			}

			//thruster fall trail
			if(Mathf.chanceDelta(0.2)){
				float offset = type.engineOffset/2f + type.engineOffset/2f * elevation;
				float range = Mathf.range(type.engineSize);
				type.fallThrusterEffect.at(
					x + Angles.trnsx(rotation + 180, offset) + Mathf.range(range),
					y + Angles.trnsy(rotation + 180, offset) + Mathf.range(range),
					Mathf.random()
				);
			}*/

		//move down
		this.elevation -= this.utype.fallSpeed /* * Time.delta*/

		/*if(this.isGrounded() || this.health <= -this.maxHealth){
			despawn();
			//Call.unitDestroy(id);
		}*/
	}
	/*
		Tile tile = tileOn();
		Floor floor = floorOn();

		if(tile != null && isGrounded() && !type.hovering){
			//unit block update
			if(tile.build != null){
				tile.build.unitOn(self());
			}

			//apply damage
			if(floor.damageTaken > 0f){
				damageContinuous(floor.damageTaken);
			}
		}

		//kill entities on tiles that are solid to them
		if(tile != null && !canPassOn()){
			//boost if possible
			if(type.canBoost){
				elevation = 1f;
			}else if(!net.client()){
				kill();
			}
		}

		//AI only updates on the server
		if(!net.client() && !dead){
			controller.updateUnit();
		}

		//clear controller when it becomes invalid
		if(!controller.isValidController()){
			resetController();
		}

		//remove units spawned by the core
		if(spawnedByCore && !isPlayer() && !dead){
			despawn(this);
		}*/
	// do physics
	this.x += this.vel.x
	this.y += this.vel.y
	this.vel.x /= (1 + this.drag)
	this.vel.y /= (1 + this.drag)
}

/** @return a preview icon for this unit. */
/*
public TextureRegion icon(){
	return type.fullIcon;
}*/

/** Actually destroys the unit, removing it and creating explosions. **/
/*
public void destroy(){
	if(!isAdded()) return;

	float explosiveness = 2f + item().explosiveness * stack().amount * 1.53f;
	float flammability = item().flammability * stack().amount / 1.9f;
	float power = item().charge * stack().amount * 150f;

	if(!spawnedByCore){
		Damage.dynamicExplosion(x, y, flammability, explosiveness, power, bounds() / 2f, state.rules.damageExplosions, item().flammability > 1, team, type.deathExplosionEffect);
	}

	float shake = hitSize / 3f;

	Effect.scorch(x, y, (int)(hitSize / 5));
	Fx.explosion.at(this);
	Effect.shake(shake, shake, this);
	type.deathSound.at(this);

	Events.fire(new UnitDestroyEvent(self()));

	if(explosiveness > 7f && (isLocal() || wasPlayer)){
		Events.fire(Trigger.suicideBomb);
	}

	//if this unit crash landed (was flying), damage stuff in a radius
	if(type.flying && !spawnedByCore){
		Damage.damage(team,x, y, Mathf.pow(hitSize, 0.94f) * 1.25f, Mathf.pow(hitSize, 0.75f) * type.crashDamageMultiplier * 5f, true, false, true);
	}

	if(!headless){
		for(int i = 0; i < type.wreckRegions.length; i++){
			if(type.wreckRegions[i].found()){
				float range = type.hitSize /4f;
				Tmp.v1.rnd(range);
				Effect.decal(type.wreckRegions[i], x + Tmp.v1.x, y + Tmp.v1.y, rotation - 90);
			}
		}
	}

	remove();
}*/

/** @return name of direct or indirect player controller. */
func (this Unit) getControllerName() string {
	return this.controller.name
}

/*@Override
func (this Unit) display(Table table){
	utype.display(self(), table);
}*/

func (this Unit) isImmune(effect string) bool {
	return false
	//return this.utype.immunities.contains(effect)
}

/*
@Override
public void draw(){
	type.draw(self());
}*/

func (this Unit) isPlayer() bool {
	return this.controller.ctype == 2
}

func (this Unit) getPlayer() interface{} {
	if this.isPlayer() {
		return this.controller.parent
	} else {
		return nil
	}
}

func (this Unit) killed() {
	//this.wasPlayer = this.getController()
	this.health = math.Min(this.health, 0)
	this.dead = true

	//don't waste time when the unit is already on the ground, just destroy it
	if !this.utype.flying {
		//despawn()
	}
}

func (this Unit) kill() {
	if this.dead {
		return
	}
	//deaths are synced; this calls killed()
	//Call.unitDeath(id);
	this.killed()
}

func moveToward(angle, to, speed float64) float64 {
	if math.Abs(math.Min(math.Mod(angle-to+360, 360), math.Mod(to-angle+360, 360))) < speed {
		return to
	}
	angle = math.Mod(angle, 360)
	to = math.Mod(to, 360)

	if (angle > to) == ((360 - math.Abs(angle-to)) > math.Abs(angle-to)) {
		angle -= speed
	} else {
		angle += speed
	}

	return angle
}
