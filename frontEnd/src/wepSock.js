let ws = null
let listeners = []

function connectWebSocket() {
  if (ws) return
  ws = new WebSocket('ws://localhost:8080/ws')

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      listeners.forEach((callback) => callback(data))
    } catch (error) {
      console.error('message error:', error)
    }
  }
}
function addListener(callback) {
  listeners.push(callback)
}
export { connectWebSocket, addListener }
