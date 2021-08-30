package runtime

import (
	"reflect"
	"runtime/types"
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

func (this Unit) moveAt(vector Vec2) {
	moveAt(this, vector, this.utype.accel)
}

func (this Unit) approach(vector Vec2) {
	vel.approachDelta(vector, this.utype.accel*this.realSpeed())
}

func (this Unit) rotateMove(vec Vec2) {
	moveAt(vec2.trns(this.rotation, vec.len1()))

	if !vec.isZero() {
		this.rotation = moveToward(this.rotation, vec.angle(), this.utype.rotateSpeed /* * math.Max(Time.delta, 1)*/)
	}
}

func (this Unit) aimLook(pos Vec2) {
	this.aim(pos)
	this.lookAt(pos)
}

func (this Unit) aimLook(float x, float y) {
	aim(x, y)
	lookAt(x, y)
}

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
	for _, unit := range Groups.unit {
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

func (this Unit) inRange(other Vec2) bool {
	return this.within(other, this.this.utype.urange)
}

func (this Unit) hasWeapons() bool {
	return this.utype.hasWeapons()
}

func (this Unit) speed() float64 {
	a := this.vel().angle()
	strafePenalty = math.Min((this.isGrounded() || !this.isPlayer()).(float64), 1+(this.utype.strafePenalty-1)*math.Min(math.Mod(a-this.rotation+360), math.Mod(this.rotation-a+360, 360))/180) // Angles.angleDist
	//limit speed to minimum formation speed to preserve formation
	if this.isCommanding() {
		ns := this.minFormationSpeed * 0.98
	} else {
		ns := this.utype.speed
	}
	return ns * this.strafePenalty
}

/** @return speed with boost multipliers factored in. */
func (this Unit) realSpeed() float64 {
	if this.utype.canBoost {
		b := this.utype.boostMultiplier
	} else {
		b := 1
	}
	return (1 + (b-1)*elevation) * this.speed() * this.floorSpeedMultiplier()
}

/** Iterates through this unit and everything it is controlling. */
// BROKEN, TODO fix
/*func (this Unit) eachGroup(cons interface{}){ // originally an arc.func.Cons
	for _, unit in formation {
		cons(unit);
	}
}*/

/** @return where the unit wants to look at. */
func (this Unit) prefRotation() {
	if this.isBuilding() {
		return angleTo(this.buildX, this.buildY)
	} else if this.mineX != nil && this.mineY != nil {
		return angleTo(this.mineX, this.mineY)
	} else if this.moving() && this.utype.omniMovement {
		return this.vel().angle()
	}
	return this.rotation
}

func (this Unit) urange() {
	return this.utype.maxRange
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
	switch(sensor){
		case "totalItems": return this.stack().amount;
		case "itemCapacity": return this.utype.itemCapacity;
		case "rotation": return this.rotation;
		case "health": return this.health;
		case "maxHealth": return this.maxHealth;
		case "ammo": return math.Max(world.rules.unitAmmo.(float64) * this.utype.ammoCapacity, this.ammo);
		case "ammoCapacity": return this.utype.ammoCapacity;
		case "x": return this.x;
		case "y": return this.y;
		case "dead": return this.dead.(float64);
		case "team": return this.team;
		case "shooting": return this.shooting.(float64);
		case "boosting": return (this.utype.canBoost && this.elevation > 0).(float64);
		case "range": this.urange();
		case "shootX": this.shootX;
		case "shootY": this.shootY;
		case "mining": this.mining.(float64);
		case "mineX": math.Max(this.mining.(float64) * this.mineX + 1, 0) - 1;
		case "mineY": math.Max(this.mining.(float64) * this.mineY + 1, 0) - 1;
		case "flag": this.flag;
		case "controlled": return getControllerType()
		case "commanded": 
			if casted, ok := this.controller.(Unit); ok && !this.dead {
				return true
			} else {
				return false // remove once v7 comes out
			}
		case "payloadCount": return len(payloads)
		case "size": return hitSize
		default: return NaN;
	};
}

func (this Unit) senseObject(sensor string) interface{} {
	switch(sensor){
		case "type": return this.utype;
		case "name": return this.getControllerName();
		case "firstItem": return nil; //stack().amount == 0 ? null : item();
		case "controller": return this.getController();
		case "payloadType": return reflect.TypeOf(payloads[len(payloads)-1])
		default: return nil;
	};
}
/*
func (this Unit) sense(Content content){ // TODO items
	if(content == stack().item) return stack().amount;
	return Float.NaN;
}*/

func (this Unit) canDrown() {
	return this.isGrounded() && !this.hovering && this.utype.canDrown
}

func (this Unit) canShoot() {
	//cannot shoot while boosting
	return !this.disarmed && !(this.utype.canBoost && this.elevation > 0)
}

func (this Unit) isCounted() {
	return this.utype.isCounted
}

func (this Unit) bounds() {
	return this.hitSize * 2
}

func (this Unit) controller(next interface{}) {
	if controller != this {
		this.controller = next
	} else {
		this.controller = nil
	}
}

func (this Unit) getController() interface{} {
	if casted, ok := this.controller.(ExecutionContext); ok {
		return casted
	} else if casted, ok := this.controller.(Player); ok {
		return casted
	} else if casted, ok := this.controller.(Unit); ok {
		return casted
	}
	return nil
}

func (this Unit) getControllerType() float64 {
	if this.dead {
		return 0
	}
	if casted, ok := this.controller.(ExecutionContext); ok {
		return 1
	} else if casted, ok := this.controller.(Player); ok {
		return 2
	} else if casted, ok := this.controller.(Unit); ok {
		return 3
	}
	return 0
}

func (this Unit) resetController() {
	this.controller = null
}

func (this Unit) set(def UnitType, controller interface{}) {
	if this.utype != def {
		setType(def)
	}
	this.controller = controller
}

/** @return pathfinder path type for calculating costs */
/*func (this Unit) pathType() int{
	return Pathfinder.costGround;
}*/

func (this Unit) lookAt(float angle) {
	rotation = moveToward(rotation, angle, this.utype.rotateSpeed*Time.delta*speedMultiplier())
}

func (this Unit) lookAt(Position pos) {
	lookAt(angleTo(pos))
}

func (this Unit) lookAt(float x, float y) {
	lookAt(angleTo(x, y))
}

func (this Unit) isAI() boolean {
	if casted, ok := this.controller.(ExecutionContext); ok {
		return true
	} else if casted, ok := this.controller.(Unit); ok {
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
	this.utype = ntype
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

	if world.rules.unitAmmo && this.ammo < this.utype.ammoCapacity-0.0001 {
		resupplyTime += 1 /*Time.delta*/

		//resupply only at a fixed interval to prevent lag
		if resupplyTime > 10 {
			utype.ammoType.resupply(self())
			resupplyTime = 0
		}
	}

	/*if(abilities.size > 0){
		for(Ability a : abilities){
			a.update(self());
		}
	}*/

	this.drag = utype.drag * (math.Min(this.isGrounded()*this.floorOn().dragMultiplier, 1)) * this.dragMultiplier

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
		elevation -= utype.fallSpeed /* * Time.delta*/

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
	x += velX
	y += velY
	velX /= drag
	velY /= drag
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
	switch this.getControllerType() {
	case 1:
		return controller.author //lastAccessed
	case 2:
		return controller.name
	case 3:
		return controller.controller.name
	}
	return nil
}

/*@Override
func (this Unit) display(Table table){
	utype.display(self(), table);
}*/

func (this Unit) isImmune(effect string) bool {
	return utype.immunities.contains(effect)
}

/*
@Override
public void draw(){
	type.draw(self());
}*/

func (this Unit) isPlayer() bool {
	return this.getControllerType == 2
}

func (this Unit) getPlayer() bool {
	if this.isPlayer() {
		return controller
	} else {
		return nil
	}
}

func (this Unit) killed() {
	wasPlayer := getController()
	health := Math.min(health, 0)
	dead := true

	//don't waste time when the unit is already on the ground, just destroy it
	if !this.utype.flying {
		destroy()
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
	if math.Abs(angleDist(angle, to)) < speed {
		return to
	}
	angle = math.Mod(angle, 360)
	to = math.Mod(to, 360)

	if angle > to == (360 - math.Abs(angle-to)) > math.Abs(angle-to) {
		angle -= speed
	} else {
		angle += speed
	}

	return angle
}
