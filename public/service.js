const cacheName = "dev-1"

self.addEventListener("install", function (event) {
    const cacheList = [
        "/favicon.ico",
        "/sounds/baron-melee.wav",
        "/sounds/baron-missile.wav",
        "/sounds/baron-death.wav",
        "/sounds/baron-pain.wav",
        "/sounds/baron-scream.wav",
        "/sounds/plasma-impact.wav",
    ]
    event.waitUntil(caches.open(cacheName).then(function (cache) {
        console.log("cache", cacheName, "opened")
        return cache.addAll(cacheList)
    }).then(function () {
        return self.skipWaiting()
    }))
})

self.addEventListener("activate", function (event) {
    console.log("cache", cacheName, "activate")
    event.waitUntil(caches.keys().then(function (keyList) {
        return Promise.all(keyList.map(function (key) {
            if (key !== cacheName) {
                console.log(cacheName, "removing old cache", key)
                return caches.delete(key)
            }
        }))
    }))
    return self.clients.claim()
})

self.addEventListener("fetch", function (event) {
    event.respondWith(caches.match(event.request).then((response) => {
        if (event.request.url === "/") {
            return fetch(event.request)
        }
        return response || fetch(event.request)
    }))
})
