import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show } from '../utils'

// const Dygraph = require('../../../dist/js/dygraphs.min.js')
var opt = 'table'

export default class extends Controller {
  static get targets () {
    return [
      'nextPageButton', 'previousPageButton', 'tableBody', 'rowTemplate',
      'totalPageCount', 'currentPage', 'btnWrapper', 'tableWrapper',
      'chartWrapper', 'viewOption', 'chartOptions', 'selectedMempoolOpt'
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
    this.exchangeTableWrapperTarget.classList.remove('d-hide')
    this.btnWrapperTarget.classList.remove('d-hide')
    this.tableWrapperTarget.innerHTML = ''
    this.nextPage = 1
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
    const _this = this
    var url
    if (display === 'table') {
      url = `/getmempool?page=${this.currentPage}`
    } else {
      url = `/getmempoolCharts?chartFilter=${this.chartFilter}`
    }
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
        console.log(response)
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
