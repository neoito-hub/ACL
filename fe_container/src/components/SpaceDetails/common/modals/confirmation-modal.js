/* eslint-disable react/prop-types */
import React from 'react'
import DeleteIcon from '../../../../assets/img/icons/delete-icon.gif'

const ConfirmModal = (props) => {
  const {
    hasConfirmModal,
    handleConfirmModal,
    handleSubmit,
    message,
    confirmMessage,
  } = props
  return (
    <div
      className={`fixed left-0 top-0 z-[10001] h-screen w-full ${
        hasConfirmModal ? 'fadeIn' : 'hidden'
      }`}
    >
      <div
        onClick={() => {
          handleConfirmModal()
        }}
        className="fixed left-0 top-0 z-[10001] h-full w-full bg-black/40"
      />
      <div
        className={`absolute top-1/2 left-1/2 z-[10002] w-full max-w-[620px] -translate-x-1/2 -translate-y-1/2 transform px-4 ${
          hasConfirmModal ? '' : 'hidden'
        }`}
      >
        <div className="relative float-left flex w-full rounded-md bg-white p-6 md:space-x-10 md:p-[60px] md-lt:flex-col">
          <img
            src={DeleteIcon}
            alt="Delete"
            className="h-20 w-20 flex-shrink-0 rounded-full md-lt:mx-2 md-lt:mb-4"
          />
          <div className="flex flex-grow flex-col overflow-hidden">
            {message && (
              <h5 className="mb-3 text-lg font-semibold text-black">
                {message}
              </h5>
            )}
            {confirmMessage && (
              <p className="text-ab-base mb-3 text-black">{confirmMessage}</p>
            )}
            <div className="float-left mt-4 flex w-full items-center">
              <button
                type="button"
                onClick={handleSubmit}
                className="btn-default bg-ab-red mr-2 min-w-[80px] rounded px-5 py-2.5 text-sm font-bold leading-tight text-white transition-all hover:opacity-80"
              >
                Confirm
              </button>
              <button
                onClick={() => {
                  handleConfirmModal()
                }}
                type="button"
                className="text-ab-disabled hover:text-ab-black rounded px-3 py-1 text-sm font-bold leading-tight text-white"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ConfirmModal
