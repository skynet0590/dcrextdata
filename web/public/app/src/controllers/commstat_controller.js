import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, setActiveOptionBtn, showLoading, hideLoading, formatDate, trimUrl, insertOrUpdateQueryParam, removeUrlParam } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')
const redditPlatform = 'Reddit'
const twitterPlatform = 'Twitter'
const githubPlatform = 'GitHub'
const youtubePlatform = 'YouTube'

export default class extends Controller {
  viewOption
  platform
  subreddit
  twitterHandle
  repository
  dataType

  static get targets () {
    return [
      'paginationWrapper', 'previousPageButton', 'totalPageCount', 'nextPageButton',
      'currentPage', 'pageSizeWrapper', 'pageSize', 'messageView',
      'viewOptionControl', 'viewOption',
      'chartWrapper', 'chartsView', 'labels', 'tableWrapper', 'loadingData', 'messageView',
      'tableWrapper', 'table', 'rowTemplate', 'tableCol1', 'tableCol2', 'tableCol3',
      'platform', 'subreddit', 'subAccountWrapper', 'dataTypeWrapper', 'dataType',
      'twitterHandle', 'repository', 'channel'
    ]
  }

  initialize () {
    this.currentPage = parseInt(this.currentPageTarget.dataset.currentPage)
    if (this.currentPage < 1) {
      this.currentPage = 1
    }

    this.pageSize = this.pageSizeTarget.value

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

    this.channel = this.channelTarget.dataset.initialValue
    if (this.channel === '') {
      this.channel = this.channelTarget.value = this.channelTarget.options[0].innerText
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
    insertOrUpdateQueryParam('view-option', this.viewOption, 'chart')
    setActiveOptionBtn(this.viewOption, this.viewOptionTargets)
    hide(this.chartWrapperTarget)
    hide(this.messageViewTarget)
    show(this.tableWrapperTarget)
    show(this.pageSizeWrapperTarget)
    show(this.paginationWrapperTarget)
    this.pageSizeTarget.value = this.pageSize
    this.updateDataTypeControl()
    this.fetchData()
  }

  setChart () {
    this.viewOption = 'chart'
    insertOrUpdateQueryParam('view-option', this.viewOption, 'chart')
    setActiveOptionBtn(this.viewOption, this.viewOptionTargets)
    hide(this.tableWrapperTarget)
    hide(this.messageViewTarget)
    show(this.chartWrapperTarget)
    hide(this.paginationWrapperTarget)
    hide(this.pageSizeWrapperTarget)
    this.updateDataTypeControl()
    this.fetchDataAndPlotGraph()
    // reset this table properties as the url params will be reset
    this.currentPage = 1
    this.pageSize = 20
  }

  trimUrlParam () {
    var baseSet = ['platform', 'view-option']
    var keepSet = []
    if (this.viewOption === 'table') {
      const tableParams = ['page', 'records-per-page', ...baseSet]
      switch (this.platform) {
        case redditPlatform:
          keepSet = ['subreddit', ...tableParams]
          break
        case youtubePlatform:
          keepSet = ['channel', ...tableParams]
          break
        case githubPlatform:
          keepSet = ['repository', ...tableParams]
          break
        case twitterPlatform:
          keepSet = ['twitter-handle', ...tableParams]
          break
      }
    } else {
      var chartParams = baseSet
      switch (this.platform) {
        case redditPlatform:
          keepSet = ['subreddit', 'data-type', ...chartParams]
          break
        case youtubePlatform:
          keepSet = ['data-type', ...chartParams]
          break
        case githubPlatform:
          keepSet = ['repository', 'data-type', ...chartParams]
          break
        case twitterPlatform:
          keepSet = ['twitter-handle', ...chartParams]
          break
      }
    }

    trimUrl(keepSet)
  }

  platformChanged (event) {
    this.platform = event.currentTarget.value
    insertOrUpdateQueryParam('platform', this.platform, this.platformTarget.options[0].value)
    this.showCurrentSubAccountWrapper()
    this.updateDataTypeControl()
    this.resetSubAccountsAndDataType()
    removeUrlParam('data-type')
    this.currentPage = 1
    if (this.viewOption === 'table') {
      this.fetchData()
    } else {
      this.fetchDataAndPlotGraph()
    }
  }

  resetSubAccountsAndDataType () {
    if (this.subredditTarget.options.length > 0) {
      this.subredditTarget.value = this.subredditTarget.options[0].value
    }
    if (this.twitterHandleTarget.options.length > 0) {
      this.twitterHandleTarget.value = this.twitterHandleTarget.options[0].value
    }
    if (this.repositoryTarget.options.length > 0) {
      this.repositoryTarget.value = this.repositoryTarget.options[0].value
    }
    if (this.channelTarget.options.length > 0) {
      this.channelTarget.value = this.channelTarget.options[0].value
    }
    if (this.dataTypeTarget.options.length > 0) {
      this.dataTypeTarget.value = this.dataTypeTarget.options[0].value
    }
  }

  subredditChanged (event) {
    this.subreddit = event.currentTarget.value
    let defaultSubreddit
    if (event.currentTarget.options.length > 0) {
      defaultSubreddit = event.currentTarget.options[0].value
    }
    insertOrUpdateQueryParam('subreddit', this.subreddit, defaultSubreddit)
    this.currentPage = 1
    if (this.viewOption === 'table') {
      this.fetchData()
    } else {
      this.fetchDataAndPlotGraph()
    }
  }

  twitterHandleChanged (event) {
    this.twitterHandle = event.currentTarget.value
    let defaultTwitterHandle
    if (event.currentTarget.options.length > 0) {
      defaultTwitterHandle = event.currentTarget.options[0].value
    }
    insertOrUpdateQueryParam('twitter-handle', this.twitterHandle, defaultTwitterHandle)
    this.currentPage = 1
    if (this.viewOption === 'table') {
      this.fetchData()
    } else {
      this.fetchDataAndPlotGraph()
    }
  }

  repositoryChanged (event) {
    this.repository = event.currentTarget.value
    let defaultRepository
    if (event.currentTarget.options.length > 0) {
      defaultRepository = event.currentTarget.options[0].value
    }
    insertOrUpdateQueryParam('repository', this.repository, defaultRepository)
    this.currentPage = 1
    if (this.viewOption === 'table') {
      this.fetchData()
    } else {
      this.fetchDataAndPlotGraph()
    }
  }

  channelChanged (event) {
    this.channel = event.currentTarget.value
    let defaultChannel
    if (event.currentTarget.options.length > 0) {
      defaultChannel = event.currentTarget.options[0].value
    }
    insertOrUpdateQueryParam('channel', this.channel, defaultChannel)
    this.currentPage = 1
    if (this.viewOption === 'table') {
      this.fetchData()
    } else {
      this.fetchDataAndPlotGraph()
    }
  }

  dataTypeChanged (event) {
    this.dataType = event.currentTarget.value
    let defaultDataType
    if (event.currentTarget.options.length > 0) {
      defaultDataType = event.currentTarget.options[0].value
    }
    insertOrUpdateQueryParam('data-type', this.dataType, defaultDataType)
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
      case redditPlatform:
        if (this.dataType !== 'subscribers' && this.dataType !== 'active_accounts') {
          this.dataType = 'subscribers'
        }
        addDataTypeOption('subscribers', 'Subscribers')
        addDataTypeOption('active_accounts', 'Active Accounts')
        show(_this.dataTypeWrapperTarget)
        break
      case githubPlatform:
        if (this.dataType !== 'folks' && this.dataType !== 'stars') {
          this.dataType = 'folks'
        }
        addDataTypeOption('folks', 'Forks')
        addDataTypeOption('stars', 'Stars')
        show(_this.dataTypeWrapperTarget)
        break
      case youtubePlatform:
        if (this.dataType !== 'subscribers' && this.dataType !== 'view_count') {
          this.dataType = 'subscribers'
        }
        addDataTypeOption('subscribers', 'Subscribers')
        addDataTypeOption('view_count', 'View Count')
        show(_this.dataTypeWrapperTarget)
        break
    }

    if (this.dataType === '' && this.dataTypeTarget.innerHTML !== '') {
      this.dataType = this.dataTypeTarget.value = this.dataTypeTarget.options[0].innerText
    }

    this.dataTypeTarget.value = this.dataType
  }

