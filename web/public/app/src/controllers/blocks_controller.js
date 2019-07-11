import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, isHidden, show } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'nextPageButton', 'previousPageButton',
      'table', 'votesTbody', 'blockTbody', 'blockTbodyTemplate', 'votesTbodyTemplate',
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
    axios.get(`/getblocks?page=${page}`).then(function (response) {
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

      _this.displayBlock(result.records)
    }).catch(function (e) {
      console.log(e) // todo: handle error
    })
  }

  displayBlock (data) {
    const _this = this
    const tableHead = this.tableTarget.querySelector('thead')
    this.tableTarget.innerHTML = ''
    this.tableTarget.appendChild(tableHead)

    data.forEach(item => {
      const blockTbodyTemplate = document.importNode(_this.blockTbodyTemplateTarget.content, true)
      const fields = blockTbodyTemplate.querySelectorAll('td')

      fields[0].innerText = item.block_receive_time
      fields[1].innerText = item.block_internal_time
      fields[2].innerText = item.delay
      fields[3].innerText = item.block_height
      fields[4].innerText = item.block_hash

      // TODO: set the data-hash and populate the vote
      blockTbodyTemplate.firstElementChild.setAttribute('data-block-hash', item.block_hash)
      _this.tableTarget.appendChild(blockTbodyTemplate)

      const votesTbody = document.importNode(_this.votesTbodyTemplateTarget.content, true)
      votesTbody.firstElementChild.setAttribute('data-block-hash', item.block_hash)
      // add actual data
      _this.tableTarget.appendChild(votesTbody)
    })
  }

  showVotes (event) {
    const blockHash = event.currentTarget.getAttribute('data-block-hash')
    this.blockTbodyTargets.forEach(el => {
      el.classList.remove('labels')
    })
    this.votesTbodyTargets.forEach(el => {
      if (el.getAttribute('data-block-hash') === blockHash) {
        if (isHidden(el)) {
          show(el)
          event.currentTarget.classList.add('labels')
        } else {
          hide(el)
        }
        return
      }
      hide(el)
    })
  }
}
