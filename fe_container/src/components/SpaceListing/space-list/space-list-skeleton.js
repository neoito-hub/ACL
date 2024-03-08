/* eslint-disable react/no-array-index-key */
import React from 'react'
import Skeleton from 'react-loading-skeleton'
import 'react-loading-skeleton/dist/skeleton.css'

const SpaceListCardSkeleton = () => (
  <div className="float-left p-2">
    <div className="shadow-box hover:shadow-box-md flex cursor-pointer flex-col rounded-sm border border-[#E9E9E9] bg-white p-5 transition-shadow duration-300 w-64 h-[228px]">
      <div className="float-left flex w-full items-center justify-between">
        <Skeleton className="!h-12 !w-12 !rounded-full" />
      </div>
      <div className="float-left mt-2 flex w-full flex-grow flex-col justify-between space-y-3">
        <div className="float-left w-full">
          <Skeleton
            className="mt-1.5 h-3 !rounded first-of-type:mt-3"
            width="90%"
          />
          <Skeleton
            className="mt-1.5 h-3 !rounded first-of-type:mt-3"
            width="90%"
          />
          <div className="float-left mt-[18px] flex w-full items-center justify-between space-x-4">
            <div className="flex items-center overflow-hidden">
              <Skeleton className="!h-6 !w-20 !rounded-full" />
            </div>
          </div>
        </div>
        <div className="float-left flex w-full flex-col border-t border-ab-dark">
          <div className="float-left mt-2 flex w-full items-center space-x-4">
            <Skeleton className="!h-6 !w-20 !rounded-full" />
            <Skeleton className="!h-6 !w-20 !rounded-full" />
          </div>
        </div>
      </div>
    </div>
  </div>
)

export default SpaceListCardSkeleton
