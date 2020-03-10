import { Controller } from 'stimulus'
import axios from 'axios'
import {
  hide,
  hideLoading,
  insertOrUpdateQueryParam, legendFormatter, options,
  selectedOption,
  setActiveOptionBtn,
  show,
  showLoading,
  updateQueryParam,
  hideAll,
  barChartPlotter
} from '../utils'

import { animationFrame } from '../helpers/animation_helper'
import Zoom from '../helpers/zoom_helper'
import humanize from '../helpers/humanize_helper'

const Dygraph = require('../../../dist/js/dygraphs.min.js')

const dataTypeNodes = 'nodes'
const dataTypeVersion = 'version'
const dataTypeLocation = 'location'

export default class extends Controller {
  timestamp
  nextTimestamp
  previousTimestamp
  height
  currentPage
  pageSize
  userAgentsPage
  query
  selectedViewOption

  static get targets () {
    return [
      'timestamp', 'snapshotHeader', 'viewOptionControl', 'viewOption', 'chartDataTypeSelector', 'chartDataType',
      'numPageWrapper', 'pageSize', 'messageView', 'chartWrapper', 'chartsView', 'labels', 'previousTimestampBtn',
      'nextTimestampBtn',
      'btnWrapper', 'nextPageButton', 'previousPageButton', 'tableTitle', 'tableWrapper', 'tableHeader', 'tableBody',
      'snapshotRowTemplate', 'userAgentRowTemplate', 'countriesRowTemplate', 'totalPageCount', 'currentPage', 'loadingData',
      'dataTypeSelector', 'dataType'
    ]
  }

  async initialize () {
    this.timestamp = parseInt(this.data.get('timestamp'))
    this.timestampTargets.forEach(el => {
      el.innerHTML = humanize.date(this.timestamp * 1000)
    })
    this.nextTimestamp = parseInt(this.data.get('nextTimestamp'))
    if (this.nextTimestamp === 0) {
      hide(this.nextTimestampBtnTarget)
    }
    this.previousTimestamp = parseInt(this.data.get('previousTimestamp'))
    if (this.previousTimestamp === 0) {
      hide(this.previousTimestampBtnTarget)
    }

    this.height = parseInt(this.data.get('height'))
    this.currentPage = parseInt(this.currentPageTarget.dataset.initialValue) || 1
    this.pageSize = parseInt(this.data.get('pageSize')) || 20
    this.selectedViewOption = this.data.get('viewOption')
    this.dataType = this.data.get('dataType') || dataTypeNodes
    setActiveOptionBtn(this.dataType, this.dataTypeTargets)

    this.zoomCallback = this._zoomCallback.bind(this)
    this.drawCallback = this._drawCallback.bind(this)

    this.userAgentsPage = 1
    this.countriesPage = 1
    this.updateView()
  }

  gotoPreviousTimestamp () {
    const urlParams = new URLSearchParams(window.location.search)
    const baseUrl = window.location.href.replace(window.location.search, '')
    if (urlParams.has('timestamp')) {
      urlParams.set('timestamp', this.previousTimestamp)
    } else {
      urlParams.append('timestamp', this.previousTimestamp)
    }
    window.location.href = `${baseUrl}?${urlParams.toString()}`
  }

  gotoNextTimestamp () {
    const urlParams = new URLSearchParams(window.location.search)
    const baseUrl = window.location.href.replace(window.location.search, '')
    if (urlParams.has('timestamp')) {
      urlParams.set('timestamp', this.nextTimestamp)
    } else {
      urlParams.append('timestamp', this.nextTimestamp)
    }
    window.location.href = `${baseUrl}?${urlParams.toString()}`
  }

  updateView () {
    if (this.dataType === dataTypeNodes) {
      hide(this.snapshotHeaderTarget)
    } else {
      show(this.snapshotHeaderTarget)
    }
    if (this.selectedViewOption === 'table') {
      this.setTable()
    } else {
      this.setChart()
    }
  }

  async setTable () {
    this.selectedViewOption = 'table'
    setActiveOptionBtn(this.selectedViewOption, this.viewOptionTargets)
    hide(this.chartWrapperTarget)
    hide(this.messageViewTarget)
    show(this.tableWrapperTarget)
    show(this.numPageWrapperTarget)
    insertOrUpdateQueryParam('view-option', this.selectedViewOption)
    this.reloadTable()
  }

