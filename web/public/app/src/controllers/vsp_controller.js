import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, date } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter', 'vspTicksTable',
      'previousPageButton', 'totalPageCount', 'nextPageButton',
      'vspRowTemplate', 'currentPage', 'selectedNum'
    ]
  }

  loadPreviousPage () {
    this.nextPage = this.previousPageButtonTarget.getAttribute('data-next-page')
    if (this.nextPage <= 1) {
      hide(this.previousPageButtonTarget)
    }
    this.vspTicksTableTarget.innerHTML = ''
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
    this.vspTicksTableTarget.innerHTML = ''
    this.fetchExchange()
  }

  selectedFilterChanged () {
    this.vspTicksTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange()
  }

  NumberOfRowsChanged () {
    this.vspTicksTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange()
  }

  fetchExchange () {
    const selectedFilter = this.selectedFilterTarget.value
    const numberOfRows = this.selectedNumTarget.value

    const _this = this
    axios.get(`/filteredvspticks?page=${this.nextPage}&filter=${selectedFilter}&recordsPerPage=${numberOfRows}`)
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
        _this.displayVSPs(result.vspData)
      }).catch(function (e) {
        console.log(e)
      })
  }

  displayVSPs (vsps) {
    const _this = this

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
      fields[10].innerHTML = date(vsp.time)

      _this.vspTicksTableTarget.appendChild(vspRow)
    })
  }
}
