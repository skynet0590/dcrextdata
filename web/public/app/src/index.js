import { Application } from 'stimulus'
import { definitionsFromContext } from 'stimulus/webpack-helpers'
import './css/style.scss'
import '../node_modules/bootstrap4-toggle/css/bootstrap4-toggle.css'
import '../node_modules/bootstrap4-toggle/js/bootstrap4-toggle.js'
import ws from './services/messagesocket_service'

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