  loadPreviousPage () {
    this.currentPage -= 1
    if (this.currentPage < 1) {
      this.currentPage = 1
    }
    insertOrUpdateQueryParam('page', this.currentPage, 1)
    this.fetchData()
  }

  loadNextPage () {
    this.currentPage += 1
    insertOrUpdateQueryParam('page', this.currentPage, 1)
    this.fetchData()
  }

  pageSizeChanged (event) {
    this.currentPage = 1
    this.pageSize = event.currentTarget.value
    let defaultPageSize
    if (event.currentTarget.options.length > 0) {
      defaultPageSize = event.currentTarget.options[0].value
    }
    insertOrUpdateQueryParam('page', this.currentPage, 1)
    insertOrUpdateQueryParam('records-per-page', this.pageSize, defaultPageSize)
    this.fetchData()
  }

  fetchData () {
    let elementsToToggle = [this.tableWrapperTarget]
    showLoading(this.loadingDataTarget, elementsToToggle)
    const _this = this
    const queryString = `page=${_this.currentPage}&records-per-page=${this.pageSize}&view-option=` +
      `${_this.viewOption}&platform=${this.platform}&subreddit=${this.subreddit}&twitter-handle=${this.twitterHandle}` +
      `&repository=${this.repository}&channel=${this.channel}`
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
          hide(_this.paginationWrapperTarget)
          _this.totalPageCountTarget.textContent = 0
          _this.currentPageTarget.textContent = 0
          _this.trimUrlParam()
        } else {
          show(_this.tableTarget)
          show(_this.paginationWrapperTarget)
          hide(_this.messageViewTarget)
          _this.trimUrlParam()

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
      show(this.tableCol3Target)
    } else {
      hide(this.tableCol3Target)
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
        case 'GitHub':
          _this.displayGithubData(stat, fields)
          break
        case 'YouTube':
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
    fields[2].innerText = stat.view_count
  }

  fetchDataAndPlotGraph () {
    let elementsToToggle = [this.chartWrapperTarget]
    showLoading(this.loadingDataTarget, elementsToToggle)

    const _this = this
    const queryString = `data-type=${this.dataType}&platform=${this.platform}&subreddit=${_this.subreddit}` +
      `&twitter-handle=${this.twitterHandle}&view-option=${this.viewOption}&repository=${this.repository}&channel=${this.channel}`
    _this.trimUrlParam()

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
