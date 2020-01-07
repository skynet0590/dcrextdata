import { Controller } from 'stimulus'
import axios from 'axios'
import moment from 'moment'
import {
  setAllValues, insertOrUpdateQueryParam, hideAll, showAll, hide, show, getNumberOfPages
} from '../utils'

export default class extends Controller {
  timestamp
  height
  currentPage
  userAgentsPage
  query

  static get targets () {
    return [
      'timestamp', 'height', 'queryInput', 'peerCount', 'userAgents', 'previousUserAgentsButton', 'nextUserAgentsButton',
      'userAgentRowTemplate', 'nextCountriesButton', 'previousCountriesButton', 'countries', 'countryRowTemplate',
      'nextPageButton', 'previousPageButton', 'tableBody', 'rowTemplate', 'totalPageCount', 'currentPage', 'loadingData'
    ]
  }

  async initialize () {
    this.timestamp = parseInt(this.timestampTarget.dataset.initialValue)
    this.height = parseInt(this.heightTarget.dataset.initialValue)
    this.currentPage = parseInt(this.currentPageTarget.dataset.initialValue) || 1
    this.query = this.queryInputTarget.value

    this.userAgentsPage = 1
    await this.loadUserAgents()
    this.displayUserAgents()

    this.countriesPage = 1
    await this.loadCountries()
    this.displayCountries()

    this.loadNetworkPeers()
  }

  loadPreviousPage (e) {
    e.preventDefault()
    this.currentPage = this.currentPage - 1
    this.loadNetworkPeers()
    insertOrUpdateQueryParam('page', this.currentPage)
  }

  loadNextPage (e) {
    e.preventDefault()
    this.currentPage = this.currentPage + 1
    this.loadNetworkPeers()
    insertOrUpdateQueryParam('page', this.currentPage)
  }

  queryLinkClicked (e) {
    e.preventDefault()
    this.query = this.queryInputTarget.value = e.currentTarget.dataset.query
    this.search()
  }

  searchFormSubmitted (e) {
    e.preventDefault()
    this.query = this.queryInputTarget.value
    this.search()
  }

  search () {
    if (this.query === 'Unkown') {
      this.query = ''
    }
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
    if (!nodes) {
      // todo show error message
      return
    }
    nodes.forEach(node => {
      const exRow = document.importNode(_this.rowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      let lastSeen = node.last_seen > 0 ? moment.unix(node.last_seen).fromNow() : 'N/A'
      let connectionTime = node.connection_time > 0 ? moment.unix(node.connection_time).fromNow() : 'N/A'

      fields[0].innerHTML = `<a href="/nodes/view/${node.address}" title="Node status">${node.address}</a><br>
        <span class="text-muted">Connected ${connectionTime}</span> | 
        <span class="text-muted">Seen ${lastSeen}</span><br>`
      fields[1].innerHTML = `${node.user_agent} (${node.protocol_version})<br>
        <span class="text-muted">${node.services}</span>`
      fields[2].innerHTML = `${node.current_height || 'Unknown'}
        <div class="progress"><div class="progress-bar" style="width: ${(100 * node.current_height / _this.height).toFixed(2)}%;"></div></div>`
      fields[3].innerHTML = node.country

      _this.tableBodyTarget.appendChild(exRow)
    })
  }

  loadPreviousUserAgents () {
    this.userAgentsPage -= 1
    this.displayUserAgents()
  }

  loadNextUserAgents () {
    this.userAgentsPage += 1
    this.displayUserAgents()
  }

  async loadUserAgents () {
    const that = this
    const url = `/api/snapshot/${this.timestamp}/user-agents`
    const response = await axios.get(url)
    that.userAgents = response.data.userAgents
  }

  loadPreviousCountries () {
    this.countriesPage -= 1
    this.displayCountries()
  }

  loadNextCountries () {
    this.countriesPage += 1
    this.displayCountries()
  }

  displayUserAgents () {
    if (!this.userAgents) return
    let pageCount = getNumberOfPages(this.userAgents.length, 6)
    if (this.userAgentsPage >= pageCount) {
      hide(this.nextUserAgentsButtonTarget)
    } else {
      show(this.nextUserAgentsButtonTarget)
    }

    if (this.userAgentsPage <= 1) {
      hide(this.previousUserAgentsButtonTarget)
    } else {
      show(this.previousUserAgentsButtonTarget)
    }

    this.userAgentsTarget.innerHTML = ''
    const that = this
    const offset = (this.userAgentsPage - 1) * 6
    const userAgents = this.userAgents.slice(offset, offset + 6)
    userAgents.forEach((userAgent, i) => {
      const exRow = document.importNode(that.userAgentRowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerHTML = 1 + i + offset
      fields[1].innerHTML = `<a data-query="${userAgent.user_agent}" data-action="click->nodes#queryLinkClicked" 
                                  href="#network-snapshot">${userAgent.user_agent}</a>`
      fields[2].innerHTML = `${userAgent.nodes}(${userAgent.percentage}%)`

      that.userAgentsTarget.appendChild(exRow)
    })
  }

  async loadCountries () {
    const that = this
    const url = `/api/snapshot/${this.timestamp}/countries`
    const response = await axios.get(url)
    that.countries = response.data.countries
  }

  displayCountries () {
    if (!this.countries) return
    let pageCount = getNumberOfPages(this.countries.length, 6)
    if (this.countriesPage >= pageCount) {
      hide(this.nextCountriesButtonTarget)
    } else {
      show(this.nextCountriesButtonTarget)
    }

    if (this.countriesPage <= 1) {
      hide(this.previousCountriesButtonTarget)
    } else {
      show(this.previousCountriesButtonTarget)
    }

    this.countriesTarget.innerHTML = ''
    const that = this
    const offset = (this.countriesPage - 1) * 6
    const countries = this.countries.slice(offset, offset + 6)
    countries.forEach((country, i) => {
      const exRow = document.importNode(that.countryRowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerHTML = 1 + i + offset
      fields[1].innerHTML = `<a data-query="${country.country}" data-action="click->nodes#queryLinkClicked" 
                                href="#network-snapshot">${country.country}</a>`
      fields[2].innerHTML = `${country.nodes}(${country.percentage}%)`

      that.countriesTarget.appendChild(exRow)
    })
  }
}
