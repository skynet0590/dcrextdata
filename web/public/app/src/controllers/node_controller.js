import { Controller } from 'stimulus'
import moment from 'moment'

export default class extends Controller {
  timestamp
  height
  currentPage
  userAgentsPage
  query

  static get targets () {
    return [
      'moment'
    ]
  }

  initialize () {
    this.momentTargets.forEach(el => {
      el.textContent = moment.unix(el.dataset.value).fromNow()
    })
  }
}
