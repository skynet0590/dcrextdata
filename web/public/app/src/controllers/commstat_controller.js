import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, setActiveOptionBtn, showLoading, hideLoading } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')

export default class extends Controller {
  viewOption

  static get targets () {
    return [
      'pageSizeWrapper', 'previousPageButton', 'totalPageCount', 'nextPageButton',
      'currentPage', 'numPageWrapper', 'selectedNum', 'messageView',
      'viewOptionControl', 'viewOption',
      'chartWrapper', 'chartsView', 'tableWrapper', 'loadingData', 'messageView',
      'tableWrapper', 'table', 'rowTemplate'
    ]
  }

  initialize () {
    this.currentPage = parseInt(this.currentPageTarget.dataset.currentPage)
    if (this.currentPage < 1) {
      this.currentPage = 1
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
    axios.get(`/getCommunityStat?page=${_this.nextPage}&records-per-page=${numberOfRows}&view-option=${_this.viewOption}`)
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
          window.history.pushState(window.history.state, _this.addr, `/communityStat?page=${_this.nextPage}&records-per-page=${numberOfRows}&view-option=${_this.viewOption}`)
        } else {
          show(_this.tableTarget)
          show(_this.pageSizeWrapperTarget)
          hide(_this.messageViewTarget)
          const pageUrl = `/communityStat?page=${result.currentPage}&records-per-page=${result.selectedNum}&view-option=${_this.viewOption}`
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

          _this.displayRecord(result.stats)
        }
      }).catch(function (e) {
        console.log(e)
      })
  }

  displayRecord (stats) {
    hide(this.messageViewTarget)
    show(this.tableWrapperTarget)
    const _this = this
    this.tableTarget.innerHTML = ''

    stats.forEach(stat => {
      const exRow = document.importNode(_this.rowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerHTML = stat.date
      fields[1].innerText = stat.youtube_subscribers
      fields[2].innerText = stat.twitter_followers
      fields[3].innerText = stat.github_stars
      fields[4].innerHTML = stat.github_folks

      const redditStat = stat.reddit_stats
      for (let subreddit in redditStat) {
        const currSubreddit = subreddit
        if (redditStat.hasOwnProperty(currSubreddit)) {
          fields[5].innerHTML += `<p>${currSubreddit} Subscriber: ${redditStat[currSubreddit].subscribers}, 
                                        Active: ${redditStat[currSubreddit].active_user_count}</p>`
        }
      }
      _this.tableTarget.appendChild(exRow)
    })
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
