import React from 'react'
import ClipLoader from 'react-spinners/ClipLoader'

const Loader = () => (
  <div className="h-full w-full flex items-center justify-center">
    <ClipLoader color="#5E5EDD" size={50} />
  </div>
)

export default Loader
