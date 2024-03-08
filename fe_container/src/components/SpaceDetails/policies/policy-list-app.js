/* eslint-disable no-unsafe-optional-chaining */
/* eslint-disable react/no-unstable-nested-components */
/* eslint-disable react/prop-types */
import React, { useState, useEffect, useCallback, useContext } from 'react'
import { debounce } from 'lodash'
import Popup from 'reactjs-popup'
import useOnclickOutside from 'react-cool-onclickoutside'
import apiHelper from '../common/helpers/apiGetters'
import MyContext from '../common/my-context'
import Pagination from '../../Layout/pagination/pagination'
import 'reactjs-popup/dist/index.css'
import DownArrow from '../../../assets/img/icons/down-arrow.svg'
import Chip from './chip'

const LIST_API_MAPPING = {
  Member: process.env.USER_LIST_ENTITIES_URL,
  Team: process.env.TEAMS_LIST_ENTITIES_URL,
  Role: process.env.ROLES_LIST_ENTITIES_URL,
}

const ADD_API_MAPPING = {
  Member: process.env.USER_ADD_ENTITIES,
  Team: process.env.TEAMS_ADD_ENTITIES,
  Role: process.env.ROLES_ADD_ENTITIES,
}

const PolicyListApp = (props) => {
  const { current, type } = props
  const { spaceId, entityTypeList } = useContext(MyContext)
  const pageLimit = Number(process.env.PAGE_LIMIT)
  const chipsToDisplay = 2
  const payloadType = {
    ...(type === 'Member' && { user_id: current?.user_id }),
    ...(type === 'Team' && { team_id: current?.team_id }),
    ...(type === 'Role' && { role_id: current?.role_id }),
  }

  const [appData, setAppData] = useState(null)
  const [selectedApp, setSelectedApp] = useState(null)
  const [flag, setFlag] = useState(false)
  const [totalCount, setTotalCount] = useState(null)
  const [selectedPage, setSelectedPage] = useState(0)
  const [selectedPolicies, setSelectedPolicies] = useState({
    ...payloadType,
    deleted_entity_mappings: [],
    new_entity_mappings: [],
  })

  const [entityTypeDropdown, setEntityTypeDropdown] = useState(false)
  const [entityType, setEntityType] = useState(entityTypeList[0])
  const entityTypeDropContainer = useOnclickOutside(() => {
    setEntityTypeDropdown(false)
  })

  const [loader, setLoader] = useState(false)

  const getFilterData = () => ({
    ...payloadType,
    offset: 0,
    limit: pageLimit,
    sort_column: 'updatedAt',
    sort_direction: 'desc',
    type_id: Number(entityType?.id),
  })

  const getListEntities = async (filterData) => {
    setLoader(true)
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: LIST_API_MAPPING[type],
      value: filterData,
      spaceId,
    })
    setAppData(res?.entities || null)
    setTotalCount(res.count || 0)
    setLoader(false)
  }

  useEffect(async () => {
    setEntityTypeDropdown(false)
    getListEntities(getFilterData())
  }, [flag, entityType])

  const onSubmit = async () => {
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: ADD_API_MAPPING[type],
      value: { ...selectedPolicies, type_id: Number(entityType?.id) },
      spaceId,
    })
    setSelectedApp(null)
    setSelectedPolicies({
      ...payloadType,
      deleted_entity_mappings: [],
      new_entity_mappings: [],
    })
    setFlag((x) => !x)
  }

  const onPolicySelect = (policyData) => {
    const filteredPolicies = selectedPolicies
    if (policyData.entity_mapping_id) {
      if (
        filteredPolicies.deleted_entity_mappings.some(
          (obj) => obj.entity_mapping_id === policyData.entity_mapping_id
        )
      ) {
        filteredPolicies.deleted_entity_mappings =
          filteredPolicies.deleted_entity_mappings.filter(
            (item) => item.entity_mapping_id !== policyData.entity_mapping_id
          )
      } else
        filteredPolicies.deleted_entity_mappings.push(
          policyData.entity_mapping_id
        )
    } else if (
      filteredPolicies.new_entity_mappings.some(
        (obj) => obj.ac_pol_grp_id === policyData.pol_grp_id
      )
    ) {
      filteredPolicies.new_entity_mappings =
        filteredPolicies.new_entity_mappings.filter(
          (item) => item.ac_pol_grp_id !== policyData.pol_grp_id
        )
    } else
      filteredPolicies.new_entity_mappings.push({
        entity_id: selectedApp.entity_id,
        ac_pol_grp_id: policyData.pol_grp_id,
      })
    setSelectedPolicies(filteredPolicies)
  }

  const handler = useCallback(
    debounce((searchtext) => {
      setAppData(null)
      const arg = {
        ...getFilterData(),
        conditions: {
          is_keyword_search: !!searchtext,
          keyword: searchtext,
        },
      }
      getListEntities(arg)
    }, 1000),
    []
  )

  const onSearchTextChange = (e) => {
    handler(e.target.value)
  }

  const handlePageChange = (event) => {
    const { selected } = event
    setSelectedPage(selected)
    const arg = {
      ...getFilterData(),
      offset: pageLimit * selected,
    }
    getListEntities(arg)
  }

  return (
    <div className="float-left w-full">
      <div className="text-ab-black float-left mb-3 w-full">
        <label className="float-left text-xs font-medium">Entity Type</label>
        <div
          className="relative float-left my-1.5 w-full"
          ref={entityTypeDropContainer}
        >
          <div
            onClick={() => setEntityTypeDropdown((x) => !x)}
            className={`text-ab-sm bg-ab-gray-light focus:border-primary float-left flex w-full cursor-pointer select-none items-center justify-between rounded-md border py-0.5 px-2.5 focus:outline-none ${
              entityTypeDropdown ? 'border-primary' : 'border-ab-gray-light'
            }`}
          >
            <div className="text-ab-black text-ab-sm flex-grow overflow-hidden">
              <p className="flex h-[32px] items-center px-1 text-xs font-medium">
                {entityType?.name || 'Select Entity Type'}
              </p>
            </div>
            <img
              className={`ml-3 flex-shrink-0 transform transition-transform duration-300 ${
                entityTypeDropdown && 'rotate-180'
              }`}
              src={DownArrow}
              alt=""
            />
          </div>
          <div
            className={`shadow-box dropDownFade border-ab-gray-medium position-inverse-2 absolute top-full left-0 z-10 mt-1 w-full border bg-white p-1 py-3 ${
              entityTypeDropdown ? '' : 'hidden'
            }`}
          >
            <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-2">
              {entityTypeList?.map((entityTyp) => (
                <li
                  key={entityTyp.id}
                  onClick={() => setEntityType(entityTyp)}
                  className={`float-left mb-4 w-full last-of-type:mb-0 cursor-pointer items-center leading-normal truncate text-xs font-medium tracking-tight hover:text-primary ${
                    entityTyp.id === entityType.id
                      ? 'text-primary'
                      : 'text-ab-black'
                  }`}
                >
                  {entityTyp.name}
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
      <input
        type="text"
        onChange={onSearchTextChange}
        className="search-input border-ab-gray-dark text-ab-sm h-10 w-full rounded-md !bg-[length:14px_14px] px-2 pl-9 focus:outline-none"
        placeholder={`Search for ${entityType.name}`}
      />
      <div className="border-ab-gray-dark float-left mt-3.5 w-full border">
        <div className="text-ab-black bg-ab-gray-light text-ab-sm float-left w-full">
          <p className="float-left w-5/12 p-3 font-medium">Entity Name</p>
          <p className="float-left w-7/12 p-3 font-medium">Policies</p>
        </div>
        {!loader && !appData && (
          <span className="text-ab-black float-left w-full py-10 text-center text-sm">
            No Entities Found
          </span>
        )}
        {appData?.map((item) => (
          <div key={item?.entity_id} className="float-left w-full">
            <div className="border-ab-gray-dark float-left w-full border-t">
              <div className="float-left w-5/12 px-3 py-2">
                <p className="text-ab-black text-ab-sm my-1 py-1 font-medium">
                  {item.label}
                </p>
                {/* <p className='text-primary bg-ab-disabled-yellow float-left clear-both my-1 max-w-full truncate rounded-full py-1 px-2 text-xs font-medium'>
                  Team Name
                </p> */}
              </div>
              <div className="float-left w-7/12">
                <div className="flex cursor-pointer items-center justify-between space-x-2 px-3 py-2">
                  <div className="flex-grow overflow-hidden">
                    {item.policy_groups
                      ?.filter((policy) => policy.entity_mapping_id !== null)
                      ?.slice(0, chipsToDisplay)
                      ?.map(
                        (policy, index) =>
                          policy.entity_mapping_id !== null && (
                            // eslint-disable-next-line react/no-array-index-key
                            <Chip key={index} data={policy.pol_grp_name} />
                          )
                      )}

                    {item?.policy_groups?.filter(
                      (policy) => policy.entity_mapping_id !== null
                    ).length > chipsToDisplay && (
                      <div className="float-left">
                        <Popup
                          arrow={false}
                          keepTooltipInside
                          trigger={() => (
                            <span className="bg-primary float-left my-1 flex h-6 min-w-[24px] cursor-pointer items-center justify-center rounded-full py-1 px-1 text-[10px] font-medium leading-tight text-white">
                              +
                              {item?.policy_groups?.filter(
                                (policy) => policy.entity_mapping_id !== null
                              ).length - chipsToDisplay}
                            </span>
                          )}
                          position="bottom right"
                          className="ab-dropdown-popup max-w-270"
                        >
                          <div className="border-ab-gray-dark shadow-box dropdownFade border-ab-gray-dark shadow-box float-left mt-3 w-full max-w-[270px] rounded-b-md border bg-white px-2 py-2.5">
                            <div className="float-left w-full bg-white">
                              <div className="float-left w-full p-2">
                                <p className="text-ab-black border-ab-gray-dark float-left w-full border-b pb-2 text-xs font-medium">
                                  Policies
                                </p>
                              </div>
                              <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-2">
                                {item?.policy_groups
                                  ?.slice(
                                    chipsToDisplay,
                                    item?.policy_groups?.length
                                  )
                                  ?.map(
                                    (policy) =>
                                      policy.entity_mapping_id !== null && (
                                        <li
                                          key={policy?.pol_grp_id}
                                          className="text-ab-black float-left mb-4 w-full truncate text-xs font-medium tracking-tight last-of-type:mb-0"
                                        >
                                          {policy.pol_grp_name}
                                        </li>
                                      )
                                  )}
                              </ul>
                            </div>
                          </div>
                        </Popup>
                      </div>
                    )}
                  </div>
                  <Popup
                    open={selectedApp === item}
                    arrow={false}
                    keepTooltipInside
                    onOpen={() => setSelectedApp(item)}
                    onClose={() => {
                      setSelectedApp(null)
                      setSelectedPolicies({
                        ...payloadType,
                        deleted_entity_mappings: [],
                        new_entity_mappings: [],
                      })
                    }}
                    trigger={(open) => (
                      <div className="float-left flex-shrink-0 py-2 px-1">
                        <img
                          className={`flex-shrink-0 transform transition-all ${
                            open ? 'rotate-180' : ''
                          }`}
                          alt=""
                          src={DownArrow}
                        />
                      </div>
                    )}
                    position="bottom right"
                    className="ab-dropdown-popup max-w-270"
                  >
                    <div className="border-ab-gray-dark shadow-box dropdownFade border-ab-gray-dark shadow-box float-left mt-3 w-full max-w-[270px] rounded-b-md border bg-white px-2 py-2.5">
                      <div className="float-left w-full bg-white">
                        <div className="float-left w-full p-2">
                          <p className="text-ab-black border-ab-gray-dark float-left w-full border-b pb-2 text-xs font-medium">
                            Attach Policies
                          </p>
                        </div>
                        {/* <input
                        type='text'
                        className='search-input border-ab-gray-dark text-ab-sm h-9 w-full rounded-md !bg-[length:14px_14px] px-2 pl-9 focus:outline-none'
                        placeholder='Search for Policies'
                      /> */}
                        <ul className="custom-scroll-primary float-left mt-2 max-h-[150px] w-full overflow-y-auto p-2">
                          {item?.policy_groups?.map((policy) => (
                            <li
                              key={policy.pol_grp_id}
                              className="float-left mb-4 w-full last-of-type:mb-0"
                            >
                              <label className="float-left flex max-w-full cursor-pointer items-center leading-normal">
                                <input
                                  className="peer hidden"
                                  type="checkbox"
                                  onChange={() => onPolicySelect(policy)}
                                  disabled={
                                    type === 'Member' &&
                                    (policy?.role_id !== null ||
                                      policy?.team_id !== null)
                                  }
                                  defaultChecked={
                                    policy.entity_mapping_id !== null
                                  }
                                />
                                <span className="chkbox-icon border-ab-disabled float-left mr-2 h-5 w-5 flex-shrink-0 cursor-pointer rounded border bg-white" />
                                <p className="text-ab-black truncate text-xs font-medium tracking-tight">
                                  {policy.pol_grp_name}
                                </p>
                                {(policy?.team_name || policy?.role_name) && (
                                  <span
                                    className={`text-primary inline m-1 whitespace-nowrap rounded-full py-1 px-2 text-xs font-medium ${
                                      policy?.team_id
                                        ? 'bg-ab-orange'
                                        : 'bg-ab-yellow'
                                    }`}
                                  >
                                    {policy?.team_name || policy?.role_name}
                                  </span>
                                )}
                              </label>
                            </li>
                          ))}
                        </ul>
                        <div className="float-left my-2 flex w-full items-center px-2">
                          <button
                            type="button"
                            onClick={onSubmit}
                            className="btn-secondary disabled:bg-ab-disabled mr-3 rounded px-3 py-2 text-xs font-bold leading-tight text-white transition-all"
                          >
                            Save Changes
                          </button>
                        </div>
                      </div>
                    </div>
                  </Popup>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
      <Pagination
        Padding={0}
        marginTop={1}
        pageCount={Math.ceil(totalCount / pageLimit)}
        handlePageChange={handlePageChange}
        selected={selectedPage}
      />
    </div>
  )
}

export default PolicyListApp
