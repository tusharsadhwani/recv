const appName = 'recv-app'
const assets = ['/', '/index.html', '/style.css']
self.addEventListener('install', installEvent => {
  installEvent.waitUntil(
    caches.open(appName).then(cache => {
      cache.addAll(assets)
    }),
  )
})

self.addEventListener('fetch', event => {
  event.respondWith(
    fetch(event.request)
      .then(response => {
        return caches.open(appName).then(cache => {
          cache.put(event.request, response.clone())
          return response
        })
      })
      .catch(error => {
        return caches.open(appName).then(cache => {
          return cache.match(event.request).then(response => {
            if (response != null) {
              return response
            }
            throw error
          })
        })
      }),
  )
})
