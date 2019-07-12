import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, date } from '../utils'
import dompurify from 'dompurify'

const Dygraph = require('../../../dist/js/dygraphs.min.js')
var opt = 'table'

function legendFormatter (data) {
  var html = ''
  if (data.x == null) {
    let dashLabels = data.series.reduce((nodes, series) => {
      return `${nodes} <div class="pr-2">${series.dashHTML} ${series.labelHTML}</div>`
    }, '')
    html = `<div class="d-flex flex-wrap justify-content-center align-items-center">
              <div class="pr-3">${this.getLabels()[0]}: N/A</div>
              <div class="d-flex flex-wrap">${dashLabels}</div>
            </div>`
  } else {
    data.series.sort((a, b) => a.y > b.y ? -1 : 1)
    var extraHTML = ''
    // The circulation chart has an additional legend entry showing percent
    // difference.
    if (data.series.length === 2 && data.series[1].label.toLowerCase() === 'coin supply') {
      let predicted = data.series[0].y
      let actual = data.series[1].y
      let change = (((actual - predicted) / predicted) * 100).toFixed(2)
      extraHTML = `<div class="pr-2">&nbsp;&nbsp;Change: ${change} %</div>`
    }

    let yVals = data.series.reduce((nodes, series) => {
      if (!series.isVisible) return nodes
      let yVal = series.yHTML
      yVal = series.y

      return `${nodes} <div class="pr-2">${series.dashHTML} ${series.labelHTML}: ${yVal}</div>`
    }, '')

    html = `<div class="d-flex flex-wrap justify-content-center align-items-center">
                <div class="pr-3">${this.getLabels()[0]}: ${data.xHTML}</div>
                <div class="d-flex flex-wrap"> ${yVals}</div>
            </div>${extraHTML}`
  }

  dompurify.sanitize(html)
  return html
}

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter', 'powTable', 'numPageWrapper',
      'previousPageButton', 'totalPageCount', 'nextPageButton',
      'powRowTemplate', 'currentPage', 'selectedNum', 'powTableWrapper',
      'chartWrapper', 'labels', 'chartsView', 'viewOption', 'chartWrapper'
    ]
  }

  setTable () {
    opt = 'table'
    this.setActiveOptionBtn(opt, this.viewOptionTargets)
    this.chartWrapperTarget.classList.add('d-hide')
    this.powTableWrapperTarget.classList.remove('d-hide')
    this.numPageWrapperTarget.classList.remove('d-hide')
    this.powTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange('table')
  }

  setChart () {
    opt = 'chart'
    this.numPageWrapperTarget.classList.add('d-hide')
    this.powTableWrapperTarget.classList.add('d-hide')
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
    this.powTableTarget.innerHTML = ''
    this.fetchExchange('table')
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
    this.fetchExchange('table')
  }

  selectedFilterChanged () {
    this.powTableTarget.innerHTML = ''
    this.nextPage = 1
    console.log(opt)
    console.log(this.opt)
    this.fetchExchange(opt)
  }

  NumberOfRowsChanged () {
    this.powTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange('table')
  }

  fetchExchange (display) {
    const selectedFilter = this.selectedFilterTarget.value
    var numberOfRows
    if (display === 'chart') {
      numberOfRows = 3000
    } else {
      numberOfRows = this.selectedNumTarget.value
    }

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

        if (display === 'table') {
          _this.displayPoW(result.powData)
        } else {
          _this.plotGraph(result.powData)
        }
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
  plotGraph (pows) {
    var options = {
      axes: { y: { axisLabelWidth: 70 }, y2: { axisLabelWidth: 70 } },
      axisLabelFontSize: 12,
      labels: ['Date', 'Pool Hashrate', 'Network Difficulty', 'Workers', 'Network Hashrate'],
      colors: ['#2971FF', '#FF8C00', '#006600', '#ff0090'],
      digitsAfterDecimal: 3,
      retainDateWindow: false,
      showRangeSelector: true,
      rangeSelectorHeight: 40,
      drawPoints: true,
      sigFigs: 1,
      legendFormatter: legendFormatter,
      pointSize: 0.25,
      legend: 'always',
      labelsDiv: this.labelsTarget,
      labelsSeparateLines: true,
      highlightCircleSize: 4,
      ylabel: 'Pool Hashrate',
      y2label: 'Network Difficulty',
      yLabelWidth: 20
    }

    const _this = this

    var data = []
    var dataSet = []
    pows.forEach(pow => {
      data.push(new Date(pow.Time))
      data.push(pow.PoolHashrate)
      data.push(pow.NetworkDifficulty)
      data.push(pow.Workers)
      data.push(pow.NetworkHashrate)

      dataSet.push(data)
      data = []
    })

    _this.chartsView = new Dygraph(
      _this.chartsViewTarget,
      dataSet, options
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
