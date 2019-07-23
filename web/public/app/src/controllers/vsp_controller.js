import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, legendFormatter, barChartPlotter } from '../utils'

const Dygraph = require('../../../dist/js/dygraphs.min.js')

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
    this.viewOption = 'table'
  }

  setTable () {
    this.viewOption = 'table'
    this.setActiveOptionBtn(this.viewOption, this.viewOptionTargets)
    hide(this.chartWrapperTarget)
    hide(this.graphTypeWrapperTarget)
    show(this.vspTableWrapperTarget)
    show(this.numPageWrapperTarget)
    this.vspTicksTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange('table')
  }

  setChart () {
    this.viewOption = 'chart'
    hide(this.numPageWrapperTarget)
    hide(this.vspTableWrapperTarget)
    show(this.graphTypeWrapperTarget)
    show(this.chartWrapperTarget)
    this.setActiveOptionBtn(this.viewOption, this.viewOptionTargets)
    this.nextPage = 1
    if (this.selectedFilterTarget.selectedIndex === 0) {
      this.selectedFilterTarget.selectedIndex = 1
    }
    this.fetchExchange('chart')
  }

  loadPreviousPage () {
    this.nextPage = this.previousPageButtonTarget.getAttribute('data-next-page')
    this.fetchExchange(this.viewOption)
  }

  loadNextPage () {
    this.nextPage = this.nextPageButtonTarget.getAttribute('data-next-page')
    this.totalPages = (this.nextPageButtonTarget.getAttribute('data-total-page'))
    this.fetchExchange(this.viewOption)
  }

  selectedFilterChanged () {
    if (this.viewOption === 'table') {
      this.nextPage = 1
      this.fetchExchange(this.viewOption)
    } else {
      if (this.selectedFilterTarget.selectedIndex === 0) {
        this.selectedFilterTarget.selectedIndex = 1
      }
      this.fetchDataAndGraph()
    }
  }

  numberOfRowsChanged () {
    this.nextPage = 1
    this.fetchExchange(this.viewOption)
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
    let yLabel = this.graphTypeTarget.value.split('_').join(' ')
    // let labels = ['Date', this.selectedFilterTarget.value]
    // let colors = ['#007BFF']

    let csv = `Date,${yLabel}\n`
    let minDate, maxDate

    for (let i = 0; i < dataSet.length; i++) {
      if (!Array.isArray(dataSet[i])) continue
      for (let j = 0; j < dataSet[i].length; j++) {
        if (j === 0) {
          const date = new Date(dataSet[i][j])
          if (minDate === undefined || date < minDate) {
            minDate = date
          }
          if (maxDate === undefined || date > maxDate) {
            maxDate = date
          }
          csv += `${date},`
        } else if (!isNaN(dataSet[i][j])) {
          csv += `${dataSet[i][j]}\n`
        }
      }
    }
    console.log(minDate, maxDate)
    const _this = this
    _this.chartsView = new Dygraph(
      _this.chartsViewTarget,
      csv,
      {
        legend: 'always',
        // title: title,
        includeZero: true,
        dateWindow: [minDate, maxDate],
        animatedZooms: true,
        legendFormatter: legendFormatter,
        plotter: barChartPlotter,
        labelsDiv: _this.labelsTarget,
        ylabel: yLabel,
        xlabel: 'Date',
        labelsUTC: true,
        labelsKMB: true,
        maxNumberWidth: 10,
        axes: {
          x: {
            drawGrid: false
          },
          y: {
            axisLabelWidth: 90
          }
        }
      }
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
