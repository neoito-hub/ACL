import React, { useState } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import Header from '../header/header'
import Sidebar from '../sidebar/sidebar'
import FallbackUI from './fallback-ui'

// eslint-disable-next-line react/prop-types
const Layout = ({ children }) => {
  const [activePage, SetActivePage] = useState('home')

  const handleError = (error, errorInfo) => {
    console.log('Error occured in ', errorInfo.componentStack.split(' ')[5])
  }

  return (
    <div className="float-left w-full">
      <Header />
      <div className="float-left w-full pt-16">
        <div className="float-left w-full max-w-[1920px] md:px-6 xl:px-12">
          <div className="fixed left-0 z-[20] float-left w-full overflow-x-hidden border-r bg-white py-3 pr-2 md:top-16 md:left-6 md:block md:h-[calc(100vh-64px)] md:w-[280px] md:pt-6 md:pb-4 xl:left-12 md-lt:z-[9999] md-lt:hidden md-lt:pl-4">
            <Sidebar
              activePage={activePage}
              handleActivePage={(name) => SetActivePage(name)}
            />
          </div>
          <div className="float-left w-full md:pl-[320px]">
            <ErrorBoundary
              FallbackComponent={FallbackUI}
              onError={handleError}
              onReset={() => {
                // reset the state of your app so the error doesn't happen again
              }}
            >
              {children}
            </ErrorBoundary>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Layout
