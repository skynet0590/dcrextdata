import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, options } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')
var opt = 'table'

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter', 'vspTicksTable', 'numPageWrapper',
      'previousPageButton', 'totalPageCount', 'nextPageButton',
      'vspRowTemplate', 'currentPage', 'selectedNum', 'vspTableWrapper',
      'graphTypeWrapper', 'graphType',
      'chartWrapper', 'labels', 'chartsView', 'viewOption'
    ]
  }

  initialize () {
    this.currentPage = parseInt(this.currentPageTarget.getAttribute('data-current-page'))
    if (this.currentPage < 1) {
      this.currentPage = 1
    }
  }

  setTable () {
    this.opt = 'table'
    this.setActiveOptionBtn(this.opt, this.viewOptionTargets)
    hide(this.chartWrapperTarget)
    hide(this.graphTypeWrapperTarget)
    show(this.vspTableWrapperTarget)
    show(this.numPageWrapperTarget)
    this.vspTicksTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange('table')
  }

  setChart () {
    this.opt = 'chart'
    hide(this.numPageWrapperTarget)
    hide(this.vspTableWrapperTarget)
    show(this.graphTypeWrapperTarget)
    show(this.chartWrapperTarget)
    this.setActiveOptionBtn(this.opt, this.viewOptionTargets)
    this.nextPage = 1
    if (this.selectedFilterTarget.selectedIndex === 0) {
      this.selectedFilterTarget.value = this.selectedFilterTarget.options[1].text
    }
    this.fetchExchange('chart')
  }

  loadPreviousPage () {
    this.nextPage = this.previousPageButtonTarget.getAttribute('data-next-page')
    this.fetchExchange(opt)
  }

  loadNextPage () {
    this.nextPage = this.nextPageButtonTarget.getAttribute('data-next-page')
    this.totalPages = (this.nextPageButtonTarget.getAttribute('data-total-page'))
    this.fetchExchange(opt)
  }

  selectedFilterChanged () {
    if (this.opt === 'table') {
      this.nextPage = 1
      this.fetchExchange(opt)
    } else {
      if (this.selectedFilterTarget.value === 'All') {
        this.selectedFilterTarget.value = this.selectedFilterTarget.options[1].text
      }
      this.fetchDataAndGraph()
    }
  }

  numberOfRowsChanged () {
    this.nextPage = 1
    this.fetchExchange(opt)
  }

  fetchExchange (display) {
    const selectedFilter = this.selectedFilterTarget.value
    var numberOfRows
    if (display === 'chart') {
      this.fetchDataAndGraph()
      return
    } else {
      numberOfRows = this.selectedNumTarget.value
    }

    const _this = this
    axios.get(`/filteredvspticks?page=${this.nextPage}&filter=${selectedFilter}&recordsPerPage=${numberOfRows}`)
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

          _this.displayVSPs(result.vspData)
        } else {
          _this.plotGraph(result.vspData)
        }
      }).catch(function (e) {
        console.log(e)
      })
  }

  displayVSPs (vsps) {
    const _this = this
    this.vspTicksTableTarget.innerHTML = ''

    vsps.forEach(vsp => {
      const vspRow = document.importNode(_this.vspRowTemplateTarget.content, true)
      const fields = vspRow.querySelectorAll('td')

      fields[0].innerText = vsp.vsp
      fields[1].innerText = vsp.immature
      fields[2].innerText = vsp.live
      fields[3].innerHTML = vsp.voted
      fields[4].innerHTML = vsp.missed
      fields[5].innerHTML = vsp.pool_fees
      fields[6].innerText = vsp.proportion_live
      fields[7].innerHTML = vsp.proportion_missed
      fields[8].innerHTML = vsp.user_count
      fields[9].innerHTML = vsp.users_active
      fields[10].innerHTML = vsp.time

      _this.vspTicksTableTarget.appendChild(vspRow)
    })
  }

  fetchDataAndGraph () {
    let _this = this
    let url = `/vspchartdata?selectedAttribute=${this.graphTypeTarget.value}&sources=${this.selectedFilterTarget.value}`
    axios.get(url).then(function (response) {
      _this.plotGraph(response.data)
    })
  }

  graphTypeChanged () {
    this.fetchDataAndGraph()
  }

  // vsp chart
  plotGraph (dataSet) {
    dataSet = Object.values(dataSet)
    for (let i = 0; i < dataSet.length; i++) {
      if (!Array.isArray(dataSet[i])) continue
      for (let j = 0; j < dataSet[i].length; j++) {
        if (j === 0) {
          dataSet[i][j] = new Date(dataSet[i][j])
        } else if (!isNaN(dataSet[i][j])) {
          dataSet[i][j] = parseFloat(dataSet[i][j])
        }
      }
    }

    let labels = ['Date', this.selectedFilterTarget.value]
    let colors = ['#007BFF']

    let yLabel = this.graphTypeTarget.value.split('_').join(' ')
    var extra = {
      legendFormatter: legendFormatter,
      labelsDiv: this.labelsTarget,
      ylabel: yLabel,
      sigFigs: 8,
      maxNumberWidth: 8,
      labels: labels,
      colors: colors
    }

    const _this = this

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
