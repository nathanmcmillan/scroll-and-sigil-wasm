const go = new Go()

WebAssembly.instantiateStreaming(fetch(location.origin + "/ss.wasm"), go.importObject).then(obj => {
    go.run(obj.instance)
})
