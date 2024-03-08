/* eslint-disable react/destructuring-assignment */
/* eslint-disable react/prop-types */
import React, { useEffect, useRef, useState } from 'react'
// import CloseIcon from '../../../../assets/img/icons/close-icon.svg';

const AclSidebarModal = (props) => {
  const aclContainer = useRef()
  const { hasAclModal, handleAclModal, modalTitle } = props
  const [animationStatus, setAnimationStatus] = useState(true)
  useEffect(() => {
    if (hasAclModal) {
      setAnimationStatus(false)
    }
    aclContainer.current.onanimationend = () => {
      if (!hasAclModal) {
        setAnimationStatus(true)
      }
    }
  }, [hasAclModal])
  return (
    <>
      <div
        onClick={() => {
          handleAclModal()
        }}
        className={`fixed left-0 top-0 z-[1000] h-screen w-full bg-black/40 ${
          hasAclModal ? 'fadeIn' : 'hidden'
        }`}
      />
      <div
        ref={aclContainer}
        className={`fixed top-0 z-[1001] h-screen w-full max-w-xl bg-white py-6 px-4 md:px-12 md:py-10 ${
          animationStatus && 'hidden'
        } ${
          hasAclModal ? 'sidebar-open-animation' : 'sidebar-close-animation'
        }`}
      >
        <div className="relative float-left flex h-full w-full flex-col">
          <div className="flex w-full flex-shrink-0 flex-col space-y-2 px-2">
            <svg
              onClick={() => {
                handleAclModal()
              }}
              className="flex-shrink-0 cursor-pointer"
              width="20"
              height="20"
              viewBox="0 0 20 20"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M5.25033 15.2917L4.72949 14.75L9.47949 10L4.72949 5.22917L5.25033 4.6875L10.0212 9.45833L14.7503 4.6875L15.2712 5.22917L10.542 10L15.2712 14.75L14.7503 15.2917L10.0212 10.5208L5.25033 15.2917Z"
                fill="black"
              />
            </svg>
            <p className="text-lg font-medium text-[#24292E] capitalize">
              {modalTitle}
            </p>
          </div>
          <div className="custom-scroll-primary float-left w-full flex-grow overflow-y-auto px-2">
            {props.children}
          </div>
        </div>
      </div>
    </>
  )
}

export default AclSidebarModal
