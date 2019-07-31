import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, setActiveOptionBtn } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')

export default class extends Controller {
  static get targets () {
    return [
      'vspFilterWrapper', 'selectedFilter', 'powTable', 'numPageWrapper',
      'previousPageButton', 'totalPageCount', 'nextPageButton',
      'powRowTemplate', 'currentPage', 'selectedNum', 'powTableWrapper',
      'chartSourceWrapper', 'pool', 'chartWrapper', 'chartDataTypeSelector', 'dataType', 'labels',
      'chartsView', 'viewOption', 'pageSizeWrapper'
    ]
  }

  initialize () {
    this.dataType = 'pool_hashrate'
    this.setChart()
  }

  setTable () {
    this.viewOption = 'table'
    var filter = this.selectedFilterTarget.options
    var num = this.selectedNumTarget.options
    this.selectedFilterTarget.value = filter[0].text
    this.selectedNumTarget.value = num[0].text
    setActiveOptionBtn(this.viewOption, this.viewOptionTargets)
    hide(this.chartWrapperTarget)
    hide(this.chartSourceWrapperTarget)
    show(this.vspFilterWrapperTarget)
    show(this.powTableWrapperTarget)
    show(this.numPageWrapperTarget)
    show(this.pageSizeWrapperTarget)
    hide(this.chartDataTypeSelectorTarget)
    this.nextPage = 1
    this.fetchData()
  }

  setChart () {
    this.viewOption = 'chart'
    setActiveOptionBtn(this.viewOption, this.viewOptionTargets)
    hide(this.numPageWrapperTarget)
    hide(this.vspFilterWrapperTarget)
    hide(this.powTableWrapperTarget)
    show(this.chartSourceWrapperTarget)
    show(this.chartWrapperTarget)
    hide(this.pageSizeWrapperTarget)
    show(this.chartDataTypeSelectorTarget)
    this.nextPage = 1
    this.fetchDataAndPlotGraph()
  }

  poolCheckChanged (event) {
    this.fetchDataAndPlotGraph()
  }

  selectedFilterChanged () {
    this.nextPage = 1
    this.fetchData(this.viewOption)
  }

  loadPreviousPage () {
    this.nextPage = this.currentPage - 1
    this.fetchExchange(this.viewOption)
  }

  loadNextPage () {
    this.nextPage = this.currentPage + 1
    this.fetchExchange(this.viewOption)
  }

  numberOfRowsChanged () {
    this.nextPage = 1
    this.fetchData(this.viewOption)
  }

  fetchData () {
    const selectedFilter = this.selectedFilterTarget.value
    var numberOfRows = this.selectedNumTarget.value

    const _this = this
    axios.get(`/filteredpow?page=${this.nextPage}&filter=${selectedFilter}&recordsPerPage=${numberOfRows}`)
      .then(function (response) {
        let result = response.data
        window.history.pushState(window.history.state, _this.addr, `pow?page=${result.previousPage}&filter=${selectedFilter}&recordsPerPage=${result.selectedNum}`)

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
        _this.previousPageButtonTarget.setAttribute('href', `?page=${result.previousPage}&filter=${result.selectedFilter}&recordsPerPage=${result.selectedNum}`)
        _this.nextPageButtonTarget.setAttribute('href', `?page=${result.nextPage}&filter=${result.selectedFilter}&recordsPerPage=${result.selectedNum}`)

        _this.displayPoW(result.powData)
      }).catch(function (e) {
        console.log(e)
      })
  }

  displayPoW (pows) {
    const _this = this
    this.powTableTarget.innerHTML = ''

    pows.forEach(pow => {
      const powRow = document.importNode(_this.powRowTemplateTarget.content, true)
      const fields = powRow.querySelectorAll('td')

      fields[0].innerText = pow.source
      fields[1].innerText = pow.pool_hashrate_th
      fields[2].innerHTML = pow.workers
      fields[4].innerHTML = pow.time

      _this.powTableTarget.appendChild(powRow)
    })
  }

  setDataType (event) {
    this.dataType = event.currentTarget.getAttribute('data-option')
    setActiveOptionBtn(this.dataType, this.dataTypeTargets)
    this.fetchDataAndPlotGraph()
  }

  fetchDataAndPlotGraph () {
    let selectedPools = []
    this.poolTargets.forEach(el => {
      if (el.checked) {
        selectedPools.push(el.value)
      }
    })

    const _this = this
    axios.get(`/powchartdata?pools=${selectedPools.join('|')}&datatype=${this.dataType}`).then(function (response) {
      let result = response.data
      if (result.error) {
        console.log(result.error) // todo show error page fron front page
        return
      }

      _this.plotGraph(result)
    }).catch(function (e) {
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
      // plotter: barChartPlotter,
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
  }
}
