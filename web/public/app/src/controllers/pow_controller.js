import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, date } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter', 'powTable',
      'previousPageButton', 'totalPageCount', 'nextPageButton',
      'powRowTemplate', 'currentPage', 'selectedNum',
      'chartWrapper', 'labels', 'chartsView', 'viewOption'
    ]
  }

  setTable () {
    var opt = 'table'
    this.setActiveOptionBtn(opt, this.viewOptionTargets)
    this.powTableTarget.innerHTML = ''
    this.selectedFilter = 'All'
    this.selectedNum = 20
    this.nextPage = 1
    this.fetchExchange()
  }

  setChart () {
    var opt = 'chart'
    this.setActiveOptionBtn(opt, this.viewOptionTargets)
    this.nextPage = 1
    this.plotGraph()
  }

  loadPreviousPage () {
    this.nextPage = this.previousPageButtonTarget.getAttribute('data-next-page')
    if (this.nextPage <= 1) {
      hide(this.previousPageButtonTarget)
    }
    this.powTableTarget.innerHTML = ''
    this.fetchExchange()
  }

  loadNextPage () {
    this.nextPage = this.nextPageButtonTarget.getAttribute('data-next-page')
    this.totalPages = (this.nextPageButtonTarget.getAttribute('data-total-page'))
    if (this.nextPage > 1) {
      show(this.previousPageButtonTarget)
    }
    if (this.totalPages === this.nextPage) {
      hide(this.nextPageButtonTarget)
    }
    this.powTableTarget.innerHTML = ''
    this.fetchExchange()
  }

  selectedFilterChanged () {
    this.powTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange()
  }

  NumberOfRowsChanged () {
    this.powTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange()
  }

  fetchExchange () {
    const selectedFilter = this.selectedFilterTarget.value
    const numberOfRows = this.selectedNumTarget.value

    const _this = this
    axios.get(`/filteredpow?page=${this.nextPage}&filter=${selectedFilter}&recordsPerPage=${numberOfRows}`)
      .then(function (response) {
      // since results are appended to the table, discard this response
      // if the user has changed the filter before the result is gotten
        if (_this.selectedFilterTarget.value !== selectedFilter) {
          return
        }

        console.log(response.data)

        let result = response.data
        _this.totalPageCountTarget.textContent = result.totalPages
        _this.currentPageTarget.textContent = result.currentPage
        _this.previousPageButtonTarget.setAttribute('data-next-page', `${result.previousPage}`)
        _this.nextPageButtonTarget.setAttribute('data-next-page', `${result.nextPage}`)
        _this.nextPageButtonTarget.setAttribute('data-total-page', `${result.totalPages}`)
        _this.displayPoW(result.powData)
      }).catch(function (e) {
        console.log(e)
      })
  }

  displayPoW (pows) {
    const _this = this

    pows.forEach(pow => {
      const powRow = document.importNode(_this.powRowTemplateTarget.content, true)
      const fields = powRow.querySelectorAll('td')

      fields[0].innerText = pow.Source
      fields[1].innerText = pow.NetworkHashrate
      fields[2].innerText = pow.PoolHashrate
      fields[3].innerHTML = pow.Workers
      fields[4].innerHTML = pow.NetworkDifficulty
      fields[5].innerHTML = date(pow.Time)

      _this.powTableTarget.appendChild(powRow)
    })
  }

  // pow chart
  plotGraph () {
    var options = {
      axes: { y: { axisLabelWidth: 70 }, y2: { axisLabelWidth: 70 } },
      labels: ['Date', 'Network Difficulty', 'pool hash'],
      digitsAfterDecimal: 2,
      showRangeSelector: true,
      rangeSelectorPlotFillColor: '#8997A5',
      rangeSelectorAlpha: 0.4,
      rangeSelectorHeight: 40,
      drawPoints: true,
      pointSize: 0.25,
      legend: 'always',
      labelsSeparateLines: true,
      highlightCircleSize: 4,
      ylabel: 'hash',
      y2label: 'diff',
      labelsUTC: true
    }

    const _this = this
    axios.get(`/getChartPowData?page=${this.nextPage}`)
      .then(function (response) {
        console.log(response.data)
        let result = response.data

        var data = []
        var dataSet = []
        result.powData.forEach(pow => {
          data.push(new Date(pow.Time))
          data.push(pow.PoolHashrate)
          data.push(pow.NetworkDifficulty)

          dataSet.push(data)
          data = []
        })
        console.log('...java Script Array... \n' + JSON.stringify(dataSet))
        _this.chartsView = new Dygraph(
          _this.chartsViewTarget,
          dataSet, options
        )
      }).catch(function (e) {
        console.log(e)
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
