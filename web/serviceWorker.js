const appName = 'recv-app'
const assets = ['/', '/index.html', '/style.css']

self.addEventListener('install', event => {
  event.waitUntil(caches.open(appName).then(cache => cache.addAll(assets)))
})

self.addEventListener('fetch', event => {
  event.respondWith(
    (async () => {
      try {
        const response = await fetch(event.request)
        const cache = await caches.open(appName)
        cache.put(event.request, response.clone())
        return response
      } catch (error) {
        const cache = await caches.open(appName)
        const response = await cache.match(event.request)
        if (response != null) {
          return response
        }
        throw error
      }
    })(),
  )
})
