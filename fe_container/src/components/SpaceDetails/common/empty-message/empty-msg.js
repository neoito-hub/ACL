/* eslint-disable react/prop-types */
import React from 'react'

const EmptyMsg = (props) => {
  const {
    border,
    img,
    text,
    buttonText,
    // buttonType,
    paddingY,
    marginY,
    handleClick,
  } = props
  return (
    <div
      className={`border-ab-dark float-left flex w-full flex-col items-center justify-center space-y-2.5 px-4  my-${marginY} ${
        border && 'border'
      } py-${paddingY}`}
    >
      <img src={img} alt="" />
      <p className="text-ab-sm text-ab-disabled font-medium">{text}</p>
      <button
        type="button"
        onClick={handleClick}
        className="btn-primary text-ab-sm !focus:outline-none rounded px-5 py-2.5 font-bold leading-tight text-white transition-all hover:opacity-90"
      >
        {buttonText}
      </button>
    </div>
  )
}

export default EmptyMsg
