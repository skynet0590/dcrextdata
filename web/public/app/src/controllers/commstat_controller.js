import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, setActiveOptionBtn, showLoading, hideLoading, formatDate } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')

export default class extends Controller {
  viewOption
  platform
  subreddit
  twitterHandle
  repository
  dataType

  static get targets () {
    return [
      'pageSizeWrapper', 'previousPageButton', 'totalPageCount', 'nextPageButton',
      'currentPage', 'numPageWrapper', 'selectedNum', 'messageView',
      'viewOptionControl', 'viewOption',
      'chartWrapper', 'chartsView', 'labels', 'tableWrapper', 'loadingData', 'messageView',
      'tableWrapper', 'table', 'rowTemplate', 'tableCol1', 'tableCol2', 'tableCol3',
      'platform', 'subreddit', 'subAccountWrapper', 'dataTypeWrapper', 'dataType', 'twitterHandle', 'repository'
    ]
  }

  initialize () {
    this.currentPage = parseInt(this.currentPageTarget.dataset.currentPage)
    if (this.currentPage < 1) {
      this.currentPage = 1
    }

    this.platform = this.platformTarget.dataset.initialValue
    if (this.platform === '') {
      this.platform = this.platformTarget.value = this.platformTarget.options[0].innerText
    }

    this.showCurrentSubAccountWrapper()

    this.subreddit = this.subredditTarget.dataset.initialValue
    if (this.subreddit === '') {
      this.subreddit = this.subredditTarget.value = this.subredditTarget.options[0].innerText
    }

    this.twitterHandle = this.twitterHandleTarget.dataset.initialValue
    if (this.twitterHandle === '') {
      this.twitterHandle = this.twitterHandleTarget.value = this.twitterHandleTarget.options[0].innerText
    }

    this.repository = this.repositoryTarget.dataset.initialValue
    if (this.repository === '') {
      this.repository = this.repositoryTarget.value = this.repositoryTarget.options[0].innerText
    }

    this.dataType = this.dataTypeTarget.dataset.initialValue

    this.viewOption = this.viewOptionControlTarget.dataset.initialValue
    if (this.viewOption === 'chart') {
      this.setChart()
    } else {
      this.setTable()
    }
  }

  setTable () {
    this.viewOption = 'table'
    setActiveOptionBtn(this.viewOption, this.viewOptionTargets)
    hide(this.chartWrapperTarget)
    hide(this.messageViewTarget)
    show(this.tableWrapperTarget)
    show(this.numPageWrapperTarget)
    show(this.pageSizeWrapperTarget)
    this.updateDataTypeControl()
    this.nextPage = this.currentPage
    this.fetchData()
  }

  setChart () {
    this.viewOption = 'chart'
    setActiveOptionBtn(this.viewOption, this.viewOptionTargets)
    hide(this.tableWrapperTarget)
    hide(this.messageViewTarget)
    show(this.chartWrapperTarget)
    hide(this.pageSizeWrapperTarget)
    hide(this.numPageWrapperTarget)
    this.updateDataTypeControl()
    this.fetchDataAndPlotGraph()
  }

  platformChanged (event) {
    this.platform = event.currentTarget.value
    this.showCurrentSubAccountWrapper()

    this.updateDataTypeControl()
    this.currentPage = 1
    if (this.viewOption === 'table') {
      this.fetchData()
    } else {
      this.fetchDataAndPlotGraph()
    }
  }

  subredditChanged (event) {
    this.subreddit = event.currentTarget.value
    this.currentPage = 1
    if (this.viewOption === 'table') {
      this.fetchData()
    } else {
      this.fetchDataAndPlotGraph()
    }
  }

  twitterHandleChanged (event) {
    this.twitterHandle = event.currentTarget.value
    this.currentPage = 1
    if (this.viewOption === 'table') {
      this.fetchData()
    } else {
      this.fetchDataAndPlotGraph()
    }
  }

  repositoryChanged (event) {
    this.repository = event.currentTarget.value
    this.currentPage = 1
    if (this.viewOption === 'table') {
      this.fetchData()
    } else {
      this.fetchDataAndPlotGraph()
    }
  }

  dataTypeChanged (event) {
    this.dataType = event.currentTarget.value
    this.fetchDataAndPlotGraph()
  }

  showCurrentSubAccountWrapper () {
    const that = this
    this.subAccountWrapperTargets.forEach(el => {
      if (el.dataset.platform === that.platform) {
        show(el)
      } else {
        hide(el)
      }
    })
  }

  updateDataTypeControl () {
    this.dataTypeTarget.innerHTML = ''
    hide(this.dataTypeWrapperTarget)
    if (this.viewOption !== 'chart') {
      return
    }

    const _this = this
    const addDataTypeOption = function (value, label) {
      let selected = _this.dataType === value ? 'selected' : ''
      _this.dataTypeTarget.innerHTML += `<option ${selected} value="${value}">${label}</option>`
    }
    switch (this.platform) {
      case 'Reddit':
        if (this.dataType !== 'subscribers' && this.dataType !== 'active_accounts') {
          this.dataType = 'subscribers'
        }
        addDataTypeOption('subscribers', 'Subscribers')
        addDataTypeOption('active_accounts', 'Active Accounts')
        show(_this.dataTypeWrapperTarget)
        break
      case 'Github':
        if (this.dataType !== 'folks' && this.dataType !== 'stars') {
          this.dataType = 'folks'
        }
        addDataTypeOption('folks', 'Folks')
        addDataTypeOption('stars', 'Stars')
        show(_this.dataTypeWrapperTarget)
        break
    }

    if (this.dataType === '' && this.dataTypeTarget.innerHTML !== '') {
      this.dataType = this.dataTypeTarget.value = this.dataTypeTarget.options[0].innerText
    }

    this.dataTypeTarget.value = this.dataType
  }

