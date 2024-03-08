import React from 'react'
import ReactDOM from 'react-dom'
import { BrowserRouter } from 'react-router-dom'
import App from './App'
import { ACLProvider } from './context/ACLContext'

ReactDOM.render(
  <BrowserRouter>
    <ACLProvider>
      <App />
    </ACLProvider>
  </BrowserRouter>,
  document.getElementById('root')
)
