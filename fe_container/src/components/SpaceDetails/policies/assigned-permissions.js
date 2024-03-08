/* eslint-disable no-unsafe-optional-chaining */
/* eslint-disable no-unused-vars */
/* eslint-disable react/no-unstable-nested-components */
/* eslint-disable camelcase */
/* eslint-disable no-shadow */
/* eslint-disable react/prop-types */
import React, { useState, useEffect, useContext } from 'react'
import Popup from 'reactjs-popup'
import 'reactjs-popup/dist/index.css'
import DownArrow from '../../../assets/img/icons/down-arrow.svg'
import apiHelper from '../common/helpers/apiGetters'
import MyContext from '../common/my-context'
import Entities from './entities'
// import EntityPermission from './entity-permission';

const LIST_PERMISSIONS_API_MAPPING = {
  Role: process.env.ROLES_LIST_PERMISSIONS,
  Team: process.env.TEAMS_LIST_PERMISSIONS,
  Member: process.env.USER_LIST_PERMISSIONS,
}
const ADD_PERMISSIONS_API_MAPPING = {
  Role: process.env.ROLES_ADD_PERMISSIONS,
  Team: process.env.TEAMS_ADD_PERMISSIONS,
  Member: process.env.USER_ADD_PERMISSIONS,
}

const AssignedPermissions = (props) => {
  const { current, type } = props
  const { spaceId, entityTypeList } = useContext(MyContext)
  const chipesToDisplay = 5

  const [permissionList, setPermissionList] = useState(null)
  const [flag, setFlag] = useState(false)
  const [loading, setLoading] = useState(false)
  const [accordianID, setAccordianID] = useState(null)

  const payloadType = {
    ...(type === 'Member' && { user_id: current?.user_id }),
    ...(type === 'Team' && { team_id: current?.team_id }),
    ...(type === 'Role' && { role_id: current?.role_id }),
  }

  const filterDataStructure = () => ({
    ...payloadType,
    search_keyword: '',
    page_limit: 100,
    offset: 0,
    sort_column: 'PolicyCount',
    sort_direction: 'ASC',
  })

  const getUserListPermissions = async () => {
    setLoading(true)
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: LIST_PERMISSIONS_API_MAPPING[type],
      value: filterDataStructure(),
      spaceId,
    })
    setPermissionList(res?.data)
    setLoading(false)
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

  const handlePermissionChange = async (e) => {
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: ADD_PERMISSIONS_API_MAPPING[type],
      value: {
        ...payloadType,
        permissions: [
          {
            permission_id: e.target.value,
            added_entities: [],
            added_space_access_entities: [],
            deleted_entities: [],
            deleted_space_access_entities: [],
            is_delete: true,
          },
        ],
      },
      showSuccessMessage: true,
      spaceId,
    })
    setFlag((flag) => !flag)
  }

  return (
    <div className="float-left w-full">
      <div className="float-left mt-1 w-full">
        {permissionList?.map((item) => (
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
                    defaultChecked
                    onClick={handlePermissionChange}
                  />
                  <span className="chkbox-icon border-ab-disabled float-left h-5 w-5 flex-shrink-0 cursor-pointer rounded border bg-white" />
                </label>
                <div className="float-left flex flex-col overflow-hidden">
                  <p className="max-w-full truncate text-sm font-medium">
                    {item.name}
                  </p>
                  <div className="app-info-group text-ab-black/50 mt-1.5 flex flex-wrap text-xs font-medium">
                    {Object.keys(item?.entity_types).map((key) => (
                      <span key={key}>
                        {item?.attached_entities?.SpaceAccessEntities[key]
                          ?.length
                          ? 'All'
                          : item?.attached_entities?.AddedEntities[key]
                              ?.length || '0'}{' '}
                        {entityTypeList.find((x) => x.id === key)?.display_name}
                      </span>
                    ))}
                  </div>
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
                  alt=" "
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
                  {item?.policy_groups
                    ?.slice(0, chipesToDisplay)
                    .map((policy_group) => (
                      <p
                        key={policy_group?.subs_id}
                        className="text-primary bg-ab-disabled-yellow float-left my-1.5 mr-2.5 max-w-full truncate rounded-full py-1 px-2 text-xs font-medium"
                      >
                        {policy_group?.policy_group_name}
                      </p>
                    ))}
                  {item?.policy_groups?.length > chipesToDisplay && (
                    <div className="float-left">
                      <Popup
                        arrow={false}
                        keepTooltipInside
                        trigger={(open) => (
                          <span className="float-left my-1.5 flex h-6 min-w-[24px] cursor-pointer items-center justify-center rounded-full bg-[#9747FF] py-1 px-1 text-xs font-medium leading-tight text-white">
                            +{item?.policy_groups?.length - chipesToDisplay}
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
                                ?.slice(
                                  chipesToDisplay,
                                  item?.policy_groups?.length,
                                )
                                .map((policy_group) => (
                                  <li
                                    key={policy_group?.subs_id}
                                    className="text-ab-black float-left mb-4 w-full truncate text-xs font-medium tracking-tight last-of-type:mb-0"
                                  >
                                    {policy_group?.policy_group_name}
                                  </li>
                                ))}
                            </ul>
                          </div>
                        </div>
                      </Popup>
                    </div>
                  )}
                </div>
                {item?.entity_types &&
                  Object.values(item?.entity_types).some((value) => value) &&
                  accordianID === item.permission_id && (
                    <Entities
                      type={type}
                      spaceId={spaceId}
                      current={current}
                      currentPage="assigned-permissions"
                      updateList={() => setFlag((flag) => !flag)}
                      currentPermission={item}
                      closeAccordian={() => handleAccordian(item.permission_id)}
                    />
                  )}
              </div>
            </div>
          </div>
        ))}
        {permissionList?.length && (
          <p className="text-ab-black/60 float-left mt-4 mb-2 w-full text-xs font-medium">
            Note: Uncheck the item to unassign the permission
          </p>
        )}
        {!permissionList?.length && !loading && (
          <p className="text-ab-black float-left w-full py-10 text-center text-sm">
            No Permissions Assigned
          </p>
        )}
      </div>
    </div>
  )
}

export default AssignedPermissions
