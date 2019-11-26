import { Controller } from 'stimulus'
export default class extends Controller {
  static get targets () {
    return [ 'nextPageButton', 'previousPageButton', 'tableBody', 'rowTemplate', 'totalPageCount', 'currentPage', 'loadingData' ]
  }
}
