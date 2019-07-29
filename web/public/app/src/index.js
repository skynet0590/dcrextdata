import { Application } from 'stimulus'
import { definitionsFromContext } from 'stimulus/webpack-helpers'
import ws from './services/messagesocket_service'
import './css/style.scss'

function getSocketURI () {
  let protocol = (window.location.protocol === 'https:') ? 'wss' : 'ws'
  return `${protocol}://${window.location.host}/ws`
}

function createWebSocket () {
  setTimeout(() => {
    // wait a bit to prevent websocket churn from drive by page loads
    let uri = getSocketURI()
    ws.connect(uri)
  }, 1000)
}

createWebSocket()

const application = Application.start()
const context = require.context('./controllers', true, /\.js$/)
application.load(definitionsFromContext(context))
