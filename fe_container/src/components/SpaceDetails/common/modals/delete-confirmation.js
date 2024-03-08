/* eslint-disable react/prop-types */
import React from 'react'

const DeleteConfirmModal = (props) => {
  const { hasDeleteConfirmModal, handleDeleteConfirmModal } = props
  return (
    <div
      className={`fixed left-0 top-0 z-[10001] h-screen w-full ${
        hasDeleteConfirmModal ? 'fadeIn' : 'hidden'
      }`}
    >
      <div
        onClick={() => {
          handleDeleteConfirmModal()
        }}
        className="fixed left-0 top-0 z-[10001] h-full w-full bg-black/40"
      />
      <div
        className={`absolute top-1/2 left-1/2 z-[10002] w-full max-w-sm -translate-x-1/2 -translate-y-1/2 transform px-4 md:max-w-xl ${
          hasDeleteConfirmModal ? '' : 'hidden'
        }`}
      >
        <div className="relative float-left flex w-full rounded-md bg-white p-6 md:space-x-10 md:p-[60px] md-lt:flex-col">
          <div className="flex flex-grow flex-col overflow-hidden">
            <h5 className="text-ab-black mb-2 text-base font-semibold">
              Delete This Block
            </h5>
            <p className="text-ab-black mb-4 text-sm">
              Delisting this block will remove it from the store. You can bring
              it back on sale by finding it in Private blocks table and clicking
              on the &apos;Sell the Block&apos; button. Are you sure you want to
              delist this block now?
            </p>
            <div className="float-left mt-4 flex w-full items-center">
              <button
                type="button"
                className="btn-default bg-ab-red mr-2 min-w-[80px] rounded px-4 py-2 text-sm font-bold leading-tight text-white transition-all hover:opacity-80"
              >
                Yes
              </button>
              <button
                onClick={() => {
                  handleDeleteConfirmModal()
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

export default DeleteConfirmModal
