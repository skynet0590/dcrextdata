import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, options } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')
var opt = 'table'

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter', 'exchangeTable', 'selectedCpair', 'numPageWrapper', 'intervalWapper',
      'previousPageButton', 'totalPageCount', 'nextPageButton', 'selectedDticks', 'selectedInterval',
      'exRowTemplate', 'currentPage', 'selectedNum', 'exchangeTableWrapper', 'tickWapper',
      'chartWrapper', 'labels', 'chartsView', 'viewOption', 'hideOption', 'sourceWrapper'
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

    this.intervalWapperTarget.classList.add('d-hide')
    this.selectedDticksTarget.value = 'close'
    this.tickWapperTarget.classList.add('d-hide')
    this.sourceWrapperTarget.classList.remove('d-hide')
    this.hideOptionTarget.classList.remove('d-hide')
    this.selectedCpair = this.selectedCpairTarget.value
    this.setActiveOptionBtn(opt, this.viewOptionTargets)
    this.chartWrapperTarget.classList.add('d-hide')
    this.exchangeTableWrapperTarget.classList.remove('d-hide')
    this.numPageWrapperTarget.classList.remove('d-hide')
    this.exchangeTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange(opt)
  }

  setChart () {
    opt = 'chart'

    this.intervalWapperTarget.classList.remove('d-hide')
    this.tickWapperTarget.classList.remove('d-hide')
    this.sourceWrapperTarget.classList.add('d-hide')
    this.hideOptionTarget.classList.add('d-hide')
    this.numPageWrapperTarget.classList.add('d-hide')
    this.exchangeTableWrapperTarget.classList.add('d-hide')
    this.setActiveOptionBtn(opt, this.viewOptionTargets)
    this.chartWrapperTarget.classList.remove('d-hide')
    var y = this.selectedIntervalTarget.options
    this.selectedInterval = this.selectedIntervalTarget.value = y[0].text
    this.selectedDtick = this.selectedDticksTarget.value = 'close'
    this.selectedCpair = this.selectedCpairTarget.value = 'BTC/DCR'
    this.fetchExchange(opt)
  }

  selectedIntervalChanged () {
    this.selectedInterval = this.selectedIntervalTarget.value
    this.fetchExchange(opt)
  }

  selectedDticksChanged () {
    this.selectedDtick = this.selectedDticksTarget.value
    this.fetchExchange(opt)
  }

  loadPreviousPage () {
    this.selectedCpair = this.selectedCpairTarget.value
    this.nextPage = this.previousPageButtonTarget.getAttribute('data-previous-page')
    this.fetchExchange(opt)
  }

  loadNextPage () {
    this.selectedCpair = this.selectedCpairTarget.value
    this.nextPage = this.nextPageButtonTarget.getAttribute('data-next-page')
    this.fetchExchange(opt)
  }

  selectedFilterChanged () {
    this.nextPage = 1
    this.selectedCpair = this.selectedCpairTarget.value
    this.fetchExchange(opt)
  }

  selectedCpairChanged () {
    this.nextPage = 1
    this.selectedCpair = this.selectedCpairTarget.value
    this.fetchExchange(opt)
  }

  NumberOfRowsChanged () {
    this.nextPage = 1
    this.selectedCpair = this.selectedCpairTarget.value
    this.fetchExchange(opt)
  }

  fetchExchange (display) {
    const _this = this
    var url
    var selectedFilter
    if (display === 'table') {
      const numberOfRows = this.selectedNumTarget.value
      selectedFilter = this.selectedFilterTarget.value

      url = `/filteredEx?page=${this.nextPage}&filter=${selectedFilter}&recordsPerPage=${numberOfRows}&selectedCpair=${this.selectedCpair}`
    } else {
      url = `/chartExchange?selectedDtick=${this.selectedDtick}&selectedCpair=${this.selectedCpair}&selectedInterval=${this.selectedInterval}`
    }

    axios.get(url)
      .then(function (response) {
        let result = response.data
        if (display === 'table') {
          if (_this.selectedFilterTarget.value !== selectedFilter) {
            return
          }

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
          _this.previousPageButtonTarget.setAttribute('data-previous-page', `${result.previousPage}`)
          _this.nextPageButtonTarget.setAttribute('data-next-page', `${result.nextPage}`)

          _this.displayExchange(result.exData)
        } else {
          console.log(result)
          _this.plotGraph(result.chartData)
        }
      }).catch(function (e) {
        console.log(e)
      })
  }

  displayExchange (exs) {
    const _this = this
    this.exchangeTableTarget.innerHTML = ''

    exs.forEach(ex => {
      const exRow = document.importNode(_this.exRowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerHTML = ex.time
      fields[1].innerText = ex.exchange_name
      fields[2].innerText = ex.high
      fields[3].innerText = ex.low
      fields[4].innerHTML = ex.open
      fields[5].innerHTML = ex.close
      fields[6].innerHTML = ex.volume
      fields[7].innerText = ex.interval
      fields[8].innerHTML = ex.currency_pair

      _this.exchangeTableTarget.appendChild(exRow)
    })
  }

  // exchange chart
  plotGraph (exs) {
    var extra = {
      legendFormatter: legendFormatter,
      labelsDiv: this.labelsTarget,
      ylabel: 'Price',
      labels: ['Date', 'bittrex', 'binance', 'bleutrade', 'poloniex'],
      colors: ['#2971FF', '#00FF30', '#8F00FF', '#ff1212', '#8ff090']
    }

    var data = [0, 0, 0, 0, 0]
    var dataSet = []

    const _this = this
    exs.forEach(ex => {
      data[0] = new Date(ex.time)
      data.splice(ex.exchange_id, 1, ex.filter)

      dataSet.push(data)
      data = [0, 0, 0, 0, 0]
    })

    var hash = {}
    var i, j,
      result,
      item,
      key

    for (i = 0; i < dataSet.length; i++) {
      item = dataSet[i]
      key = item[0].toString()
      if (!hash[key]) {
        hash[key] = item.slice()
        continue
      }
      for (j = 1; j < item.length; j++) hash[key][j] += item[j]
    }

    result = Object.values(hash)

    console.log(result)
    _this.chartsView = new Dygraph(
      _this.chartsViewTarget,
      result, { ...options, ...extra }
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
