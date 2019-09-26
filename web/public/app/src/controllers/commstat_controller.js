import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, setActiveOptionBtn, showLoading, hideLoading } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')

export default class extends Controller {
  viewOption
  platform
  subreddit

  static get targets () {
    return [
      'pageSizeWrapper', 'previousPageButton', 'totalPageCount', 'nextPageButton',
      'currentPage', 'numPageWrapper', 'selectedNum', 'messageView',
      'viewOptionControl', 'viewOption',
      'chartWrapper', 'chartsView', 'tableWrapper', 'loadingData', 'messageView',
      'tableWrapper', 'table', 'rowTemplate', 'tableCol1', 'tableCol2', 'tableCol3',
      'platform', 'subreddit', 'subredditWrapper'
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

    if (this.platform === 'Reddit') {
      show(this.subredditWrapperTarget)
    } else {
      hide(this.subredditWrapperTarget)
    }

    this.subreddit = this.subredditTarget.dataset.initialValue
    if (this.subreddit === '') {
      this.subreddit = this.subredditTarget.value = this.subredditTarget.options[0].innerText
    }

    this.viewOption = this.viewOptionControlTarget.dataset.viewOption
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
    this.fetchDataAndPlotGraph()
  }

  platformChanged (event) {
    this.platform = event.currentTarget.value
    if (this.platform === 'Reddit') {
      show(this.subredditWrapperTarget)
    } else {
      hide(this.subredditWrapperTarget)
    }
    this.currentPage = 1
    this.fetchData()
  }

  subredditChanged (event) {
    this.subreddit = event.currentTarget.value
    this.currentPage = 1
    this.fetchData()
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
    axios.get(`/getCommunityStat?page=${_this.nextPage}&records-per-page=${numberOfRows}&view-option=${_this.viewOption}&platform=${this.platform}&subreddit=${this.subreddit}`)
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
          window.history.pushState(window.history.state, _this.addr, `/communityStat?page=${_this.nextPage}&records-per-page=${numberOfRows}&view-option=${_this.viewOption}&platform=${_this.platform}&subreddit=${_this.subreddit}`)
        } else {
          show(_this.tableTarget)
          show(_this.pageSizeWrapperTarget)
          hide(_this.messageViewTarget)
          const pageUrl = `/communityStat?page=${result.currentPage}&records-per-page=${result.selectedNum}&view-option=${_this.viewOption}&platform=${_this.platform}&subreddit=${_this.subreddit}`
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
    } else {
      hide(this.tableCol2Target)
    }

    stats.forEach(stat => {
      const exRow = document.importNode(_this.rowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerHTML = stat.date
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
    let selectedPools = []
    this.poolTargets.forEach(el => {
      if (el.checked) {
        selectedPools.push(el.value)
      }
    })

    let elementsToToggle = [this.chartWrapperTarget]
    showLoading(this.loadingDataTarget, elementsToToggle)

    const _this = this
    const queryString = `data-type=${this.dataType}&pools=${selectedPools.join('|')}&view-option=${_this.selectedViewOption}`
    window.history.pushState(window.history.state, _this.addr, `/pow?${queryString}`)

    axios.get(`/powchart?${queryString}`).then(function (response) {
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
    let dataTypeLabel = 'Pool Hashrate (Th/s)'
    if (_this.dataType === 'workers') {
      dataTypeLabel = 'Workers'
    }

    let options = {
      legend: 'always',
      includeZero: true,
      legendFormatter: legendFormatter,
      labelsDiv: _this.labelsTarget,
      ylabel: dataTypeLabel,
      xlabel: 'Date',
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

    _this.chartsView = new Dygraph(_this.chartsViewTarget, dataSet.csv, options)
    _this.validateZoom()
  }
}
