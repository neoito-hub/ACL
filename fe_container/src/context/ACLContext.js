import React, { createContext, useState } from 'react'

export const ACLContext = createContext({
  spaceId: null,
  setSpaceId: () => {},
  spaceDetails: null,
  setSpaceDetails: () => {},
  userDetails: null,
  setUserDetails: () => {},
  spaceData: null,
  setSpaceData: () => {},
})

export const ACLProvider = ({ children }) => {
  const [spaceId, setSpaceId] = useState(null)
  const [spaceDetails, setSpaceDetails] = useState(null)
  const [userDetails, setUserDetails] = useState(null)
  const [spaceData, setSpaceData] = useState(null)

  const contextValue = {
    spaceId,
    spaceDetails,
    userDetails,
    spaceData,
    setSpaceId,
    setSpaceDetails,
    setUserDetails,
    setSpaceData,
  }

  return (
    <ACLContext.Provider value={contextValue}>{children}</ACLContext.Provider>
  )
}