  setChart () {
    this.selectedViewOption = 'chart'
    hide(this.tableWrapperTarget)
    hide(this.messageViewTarget)
    setActiveOptionBtn(this.selectedViewOption, this.viewOptionTargets)
    setActiveOptionBtn(this.dataType, this.chartDataTypeTargets)
    hide(this.numPageWrapperTarget)
    show(this.chartWrapperTarget)
    updateQueryParam('view-option', this.selectedViewOption)
    this.reloadChat()
  }

  setDataType (e) {
    this.dataType = e.currentTarget.getAttribute('data-option')
    if (this.dataType === selectedOption(this.dataTypeTargets)) {
      return
    }
    this.currentPage = 1
    insertOrUpdateQueryParam('page', this.currentPage, 1)
    setActiveOptionBtn(this.dataType, this.dataTypeTargets)
    insertOrUpdateQueryParam('data-type', this.dataType)
    this.updateView()
  }

  loadNextPage () {
    this.currentPage += 1
    insertOrUpdateQueryParam('page', this.currentPage, 1)
    this.reloadTable()
  }

  loadPreviousPage () {
    this.currentPage -= 1
    if (this.currentPage < 1) {
      this.currentPage = 1
    }
    insertOrUpdateQueryParam('page', this.currentPage, 1)
    this.reloadTable()
  }

  reloadTable () {
    let url
    let displayFn
    switch (this.dataType) {
      case dataTypeVersion:
        url = '/api/snapshots/user-agents'
        displayFn = this.displayUserAgents
        break
      case dataTypeLocation:
        url = '/api/snapshots/countries'
        displayFn = this.displayCountries
        break
      case dataTypeNodes:
      default:
        url = '/api/snapshots'
        displayFn = this.displaySnapshotTable
        break
    }
    const _this = this
    showLoading(this.loadingDataTarget, [_this.tableWrapperTarget])
    url += `?page=${this.currentPage}&page-size=${this.pageSize}&timestamp=${this.timestamp}`
    axios.get(url).then(function (response) {
      let result = response.data
      hideLoading(_this.loadingDataTarget, [_this.tableWrapperTarget])
      if (result.error) {
        let messageHTML = `<div class="alert alert-primary"><strong>${result.error}</strong></div>`
        _this.messageViewTarget.innerHTML = messageHTML
        show(_this.messageViewTarget)
        hide(_this.tableBodyTarget)
        hide(_this.btnWrapperTarget)
        return
      }
      hide(_this.messageViewTarget)
      show(_this.tableBodyTarget)
      show(_this.btnWrapperTarget)
      _this.totalPageCountTarget.textContent = result.totalPages
      _this.currentPageTarget.textContent = _this.currentPage

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
      displayFn = displayFn.bind(_this)
      displayFn(result)
    }).catch(function (e) {
      hideLoading(_this.loadingDataTarget)
      console.log(e) // todo: handle error
    })
  }

