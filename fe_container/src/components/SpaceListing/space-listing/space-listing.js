/* eslint-disable no-unused-expressions */
import React, { useState, useEffect, useRef, useCallback } from 'react'
import { debounce } from 'lodash'
import { Link } from 'react-router-dom'
import apiHelper from '../common/helpers/apiGetters'
import SpaceList from '../space-list/space-list'

const SpaceListing = () => {
  const inputRef = useRef(null)

  const [tabActive, setTabActive] = useState(0)
  const [loader, setLoader] = useState(false)
  const [spaceList, setSpaceList] = useState(null)
  const [filterData, setFilterData] = useState({
    state: 0,
    search_keyword: null,
    page_limit: 50,
    offset: 0,
  })

  useEffect(async () => {
    setLoader(true)
    setSpaceList(null)
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.LIST_SPACES_DETAILED_URL,
      value: filterData,
    })
    res && setSpaceList(res?.data)
    setLoader(false)
  }, [filterData])

  const handler = useCallback(
    debounce((text, activeTab) => {
      const arg = {
        ...filterData,
        search_keyword: text,
        state: activeTab,
      }
      setFilterData(arg)
    }, 1000),
    []
  )

  const onSearchTextChange = (e) => {
    handler(e.target.value, tabActive)
  }

  const handlePageChange = (e, selected) => {
    e.preventDefault()
    inputRef.current.value = ''
    setTabActive(selected)
    const arg = {
      ...filterData,
      state: selected,
      search_keyword: null,
    }
    setFilterData(arg)
  }

  return (
    <div className="float-left w-full max-w-6xl py-6 md-lt:px-4">
      <div className="float-left mt-5 flex w-full items-center justify-between">
        <p className="float-left text-2xl font-semibold text-black/80">
          Spaces
        </p>
      </div>
      <div className="float-left mt-7 w-full">
        <div className="md-h-scroll-primary float-left flex flex-col w-full overflow-x-auto">
          <div className="border-ab-gray-dark float-left flex w-full space-x-3 border-b">
            <div
              onClick={(e) => handlePageChange(e, 0)}
              className={`text-ab-sm relative -bottom-px flex cursor-pointer items-center justify-center border-b-2 px-3 py-2.5 font-medium ${
                tabActive === 0
                  ? 'text-primary border-primary'
                  : 'text-ab-black hover:text-primary border-transparent'
              }`}
            >
              <p className="whitespace-nowrap">All Spaces</p>
            </div>
            <div
              onClick={(e) => handlePageChange(e, 1)}
              className={`text-ab-sm relative -bottom-px flex cursor-pointer items-center justify-center border-b-2 px-3 py-2.5 font-medium ${
                tabActive === 1
                  ? 'text-primary border-primary'
                  : 'text-ab-black hover:text-primary border-transparent'
              }`}
            >
              <p className="whitespace-nowrap">Personal</p>
            </div>
            <div
              onClick={(e) => handlePageChange(e, 2)}
              className={`text-ab-sm relative -bottom-px flex cursor-pointer items-center justify-center border-b-2 px-3 py-2.5 font-medium ${
                tabActive === 2
                  ? 'text-primary border-primary'
                  : 'text-ab-black hover:text-primary border-transparent'
              }`}
            >
              <p className="whitespace-nowrap">Business</p>
            </div>
            {/* <div
              className={`text-sm relative -bottom-px flex cursor-pointer items-center justify-center px-3 py-2.5 font-medium text-ab-black`}
            >
              <p className="whitespace-nowrap">All spaces</p>
            </div> */}
          </div>
          <div className="float-left flex items-center space-x-3 mt-5">
            <input
              ref={inputRef}
              placeholder={
                tabActive === 0
                  ? 'Search Spaces'
                  : tabActive === 1
                  ? 'Search personal spaces'
                  : tabActive === 2
                  ? 'Search business spaces'
                  : 'Search Spaces'
              }
              onChange={onSearchTextChange}
              className="search-input-white border-ab-gray-dark text-ab-sm h-9 w-full rounded-md border !bg-[length:14px_14px] px-2 pl-9 focus:outline-none"
            />
            <Link
              to="/create-space"
              className="btn-primary text-ab-sm flex flex-shrink-0 items-center space-x-2.5 rounded px-5 py-2.5 font-bold leading-tight text-white transition-all hover:opacity-90"
            >
              <svg
                width="14"
                height="14"
                viewBox="0 0 14 14"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M6.99974 0.0664062V13.9331M0.0664062 6.99974H13.9331"
                  stroke="white"
                  strokeWidth="1.5"
                />
              </svg>
              <span>New Space</span>
            </Link>
          </div>
        </div>
        <div className="float-left w-full overflow-x-hidden py-3">
          <SpaceList loader={loader} spaceList={spaceList} />
        </div>
      </div>
    </div>
  )
}

export default SpaceListing
