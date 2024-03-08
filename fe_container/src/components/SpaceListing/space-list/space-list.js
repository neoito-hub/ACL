/* eslint-disable react/no-array-index-key */
/* eslint-disable camelcase */
/* eslint-disable react/prop-types */
import React, { useContext } from 'react'
import dayjs from 'dayjs'
import { Link } from 'react-router-dom'
import SpaceListCardSkeleton from './space-list-skeleton'
import CalendarIcon from '../../../assets/img/icons/spaces-calendar.svg'
import UserIcon from '../../../assets/img/icons/spaces-user-icon.svg'

import { ACLContext } from '../../../context/ACLContext'

const SpaceList = (props) => {
  const { loader, spaceList } = props
  const loaderData = Array(6).fill(null)

  const { setSpaceId, setSpaceDetails } = useContext(ACLContext)

  const handleSpaceChange = (space_id) => {
    setSpaceDetails(null)
    setSpaceId(space_id)
  }

  return (
    <div className="float-left w-full">
      <div className="-mx-2 flex flex-wrap">
        {loader &&
          !spaceList &&
          loaderData.map((_, index) => <SpaceListCardSkeleton key={index} />)}
        {!loader &&
          spaceList?.map((space) => (
            <Link
              to="/spaces/my-entities"
              key={space?.space_id}
              onClick={() => handleSpaceChange(space?.space_id)}
              className="float-left p-2"
            >
              <div className="border-ab-gray-dark/60 float-left flex w-64 flex-col border bg-white p-5 cursor-pointer">
                <div className="text-primary float-left flex h-12 w-12 items-center justify-center rounded-full bg-[#F2EBFF] text-2xl font-medium capitalize">
                  {space?.space_name[0]}
                </div>
                {/* <img src={} className="w-12 h-12 object-cover object-center"></img> */}
                <div className="float-left w-full py-4">
                  <p className="text-ab-base text-ab-black overflow-hidden text-ellipsis font-medium">
                    {space?.space_name}
                  </p>
                  <div className="text-ab-black/60 spaces-info-group mt-1.5 flex flex-wrap text-xs font-medium">
                    <span>{space?.entity_count} Entities</span>
                  </div>
                </div>
                <div className="float-left flex w-full items-center space-x-2">
                  {space?.my_space ? (
                    <p className="bg-primary/20 text-primary max-w-full truncate rounded-full py-1 px-2 text-xs font-medium">
                      My Space
                    </p>
                  ) : (
                    <>
                      <div className="bg-secondary float-left flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-full text-xs font-medium text-white capitalize">
                        {space?.user_name[0]}
                      </div>
                      <p className="max-w-full truncate rounded-full bg-[#FF9364]/20 py-1 px-2 text-xs font-medium text-[#FF9364]">
                        Shared Space
                      </p>
                    </>
                  )}
                </div>
                <div className="border-ab-dark text-ab-black/80 float-left mt-3 flex w-full flex-wrap items-center border-t pt-1 font-medium">
                  <p className="mt-2 mr-3 flex items-center text-[10px]">
                    <img className="mr-1" src={CalendarIcon} alt="" />
                    {dayjs(space?.created_at).format('DD MMM YYYY')}
                  </p>
                  <p className="mt-2 flex items-center text-[10px]">
                    <img className="mr-1" src={UserIcon} alt="" />
                    {`${space?.member_count} ${
                      space?.member_count > 1 ? 'Members' : 'Member'
                    }`}
                  </p>
                </div>
              </div>
            </Link>
          ))}
        {!loader && !spaceList && (
          <p className="text-ab-black float-left w-full py-10 text-center text-sm">
            No Spaces Found!
          </p>
        )}
      </div>
      {/* <div className='float-left mt-4 w-full'>
        <p className='text-primary text-ab-sm font-medium underline underline-offset-1'>
          See more
        </p>
      </div> */}
    </div>
  )
}

export default SpaceList
