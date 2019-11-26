import { Controller } from 'stimulus'
import axios from 'axios'
import moment from 'moment'
import {
  setAllValues, insertOrUpdateQueryParam, hideAll, showAll
} from '../utils'

export default class extends Controller {
  timestamp
  height
  currentPage
  query

  static get targets () {
    return [
      'timestamp', 'height', 'queryInput', 'peerCount',
      'nextPageButton', 'previousPageButton', 'tableBody', 'rowTemplate', 'totalPageCount', 'currentPage', 'loadingData'
    ]
  }

  initialize () {
    this.timestamp = parseInt(this.timestampTarget.dataset.initialValue)
    this.height = parseInt(this.heightTarget.dataset.initialValue)
    this.currentPage = parseInt(this.currentPageTarget.dataset.initialValue) || 1
    this.query = this.queryInputTarget.value
    this.loadNetworkPeers()
  }

  loadPreviousPage (e) {
    e.preventDefault()
    this.currentPage = this.currentPage - 1
    this.loadNetworkPeers()
    insertOrUpdateQueryParam('page', this.currentPage)
  }

  laodNextPage (e) {
    e.preventDefault()
    this.currentPage = this.currentPage + 1
    this.loadNetworkPeers()
    insertOrUpdateQueryParam('page', this.currentPage)
  }

  queryLinkClicked (e) {
    e.preventDefault()
    this.query = e.currentTarget.dataset.query
    this.search()
  }

  searchFormSubmitted (e) {
    e.preventDefault()
    this.query = this.queryInputTarget.value
    this.search()
  }

  search () {
    insertOrUpdateQueryParam('q', this.query)
    this.currentPage = 1
    insertOrUpdateQueryParam('page', 1)
    this.loadNetworkPeers()
  }

  loadNetworkPeers () {
    const _this = this
    const url = `/api/snapshot/${this.timestamp}/nodes?q=${this.query}&page=${this.currentPage}`
    axios.get(url).then(response => {
      let result = response.data
      if (_this.currentPage <= 1) {
        hideAll(_this.previousPageButtonTargets)
      } else {
        showAll(_this.previousPageButtonTargets)
      }

      if (_this.currentPage >= result.pageCount) {
        hideAll(_this.nextPageButtonTargets)
      } else {
        showAll(_this.nextPageButtonTargets)
      }

      setAllValues(_this.currentPageTargets, result.page)
      setAllValues(_this.totalPageCountTargets, result.pageCount)
      setAllValues(_this.peerCountTargets, result.peerCount)

      _this.displayNodes(result.nodes)
    })
  }

  displayNodes (nodes) {
    const _this = this
    this.tableBodyTarget.innerHTML = ''

    nodes.forEach(node => {
      const exRow = document.importNode(_this.rowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerHTML = `<a href="/nodes/view/${node.address}" title="Node status">${node.address}</a><br>
        <span class="text-muted">Since ${moment.unix(node.last_seen).fromNow()}</span><br>`
      fields[1].innerHTML = `${node.user_agent} (${node.protocol_version})<br>
        <span class="text-muted">services</span>`
      fields[2].innerHTML = `${node.current_height}
        <div class="progress"><div class="progress-bar" style="width: ${(100 * node.current_height / _this.height).toFixed(2)}%;"></div></div>`
      fields[3].innerHTML = node.country

      _this.tableBodyTarget.appendChild(exRow)
    })
  }
}
