import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, options } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')
var opt = 'table'
// var ylabel

export default class extends Controller {
  static get targets () {
    return [
      'nextPageButton', 'previousPageButton', 'tableBody', 'rowTemplate',
      'totalPageCount', 'currentPage', 'btnWrapper', 'tableWrapper', 'chartsView',
      'chartWrapper', 'viewOption', 'chartOptions', 'labels', 'selectedMempoolOpt'
    ]
  }

  initialize () {
    this.currentPage = parseInt(this.currentPageTarget.getAttribute('data-current-page'))
    if (this.currentPage < 1) {
      this.currentPage = 1
    }
  }

  setTable () {
    opt = 'table'
    this.chartOptionsTarget.classList.add('d-hide')
    this.setActiveOptionBtn(opt, this.viewOptionTargets)
    this.chartWrapperTarget.classList.add('d-hide')
    this.tableWrapperTarget.classList.remove('d-hide')
    this.btnWrapperTarget.classList.remove('d-hide')
    this.currentPage = this.currentPage
    this.fetchData(opt)
  }

  setChart () {
    opt = 'chart'
    var y = this.selectedMempoolOptTarget.options
    this.chartFilter = this.selectedMempoolOptTarget.value = y[0].value
    this.chartOptionsTarget.classList.remove('d-hide')
    this.btnWrapperTarget.classList.add('d-hide')
    this.tableWrapperTarget.classList.add('d-hide')
    this.setActiveOptionBtn(opt, this.viewOptionTargets)
    this.chartWrapperTarget.classList.remove('d-hide')
    this.nextPage = 1
    this.fetchData(opt)
  }

  MempoolOptionChanged () {
    this.chartFilter = this.selectedMempoolOptTarget.value
    this.fetchData(opt)
  }

  gotoPreviousPage () {
    this.currentPage = this.currentPage - 1
    this.fetchData(opt)
  }

  gotoNextPage () {
    this.currentPage = this.currentPage + 1
    this.fetchData(opt)
  }

  fetchData (display) {
    var url
    if (display === 'table') {
      url = `/getmempool?page=${this.currentPage}`
    } else {
      url = `/getmempoolCharts?chartFilter=${this.chartFilter}`
    }

    const _this = this
    axios.get(url).then(function (response) {
      let result = response.data
      if (display === 'table') {
        _this.totalPageCountTarget.textContent = result.totalPages
        _this.currentPageTarget.textContent = result.currentPage

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

        _this.displayMempool(result.mempoolData)
      } else {
        console.log(result)
        _this.plotGraph(result)
      }
    }).catch(function (e) {
      console.log(e) // todo: handle error
    })
  }

  displayMempool (data) {
    const _this = this
    this.tableBodyTarget.innerHTML = ''

    data.forEach(item => {
      const exRow = document.importNode(_this.rowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerText = item.time
      fields[1].innerText = item.number_of_transactions
      fields[2].innerText = item.size
      fields[3].innerHTML = item.total_fee.toFixed(8)

      _this.tableBodyTarget.appendChild(exRow)
    })
  }

  // exchange chart
  plotGraph (exs) {
    var chartData = exs.mempoolchartData
    var mempool = exs.chartFilter
    console.log(mempool)
    var extra = {
      legendFormatter: legendFormatter,
      labelsDiv: this.labelsTarget,
      ylabel: mempool,
      labels: ['Date', mempool],
      colors: ['#2971FF', '#FF8C00']
    }

    const _this = this

    var data = []
    var dataSet = []
    chartData.forEach(mp => {
      data.push(new Date(mp.time))
      if (mempool === 'size') {
        data.push(mp.size)
      } else if (mempool === 'total_fee') {
        data.push(mp.total_fee)
      } else {
        data.push(mp.number_of_transactions)
      }

      dataSet.push(data)
      data = []
    })
    console.log(dataSet)
    _this.chartsView = new Dygraph(
      _this.chartsViewTarget,
      dataSet, { ...options, ...extra }
    )
  }

  setActiveOptionBtn (opt, optTargets) {
    optTargets.forEach(li => {
      if (li.dataset.option === opt) {
        li.classList.add('active')
      } else {
        li.classList.remove('active')
      }
    })
  }
}
