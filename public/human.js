const HumanAnimationIdle = []
const HumanAnimationWalk = []
const HumanAnimationMelee = []
const HumanAnimationMissile = []
const HumanAnimationDeath = []

const HumanDead = 0
const HumanIdle = 1
const HumanWalk = 2
const HumanMelee = 3
const HumanMissile = 4

class Human extends Thing {
    constructor(world, nid, x, y, z, angle, health, status) {
        super()
        this.World = world
        this.UID = HumanUID
        this.SID = "baron"
        this.NID = nid
        this.Animation = HumanAnimationWalk
        this.X = x
        this.Y = y
        this.Z = z
        this.Angle = angle
        this.OldX = x
        this.OldY = y
        this.OldZ = z
        this.Radius = 0.4
        this.Height = 1.0
        this.Speed = 0.1
        this.Health = health
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
    NetUpdateState(status) {
        if (this.Status === status)
            return
        this.AnimationMod = 0
        this.AnimationFrame = 0
        switch (status) {
            case HumanDead:
                this.Animation = HumanAnimationDeath
                break
            case HumanMissile:
                this.Animation = HumanAnimationMissile
                PlaySound("baron-missile")
                break
            case HumanIdle:
                this.Animation = HumanAnimationIdle
            default:
                this.Animation = HumanAnimationWalk
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
    Missile() {
        if (this.UpdateAnimation() === AnimationDone) {
            this.AnimationFrame = 0
            this.Animation = HumanAnimationIdle
            this.Status = HumanIdle
        }
    }
    Walk() {
        if (this.UpdateAnimation() === AnimationDone)
            this.AnimationFrame = 0
    }
    Update() {
        switch (this.Status) {
            case HumanDead:
                this.Dead()
            case HumanMissile:
                this.Missile()
                break
            case HumanIdle:
                break
            default:
                this.Walk()
                break
        }
        this.LerpNetCode()
    }
    EmptyUpdate() {}
}
