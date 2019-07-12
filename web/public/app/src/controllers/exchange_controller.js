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
    this.nextPage = 1
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
    var numberOfRows
    if (display === 'chart') {
      numberOfRows = 3000
    } else {
      numberOfRows = this.selectedNumTarget.value
    }
    const selectedFilter = this.selectedFilterTarget.value
    const selectedCpair = this.selectedCpairTarget.value

    const _this = this
    axios.get(`/filteredEx?page=${this.nextPage}&filter=${selectedFilter}&recordsPerPage=${numberOfRows}&selectedCpair=${selectedCpair}`)
      .then(function (response) {
      // since results are appended to the table, discard this response
      // if the user has changed the filter before the result is gotten
        if (_this.selectedFilterTarget.value !== selectedFilter) {
          return
        }

        let result = response.data
        _this.totalPageCountTarget.textContent = result.totalPages
        _this.currentPageTarget.textContent = result.currentPage
        _this.previousPageButtonTarget.setAttribute('data-next-page', `${result.previousPage}`)
        _this.nextPageButtonTarget.setAttribute('data-next-page', `${result.nextPage}`)
        _this.nextPageButtonTarget.setAttribute('data-total-page', `${result.totalPages}`)

        if (display === 'table') {
          _this.displayExchange(result.exData)
        } else {
          _this.plotGraph(result.exData)
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
      labels: ['Date', 'interval', 'high', 'low', 'open', 'close', 'volume'],
      colors: ['#2971FF', '#FF8C00', '#006ed0', '#ff0090', '#8ff090', '#ffee90', '#dab390']
    }

    const _this = this

    var data = []
    var dataSet = []
    exs.forEach(ex => {
      data.push(new Date(ex.time))
      data.push(ex.interval)
      data.push(ex.high)
      data.push(ex.low)
      data.push(ex.open)
      data.push(ex.close)
      data.push(ex.volume)

      dataSet.push(data)
      data = []
    })

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