  loadPreviousPage () {
    this.nextPage = this.currentPage - 1
    this.fetchData()
  }

  loadNextPage () {
    this.nextPage = this.currentPage + 1
    this.fetchData()
  }

  numberOfRowsChanged () {
    this.nextPage = 1
    this.fetchData()
  }

  fetchData () {
    const numberOfRows = this.selectedNumTarget.value

    let elementsToToggle = [this.tableWrapperTarget]
    showLoading(this.loadingDataTarget, elementsToToggle)

    const _this = this
    const queryString = `page=${_this.nextPage}&records-per-page=${numberOfRows}&view-option=` +
      `${_this.viewOption}&platform=${this.platform}&subreddit=${this.subreddit}&twitter-handle=${this.twitterHandle}` +
      `&repository=${this.repository}`
    axios.get(`/getCommunityStat?${queryString}`)
      .then(function (response) {
        hideLoading(_this.loadingDataTarget, elementsToToggle)
        let result = response.data
        if (result.message) {
          let messageHTML = ''
          messageHTML += `<div class="alert alert-primary">
                           <strong>${result.message}</strong>
                      </div>`

          _this.messageViewTarget.innerHTML = messageHTML
          show(_this.messageViewTarget)
          hide(_this.tableTarget)
          hide(_this.pageSizeWrapperTarget)
          _this.totalPageCountTarget.textContent = 0
          _this.currentPageTarget.textContent = 0
          window.history.pushState(window.history.state, _this.addr, `/community?${queryString}`)
        } else {
          show(_this.tableTarget)
          show(_this.pageSizeWrapperTarget)
          hide(_this.messageViewTarget)
          const pageUrl = `/community?${queryString}`
          window.history.pushState(window.history.state, _this.addr, pageUrl)

          _this.currentPage = result.currentPage
          if (_this.currentPage <= 1) {
            hide(_this.previousPageButtonTarget)
          } else {
            show(_this.previousPageButtonTarget)
          }

          if (_this.currentPage >= result.totalPages) {
            hide(_this.nextPageButtonTarget)
          } else {
            show(_this.nextPageButtonTarget)
          }

          _this.totalPageCountTarget.textContent = result.totalPages
          _this.currentPageTarget.textContent = result.currentPage

          _this.displayRecord(result.stats, result.columns)
        }
      }).catch(function (e) {
        console.log(e)
      })
  }

  displayRecord (stats, columns) {
    hide(this.messageViewTarget)
    show(this.tableWrapperTarget)
    const _this = this
    this.tableTarget.innerHTML = ''

    this.tableCol1Target.innerText = columns[0]
    this.tableCol2Target.innerText = columns[1]
    if (columns.length > 2) {
      this.tableCol3Target.innerText = columns[2]
      show(this.tableCol2Target)
    } else {
      hide(this.tableCol2Target)
    }

    if (!stats) {
      return
    }

    stats.forEach(stat => {
      const exRow = document.importNode(_this.rowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerHTML = formatDate(new Date(stat.date))
      switch (_this.platform) {
        case 'Reddit':
          _this.displayRedditData(stat, fields)
          break
        case 'Twitter':
          _this.displayTwitterStat(stat, fields)
          break
        case 'Github':
          _this.displayGithubData(stat, fields)
          break
        case 'Youtube':
          _this.displayYoutubeData(stat, fields)
          break
      }

      _this.tableTarget.appendChild(exRow)
    })
  }

  displayRedditData (stat, fields) {
    fields[1].innerHTML = stat.subscribers
    fields[2].innerText = stat.active_user_count
  }

  displayTwitterStat (stat, fields) {
    fields[1].innerHTML = stat.followers
    hide(fields[2])
  }

  displayGithubData (stat, fields) {
    fields[1].innerHTML = stat.stars
    fields[2].innerText = stat.folks
  }

  displayYoutubeData (stat, fields) {
    fields[1].innerHTML = stat.subscribers
    hide(fields[2])
  }

  fetchDataAndPlotGraph () {
    let elementsToToggle = [this.chartWrapperTarget]
    showLoading(this.loadingDataTarget, elementsToToggle)

    const _this = this
    const queryString = `data-type=${this.dataType}&platform=${this.platform}&subreddit=${_this.subreddit}` +
      `&twitter-handle=${this.twitterHandle}&view-option=${this.viewOption}&repository=${this.repository}`
    window.history.pushState(window.history.state, _this.addr, `/community?${queryString}`)

    axios.get(`/communitychat?${queryString}`).then(function (response) {
      hideLoading(_this.loadingDataTarget, elementsToToggle)
      let result = response.data
      if (result.error) {
        console.log(result.error) // todo show error page from front page
        return
      }

      _this.plotGraph(result)
    }).catch(function (e) {
      hideLoading(_this.loadingDataTarget, elementsToToggle)
      console.log(e)
    })
  }

  // vsp chart
  plotGraph (dataSet) {
    const _this = this

    let options = {
      legend: 'always',
      includeZero: true,
      legendFormatter: legendFormatter,
      labelsDiv: _this.labelsTarget,
      ylabel: dataSet.ylabel,
      xlabel: 'Date',
      labels: ['Date', dataSet.ylabel],
      labelsUTC: true,
      labelsKMB: true,
      connectSeparatedPoints: true,
      showRangeSelector: true,
      axes: {
        x: {
          drawGrid: false
        }
      }
    }

    _this.chartsView = new Dygraph(_this.chartsViewTarget, dataSet.stats, options)
  }
}
