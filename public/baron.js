const BaronAnimationIdle = []
const BaronAnimationWalk = []
const BaronAnimationMelee = []
const BaronAnimationMissile = []
const BaronAnimationDeath = []

const BaronSleep = 0
const BaronDead = 1
const BaronLook = 2
const BaronChase = 3
const BaronMelee = 4
const BaronMissile = 5

class Baron extends Thing {
    constructor(world, nid, x, y, z, direction, heatlh, status) {
        super()
        this.World = world
        this.UID = BaronUID
        this.SID = "baron"
        this.NID = nid
        this.Update = this.BaronUpdate
        this.Animation = this.GetAnimation(status)
        this.X = x
        this.Y = y
        this.Z = z
        this.Angle = DirectionToAngle[direction]
        this.OldX = x
        this.OldY = y
        this.OldZ = z
        this.Radius = 0.4
        this.Height = 1.0
        this.Speed = 0.1
        this.Health = heatlh
        this.Status = status
        world.AddThing(this)
        this.BlockBorders()
        this.AddToBlocks()
    }
    Save() {
        let data = "{u:" + this.UID
        data += ",x:" + this.X
        data += ",y:" + this.Y
        data += ",z:" + this.Z
        data += ",a:" + this.Angle
        data += ",h:" + this.Health
        data += "}"
        return data
    }
    GetAnimation(status) {
        switch (status) {
            case BaronDead:
                return BaronAnimationDeath
            case BaronMelee:
                return BaronAnimationMelee
            case BaronMissile:
                return BaronAnimationMissile
            default:
                return BaronAnimationWalk
        }
    }
    NetUpdateState(status) {
        if (this.Status === status)
            return
        this.AnimationMod = 0
        this.AnimationFrame = 0
        switch (status) {
            case BaronDead:
                this.Animation = BaronAnimationDeath
                break
            case BaronMelee:
                this.Animation = BaronAnimationMelee
                PlaySound("baron-melee")
                break
            case BaronMissile:
                this.Animation = BaronAnimationMissile
                PlaySound("baron-missile")
                break
            case BaronChase:
                if (Math.random() < 0.1)
                    PlaySound("baron-scream")
            default:
                this.Animation = BaronAnimationWalk
                break
        }
        this.Status = status
    }
    NetUpdateHealth(health) {
        if (health < this.Health) {
            if (health < 1) {
                PlaySound("baron-death")
            } else {
                PlaySound("baron-pain")
            }
            for (let i = 0; i < 20; i++) {
                let spriteName = "blood-" + Math.floor(Math.random() * 3)
                let x = this.X + this.Radius * (1 - Math.random() * 2)
                let y = this.Y + this.Height * Math.random()
                let z = this.Z + this.Radius * (1 - Math.random() * 2)
                const spread = 0.2
                let dx = spread * (1 - Math.random() * 2)
                let dy = spread * Math.random()
                let dz = spread * (1 - Math.random() * 2)
                new Blood(this.World, x, y, z, dx, dy, dz, spriteName)
            }

        }
        this.Health = health
    }
    Dead() {
        if (this.AnimationFrame === this.Animation.length - 1) {
            this.Update = this.EmptyUpdate
        } else {
            this.UpdateAnimation()
        }
    }
    Look() {
        if (this.UpdateAnimation() === AnimationDone) {
            this.AnimationFrame = 0
        }
    }
    Melee() {
        if (this.UpdateAnimation() === AnimationDone) {
            this.AnimationFrame = 0
            this.Animation = BaronAnimationWalk
        }
    }
    Missile() {
        if (this.UpdateAnimation() === AnimationDone) {
            this.AnimationFrame = 0
            this.Animation = BaronAnimationWalk
        }
    }
    Chase() {
        if (this.UpdateAnimation() === AnimationDone) {
            this.AnimationFrame = 0
        }
    }
    BaronUpdate() {
        switch (this.Status) {
            case BaronDead:
                this.Dead()
                break
            case BaronLook:
                this.Look()
                break
            case BaronMelee:
                this.Melee()
                break
            case BaronMissile:
                this.Missile()
                break
            case BaronChase:
                this.Chase()
                break
        }
        this.LerpNetCode()
    }
    EmptyUpdate() {}
}
