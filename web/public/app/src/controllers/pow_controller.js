import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, options, setActiveOptionBtn } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')

export default class extends Controller {
  static get targets () {
    return [
      'vspFilterWrapper', 'selectedFilter', 'powTable', 'numPageWrapper',
      'previousPageButton', 'totalPageCount', 'nextPageButton',
      'powRowTemplate', 'currentPage', 'selectedNum', 'powTableWrapper',
      'chartSourceWrapper', 'pool', 'chartWrapper', 'chartDataTypeSelector', 'labels',
      'chartsView', 'viewOption', 'pageSizeWrapper'
    ]
  }

  initialize () {
    this.viewOption = 'table'
    this.dataType = 'pool_hashrate'
  }

  connect () {
    var filter = this.selectedFilterTarget.options
    var num = this.selectedNumTarget.options
    this.selectedFilterTarget.value = filter[0].text
    this.selectedNumTarget.value = num[0].text
  }

  setTable () {
    this.viewOption = 'table'
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

  loadPreviousPage () {
    this.nextPage = this.previousPageButtonTarget.getAttribute('data-next-page')
    this.fetchData(this.viewOption)
  }

  loadNextPage () {
    this.nextPage = this.nextPageButtonTarget.getAttribute('data-next-page')
    this.fetchData(this.viewOption)
  }

  selectedFilterChanged () {
    this.nextPage = 1
    this.fetchData(this.viewOption)
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
        _this.previousPageButtonTarget.setAttribute('data-next-page', `${result.previousPage}`)
        _this.nextPageButtonTarget.setAttribute('data-next-page', `${result.nextPage}`)

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

  selectedChartDataTypeChanged (event) {
    this.dataType = event.currentTarget.value
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

  plotGraph1 (data) {
    const _this = this

    console.log(data.csv)
    let dataTypeLabel = 'Pool Hashrate'
    if (_this.dataType === 'workers') {
      dataTypeLabel = 'Workers'
    }

    const extra = {
      includeZero: true,
      colors: ['#2971FF', '#FF8C00'],
      labelsDiv: this.labelsTarget,
      ylabel: dataTypeLabel,
      labelsKMB: true,
      legendFormatter: legendFormatter,
      dateWindow: [data.minDate, data.maxDate],
      xlabel: 'Date',
      labelsUTC: true,
      connectSeparatedPoints: true
    }

    _this.chartsView = new Dygraph(_this.chartsViewTarget, data.csv, { ...options, ...extra }
    )
  }

  // vsp chart
  plotGraph (dataSet) {
    const _this = this
    let dataTypeLabel = 'Pool Hashrate'
    if (_this.dataType === 'workers') {
      dataTypeLabel = 'Workers'
    }

    let options = {
      legend: 'always',
      includeZero: true,
      animatedZooms: true,
      legendFormatter: legendFormatter,
      // plotter: barChartPlotter,
      labelsDiv: _this.labelsTarget,
      ylabel: dataTypeLabel,
      xlabel: 'Date',
      labelsUTC: true,
      labelsKMB: true,
      connectSeparatedPoints: true,
      axes: {
        x: {
          drawGrid: false
        }
      }
    }

    _this.chartsView = new Dygraph(_this.chartsViewTarget, dataSet.csv, options)
  }
}
