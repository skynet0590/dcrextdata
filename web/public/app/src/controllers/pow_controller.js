import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, options } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter', 'powTable', 'numPageWrapper',
      'previousPageButton', 'totalPageCount', 'nextPageButton',
      'powRowTemplate', 'currentPage', 'selectedNum', 'powTableWrapper',
      'chartWrapper', 'chartDataTypeSelector', 'chartDataType', 'labels',
      'chartsView', 'viewOption', 'pageSizeWrapper'
    ]
  }

  initialize () {
    this.viewOption = 'table'
    this.dataType = 'hashrate'
  }

  connect () {
    var filter = this.selectedFilterTarget.options
    var num = this.selectedNumTarget.options
    this.selectedFilterTarget.value = filter[0].text
    this.selectedNumTarget.value = num[0].text
  }

  setTable () {
    this.viewOption = 'table'
    this.setActiveOptionBtn(this.viewOption, this.viewOptionTargets)
    hide(this.chartWrapperTarget)
    show(this.powTableWrapperTarget)
    show(this.numPageWrapperTarget)
    show(this.pageSizeWrapperTarget)
    hide(this.chartDataTypeSelectorTarget)
    this.nextPage = 1
    this.fetchExchange('table')
  }

  setChart () {
    this.viewOption = 'chart'
    this.setActiveOptionBtn(this.viewOption, this.viewOptionTargets)
    hide(this.numPageWrapperTarget)
    hide(this.powTableWrapperTarget)
    show(this.chartWrapperTarget)
    hide(this.pageSizeWrapperTarget)
    show(this.chartDataTypeSelectorTarget)
    // hide(this.numPageWrapperTarget)
    this.nextPage = 1
    this.fetchExchange('chart')
  }

  setHashrateDataType (event) {
    this.dataType = 'hashrate'
    this.chartDataTypeTargets.forEach(el => {
      el.classList.remove('active')
    })
    event.currentTarget.classList.add('active')
    this.fetchExchange('chart')
  }

  setWorkersDataType (event) {
    this.dataType = 'workers'
    this.chartDataTypeTargets.forEach(el => {
      el.classList.remove('active')
    })
    event.currentTarget.classList.add('active')
    this.fetchExchange('chart')
  }

  loadPreviousPage () {
    this.nextPage = this.previousPageButtonTarget.getAttribute('data-next-page')
    this.fetchExchange(this.viewOption)
  }

  loadNextPage () {
    this.nextPage = this.nextPageButtonTarget.getAttribute('data-next-page')
    this.fetchExchange(this.viewOption)
  }

  selectedFilterChanged () {
    this.nextPage = 1
    this.fetchExchange(this.viewOption)
  }

  numberOfRowsChanged () {
    this.nextPage = 1
    this.fetchExchange(this.viewOption)
  }

  fetchExchange (display) {
    const selectedFilter = this.selectedFilterTarget.value
    var numberOfRows = this.selectedNumTarget.value

    const _this = this
    axios.get(`/filteredpow?page=${this.nextPage}&filter=${selectedFilter}&recordsPerPage=${numberOfRows}`)
      .then(function (response) {
        let result = response.data

        if (display === 'table') {
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
        } else {
          _this.plotGraph(result.powData)
        }
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

  // pow chart
  plotGraph (pows) {
    const _this = this

    var data = []
    var dataSet = []
    pows.forEach(pow => {
      data.push(new Date(pow.time))

      if (_this.dataType === 'hashrate') {
        data.push(parseInt(pow.pool_hashrate_th))
      } else {
        data.push(parseInt(pow.workers))
      }

      dataSet.push(data)
      data = []
    })

    let dataTypeLabel = 'Pool Hashrate'
    if (_this.dataType === 'workers') {
      dataTypeLabel = 'Workers'
    }

    var extra = {
      labels: ['Date', dataTypeLabel],
      colors: ['#2971FF', '#FF8C00'],
      labelsDiv: this.labelsTarget,
      ylabel: dataTypeLabel,
      y2label: 'Network Difficulty',
      sigFigs: 1,
      legendFormatter: legendFormatter
    }

    _this.chartsView = new Dygraph(
      _this.chartsViewTarget,
      dataSet, { ...options, ...extra }
    )
  }

  setActiveOptionBtn (opt, optTargets) {
    optTargets.forEach(li => {
      if (li.dataset.option === this.viewOption) {
        li.classList.add('active')
      } else {
        li.classList.remove('active')
      }
    })
  }
}
