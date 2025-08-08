const api = {
  register: (email, password) => fetch('/api/register', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email, password }) }),
  login: (email, password) => fetch('/api/login', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email, password }) }),
  list: (token) => fetch('/api/expenses', { headers: { 'Authorization': 'Bearer ' + token } }),
  create: (token, expense) => fetch('/api/expenses', { method: 'POST', headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token }, body: JSON.stringify(expense) }),
}

const state = {
  token: localStorage.getItem('token') || '',
}

const qs = (s) => document.querySelector(s)

const toast = (t) => {
  const el = qs('#toast')
  el.textContent = t
  el.classList.remove('hidden')
  setTimeout(() => el.classList.add('hidden'), 2000)
}

const showApp = (authed) => {
  qs('#auth').classList.toggle('hidden', authed)
  qs('#appSection').classList.toggle('hidden', !authed)
  qs('#logoutBtn').classList.toggle('hidden', !authed)
}

const renderExpenses = (items) => {
  const ul = qs('#expenses')
  ul.innerHTML = ''
  items.forEach(x => {
    const li = document.createElement('li')
    const left = document.createElement('div')
    const right = document.createElement('div')
    left.textContent = `${x.category} — ${x.note || ''}`.trim()
    right.textContent = `${x.amount.toFixed(2)}`
    li.appendChild(left)
    li.appendChild(right)
    ul.appendChild(li)
  })
}

const loadExpenses = async () => {
  if (!state.token) return
  const res = await api.list(state.token)
  if (!res.ok) { toast('Ошибка загрузки'); return }
  const data = await res.json()
  renderExpenses(data)
}

const addExpense = async () => {
  const amount = parseFloat(qs('#amount').value)
  const category = qs('#category').value.trim()
  const note = qs('#note').value.trim()
  const dateVal = qs('#date').value
  const date = dateVal ? new Date(dateVal).toISOString() : new Date().toISOString()
  const res = await api.create(state.token, { amount, category, note, date })
  if (!res.ok) { toast('Ошибка добавления'); return }
  await loadExpenses()
  qs('#amount').value = ''
  qs('#category').value = ''
  qs('#note').value = ''
  qs('#date').value = ''
}

const doLogin = async () => {
  const email = qs('#loginEmail').value.trim()
  const password = qs('#loginPassword').value
  const res = await api.login(email, password)
  if (!res.ok) { toast('Неверные данные'); return }
  const data = await res.json()
  state.token = data.token
  localStorage.setItem('token', state.token)
  showApp(true)
  await loadExpenses()
}

const doRegister = async () => {
  const email = qs('#regEmail').value.trim()
  const password = qs('#regPassword').value
  const res = await api.register(email, password)
  if (!res.ok) { toast('Ошибка регистрации'); return }
  toast('Успешно')
}

const doLogout = () => {
  state.token = ''
  localStorage.removeItem('token')
  showApp(false)
}

qs('#loginBtn').addEventListener('click', doLogin)
qs('#registerBtn').addEventListener('click', doRegister)
qs('#addExpenseBtn').addEventListener('click', addExpense)
qs('#logoutBtn').addEventListener('click', doLogout)

if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => navigator.serviceWorker.register('/service-worker.js'))
}

showApp(!!state.token)
loadExpenses()