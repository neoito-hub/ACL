/* eslint-disable no-unsafe-optional-chaining */
/* eslint-disable react/no-unstable-nested-components */
/* eslint-disable no-unused-vars */
/* eslint-disable camelcase */
/* eslint-disable no-shadow */
/* eslint-disable react/prop-types */
import React, { useState, useEffect, useContext, useCallback } from 'react'
import Popup from 'reactjs-popup'
import { debounce } from 'lodash'
import 'reactjs-popup/dist/index.css'
import DownArrow from '../../../assets/img/icons/down-arrow.svg'
import apiHelper from '../common/helpers/apiGetters'
import MyContext from '../common/my-context'
// import EntityPermission from './entity-permission';
import Entities from './entities'

const LIST_AVAILABLE_PERMISSIONS_API_MAPPING = {
  Member: process.env.USER_LIST_AVAILABLE_PERMISSIONS,
  Team: process.env.TEAMS_LIST_AVAILABLE_PERMISSIONS,
  Role: process.env.ROLES_LIST_AVAILABLE_PERMISSIONS,
}

const AllPermissions = (props) => {
  const { current, type } = props
  const { spaceId } = useContext(MyContext)
  const [permissionList, setPermissionList] = useState(null)
  const [flag, setFlag] = useState(false)
  const [accordianID, setAccordianID] = useState(null)
  // const [permissionCategory, setPermissionCategory] =
  //   useState('permission-blocks');
  const [searchText, setSearchText] = useState(null)

  const payloadType = {
    ...(type === 'Member' && { user_id: current?.user_id }),
    ...(type === 'Team' && { team_id: current?.team_id }),
    ...(type === 'Role' && { role_id: current?.role_id }),
  }

  const filterDataStructure = () => ({
    ...payloadType,
    search_keyword: searchText,
    page_limit: 100,
    offset: 0,
    sort_column: 'PolicyCount',
    sort_direction: 'ASC',
  })

  const getUserListPermissions = async () => {
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: LIST_AVAILABLE_PERMISSIONS_API_MAPPING[type],
      value: filterDataStructure(),
      spaceId,
    })
    setPermissionList(res?.data)
  }

  useEffect(() => {
    getUserListPermissions()
  }, [flag])

  const handleAccordian = (id) => {
    if (id === accordianID) {
      setAccordianID(null)
    } else {
      setAccordianID(id)
    }
  }

  const handler = useCallback(
    debounce((text) => {
      setPermissionList(null)
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
      <div className="float-left mt-1 w-full">
        <input
          type="text"
          onChange={onSearchTextChange}
          className="search-input border-ab-gray-dark text-ab-black bg-ab-gray-dark h-10 w-full rounded-md border !bg-[length:14px_14px] px-2 pl-9 text-sm focus:outline-none"
          placeholder="Search permissions"
        />
        <p className="text-ab-black/60 float-left mt-4 mb-3 w-full text-xs font-medium">
          Note: Select to assign the permission
        </p>
        {!permissionList?.length ? (
          <p className="text-ab-black float-left w-full py-10 text-center text-sm">
            No Permissions Found!
          </p>
        ) : (
          permissionList?.map((item) => (
            <div
              key={item.permission_id}
              className="border-ab-gray-dark float-left mb-2 w-full border p-3 last-of-type:mb-0"
            >
              <div className="float-left flex w-full justify-between space-x-2">
                <div className="float-left flex flex-grow space-x-3 overflow-hidden">
                  <label className="float-left flex flex-shrink-0 cursor-pointer">
                    <input
                      className="peer hidden"
                      type="checkbox"
                      value={item.permission_id}
                      checked={accordianID === item.permission_id}
                      onChange={(e) => {
                        // if (
                        //   item?.entity_types &&
                        //   Object.keys(item?.entity_types).some(
                        //     (key) =>
                        //       key !== '0' && item?.entity_types[key] === true
                        //   )
                        // )
                        handleAccordian(item.permission_id)
                        // else handlePermissionChange(e)
                      }}
                    />
                    <span className="chkbox-icon border-ab-disabled float-left h-5 w-5 flex-shrink-0 cursor-pointer rounded border bg-white" />
                  </label>
                  <div className="float-left flex flex-col overflow-hidden">
                    <p className="max-w-full truncate text-sm font-medium">
                      {item.name}
                    </p>
                  </div>
                </div>
                <div
                  onClick={() => {
                    handleAccordian(item.permission_id)
                  }}
                  className="float-left flex-shrink-0 cursor-pointer py-2 px-1"
                >
                  <img
                    className={`flex-shrink-0 transform transition-all ${
                      accordianID === item.permission_id ? 'rotate-180' : ''
                    }`}
                    src={DownArrow}
                    alt=""
                  />
                </div>
              </div>
              <div
                className={`float-left w-full ${
                  accordianID !== item.permission_id && 'hidden'
                }`}
              >
                <div className="order-ab-gray-dark float-left mt-3 w-full border-t p-3">
                  <p className="float-left mb-1.5 text-sm font-medium text-black">
                    Associated Policies
                  </p>
                  <div className="float-left w-full">
                    {item?.policy_groups?.slice(0, 5).map((policy_group) => (
                      <p
                        key={policy_group?.policy_group_id}
                        className="text-primary bg-ab-disabled-yellow float-left my-1.5 mr-2.5 max-w-full truncate rounded-full py-1 px-2 text-xs font-medium"
                      >
                        {policy_group?.policy_group_display_name}
                      </p>
                    ))}
                    {item?.policy_groups?.length > 5 && (
                      <div className="float-left">
                        <Popup
                          arrow={false}
                          keepTooltipInside
                          trigger={(_open) => (
                            <span className="float-left my-1.5 flex h-6 min-w-[24px] cursor-pointer items-center justify-center rounded-full bg-[#9747FF] py-1 px-1 text-xs font-medium leading-tight text-white">
                              +{item?.policy_groups?.length - 5}
                            </span>
                          )}
                          position="bottom right"
                          className="ab-dropdown-popup max-w-270"
                        >
                          <div className="border-ab-gray-dark shadow-box dropdownFade border-ab-gray-dark shadow-box float-left mt-3 w-full max-w-[270px] rounded-b-md border bg-white px-2 py-2.5">
                            <div className="float-left w-full bg-white">
                              <div className="float-left w-full p-2">
                                <p className="text-ab-black border-ab-gray-dark float-left w-full border-b pb-2 text-xs font-medium">
                                  Associated Policies
                                </p>
                              </div>
                              <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-2">
                                {item?.policy_groups
                                  ?.slice(3, item?.policy_groups?.length)
                                  .map((policy_group) => (
                                    <li
                                      key={policy_group?.policy_group_id}
                                      className="text-ab-black float-left mb-4 w-full truncate text-xs font-medium tracking-tight last-of-type:mb-0"
                                    >
                                      {policy_group?.policy_group_display_name}
                                    </li>
                                  ))}
                              </ul>
                            </div>
                          </div>
                        </Popup>
                      </div>
                    )}
                  </div>
                  {
                    // Object.values(item?.entity_types).some((value) => value) &&
                    accordianID === item.permission_id && (
                      <Entities
                        type={type}
                        current={current}
                        currentPage="all-permissions"
                        updateList={() => setFlag((flag) => !flag)}
                        currentPermission={item}
                        closeAccordian={() =>
                          handleAccordian(item.permission_id)
                        }
                      />
                    )
                  }
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  )
}

export default AllPermissions
