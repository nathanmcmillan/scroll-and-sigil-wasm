if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register("service.js").then(() => {}).catch((error) => {
        console.log("failed to register service worker", error)
    })
}

class Game {
    constructor() {}
    async init(app) {
        let player = app.world.netLookup.get(app.world.PID)
        app.camera = new Camera(app.world, player, 10.0)
        player.camera = app.camera
    }
}

let game = new Game()
let app = new App(game)

app.run()

function loop() {
    app.loop()
}
