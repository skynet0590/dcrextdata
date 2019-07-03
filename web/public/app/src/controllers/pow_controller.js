import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, date } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter', 'powTable',
      'previousPageButton', 'totalPageCount', 'nextPageButton',
      'powRowTemplate', 'currentPage', 'selectedNum'
    ]
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
}
