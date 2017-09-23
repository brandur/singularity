/**
 * Expand the right TOC item based on the content which the user is reading.
 */
function activate() {
  const headers = document.querySelectorAll('h2')
  let i = 0

  for (; i < headers.length; i++) {
    if (!isInView(headers[i])) {
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

function isInView(el) {
  return el.getBoundingClientRect().top <= 0
}

window.addEventListener('load', e => activateTOCItem(-1))
window.addEventListener('scroll', e => activate())
