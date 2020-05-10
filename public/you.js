const InputOpNewMove = 0
const InputOpContinueMove = 1
const InputOpMissile = 2
const InputOpSearch = 3
const InputOpChat = 4

class You extends Human {
    constructor(world, nid, x, y, z, angle, health, status) {
        super(world, nid, x, y, z, angle, health, status)
        this.camera = null
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
                break
            case HumanIdle:
                this.Animation = HumanAnimationIdle
            default:
                this.Animation = HumanAnimationWalk
                break
        }
        this.Status = status
    }
    Walk() {
        let direction = null
        let goal = null

        if (Input.KeyPress("p")) {
            SocketSendSet.set(InputOpSearch, true)
            SocketSendSet.set(InputOpChat, "todo test chat")
        }

        if (Input.KeyDown(" ")) {
            SocketSendSet.set(InputOpMissile, true)
            this.Status = HumanMissile
            this.AnimationMod = 0
            this.AnimationFrame = 0
            this.Animation = HumanAnimationMissile
            PlaySound("baron-missile")
            return
        }

        if (Input.KeyDown("w")) {
            direction = "w"
            goal = this.camera.ry
        }

        if (Input.KeyDown("s")) {
            if (direction === null) {
                direction = "s"
                goal = this.camera.ry + Math.PI
            } else {
                direction = null
                goal = null
            }
        }

        if (Input.KeyDown("a")) {
            if (direction === null) {
                direction = "a"
                goal = this.camera.ry - HalfPi
            } else if (direction === "w") {
                direction = "wa"
                goal -= QuarterPi
            } else if (direction === "s") {
                direction = "sa"
                goal += QuarterPi
            }
        }

        if (Input.KeyDown("d")) {
            if (direction === null)
                goal = this.camera.ry + HalfPi
            else if (direction === "a")
                goal = null
            else if (direction === "wa")
                goal = this.camera.ry
            else if (direction === "sa")
                goal = this.camera.ry + Math.PI
            else if (direction === "w")
                goal += QuarterPi
            else if (direction === "s")
                goal -= QuarterPi
        }

        if (goal === null) {
            this.AnimationMod = 0
            this.AnimationFrame = 0
            this.Animation = HumanAnimationIdle
        } else {
            if (goal < 0)
                goal += Tau
            else if (goal >= Tau)
                goal -= Tau

            if (this.Angle !== goal) {
                this.Angle = goal
                SocketSendSet.set(InputOpNewMove, goal)
            } else {
                SocketSendSet.set(InputOpContinueMove, true)
            }

            // TODO improve
            // this.X += Math.sin(this.Angle) * this.Speed * InverseNetRate
            // this.Z -= Math.cos(this.Angle) * this.Speed * InverseNetRate

            if (this.Animation !== HumanAnimationWalk) {
                this.Animation = HumanAnimationWalk
            } else if (this.UpdateAnimation() === AnimationDone) {
                this.AnimationFrame = 0
            }
        }
    }
    Update() {
        switch (this.Status) {
            case HumanDead:
                this.Dead()
            case HumanMissile:
                this.Missile()
                break
            default:
                this.Walk()
                break
        }
        this.LerpNetCode()
    }
}
