import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter', 'exchangeTable',
      'previousPageButton', 'pageReport', 'nextPageButton',
      'exRowTemplate'
    ]
  }

  initialize () {
    // hide next page button to use infinite scroll
    // hide(this.previousPageButtonTarget)
    // hide(this.nextPageButtonTarget)
    // hide(this.pageReportTarget)

    this.nextPage = this.nextPageButtonTarget.getAttribute('data-next-page')
    this.selectedFilter = this.nextPageButtonTarget.getAttribute('data-filter')

    if (this.nextPage) {
      // check if there is space at the bottom to load more now
      this.fetchExchange()
    }
  }

  selectedFilterChanged () {
    this.exchangeTableTarget.innerHTML = ''
    this.nextPage = 1
    this.fetchExchange()
  }

  fetchExchange () {
    const selectedFilter = this.selectedFilterTarget.value

    const _this = this
    axios.get(`/filteredEx?page=${this.nextPage}&filter=${selectedFilter}`)
      .then(function (response) {
      // since results are appended to the table, discard this response
      // if the user has changed the filter before the result is gotten
        if (_this.selectedFilterTarget.value !== selectedFilter) {
          return
        }
        console.log(response.data)

        let result = response.data
        // _this.nextPageButtonTarget = result.selectedFilter
        // _this.previousPageButtonTarget.textContent = result.selectedFilter
        _this.nextPage = result.nextPage
        _this.selectedFilter = result.selectedFilter
        if (result.nextPage >= 2) {
          _this.previousPageButtonTarget.setAttribute('href', `?page=${result.nextPage - 1}&exchange=${_this.selectedFilterTarget.value}`)
        } else {
          _this.previousPageButtonTarget.setAttribute('href', `?page=1&exchange=${_this.selectedFilterTarget.value}`)
        }
        _this.nextPageButtonTarget.setAttribute('href', `?page=${result.nextPage}&exchange=${_this.selectedFilterTarget.value}`)
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
      fields[8].innerHTML = ex.time

      _this.exchangeTableTarget.appendChild(exRow)
    })
  }
}
