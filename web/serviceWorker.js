const appName = 'recv-app'
const assets = ['/', '/index.html', '/style.css']
self.addEventListener('install', installEvent => {
  installEvent.waitUntil(
    caches.open(appName).then(cache => {
      cache.addAll(assets)
    }),
  )
})

self.addEventListener('fetch', fetchEvent => {
  fetchEvent.respondWith(
    caches.match(fetchEvent.request).then(res => {
      return res || fetch(fetchEvent.request)
    }),
  )
})
