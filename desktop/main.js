const {
    app,
    BrowserWindow
} = require("electron")

const path = require("path")

let window

function createWindow() {
    window = new BrowserWindow({
        width: 1200,
        height: 1100
    })

    window.setMenuBarVisibility(false)
    window.webContents.openDevTools()

    window.on("closed", () => {
        window = null
    })

    window.loadFile(path.join(__dirname, "..", "public", "game.html"))
}

app.on("ready", createWindow)

app.on("window-all-closed", () => {
    if (process.platform !== "darwin") {
        app.quit()
    }
})

app.on("activate", () => {
    if (window === null) {
        createWindow()
    }
})