  displayUserAgents (result) {
    this.tableTitleTarget.innerHTML = 'User Agents'
    this.showHeader(dataTypeVersion)
    this.tableBodyTarget.innerHTML = ''

    const _this = this
    let offset = (_this.currentPage - 1) * _this.pageSize
    let top = Math.min(offset + this.pageSize, result.userAgents.length)
    for (let i = offset; i < top; i++) {
      let item = result.userAgents[i]
      const exRow = document.importNode(_this.userAgentRowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerText = i + offset + 1
      fields[1].innerText = item.user_agent
      fields[2].innerText = item.nodes

      _this.tableBodyTarget.appendChild(exRow)
    }
  }

  displayCountries (result) {
    this.tableTitleTarget.innerHTML = 'Countries'
    this.showHeader(dataTypeLocation)
    this.tableBodyTarget.innerHTML = ''

    const _this = this
    let offset = (_this.currentPage - 1) * _this.pageSize
    let top = Math.min(offset + this.pageSize, result.countries.length)
    for (let i = offset; i < top; i++) {
      let item = result.countries[i]
      const exRow = document.importNode(_this.countriesRowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerText = i + offset + 1
      fields[1].innerText = item.country || 'Unknown'
      fields[2].innerText = item.nodes

      _this.tableBodyTarget.appendChild(exRow)
    }
  }

  displaySnapshotTable (result) {
    this.tableTitleTarget.innerHTML = 'Network Snapshots'
    this.showHeader(dataTypeNodes)
    this.tableBodyTarget.innerHTML = ''

    const _this = this
    result.data.forEach(item => {
      const exRow = document.importNode(_this.snapshotRowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerText = humanize.date(item.timestamp * 1000)
      fields[1].innerText = item.height
      fields[2].innerText = item.node_count
      fields[3].innerText = item.oldest_node_timestamp <= 0 ? 'N/A' : humanize.timeSince(item.oldest_node_timestamp)

      _this.tableBodyTarget.appendChild(exRow)
    })
  }

  showHeader (dataType) {
    hideAll(this.tableHeaderTargets)
    this.tableHeaderTargets.forEach(el => {
      if (el.getAttribute('data-for') === dataType) {
        show(el)
      }
    })
  }

  changePageSize (e) {
    this.pageSize = parseInt(e.currentTarget.value)
    insertOrUpdateQueryParam('page-size', this.pageSize, 20)
    this.reloadTable()
  }

  // chart
  selectedZoom () { return selectedOption(this.zoomOptionTargets) }

  setZoom (e) {
    var target = e.srcElement || e.target
    var option
    if (!target) {
      let ex = this.chartsView.xAxisExtremes()
      option = Zoom.mapKey(e, ex, 1)
    } else {
      option = target.dataset.option
    }
    setActiveOptionBtn(option, this.zoomOptionTargets)
    if (!target) return // Exit if running for the first time
    this.validateZoom()
  }

  async validateZoom () {
    await animationFrame()
    await animationFrame()
    let oldLimits = this.limits || this.chartsView.xAxisExtremes()
    this.limits = this.chartsView.xAxisExtremes()
    var selected = this.selectedZoom()
    if (selected) {
      this.lastZoom = Zoom.validate(selected, this.limits, 1, 1)
    } else {
      this.lastZoom = Zoom.project(this.settings.zoom, oldLimits, this.limits)
    }
    if (this.lastZoom) {
      this.chartsView.updateOptions({
        dateWindow: [this.lastZoom.start, this.lastZoom.end]
      })
    }
    if (selected !== this.settings.zoom) {
      this._zoomCallback(this.lastZoom.start, this.lastZoom.end)
    }
    await animationFrame()
    this.chartsView.updateOptions({
      zoomCallback: this.zoomCallback,
      drawCallback: this.drawCallback
    })
  }

  _zoomCallback (start, end) {
    this.lastZoom = Zoom.object(start, end)
    this.settings.zoom = Zoom.encode(this.lastZoom)
    let ex = this.chartsView.xAxisExtremes()
    let option = Zoom.mapKey(this.settings.zoom, ex, 1)
    setActiveOptionBtn(option, this.zoomOptionTargets)
  }

  _drawCallback (graph, first) {
    if (first) return
    var start, end
    [start, end] = this.chartsView.xAxisRange()
    if (start === end) return
    if (this.lastZoom.start === start) return // only handle slide event.
    this._zoomCallback(start, end)
  }

  async reloadChat () {
    let url
    let drawChartFn

    switch (this.dataType) {
      case dataTypeVersion:
        url = `/api/snapshots/user-agents?chart=1&timestamp=${this.timestamp}`
        drawChartFn = this.drawUserAgentsChart
        break
      case dataTypeLocation:
        url = `/api/snapshots/countries?chart=1&timestamp=${this.timestamp}`
        drawChartFn = this.drawCountriesChart
        break
      case dataTypeNodes:
      default:
        url = '/api/snapshots/chart'
        drawChartFn = this.drawSnapshotChart
        break
    }

    this.drawInitialGraph()
    showLoading(this.loadingDataTarget)
    const response = await axios.get(url)
    const result = response.data
    if (result.error) {
      this.messageViewTarget.innerHTML = `<div class="alert alert-primary"><strong>${result.error}</strong></div>`
      show(this.messageViewTarget)
      hideLoading(this.loadingDataTarget)
      return
    }
    hide(this.messageViewTarget)
    drawChartFn = drawChartFn.bind(this)
    drawChartFn(result)
  }

  drawSnapshotChart (result) {
    let minDate, maxDate, csv

    result.forEach(record => {
      let date = new Date(record.timestamp * 1000)
      if (minDate === undefined || date < minDate) {
        minDate = date
      }

      if (maxDate === undefined || date > maxDate) {
        maxDate = date
      }
      csv += `${date},${record.node_count}\n`
    })

    this.chartsView = new Dygraph(
      this.chartsViewTarget,
      csv,
      {
        legend: 'always',
        includeZero: true,
        dateWindow: [minDate, maxDate],
        legendFormatter: legendFormatter,
        digitsAfterDecimal: 8,
        labelsDiv: this.labelsTarget,
        ylabel: 'Node Count',
        xlabel: 'Date',
        labels: ['Date', 'Node Count'],
        labelsUTC: true,
        labelsKMB: true,
        maxNumberWidth: 10,
        showRangeSelector: true,
        axes: {
          x: {
            drawGrid: false
          },
          y: {
            axisLabelWidth: 90
          }
        }
      }
    )
    hideLoading(this.loadingDataTarget)
  }

  drawUserAgentsChart (result) {
    let csv = ''
    let i = 0
    let labelMap = []
    result.userAgents.forEach(record => {
      csv += `${i},${record.nodes}\n`
      labelMap.push(record.user_agent)
      i++
    })

    this.chartsView = new Dygraph(
      this.chartsViewTarget,
      csv,
      {
        legend: 'always',
        includeZero: true,
        legendFormatter: legendFormatter,
        plotter: barChartPlotter,
        digitsAfterDecimal: 8,
        labelsDiv: this.labelsTarget,
        ylabel: 'Node Count',
        xlabel: 'User Agent',
        labels: ['User Agent', 'Node Count'],
        labelsUTC: true,
        labelsKMB: true,
        maxNumberWidth: 10,
        showRangeSelector: true,
        axes: {
          x: {
            valueFormatter: (x) => {
              return labelMap[parseInt(x)]
            },
            axisLabelFormatter: (x) => {
              return labelMap[parseInt(x)]
            }
          },
          y: {
            axisLabelWidth: 90
          }
        }
      }
    )
    hideLoading(this.loadingDataTarget)
  }

  drawCountriesChart (result) {
    let csv = ''
    let i = 0
    let labelMap = []
    result.countries.forEach(record => {
      csv += `${i},${record.nodes}\n`
      labelMap.push(record.country)
      i++
    })

    this.chartsView = new Dygraph(
      this.chartsViewTarget,
      csv,
      {
        legend: 'always',
        includeZero: true,
        legendFormatter: legendFormatter,
        plotter: barChartPlotter,
        digitsAfterDecimal: 8,
        labelsDiv: this.labelsTarget,
        ylabel: 'Node Count',
        xlabel: 'Country',
        labels: ['Country', 'Node Count'],
        labelsUTC: true,
        labelsKMB: true,
        maxNumberWidth: 10,
        showRangeSelector: true,
        axes: {
          x: {
            valueFormatter: (x) => {
              return labelMap[parseInt(x)]
            },
            axisLabelFormatter: (x) => {
              return labelMap[parseInt(x)]
            }
          },
          y: {
            axisLabelWidth: 90
          }
        }
      }
    )
    hideLoading(this.loadingDataTarget)
  }

  drawInitialGraph () {
    var extra = {
      legendFormatter: legendFormatter,
      labelsDiv: this.labelsTarget,
      ylabel: 'Node Count',
      xlabel: 'Date',
      labels: ['Date', 'Node Count'],
      labelsUTC: true,
      labelsKMB: true,
      axes: {
        x: {
          drawGrid: false
        }
      }
    }

    this.chartsView = new Dygraph(
      this.chartsViewTarget,
      [[0, 0]],
      { ...options, ...extra }
    )
  }
}
