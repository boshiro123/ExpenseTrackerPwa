const api = {
  register: (email, password) => fetch('/api/register', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email, password }) }),
  login: (email, password) => fetch('/api/login', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email, password }) }),
  list: (token) => fetch('/api/expenses', { headers: { 'Authorization': 'Bearer ' + token } }),
  create: (token, expense) => fetch('/api/expenses', { method: 'POST', headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token }, body: JSON.stringify(expense) }),
}

const DEFAULT_CATEGORIES = ['Еда', 'Транспорт', 'Жильё', 'Развлечения', 'Здоровье', 'Покупки', 'Другое']

const state = {
  token: localStorage.getItem('token') || '',
  email: localStorage.getItem('email') || '',
  categories: JSON.parse(localStorage.getItem('categories') || 'null') || DEFAULT_CATEGORIES,
  selectedCategory: '',
  currentView: 'expenses',
}

const qs = (s) => document.querySelector(s)
const qsa = (s) => Array.from(document.querySelectorAll(s))

const snackbar = (t) => {
  const el = qs('#snackbar')
  el.textContent = t
  el.classList.remove('hidden')
  el.style.opacity = '1'
  el.style.transform = 'translateX(-50%) translateY(0)'
  clearTimeout(el._t)
  el._t = setTimeout(() => { el.style.opacity = '0'; el.style.transform = 'translateX(-50%) translateY(8px)'; setTimeout(() => el.classList.add('hidden'), 200) }, 2200)
}

const showApp = (authed) => {
  qs('#auth').classList.toggle('hidden', authed)
  qs('#appSection').classList.toggle('hidden', !authed)
  qs('#logoutBtn').classList.toggle('hidden', !authed)
  if (authed) {
    qs('#profileEmail').textContent = state.email
  }
}

const renderChips = (mount, items, active) => {
  mount.innerHTML = ''
  items.forEach(name => {
    const chip = document.createElement('button')
    chip.className = 'chip ripple' + (name === active ? ' active' : '')
    chip.textContent = name
    chip.addEventListener('click', () => {
      state.selectedCategory = name
      renderChips(mount, items, name)
    })
    mount.appendChild(chip)
  })
}

const renderExpenses = (items) => {
  const grid = qs('#expensesGrid')
  grid.innerHTML = ''
  items.forEach(x => {
    const card = document.createElement('div')
    card.className = 'card ripple'
    card.style.animation = 'fadeUp .24s ease'

    const row1 = document.createElement('div')
    row1.className = 'row'
    const left = document.createElement('div')
    left.textContent = x.note || x.category

    const right = document.createElement('div')
    right.textContent = `${Number(x.amount).toFixed(2)}`

    const chip = document.createElement('span')
    chip.className = 'chip'
    chip.textContent = x.category

    row1.appendChild(left)
    row1.appendChild(chip)

    const row2 = document.createElement('div')
    row2.className = 'row'
    const date = new Date(x.date)
    const dateEl = document.createElement('div')
    dateEl.style.color = 'var(--color-muted)'
    dateEl.textContent = date.toLocaleDateString()

    const amountEl = document.createElement('div')
    amountEl.style.fontWeight = '600'
    amountEl.textContent = right.textContent

    row2.appendChild(dateEl)
    row2.appendChild(amountEl)

    card.appendChild(row1)
    card.appendChild(row2)
    grid.appendChild(card)
  })
}

const loadExpenses = async () => {
  if (!state.token) return
  const res = await api.list(state.token)
  if (!res.ok) { snackbar('Ошибка загрузки'); return }
  const data = await res.json()
  renderExpenses(data)
}

const openSheet = () => {
  const sheet = qs('#expenseSheet')
  const backdrop = qs('#sheetBackdrop')
  sheet.classList.remove('hidden'); backdrop.classList.remove('hidden')
  requestAnimationFrame(() => { sheet.classList.add('show'); backdrop.classList.add('show') })
}

const closeSheet = () => {
  const sheet = qs('#expenseSheet')
  const backdrop = qs('#sheetBackdrop')
  sheet.classList.remove('show'); backdrop.classList.remove('show')
  setTimeout(() => { sheet.classList.add('hidden'); backdrop.classList.add('hidden') }, 200)
}

const addExpense = async () => {
  const amount = parseFloat(qs('#amount').value)
  const note = qs('#note').value.trim()
  const dateVal = qs('#date').value
  const category = state.selectedCategory || state.categories[0]
  const date = dateVal ? new Date(dateVal).toISOString() : new Date().toISOString()
  const res = await api.create(state.token, { amount, category, note, date })
  if (!res.ok) { snackbar('Ошибка добавления'); return }
  closeSheet()
  qs('#amount').value = ''
  qs('#note').value = ''
  qs('#date').value = ''
  state.selectedCategory = ''
  await loadExpenses()
  snackbar('Расход добавлен')
}

const doLogin = async () => {
  const email = qs('#loginEmail').value.trim()
  const password = qs('#loginPassword').value
  const res = await api.login(email, password)
  if (!res.ok) { snackbar('Неверные данные'); return }
  const data = await res.json()
  state.token = data.token
  state.email = email
  localStorage.setItem('token', state.token)
  localStorage.setItem('email', state.email)
  showApp(true)
  await loadExpenses()
}

const doRegister = async () => {
  const email = qs('#regEmail').value.trim()
  const password = qs('#regPassword').value
  const res = await api.register(email, password)
  if (!res.ok) { snackbar('Ошибка регистрации'); return }
  const loginRes = await api.login(email, password)
  if (!loginRes.ok) { snackbar('Не удалось войти'); return }
  const data = await loginRes.json()
  state.token = data.token
  state.email = email
  localStorage.setItem('token', state.token)
  localStorage.setItem('email', state.email)
  showApp(true)
  await loadExpenses()
  snackbar('Добро пожаловать')
}

const doLogout = () => {
  state.token = ''
  state.email = ''
  localStorage.removeItem('token')
  localStorage.removeItem('email')
  showApp(false)
}

const setView = (v) => {
  state.currentView = v
  qsa('.nav-item').forEach(b => b.classList.toggle('active', b.dataset.target === v))
  qsa('.view').forEach(vw => vw.classList.remove('active'))
  if (v === 'expenses') qs('#view-expenses').classList.add('active')
  if (v === 'categories') qs('#view-categories').classList.add('active')
  if (v === 'profile') qs('#view-profile').classList.add('active')
}

const initNav = () => {
  qsa('.nav-item').forEach(b => b.addEventListener('click', () => setView(b.dataset.target)))
}

const initChips = () => {
  renderChips(qs('#categoriesChips'), state.categories, '')
  renderChips(qs('#chipsSelect'), state.categories, state.selectedCategory)
}

const initRipple = () => {
  document.body.addEventListener('pointerdown', (e) => {
    const el = e.target.closest('.ripple')
    if (!el) return
    el.classList.add('pressed')
    setTimeout(() => el.classList.remove('pressed'), 420)
  })
}

qs('#loginBtn').addEventListener('click', doLogin)
qs('#registerBtn').addEventListener('click', doRegister)
qs('#logoutBtn').addEventListener('click', doLogout)
qs('#fab').addEventListener('click', () => { initChips(); openSheet() })
qs('#sheetClose').addEventListener('click', closeSheet)
qs('#sheetBackdrop').addEventListener('click', closeSheet)
qs('#saveExpense').addEventListener('click', addExpense)

if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => navigator.serviceWorker.register('/service-worker.js'))
}

showApp(!!state.token)
setView('expenses')
if (state.token) loadExpenses()
initNav()
initRipple()