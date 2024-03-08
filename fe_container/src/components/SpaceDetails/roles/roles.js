/* eslint-disable jsx-a11y/control-has-associated-label */
/* eslint-disable no-shadow */
/* eslint-disable no-unused-expressions */
/* eslint-disable camelcase */
import React, {
  useState,
  useMemo,
  useEffect,
  createRef,
  useContext,
  useCallback,
} from 'react'
import { debounce } from 'lodash'
import { Tooltip as ReactTooltip } from 'react-tooltip'
import Popup from 'reactjs-popup'
import apiHelper from '../common/helpers/apiGetters'
import MyContext from '../common/my-context'
import DeleteUserModal from '../common/modals/delete-user-modal'
import CreateNewModal from '../common/modals/create-new-modal'
import ExpandMoreIcon from '../../../assets/img/icons/expand-more.svg'
import InviteModal from '../common/modals/invite-modal'
import AclSidebarModal from '../common/modals/acl-sidebar-modal'
import Policies from '../policies/policies'
import EmptyMsg from '../common/empty-message/empty-msg'
import RolePlaceholderIcon from '../../../assets/img/icons/role-placeholder-icon.svg'
import RolesUserList from './roles-user-list'

const Roles = () => {
  // const page_limit = Number(process.env.PAGE_LIMIT);
  const page_limit = 50
  const { spaceId } = useContext(MyContext)

  const [hasDeleteRoleModal, setHasDeleteRoleModal] = useState(false)
  const [hasDeleteUserModal, setHasDeleteUserModal] = useState(false)
  const [hasCreateNewModal, setHasCreateNewModal] = useState(false)
  const [hasInviteModal, setHasInviteModal] = useState(false)
  const [hasAclModal, setHasAclModal] = useState(false)
  const [accordianID, setAccordianID] = useState(null)
  const [accordianHeight, setAccordianHeight] = useState(null)
  const [currentRole, setCurrentRole] = useState(null)
  // const [loading, setLoading] = useState(false);
  // const [selectedPage, setSelectedPage] = useState(0);
  const [selectedRole, setselectedRole] = useState(null)
  const [selectedUser, setSelectedUser] = useState(null)
  // const [totalCount, setTotalCount] = useState(null);
  const [usersSelectedPage, setUsersSelectedPage] = useState(0)
  const [usersTotalCount, setUsersTotalCount] = useState(null)
  const [roles, setRoles] = useState(null)
  const [users, setUsers] = useState(null)
  const [flag, setFlag] = useState(false)
  const [loader, setLoader] = useState(true)
  const [userFlag, setUserFlag] = useState(false)
  const [searchText, setSearchText] = useState(null)

  const [usersFilterData, setUsersFilterData] = useState({
    space_id: spaceId,
    role_id: null,
    is_filter: false,
    is_keyword_search: false,
    conditions: {
      search_keyword: null,
      filter: {},
    },
    page_limit,
    offset: 0,
    active: 'createdAt',
    direction: 'ASC',
  })

  const filterDataStructure = () => ({
    space_id: spaceId,
    search_keyword: searchText,
    page_limit,
    offset: 0,
  })

  const refsById = useMemo(() => {
    const refs = {}
    roles?.forEach((item) => {
      refs[item.role_id] = createRef(null)
    })
    return refs
  }, [roles])

  const getUsersData = async (arg, container) => {
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.ROLES_LIST_USERS_URL,
      value: arg,
      spaceId,
    })
    setUsers(res ? res.data : null)
    setUsersTotalCount(res ? res.total_count : 0)
    setAccordianHeight(container.current.scrollHeight)
    // handleAccordian(arg.role_id, refsById[arg.role_id]);
  }

  const handleAccordian = (id, container) => {
    setAccordianHeight(container.current.scrollHeight)
    if (accordianID === id) {
      setAccordianID(null)
      setUsers(null)
    } else {
      setAccordianID(id)
      const arg = {
        ...usersFilterData,
        role_id: id,
      }
      id && getUsersData(arg, container)
      // setUsersFilterData(arg);
    }
  }

  const getRolesList = async () => {
    setLoader(true)
    setRoles(null)
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.LIST_ROLES_URL,
      value: filterDataStructure(),
      spaceId,
    })
    setRoles(res ? res.data : null)
    // setTotalCount(res ? res.total_count : 0);
    setLoader(false)
  }

  useEffect(async () => {
    getRolesList()
  }, [flag])

  useEffect(async () => {
    const res =
      usersFilterData?.role_id &&
      (await apiHelper({
        baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
        subUrl: process.env.ROLES_LIST_USERS_URL,
        value: usersFilterData,
        spaceId,
      }))
    setUsers(res ? res.data : null)
    setUsersTotalCount(res ? res.total_count : 0)
  }, [usersFilterData, userFlag])

  const handleUsersPageChange = (event) => {
    const { selected } = event
    setUsersSelectedPage(selected)
    const arg = {
      ...usersFilterData,
      offset: usersFilterData.page_limit * selected,
    }
    setUsersFilterData(arg)
  }

  const onCreateNewRole = async (appName) => {
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.CREATE_ROLE_URL,
      value: {
        name: appName,
        description: '',
        space_id: spaceId,
      },
      spaceId,
    })
    setHasCreateNewModal((hasCreateNewModal) => !hasCreateNewModal)
    setFlag((flag) => !flag)
  }

  const onDeleteRole = async () => {
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.DELETE_ROLE_URL,
      value: {
        role_id: selectedRole,
      },
      apiType: 'delete',
      spaceId,
    })
    setselectedRole(null)
    setHasDeleteRoleModal(false)
  }

  const onDeleteUser = async () => {
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.ROLES_DELETE_USER,
      value: {
        role_id: selectedUser,
      },
      apiType: 'delete',
      spaceId,
    })
    setSelectedUser(null)
    setHasDeleteUserModal(false)
    setUserFlag((userFlag) => !userFlag)
  }

  const handler = useCallback(
    debounce((text) => {
      setSearchText(text)
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
          placeholder="Search Role"
          onChange={onSearchTextChange}
          className="search-input-white border-ab-gray-dark text-ab-sm h-9 w-full rounded-md border !bg-[length:14px_14px] px-2 pl-9 focus:outline-none"
        />
        <button
          type="button"
          onClick={() => setHasCreateNewModal(true)}
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
          <span>New Role</span>
        </button>
      </div>
      {roles?.map((item) => (
        <div
          key={item?.role_id}
          className="border-ab-gray-dark float-left mt-3 w-full overflow-hidden rounded-md border"
        >
          <div className="border-ab-gray-dark float-left flex w-full items-center justify-between border-b px-4 py-3">
            <p className="text-ab-base text-ab-black truncate pr-3 font-semibold leading-tight capitalize">
              {item?.name} ({item?.member_count})
            </p>
            <div className="float-left flex items-center space-x-2">
              <div className="float-left flex items-center space-x-2.5">
                <button
                  onClick={() => {
                    setHasInviteModal(true)
                    setselectedRole(item?.role_id)
                  }}
                  type="button"
                  className={`btn-default bg-ab-yellow text-ab-black disabled:bg-ab-disabled h-8 flex-shrink-0 rounded px-3 py-2 text-xs font-medium leading-tight transition-all hover:opacity-90 ${
                    accordianID === item?.role_id ? '' : 'hidden'
                  }`}
                >
                  Invite new member
                </button>
                {!item?.is_owner && (
                  <Popup
                    trigger={
                      <button
                        type="button"
                        onClick={() => {
                          setHasAclModal(true)
                          setCurrentRole(item)
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
                )}
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
              <div className="flex h-8 items-center">
                <img
                  onClick={() =>
                    handleAccordian(item?.role_id, refsById[item?.role_id])
                  }
                  className={`flex-shrink-0 cursor-pointer transition-all duration-100 ${
                    accordianID === item?.role_id && 'rotate-180'
                  }`}
                  src={ExpandMoreIcon}
                  alt=""
                />
              </div>
            </div>
          </div>
          <div className="float-left w-full px-4">
            <div
              className={`float-left flex w-full flex-wrap py-2 ${
                accordianID !== null && 'dropDownFade'
              } ${accordianID === item.role_id && 'hidden'}`}
            >
              {item?.members?.map((member, index) => (
                <div
                  key={member?.user_id}
                  className={`${
                    index % 2 === 0 ? 'bg-secondary' : 'bg-ab-black '
                  } float-left my-1 mr-2 flex h-9 w-9 flex-shrink-0 items-center justify-center overflow-hidden rounded-full`}
                >
                  <span className="text-lg font-semibold leading-normal text-white capitalize">
                    {member?.full_name[0] || member?.email[0]}
                  </span>
                </div>
              ))}
              <span
                onClick={() => {
                  setHasInviteModal(true)
                  setselectedRole(item?.role_id)
                }}
                className="float-left my-1 mr-2 flex h-9 w-9 flex-shrink-0 cursor-pointer items-center justify-center overflow-hidden rounded-full bg-[#F2EBFF]"
              >
                <svg
                  width="12"
                  height="12"
                  viewBox="0 0 12 12"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M6 0.166504V11.8332M0.166672 5.99984H11.8333"
                    stroke="#5E5EDD"
                    strokeWidth="1.66667"
                  />
                </svg>
              </span>
            </div>
            <RolesUserList
              refsById={refsById}
              accordianID={accordianID}
              item={item}
              accordianHeight={accordianHeight}
              users={users}
              handleUsersPageChange={handleUsersPageChange}
              usersSelectedPage={usersSelectedPage}
              usersTotalCount={usersTotalCount}
            />
          </div>
        </div>
      ))}
      {!loader && !roles && searchText && (
        <div className="text-ab-black float-left w-full py-10 text-center text-sm">
          No Roles Found
        </div>
      )}
      {!roles && !searchText && !loader && (
        <EmptyMsg
          border={false}
          img={RolePlaceholderIcon}
          text="You don’t have any roles"
          buttonText="Create your first role"
          buttonType="role"
          paddingY="10"
          marginY="2"
          handleClick={() => setHasCreateNewModal(true)}
        />
      )}
      <ReactTooltip
        effect="solid"
        padding="4px 8px"
        className="ab-tooltip"
        arrowColor="transparent"
      />
      {hasDeleteRoleModal && (
        <DeleteUserModal
          hasDeleteUserModal={hasDeleteRoleModal}
          handleDeleteUserModal={() => {
            setHasDeleteRoleModal(false)
            setselectedRole(null)
          }}
          handleSubmit={() => onDeleteRole()}
          type="Team"
        />
      )}
      {hasDeleteUserModal && (
        <DeleteUserModal
          hasDeleteUserModal={hasDeleteUserModal}
          handleDeleteUserModal={() => {
            setHasDeleteUserModal(false)
            setSelectedUser(null)
          }}
          handleSubmit={() => onDeleteUser()}
          type="User"
        />
      )}
      {hasCreateNewModal && (
        <CreateNewModal
          hasCreateNewModal={hasCreateNewModal}
          handleCreateNewModal={() =>
            setHasCreateNewModal((hasCreateNewModal) => !hasCreateNewModal)
          }
          handleClick={(teamName) => onCreateNewRole(teamName)}
          type="Role"
          placeHolder="Enter Here"
          spaceId={spaceId}
        />
      )}
      {hasInviteModal && (
        <InviteModal
          hasInviteModal={hasInviteModal}
          handleInviteModal={() => {
            setHasInviteModal(false)
            setselectedRole(null)
          }}
          type="Role"
          updateList={() => {
            // users && setUserFlag((userFlag) => !userFlag);
            setFlag((flag) => !flag)
            setAccordianID(null)
            setUsers(null)
          }}
          current={selectedRole}
        />
      )}
      {hasAclModal && (
        <AclSidebarModal
          hasAclModal={hasAclModal}
          handleAclModal={() => {
            setHasAclModal(false)
            setCurrentRole(null)
          }}
          modalTitle={`${
            currentRole?.name ? `${currentRole?.name}'s` : 'Role'
          } Access Management`}
        >
          <Policies
            hasTabMenu
            tabMenus={['Permissions', 'Member', 'Entities']}
            getExistingPolicies={process.env.ROLES_LIST_EXISTING_POL_GRP_SUBS}
            getPoliciesToAdd={process.env.ROLES_LIST_TO_ADD_POL_GRP_SUBS}
            current={currentRole}
            type="Role"
          />
        </AclSidebarModal>
      )}
    </div>
  )
}

export default Roles
