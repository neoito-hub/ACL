/* eslint-disable no-unsafe-optional-chaining */
/* eslint-disable jsx-a11y/control-has-associated-label */
/* eslint-disable no-shadow */
/* eslint-disable no-unused-expressions */
/* eslint-disable camelcase */
import React, { useState, useEffect, useContext, useCallback } from 'react'
import { debounce } from 'lodash'
import MyContext from '../common/my-context'
import apiHelper from '../common/helpers/apiGetters'
import Pagination from '../../Layout/pagination/pagination'
import InviteModal from '../common/modals/invite-modal'
import ConfirmModal from '../common/modals/confirmation-modal'

const statuses = [
  { status_code: 1, slug: 'Pending', css: 'bg-[#FFFDEF] text-[#DBAB09]' },
  { status_code: 2, slug: 'Accepted', css: 'bg-ab-green/10 text-ab-green' },
  { status_code: 3, slug: 'Declined', css: 'bg-ab-red/10 text-ab-red' },
  { status_code: 4, slug: 'Expired', css: 'bg-ab-red/10 text-ab-red' },
]

const Invites = () => {
  const page_limit = Number(process.env.PAGE_LIMIT)
  const { spaceId } = useContext(MyContext)

  const [hasConfirmModal, setHasConfirmModal] = useState(false)
  const [current, setCurrent] = useState(null)
  const [hasInviteModal, setHasInviteModal] = useState(false)
  const [loader, setLoader] = useState(true)
  const [flag, setFlag] = useState(false)
  const [invites, setInvites] = useState(null)
  const [totalCount, setTotalCount] = useState(null)
  const [selectedPage, setSelectedPage] = useState(0)
  const [searchText, setSearchText] = useState(null)

  const filterDataStructure = () => ({
    space_id: spaceId,
    is_filter: false,
    is_keyword_search: !!searchText,
    conditions: {
      search_keyword: searchText,
      filter: {},
    },
    page_limit,
    offset: page_limit * selectedPage,
    active: 'createdAt',
    direction: 'DESC',
  })

  useEffect(async () => {
    setLoader(true)
    setInvites(null)
    const res =
      spaceId &&
      (await apiHelper({
        baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
        subUrl: process.env.LIST_INVITED_USERS,
        value: filterDataStructure(),
        spaceId,
      }))
    res && setInvites(res?.data)
    res && setTotalCount(res ? res.total_count : 0)
    setLoader(false)
  }, [flag])

  const handlePageChange = (event) => {
    const { selected } = event
    setSelectedPage(selected)
    setFlag((flag) => !flag)
  }

  const handler = useCallback(
    debounce((text) => {
      setSearchText(text)
      setSelectedPage(0)
      setFlag((flag) => !flag)
    }, 1000),
    []
  )

  const onSearchTextChange = (e) => {
    handler(e.target.value)
  }

  const onRevokeInvite = async () => {
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.REVOKE_INVITE,
      value: {
        invite_ids: [current?.invite_id],
        invite_details_id: [],
      },
      showSuccessMessage: true,
      spaceId,
    })
    setHasConfirmModal(false)
    setCurrent(null)
    setFlag((flag) => !flag)
  }
  const onResendInvite = async (invite) => {
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.RESEND_INVITE,
      value: {
        invite_ids: [invite?.invite_id],
      },
      showSuccessMessage: true,
      spaceId,
    })
    setFlag((flag) => !flag)
  }

  return (
    <div className="float-left w-full">
      <div className="float-left flex items-center space-x-3 pb-3 w-full">
        <input
          placeholder="Search Invite"
          onChange={onSearchTextChange}
          className="search-input-white border-ab-gray-dark text-ab-sm h-9 w-full rounded-md border !bg-[length:14px_14px] px-2 pl-9 focus:outline-none"
        />
        <button
          type="button"
          onClick={() => setHasInviteModal(true)}
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
          <span>Invite Member</span>
        </button>
      </div>
      <div className="float-left w-full py-3">
        <div className="border-ab-gray-dark custom-h-scroll-primary float-left w-full overflow-x-auto border">
          <table className="min-w-full">
            <thead>
              <tr className="text-ab-black bg-ab-gray-light text-left text-sm">
                {/* <th className='whitespace-nowrap p-3 font-normal'>Name</th> */}
                <th className="p-3 font-normal">Email</th>
                <th className="p-3 font-normal">Status</th>
                <th className="p-3 font-normal float-right" />
              </tr>
            </thead>
            <tbody>
              {!loader &&
                invites?.map((invite, index) => (
                  <tr
                    key={invite?.invite_id || index}
                    className="border-ab-gray-dark text-ab-black border-t text-xs"
                  >
                    {/* <td className='p-3'>
                      <div className='flex items-center space-x-2'>
                        <div className='bg-secondary float-left flex h-9 w-9 flex-shrink-0 items-center justify-center overflow-hidden rounded-full'>
                          <span className='text-lg font-semibold leading-normal text-white capitalize'>
                            {member?.email[0] || ''}
                          </span>
                        </div>
                        <p className='whitespace-nowrap capitalize'>
                          {member.full_name || '-'}
                        </p>
                      </div>
                    </td> */}
                    <td className="whitespace-nowrap p-3">
                      <div className="flex items-center space-x-2">
                        <div className="bg-secondary float-left flex h-9 w-9 flex-shrink-0 items-center justify-center overflow-hidden rounded-full">
                          <span className="text-lg font-semibold leading-normal text-white capitalize">
                            {invite?.email[0] || ''}
                          </span>
                        </div>
                        <p className="whitespace-nowrap">
                          {invite.email || '-'}
                        </p>
                      </div>
                    </td>
                    <td className="p-3">
                      <span
                        className={`${
                          statuses[
                            invite?.status === 1 && invite?.expired
                              ? 3
                              : invite?.status - 1
                          ].css
                        } flex-shrink-0 rounded-full py-1 px-2.5 text-xs font-medium`}
                      >
                        {statuses[
                          invite?.status === 1 && invite?.expired
                            ? 3
                            : invite?.status - 1
                        ].slug || '-'}
                      </span>
                    </td>
                    <td className="p-3">
                      {invite?.status === 1 && !invite?.expired && (
                        <div className="flex items-center justify-end space-x-2.5">
                          <button
                            type="button"
                            className="btn-default bg-ab-red mt-1.5 truncate rounded-md py-1 px-3 font-medium text-white first:mt-0 focus:outline-none"
                            onClick={() => {
                              setHasConfirmModal(true)
                              setCurrent(invite)
                            }}
                          >
                            Revoke
                          </button>
                        </div>
                      )}
                      {(invite?.status === 3 || invite?.expired) && (
                        <div className="flex items-center justify-end space-x-2.5">
                          <button
                            type="button"
                            className="btn-default bg-ab-green mt-1.5 truncate rounded-md py-1 px-3 font-medium text-white first:mt-0 focus:outline-none"
                            onClick={() => onResendInvite(invite)}
                          >
                            Resend
                          </button>
                        </div>
                      )}
                    </td>
                  </tr>
                ))}
            </tbody>
          </table>
          {!loader && !invites && (
            <div className="flex justify-center items-center">
              <span className="text-ab-black float-left w-full py-10 text-center text-sm">
                No Invites Found
              </span>
            </div>
          )}
        </div>
        {totalCount > page_limit && (
          <Pagination
            Padding={0}
            marginTop={1}
            pageCount={Math.ceil(totalCount / page_limit)}
            handlePageChange={handlePageChange}
            selected={selectedPage}
          />
        )}
      </div>
      {hasInviteModal && (
        <InviteModal
          hasInviteModal={hasInviteModal}
          handleInviteModal={() => setHasInviteModal(false)}
          type="Space"
          updateList={() => setFlag((flag) => !flag)}
        />
      )}
      <ConfirmModal
        hasConfirmModal={hasConfirmModal}
        handleConfirmModal={() => {
          setHasConfirmModal(false)
          setCurrent(null)
        }}
        handleSubmit={onRevokeInvite}
        message="Are you sure you want to revoke the invite?"
      />
    </div>
  )
}

export default Invites
