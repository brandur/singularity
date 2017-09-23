/**
 * Check if `el` is out out of view.
 */

function isBelowScroll(el) {
  return el.getBoundingClientRect().top > 0
}

/**
 * Activate the correct menu item for the
 * contents in the viewport.
 */

function activate() {
  const headers = document.querySelectorAll('h2')
  let i = 0

  for (; i < headers.length; i++) {
    if (isBelowScroll(headers[i])) {
      break
    }
  }

  activateTOCItem(i-1)
}

function activateTOCItem(i) {
  const tocItems = document.querySelectorAll('.toc > ol > li')
  tocItems.forEach(e => e.classList.add('collapsed'))

  if (i == -1) {
    return
  }

  tocItems[i].classList.remove('collapsed')
}

window.addEventListener('load', e => activateTOCItem(-1))
window.addEventListener('scroll', e => activate())
