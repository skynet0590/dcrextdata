import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'nextPageButton', 'previousPageButton', 'tableBody', 'rowTemplate',
      'totalPageCount', 'currentPage'
    ]
  }

  initialize () {
    this.currentPage = parseInt(this.currentPageTarget.getAttribute('data-current-page'))
    if (this.currentPage < 1) {
      this.currentPage = 1
    }
  }

  gotoPreviousPage () {
    this.fetchData(this.currentPage - 1)
  }

  gotoNextPage () {
    this.fetchData(this.currentPage + 1)
  }

  fetchData (page) {
    const _this = this
    axios.get(`/getvotes?page=${page}`).then(function (response) {
      let result = response.data
      _this.totalPageCountTarget.textContent = result.totalPages
      _this.currentPageTarget.textContent = result.currentPage

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

      _this.displayVotes(result.records)
    }).catch(function (e) {
      console.log(e) // todo: handle error
    })
  }

  displayVotes (data) {
    const _this = this
    this.tableBodyTarget.innerHTML = ''

    data.forEach(item => {
      const exRow = document.importNode(_this.rowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerText = item.receive_time
      fields[1].innerText = item.hash
      fields[2].innerText = item.voting_on
      fields[3].innerText = item.validator_id

      _this.tableBodyTarget.appendChild(exRow)
    })
  }
}
