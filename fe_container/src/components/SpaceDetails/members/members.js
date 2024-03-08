/* eslint-disable jsx-a11y/control-has-associated-label */
/* eslint-disable no-shadow */
/* eslint-disable no-unused-expressions */
/* eslint-disable camelcase */
import React, { useState, useEffect, useContext, useCallback } from 'react'
import { Tooltip as ReactTooltip } from 'react-tooltip'
// import DeleteUserModal from '../common/modals/delete-user-modal'
import { debounce } from 'lodash'
import Popup from 'reactjs-popup'
import MyContext from '../common/my-context'
import apiHelper from '../common/helpers/apiGetters'
import Pagination from '../../Layout/pagination/pagination'
import EditMemberModal from '../common/modals/edit-member-modal'
import InviteModal from '../common/modals/invite-modal'
import AclSidebarModal from '../common/modals/acl-sidebar-modal'
import Policies from '../policies/policies'

const Members = () => {
  const page_limit = Number(process.env.PAGE_LIMIT)
  const { spaceId } = useContext(MyContext)

  const chipsToDisplay = 2

  // const [hasDeleteUserModal, setHasDeleteUserModal] = useState(false);
  const [hasEditMemberModal, setHasEditMemberModal] = useState(false)
  const [hasInviteModal, setHasInviteModal] = useState(false)
  const [hasAclModal, setHasAclModal] = useState(false)
  const [loader, setLoader] = useState(true)
  const [flag, setFlag] = useState(false)
  const [members, setMembers] = useState(null)
  const [totalCount, setTotalCount] = useState(null)
  const [selectedPage, setSelectedPage] = useState(0)
  const [searchText, setSearchText] = useState(null)
  const [currentUser, setCurrentUser] = useState(null)

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
    direction: 'ASC',
  })

  useEffect(async () => {
    setLoader(true)
    setMembers(null)
    const res =
      spaceId &&
      (await apiHelper({
        baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
        subUrl: process.env.LIST_USERS_URL,
        value: filterDataStructure(),
        spaceId,
      }))
    res && setMembers(res?.data)
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
    [],
  )

  const onSearchTextChange = (e) => {
    handler(e.target.value)
  }

  return (
    <div className="float-left w-full">
      <div className="float-left flex items-center space-x-3 pb-3 w-full">
        <input
          placeholder="Search Member"
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
          <span>New Member</span>
        </button>
      </div>
      <div className="float-left w-full py-3">
        <div className="border-ab-gray-dark custom-h-scroll-primary float-left w-full overflow-x-auto border">
          <table className="min-w-full">
            <thead>
              <tr className="text-ab-black bg-ab-gray-light text-left text-sm">
                <th className="whitespace-nowrap p-3 font-normal">Name</th>
                <th className="p-3 font-normal">Email</th>
                <th className="p-3 font-normal">Role</th>
                <th className="p-3 font-normal" />
              </tr>
            </thead>
            <tbody>
              {!loader &&
                members?.map((member, index) => (
                  <tr
                    key={member?.user_id || index}
                    className="border-ab-gray-dark text-ab-black border-t text-xs"
                  >
                    <td className="p-3">
                      <div className="flex items-center space-x-2">
                        <div className="bg-secondary float-left flex h-9 w-9 flex-shrink-0 items-center justify-center overflow-hidden rounded-full">
                          <span className="text-lg font-semibold leading-normal text-white capitalize">
                            {member?.full_name[0] || member?.email[0]}
                          </span>
                        </div>
                        <p className="whitespace-nowrap capitalize">
                          {member.full_name || '-'}
                        </p>
                      </div>
                    </td>
                    <td className="whitespace-nowrap p-3">{member.email}</td>
                    <td className="p-3">
                      <div className="flex flex-wrap items-center">
                        <p className="float-left mr-2 capitalize">
                          {member?.roles
                            ?.slice(0, chipsToDisplay)
                            ?.map((role, index) => (
                              <span key={role?.id}>
                                {(index ? ', ' : '') + role.name}
                              </span>
                            ))}
                        </p>
                        {member?.roles?.length > chipsToDisplay && (
                          <div className="float-left">
                            <Popup
                              arrow={false}
                              keepTooltipInside
                              trigger={() => (
                                <span className="bg-primary float-left my-1 flex h-6 min-w-[24px] cursor-pointer items-center justify-center rounded-full py-1 px-1 text-[10px] font-medium leading-tight text-white">
                                  +{member?.roles?.length - chipsToDisplay}
                                </span>
                              )}
                              position="bottom center"
                              className="ab-dropdown-popup max-w-270"
                            >
                              <div className="border-ab-gray-dark shadow-box dropdownFade border-ab-gray-dark shadow-box float-left mt-4 w-full max-w-[270px] rounded-b-md border bg-white px-2 py-2.5">
                                <div className="float-left w-full bg-white">
                                  <div className="float-left w-full p-2">
                                    <p className="text-ab-black border-ab-gray-dark float-left w-full border-b pb-2 text-xs font-medium">
                                      Role
                                    </p>
                                  </div>
                                  <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-2 capitalize">
                                    {member?.roles
                                      ?.slice(chipsToDisplay)
                                      ?.map((role) => (
                                        <li
                                          key={role?.id}
                                          className="text-ab-black float-left mb-4 w-full truncate text-xs font-medium tracking-tight last-of-type:mb-0"
                                        >
                                          {role?.name}
                                        </li>
                                      ))}
                                  </ul>
                                </div>
                              </div>
                            </Popup>
                          </div>
                        )}
                      </div>
                    </td>
                    {!member?.is_owner && (
                      <td className="p-3">
                        <div className="flex items-center justify-end space-x-2.5">
                          <Popup
                            trigger={
                              <button
                                onClick={() => {
                                  setHasEditMemberModal(true)
                                  setCurrentUser(member)
                                }}
                                type="button"
                                className="float-left flex flex-shrink-0 items-center justify-center cursor-pointer rounded group"
                              >
                                <svg
                                  content="Edit Member"
                                  className="cursor-pointer focus:outline-none fill-ab-black/60 hover:fill-primary"
                                  width="16"
                                  height="16"
                                  viewBox="0 0 16 16"
                                  fill="none"
                                  xmlns="http://www.w3.org/2000/svg"
                                >
                                  <path d="M1.33331 16.0001V13.9834H14.6666V16.0001H1.33331ZM2.69998 12.2168V10.0001L8.91665 3.78343L11.1333 6.00009L4.91665 12.2168H2.69998ZM3.69998 11.2168H4.44998L9.69998 5.96676L8.94998 5.21676L3.69998 10.4668V11.2168ZM11.8666 5.26676L9.64998 3.05009L11.05 1.65009C11.1722 1.50565 11.3111 1.43065 11.4666 1.42509C11.6222 1.41954 11.7778 1.49454 11.9333 1.65009L13.2333 2.95009C13.3778 3.09454 13.45 3.24731 13.45 3.40843C13.45 3.56954 13.3889 3.72231 13.2666 3.86676L11.8666 5.26676Z" />
                                </svg>
                              </button>
                            }
                            className="ab-tooltip-v2"
                            on={['hover', 'focus']}
                            position="top center"
                            closeOnDocumentClick
                          >
                            Edit Member
                          </Popup>
                          <Popup
                            trigger={
                              <button
                                type="button"
                                onClick={() => {
                                  setHasAclModal(true)
                                  setCurrentUser(member)
                                }}
                                className="float-left flex flex-shrink-0 items-center justify-center cursor-pointer rounded group"
                              >
                                <svg
                                  className="focus:outline-none fill-ab-black/60 group-hover:fill-secondary"
                                  width="16"
                                  height="16"
                                  viewBox="0 0 16 16"
                                  fill="none"
                                  xmlns="http://www.w3.org/2000/svg"
                                >
                                  <path d="M11.5167 11.4997C11.8056 11.4997 12.05 11.3942 12.25 11.1831C12.45 10.972 12.55 10.722 12.55 10.4331C12.55 10.1442 12.45 9.89974 12.25 9.69974C12.05 9.49974 11.8056 9.39974 11.5167 9.39974C11.2278 9.39974 10.9778 9.49974 10.7667 9.69974C10.5556 9.89974 10.45 10.1442 10.45 10.4331C10.45 10.722 10.5556 10.972 10.7667 11.1831C10.9778 11.3942 11.2278 11.4997 11.5167 11.4997ZM11.5 13.5831C11.8667 13.5831 12.2 13.5053 12.5 13.3497C12.8 13.1942 13.0556 12.972 13.2667 12.6831C12.9778 12.5275 12.6889 12.4108 12.4 12.3331C12.1111 12.2553 11.8111 12.2164 11.5 12.2164C11.1889 12.2164 10.8861 12.2553 10.5917 12.3331C10.2972 12.4108 10.0111 12.5275 9.73335 12.6831C9.94447 12.972 10.1972 13.1942 10.4917 13.3497C10.7861 13.5053 11.1222 13.5831 11.5 13.5831ZM8.00002 14.6664C6.46669 14.3108 5.19446 13.4414 4.18335 12.0581C3.17224 10.6747 2.66669 9.08863 2.66669 7.29974V3.31641L8.00002 1.31641L13.3334 3.31641V7.81641C13.1778 7.73863 13.0111 7.66918 12.8334 7.60807C12.6556 7.54696 12.4889 7.50529 12.3334 7.48307V4.01641L8.00002 2.41641L3.66669 4.01641V7.29974C3.66669 8.14418 3.8028 8.92196 4.07502 9.63307C4.34724 10.3442 4.69446 10.9692 5.11669 11.5081C5.53891 12.047 6.00558 12.4942 6.51669 12.8497C7.0278 13.2053 7.52224 13.4608 8.00002 13.6164C8.06669 13.7497 8.16669 13.8997 8.30002 14.0664C8.43335 14.2331 8.54447 14.3608 8.63335 14.4497C8.53335 14.5053 8.4278 14.547 8.31669 14.5747C8.20558 14.6025 8.10002 14.6331 8.00002 14.6664ZM11.55 14.6664C10.6834 14.6664 9.94447 14.3581 9.33335 13.7414C8.72224 13.1247 8.41669 12.3942 8.41669 11.5497C8.41669 10.6831 8.72224 9.94141 9.33335 9.32474C9.94447 8.70807 10.6834 8.39974 11.55 8.39974C12.4056 8.39974 13.1417 8.70807 13.7584 9.32474C14.375 9.94141 14.6834 10.6831 14.6834 11.5497C14.6834 12.3942 14.375 13.1247 13.7584 13.7414C13.1417 14.3581 12.4056 14.6664 11.55 14.6664Z" />
                                </svg>
                              </button>
                            }
                            className="ab-tooltip-v2"
                            on={['hover', 'focus']}
                            position="top center"
                            closeOnDocumentClick
                          >
                            Manage Permissions
                          </Popup>
                          {/* <Popup
                          trigger={
                            <button
                              onClick={() => setHasDeleteUserModal(true)}
                              type='button'
                              className='float-left flex flex-shrink-0 items-center justify-center cursor-pointer rounded group'
                            >
                              <svg
                                className='focus:outline-none fill-ab-black/60 group-hover:fill-ab-red'
                                width='16'
                                height='16'
                                viewBox='0 0 16 16'
                                fill='none'
                                xmlns='http://www.w3.org/2000/svg'
                              >
                                <path d='M4.35002 14C4.07224 14 3.83613 13.9028 3.64169 13.7083C3.44724 13.5139 3.35002 13.2778 3.35002 13V3.5H2.66669V2.5H5.80002V2H10.2V2.5H13.3334V3.5H12.65V13C12.65 13.2667 12.55 13.5 12.35 13.7C12.15 13.9 11.9167 14 11.65 14H4.35002ZM11.65 3.5H4.35002V13H11.65V3.5ZM6.11669 11.5667H7.11669V4.91667H6.11669V11.5667ZM8.88335 11.5667H9.88335V4.91667H8.88335V11.5667ZM4.35002 3.5V13V3.5Z' />
                              </svg>
                            </button>
                          }
                          className='ab-tooltip-v2'
                          on={['hover', 'focus']}
                          position='top center'
                          closeOnDocumentClick
                        >
                          Delete
                        </Popup> */}
                        </div>
                      </td>
                    )}
                  </tr>
                ))}
            </tbody>
          </table>
          {!loader && !members && (
            <div className="flex justify-center items-center">
              <span className="text-ab-black float-left w-full py-10 text-center text-sm">
                No Users Found
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
      <ReactTooltip
        effect="solid"
        padding="4px 8px"
        className="ab-tooltip"
        arrowColor="transparent"
      />
      {/* <DeleteUserModal
        hasDeleteUserModal={hasDeleteUserModal}
        handleDeleteUserModal={() => setHasDeleteUserModal(false)}
      /> */}

      {hasEditMemberModal && (
        <EditMemberModal
          hasEditMemberModal={hasEditMemberModal}
          handleEditMemberModal={() => {
            setHasEditMemberModal(false)
            setCurrentUser(null)
          }}
          selectedUser={currentUser}
          updateList={(flag) => setFlag(!flag)}
        />
      )}
      {hasInviteModal && (
        <InviteModal
          hasInviteModal={hasInviteModal}
          handleInviteModal={() => setHasInviteModal(false)}
          type="Space"
          updateList={() => {}}
        />
      )}
      {hasAclModal && (
        <AclSidebarModal
          hasAclModal={hasAclModal}
          handleAclModal={() => {
            setHasAclModal(false)
            setCurrentUser(null)
          }}
          modalTitle={`${
            currentUser?.full_name ? `${currentUser?.full_name}'s` : 'Member'
          } Access Management`}
        >
          <Policies
            hasTabMenu
            tabMenus={['Permissions', 'Member', 'Entities']}
            getExistingPolicies={process.env.USER_LIST_EXISTING_POL_GRP_SUBS}
            getPoliciesToAdd={process.env.USER_LIST_TO_ADD_POL_GRP_SUBS}
            current={currentUser}
            type="Member"
          />
        </AclSidebarModal>
      )}
    </div>
  )
}

export default Members
