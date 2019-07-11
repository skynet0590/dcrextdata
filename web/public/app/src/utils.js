export const hide = (el) => {
  el.classList.add('d-none')
}

export const show = (el) => {
  el.classList.remove('d-none')
}

export const isHidden = (el) => {
  return el.classList.contains('d-none')
}

export function date (date) {
  var d = new Date(date)
  return `${String(d.getUTCFullYear())}-${String(d.getUTCMonth() + 1).padStart(2, '0')}-${String(d.getUTCDate()).padStart(2, '0')} ` +
          `${String(d.getUTCHours()).padStart(2, '0')}:${String(d.getUTCMinutes()).padStart(2, '0')}:${String(d.getUTCSeconds()).padStart(2, '0')} +0000 UTC`
}
