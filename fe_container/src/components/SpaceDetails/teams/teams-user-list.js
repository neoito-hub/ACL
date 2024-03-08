/* eslint-disable no-unsafe-optional-chaining */
/* eslint-disable react/no-unstable-nested-components */
/* eslint-disable jsx-a11y/control-has-associated-label */
/* eslint-disable camelcase */
/* eslint-disable react/prop-types */
import React from 'react'
import Popup from 'reactjs-popup'
import Pagination from '../../Layout/pagination/pagination'

const TeamsUserList = (props) => {
  const {
    refsById,
    accordianID,
    accordianHeight,
    item,
    users,
    usersSelectedPage,
    handleUsersPageChange,
    usersTotalCount,
  } = props
  const chipsToDisplay = 2
  const page_limit = Number(process.env.PAGE_LIMIT)

  return (
    <div
      ref={refsById[item.team_id]}
      style={{
        height:
          accordianID === item.team_id
            ? accordianHeight !== null && `${accordianHeight}px`
            : '0px',
      }}
      className="float-left w-full overflow-hidden transition-all duration-500"
    >
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
              {users?.map((user) => (
                <tr
                  key={user?.team_member_id}
                  className="border-ab-gray-dark text-ab-black border-t text-xs"
                >
                  <td className="p-3">
                    <div className="flex items-center space-x-2">
                      <div className="bg-secondary float-left flex h-9 w-9 flex-shrink-0 items-center justify-center overflow-hidden rounded-full">
                        <span className="text-lg font-semibold leading-normal text-white capitalize">
                          {user?.full_name[0] || user?.email[0]}
                        </span>
                      </div>
                      {/* <img src={Avatar} className="w-9 h-9 flex-shrink-0 float-left rounded-full border-ab-gray-medium border object-cover" alt=""/>  */}
                      <p className="whitespace-nowrap">{user?.full_name}</p>
                    </div>
                  </td>
                  <td className="whitespace-nowrap p-3">{user?.email}</td>
                  <td className="p-3">
                    <div className="flex flex-wrap items-center">
                      <p className="float-left mr-2 capitalize">
                        {user?.roles
                          ?.slice(0, chipsToDisplay)
                          ?.map((role, index) => (
                            <span key={role?.id}>
                              {(index ? ', ' : '') + role.name}
                            </span>
                          ))}
                      </p>
                      {user?.roles?.length > chipsToDisplay && (
                        <div className="float-left">
                          <Popup
                            arrow={false}
                            keepTooltipInside
                            trigger={() => (
                              <span className="bg-primary float-left my-1 flex h-6 min-w-[24px] cursor-pointer items-center justify-center rounded-full py-1 px-1 text-[10px] font-medium leading-tight text-white">
                                +{user?.roles?.length - chipsToDisplay}
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
                                  {user?.roles
                                    ?.slice(
                                      chipsToDisplay,
                                      item?.policy_groups?.length,
                                    )
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
                  {/* <td className='p-3'>
                  <div className='flex items-center justify-end space-x-2.5'>
                    <svg
                      data-tip="Edit Member"
                      onClick={() => setHasEditMemberModal(true)}
                      className="hover:fill-secondary cursor-pointer fill-black focus:outline-none"
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path d="M8.33333 9.2915C7.70833 9.2915 7.1875 9.08664 6.77083 8.67692C6.35417 8.2672 6.14583 7.74289 6.14583 7.104C6.14583 6.479 6.35417 5.95817 6.77083 5.5415C7.1875 5.12484 7.70833 4.9165 8.33333 4.9165C8.95833 4.9165 9.47917 5.12484 9.89583 5.5415C10.3125 5.95817 10.5208 6.479 10.5208 7.104C10.5208 7.74289 10.3125 8.2672 9.89583 8.67692C9.47917 9.08664 8.95833 9.2915 8.33333 9.2915ZM2.75 15.1457V14.1665C2.75 13.8054 2.84028 13.5068 3.02083 13.2707C3.20139 13.0346 3.46528 12.8332 3.8125 12.6665C4.57639 12.3193 5.28472 12.0554 5.9375 11.8748C6.59028 11.6943 7.38889 11.604 8.33333 11.604H8.48958C8.53819 11.604 8.59028 11.5971 8.64583 11.5832C8.61806 11.6804 8.59722 11.7637 8.58333 11.8332C8.56944 11.9026 8.55556 11.979 8.54167 12.0623H8.33333C7.45833 12.0623 6.69444 12.1353 6.04167 12.2811C5.38889 12.4269 4.72222 12.6873 4.04167 13.0623C3.70833 13.229 3.48611 13.3991 3.375 13.5728C3.26389 13.7464 3.20833 13.9443 3.20833 14.1665V14.6873H8.5625C8.57639 14.7568 8.59722 14.8332 8.625 14.9165C8.65278 14.9998 8.68055 15.0762 8.70833 15.1457H2.75ZM13.7292 15.4373L13.6667 14.5207C13.3889 14.479 13.125 14.3887 12.875 14.2498C12.625 14.1109 12.4167 13.9373 12.25 13.729L11.4583 14.0415L11.3333 13.854L12.0417 13.3123C11.9444 13.0901 11.8958 12.8332 11.8958 12.5415C11.8958 12.2498 11.9444 11.9859 12.0417 11.7498L11.3542 11.1873L11.4792 10.9998L12.25 11.3332C12.4028 11.1109 12.6076 10.9339 12.8646 10.8019C13.1215 10.67 13.3889 10.5832 13.6667 10.5415L13.7292 9.62484H13.9792L14.0417 10.5415C14.3056 10.5832 14.566 10.67 14.8229 10.8019C15.0799 10.9339 15.2917 11.104 15.4583 11.3123L16.2292 10.9998L16.3333 11.1665L15.6458 11.729C15.7569 11.979 15.8125 12.2464 15.8125 12.5311C15.8125 12.8158 15.7569 13.0762 15.6458 13.3123L16.3542 13.854L16.25 14.0415L15.4583 13.729C15.2778 13.9373 15.0625 14.1109 14.8125 14.2498C14.5625 14.3887 14.3056 14.479 14.0417 14.5207L13.9792 15.4373H13.7292ZM13.8333 13.9998C14.2639 13.9998 14.6146 13.8609 14.8854 13.5832C15.1562 13.3054 15.2917 12.9582 15.2917 12.5415C15.2917 12.1109 15.1562 11.7603 14.8854 11.4894C14.6146 11.2186 14.2639 11.0832 13.8333 11.0832C13.4167 11.0832 13.0694 11.2186 12.7917 11.4894C12.5139 11.7603 12.375 12.1109 12.375 12.5415C12.375 12.9582 12.5139 13.3054 12.7917 13.5832C13.0694 13.8609 13.4167 13.9998 13.8333 13.9998ZM8.33333 8.83317C8.81944 8.83317 9.22917 8.66998 9.5625 8.34359C9.89583 8.0172 10.0625 7.604 10.0625 7.104C10.0625 6.61789 9.89583 6.20817 9.5625 5.87484C9.22917 5.5415 8.81944 5.37484 8.33333 5.37484C7.84722 5.37484 7.4375 5.5415 7.10417 5.87484C6.77083 6.20817 6.60417 6.61789 6.60417 7.104C6.60417 7.604 6.77083 8.0172 7.10417 8.34359C7.4375 8.66998 7.84722 8.83317 8.33333 8.83317Z" />
                    </svg>
                    <svg
                      data-tip="ACL"
                      onClick={() => setHasAclModal(true)}
                      className="hover:fill-secondary cursor-pointer fill-black focus:outline-none"
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path d="M10.5417 16.9582V16.0207L14.8125 11.7498L15.75 12.6873L11.4792 16.9582H10.5417ZM3.04167 12.729V12.2707H9V12.729H3.04167ZM16.4375 11.9998L15.5 11.0623L16.1042 10.4582C16.1458 10.4165 16.1979 10.3957 16.2604 10.3957C16.3229 10.3957 16.3819 10.4165 16.4375 10.4582L17.0417 11.0623C17.0833 11.1179 17.1042 11.1769 17.1042 11.2394C17.1042 11.3019 17.0833 11.354 17.0417 11.3957L16.4375 11.9998ZM3.04167 9.354V8.89567H12.3958V9.354H3.04167ZM3.04167 5.99984V5.5415H12.3958V5.99984H3.04167Z" />
                    </svg>
                    <svg
                      data-tip="Delete"
                      onClick={() => {
                        setHasDeleteUserModal(true);
                        setSelectedUser(user.team_member_id);
                      }}
                      className="fill-ab-black hover:fill-ab-red cursor-pointer focus:outline-none"
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path d="M7.75 13.2083L10.0208 10.9167L12.3125 13.2083L12.8542 12.6667L10.5833 10.375L12.8542 8.0625L12.3125 7.52083L10.0208 9.8125L7.75 7.52083L7.1875 8.0625L9.47917 10.375L7.1875 12.6667L7.75 13.2083ZM6.33333 16.75C5.95833 16.75 5.63195 16.6111 5.35417 16.3333C5.07639 16.0556 4.9375 15.7292 4.9375 15.3542V4.77083H4.10417V4.04167H7.39583V3.4375H12.625V4.04167H15.9167V4.77083H15.0833V15.3542C15.0833 15.7292 14.9479 16.0556 14.6771 16.3333C14.4062 16.6111 14.0764 16.75 13.6875 16.75H6.33333ZM14.3542 4.77083H5.66667V15.3542C5.66667 15.5208 5.73611 15.6736 5.875 15.8125C6.01389 15.9514 6.16667 16.0208 6.33333 16.0208H13.6875C13.8542 16.0208 14.0069 15.9514 14.1458 15.8125C14.2847 15.6736 14.3542 15.5208 14.3542 15.3542V4.77083ZM5.66667 4.77083V16.0208V15.3542V4.77083Z" />
                    </svg>
                  </div>
                </td> */}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <Pagination
          Padding={0}
          marginTop={1}
          pageCount={Math.ceil(usersTotalCount / page_limit)}
          handlePageChange={handleUsersPageChange}
          selected={usersSelectedPage}
        />
      </div>
    </div>
  )
}

export default TeamsUserList
