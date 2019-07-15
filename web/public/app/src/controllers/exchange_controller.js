import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, date, legendFormatter, options } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')
var opt = 'table'

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter', 'exchangeTable', 'selectedCpair', 'numPageWrapper',
      'previousPageButton', 'totalPageCount', 'nextPageButton',
      'exRowTemplate', 'currentPage', 'selectedNum', 'exchangeTableWrapper',
      'chartWrapper', 'labels', 'chartsView', 'viewOption'
    ]
  }

  setTable () {
    opt = 'table'
    this.setActiveOptionBtn(opt, this.viewOptionTargets)
    this.chartWrapperTarget.classList.add('d-hide')
    this.exchangeTableWrapperTarget.classList.remove('d-hide')
    this.numPageWrapperTarget.classList.remove('d-hide')
    this.exchangeTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange('table')
  }

  setChart () {
    opt = 'chart'
    this.numPageWrapperTarget.classList.add('d-hide')
    this.exchangeTableWrapperTarget.classList.add('d-hide')
    this.setActiveOptionBtn(opt, this.viewOptionTargets)
    this.chartWrapperTarget.classList.remove('d-hide')
    this.filter = 'close'
    this.fetchExchange('chart')
  }

  loadPreviousPage () {
    this.nextPage = this.previousPageButtonTarget.getAttribute('data-next-page')
    if (this.nextPage <= 1) {
      hide(this.previousPageButtonTarget)
    }
    this.fetchExchange(opt)
  }

  loadNextPage () {
    this.nextPage = this.nextPageButtonTarget.getAttribute('data-next-page')
    this.totalPages = this.nextPageButtonTarget.getAttribute('data-total-page')
    if (this.nextPage > 1) {
      show(this.previousPageButtonTarget)
    }
    if (this.totalPages === this.nextPage) {
      hide(this.nextPageButtonTarget)
    }
    this.fetchExchange(opt)
  }

  selectedFilterChanged () {
    this.nextPage = 1
    this.fetchExchange(opt)
  }

  selectedCpairChanged () {
    this.nextPage = 1
    this.fetchExchange(opt)
  }

  NumberOfRowsChanged () {
    this.nextPage = 1
    this.fetchExchange(opt)
  }

  fetchExchange (display) {
    this.exchangeTableTarget.innerHTML = ''

    const numberOfRows = this.selectedNumTarget.value
    const selectedFilter = this.selectedFilterTarget.value
    const selectedCpair = this.selectedCpairTarget.value

    const _this = this
    var url
    if (display === 'table') {
      url = `/filteredEx?page=${this.nextPage}&filter=${selectedFilter}&recordsPerPage=${numberOfRows}&selectedCpair=${selectedCpair}`
    } else {
      url = `/chartExchange?filter=${this.filter}`
    }
    axios.get(url)
      .then(function (response) {
      // since results are appended to the table, discard this response
      // if the user has changed the filter before the result is gotten
        if (_this.selectedFilterTarget.value !== selectedFilter) {
          return
        }

        let result = response.data
        if (display === 'table') {
          console.log(result.exData)
          _this.totalPageCountTarget.textContent = result.totalPages
          _this.currentPageTarget.textContent = result.currentPage
          _this.previousPageButtonTarget.setAttribute('data-next-page', `${result.previousPage}`)
          _this.nextPageButtonTarget.setAttribute('data-next-page', `${result.nextPage}`)
          _this.nextPageButtonTarget.setAttribute('data-total-page', `${result.totalPages}`)

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

    exs.forEach(ex => {
      const exRow = document.importNode(_this.exRowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerText = ex.exchange_name
      fields[1].innerText = ex.high
      fields[2].innerText = ex.low
      fields[3].innerHTML = ex.open
      fields[4].innerHTML = ex.close
      fields[5].innerHTML = ex.volume
      fields[6].innerText = ex.interval
      fields[7].innerHTML = ex.currency_pair
      fields[8].innerHTML = date(ex.time)

      _this.exchangeTableTarget.appendChild(exRow)
    })
  }

  // exchange chart
  plotGraph (exs) {
    var extra = {
      legendFormatter: legendFormatter,
      labelsDiv: this.labelsTarget,
      ylabel: 'Interval',
      y2label: 'Volume',
      labels: ['Date', 'bittrex', 'binance', 'bleutrade', 'poloniex'],
      colors: ['#2971FF', '#FF8C00', '#006ed0', '#ff0090', '#8ff090']
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
