/* eslint-disable import/no-extraneous-dependencies */
import React from 'react'
import Skeleton from 'react-loading-skeleton'
import 'react-loading-skeleton/dist/skeleton.css'

const HomeSpaceListingSkeletonLoader = () => (
  <div className="border-ab-gray-medium flex w-full flex-col items-center border px-6 py-8 cursor-pointer">
    <Skeleton className="!h-12 !w-12 !rounded-full" />
    <Skeleton className="!h-6 !w-24 !mt-3.5" />
  </div>
)

export default HomeSpaceListingSkeletonLoader
