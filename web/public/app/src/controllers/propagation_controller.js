import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, isHidden, show } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'nextPageButton', 'previousPageButton',
      'selectedRecordSet',
      'table', 'blocksTbody', 'votesTbody',
      'blocksTable', 'blocksTableBody', 'blocksRowTemplate', 'votesTable', 'votesTableBody', 'votesRowTemplate',
      'totalPageCount', 'currentPage'
    ]
  }

  initialize () {
    this.currentPage = parseInt(this.currentPageTarget.getAttribute('data-current-page'))
    if (this.currentPage < 1) {
      this.currentPage = 1
    }
    this.selectedRecordSet = 'both'
  }

  selectedRecordSetChanged () {
    this.currentPage = 1
    this.selectedRecordSet = this.selectedRecordSetTarget.value
    this.fetchData(1)
  }

  gotoPreviousPage () {
    this.fetchData(this.currentPage - 1)
  }

  gotoNextPage () {
    this.fetchData(this.currentPage + 1)
  }

  fetchData (page) {
    const _this = this
    let uri = '/getpropagationdata'
    switch (this.selectedRecordSet) {
      case 'blocks':
        uri = 'getblocks'
        break
      case 'votes':
        uri = 'getvotes'
        break
      default:
        uri = 'getpropagationdata'
        break
    }
    axios.get(`/${uri}?page=${page}`).then(function (response) {
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

      _this.displayData(result.records)
    }).catch(function (e) {
      console.log(e) // todo: handle error
    })
  }

  displayData (data) {
    switch (this.selectedRecordSet) {
      case 'blocks':
        this.displayBlocks(data)
        break
      case 'votes':
        this.displayVotes(data)
        break
      default:
        this.displayPropagationData(data)
        break
    }
  }

  displayBlocks (data) {
    const _this = this
    this.blocksTableBodyTarget.innerHTML = ''
    data.forEach(block => {
      const exRow = document.importNode(_this.blocksRowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerHTML = `<a target="_blank" href="https://explorer.dcrdata.org/block/${block.block_height}">${block.block_height}</a>`
      fields[1].innerText = block.block_internal_time
      fields[2].innerText = block.block_receive_time
      fields[3].innerText = block.delay
      fields[4].innerHTML = `<a target="_blank" href="https://explorer.dcrdata.org/block/${block.block_height}">${block.block_hash}</a>`

      _this.blocksTableBodyTarget.appendChild(exRow)
    })

    hide(this.tableTarget)
    hide(this.votesTableTarget)
    show(this.blocksTableTarget)
  }

  displayVotes (data) {
    const _this = this
    this.votesTableBodyTarget.innerHTML = ''

    data.forEach(item => {
      const exRow = document.importNode(_this.votesRowTemplateTarget.content, true)
      const fields = exRow.querySelectorAll('td')

      fields[0].innerHTML = `<a target="_blank" href="https://explorer.dcrdata.org/block/${item.voting_on}">${item.voting_on}</a>`
      fields[1].innerText = item.validator_id
      fields[2].innerText = item.validity
      fields[3].innerText = item.receive_time
      fields[4].innerText = item.block_time_diff
      fields[5].innerText = item.block_receive_time_diff
      fields[6].innerHTML = `<a target="_blank" href="https://explorer.dcrdata.org/tx/${item.hash}">${item.hash}</a>`

      _this.votesTableBodyTarget.appendChild(exRow)
    })

    hide(this.tableTarget)
    hide(this.blocksTableTarget)
    show(this.votesTableTarget)
  }

  displayPropagationData (data) {
    let blocksHtml = ''
    data.forEach(block => {
      let votesHtml = ''
      let i = 0
      block.votes.forEach(vote => {
        votesHtml += `<tr>
                            <td>${vote.voting_on}</td>
                            <td>${vote.validator_id}</td>
                            <td>${vote.validity}</td>
                            <td>${vote.receive_time}</td>
                            <td>${vote.block_receive_time_diff}s</td>
                            <td><a target="_blank" href="https://explorer.dcrdata.org/tx/${vote.hash}">${vote.hash}</a></td>
                        </tr>`
      })

      let padding = i > 0 ? 'style="padding-top:50px"' : ''
      i++
      blocksHtml += `<tbody data-target="propagation.blockTbody"
                            data-block-hash="${block.block_hash}">
                        <tr>
                            <td colspan="7" ${padding}>
                              <b>Height</b>: ${block.block_height} &nbsp;&nbsp;&nbsp;&nbsp;
                              <b>Timestamp</b>: ${block.block_internal_time}  &nbsp;&nbsp;&nbsp;&nbsp;
                              <b>Received</b>: ${block.block_receive_time}  &nbsp;&nbsp;&nbsp;&nbsp;
                              <b>Hash</b>: <a target="_blank" href="https://explorer.dcrdata.org/block/${block.block_height}">${block.block_hash}</a>
                            </td>
                        </tr>
                        </tbody>
                        <tbody data-target="propagation.votesTbody" data-block-hash="${block.block_hash}">
                        <tr>
                            <td>Voting On</td>
                            <td>Validator ID</td>
                            <td>Validity</td>
                            <td>Received</td>
                            <td>Block Receive Time Diff</td>
                            <td>Hash</td>
                        </tr>
                        ${votesHtml}
                        </tbody>
                          <tr>
                              <td colspan="7" height="50" style="border: none !important;"></td>
                          </tr>`
    })

    this.tableTarget.innerHTML = blocksHtml

    show(this.tableTarget)
    hide(this.blocksTableTarget)
    hide(this.votesTableTarget)
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
