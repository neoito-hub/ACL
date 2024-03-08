import React, { useEffect, useState } from 'react'
import { ErrorBoundary } from 'react-error-boundary'

import { shield } from '@appblocks/js-sdk'
import Layout from './components/Layout/layout/layout'
import Loader from './components/Layout/Loader'

import './assets/css/main.scss'
import FallbackUI from './common/fallback-ui'
import AppRoute from '../routes/AppRoute'

const App = () => {
  const handleError = (error, errorInfo) => {
    console.log('Error occured in ', errorInfo.componentStack.split(' ')[5])
  }

  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [first, setFirst] = useState(false)

  useEffect(async () => {
    await shield.init(process.env.BLOCK_ENV_URL_CLIENT_ID)
    if (isLoggedIn) {
      setFirst((x) => !x)
    } else {
      const isLoggedinn = await shield.verifyLogin()
      setIsLoggedIn(isLoggedinn)
    }
  }, [isLoggedIn])

  return (
    <ErrorBoundary
      FallbackComponent={FallbackUI}
      onError={handleError}
      onReset={() => {
        // reset the state of your app so the error doesn't happen again
      }}
    >
      {isLoggedIn ? (
        <Layout>
          <AppRoute />
        </Layout>
      ) : (
        <>{Loader} </>
      )}
    </ErrorBoundary>
  )
}

export default App
