import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, date } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter', 'exchangeTable', 'selectedCpair',
      'previousPageButton', 'totalPageCount', 'nextPageButton',
      'exRowTemplate', 'currentPage', 'selectedNum'
    ]
  }

  loadPreviousPage () {
    this.nextPage = this.previousPageButtonTarget.getAttribute('data-next-page')
    if (this.nextPage <= 1) {
      hide(this.previousPageButtonTarget)
    }
    this.fetchExchange()
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
    this.fetchExchange()
  }

  selectedFilterChanged () {
    this.nextPage = 1
    this.fetchExchange()
  }

  selectedCpairChanged () {
    this.nextPage = 1
    this.fetchExchange()
  }

  NumberOfRowsChanged () {
    this.nextPage = 1
    this.fetchExchange()
  }

  fetchExchange () {
    this.exchangeTableTarget.innerHTML = ''
    const selectedFilter = this.selectedFilterTarget.value
    const numberOfRows = this.selectedNumTarget.value
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
        _this.displayExchange(result.exData)
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
}
