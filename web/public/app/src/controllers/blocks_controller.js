import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, isHidden, show } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'nextPageButton', 'previousPageButton',
      'table', 'votesTbody', 'blockTbody', 'blockTbodyTemplate', 'votesTbodyTemplate', 'voteRowTemplate',
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
    const tableHeadHtml = this.tableTarget.querySelector('thead').innerHTML

    let blocksHtml = ''
    data.forEach(block => {
      let votesHtml = ''
      let i = 0
      block.votes.forEach(vote => {
        votesHtml += `<tr>
                            <td>${vote.voting_on}</td>
                            <td>${vote.validator_id}</td>
                            <td>${vote.receive_time}</td>
                            <td>${vote.block_time_diff}</td>
                            <td>${vote.block_receive_time_diff}</td>
                            <td><a target="_blank" href="https://explorer.dcrdata.org/tx/${vote.hash}">${vote.hash}</a></td>
                        </tr>`
      })

      let padding = i > 0 ? 'style="padding-top:50px"' : ''
      i++
      blocksHtml += `<tbody data-target="blocks.blockTbody"
                            data-block-hash="${block.block_hash}" class="clickable">
                        <tr>
                            <td colspan="6" ${padding}>
                              Height: ${block.block_height}
                              Timestamp: ${block.block_internal_time}
                              Received: ${block.block_receive_time}
                              Delay: ${block.delay}
                              Hash: <a target="_blank" href="https://explorer.dcrdata.org/block/${block.block_height}">${block.block_hash}</a>
                            </td>
                        </tr>
                        </tbody>
                        <tbody data-target="blocks.votesTbody" data-block-hash="${block.block_hash}" style="margin-bottom: 20px;">
                        <tr>
                            <th>Voting On</th>
                            <th>Validator ID</th>
                            <th>Received</th>
                            <th>Block Time Diff</th>
                            <th>Block Receive Time Diff</th>
                            <th>Hash</th>
                        </tr>
                        ${votesHtml}
                        </tbody>`
    })

    this.tableTarget.innerHTML = `${tableHeadHtml} ${blocksHtml}`
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
