/* eslint-disable react/prop-types */
import React from 'react'
import './policies.scss'
import Popup from 'reactjs-popup'

const Chip = ({ data }) => (
  <Popup
    trigger={
      <div className="max-w-[70px] text-primary bg-ab-disabled-yellow float-left my-1 mr-1 truncate rounded-full py-1 px-2 text-xs font-medium z-[1000]">
        {data}
      </div>
    }
    className="ab-tooltip-v2"
    on={['hover', 'focus']}
    position="top center"
    closeOnDocumentClick
  >
    {data}
  </Popup>
)

export default Chip
