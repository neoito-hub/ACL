/* eslint-disable react/no-unstable-nested-components */
/* eslint-disable react/button-has-type */
/* eslint-disable react/no-array-index-key */
/* eslint-disable no-shadow */
/* eslint-disable no-unused-expressions */
/* eslint-disable camelcase */
/* eslint-disable react/prop-types */
import React, { useState, useEffect, useCallback, useContext } from 'react'
import { debounce } from 'lodash'
import Popup from 'reactjs-popup'
import apiHelper from '../common/helpers/apiGetters'
import MyContext from '../common/my-context'
import AddNewPolicy from './add-new-policy'
import Pagination from '../../Layout/pagination/pagination'

const PolicyList = (props) => {
  const { current, getExistingPolicies, getPoliciesToAdd, type } = props
  const { spaceId } = useContext(MyContext)
  const page_limit = Number(process.env.PAGE_LIMIT)

  const chipsToDisplay = 2

  const [policiesLoader, setPoliciesLoader] = useState(false)
  const [policiesToAddLoader, setPoliciesToAddLoader] = useState(false)
  const [newPolicy, setNewPolicy] = useState(false)
  const [policies, setPolicies] = useState(null)
  const [policiesTotalCount, setPoliciesTotalCount] = useState(null)
  const [policiesSelectedPage, setPoliciesSelectedPage] = useState(0)
  const [policiesToAdd, setPoliciesToAdd] = useState(null)
  const [policiesToAddSelectedPage, setPoliciesToAddSelectedPage] = useState(0)
  const [policiesToAddTotalCount, setPoliciesToAddTotalCount] = useState(null)
  const [selectedPolicies, setSelectedPolicies] = useState([])
  const [flag, setFlag] = useState(false)
  const [policiesFilterData, setPoliciesFilterData] = useState({
    space_id: spaceId,
    search_keyword: null,
    ...(type === 'Member' && { user_id: current?.user_id }),
    ...(type === 'Team' && { team_id: current?.team_id }),
    ...(type === 'Role' && { role_id: current?.role_id }),
    page_limit,
    offset: 0,
  })
  const [policiesToAddFilterData, setPoliciesToAddFilterData] = useState({
    space_id: spaceId,
    search_keyword: null,
    ...(type === 'Member' && { user_id: current?.user_id }),
    ...(type === 'Team' && { team_id: current?.team_id }),
    ...(type === 'Role' && { role_id: current?.role_id }),
    page_limit,
    offset: 0,
  })

  useEffect(async () => {
    setPoliciesLoader(true)
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: getExistingPolicies,
      value: policiesFilterData,
      spaceId,
    })
    setPolicies(res ? res.data : null)
    setSelectedPolicies(res ? res.data : [])
    setPoliciesTotalCount(res ? res.total_count : 0)
    setPoliciesLoader(false)
  }, [flag, policiesFilterData])

  useEffect(async () => {
    setPoliciesToAddLoader(true)
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: getPoliciesToAdd,
      value: policiesToAddFilterData,
      spaceId,
    })
    setPoliciesToAdd(res ? res.data : null)
    setPoliciesToAddTotalCount(res ? res.total_count : 0)
    setPoliciesToAddLoader(false)
  }, [policiesToAddFilterData])

  const handlePoliciesPageChange = (event) => {
    const { selected } = event
    setPoliciesSelectedPage(selected)
    const arg = {
      ...policiesFilterData,
      offset: policiesFilterData.page_limit * selected,
    }
    setPoliciesFilterData(arg)
  }

  const handlePoliciesToAddPageChange = (event) => {
    const { selected } = event
    setPoliciesToAddSelectedPage(selected)
    const arg = {
      ...policiesToAddFilterData,
      offset: policiesToAddFilterData.page_limit * selected,
    }
    setPoliciesToAddFilterData(arg)
  }

  const handler = useCallback(
    debounce((searchtext) => {
      setPoliciesToAdd(null)
      const arg = {
        ...policiesToAddFilterData,
        search_keyword: searchtext,
      }
      setPoliciesToAddFilterData(arg)
    }, 1000),
    [],
  )

  const onSearchTextChange = (e) => {
    handler(e.target.value)
  }

  const onPolicyAdd = async (id) => {
    setPoliciesLoader(true)
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl:
        type === 'Member'
          ? process.env.USER_ADD_POL_GRP_SUBS
          : type === 'Team'
            ? process.env.TEAMS_ADD_POL_GRP_SUBS
            : type === 'Role' && process.env.ROLES_ADD_POL_GRP_SUBS,
      value: {
        space_id: spaceId,
        ...(type === 'Member' && { user_id: current?.user_id }),
        ...(type === 'Team' && { team_id: current?.team_id }),
        ...(type === 'Role' && { role_id: current?.role_id }),
        ac_pol_grp_ids: [id],
      },
      showSuccessMessage: true,
      spaceId,
    })
    res &&
      setPoliciesToAdd((current) =>
        current.map((obj) => {
          if (obj.ac_pol_grp_id === id) {
            return { ...obj, subs_id: res[0].id }
          }
          return obj
        }),
      )
    setFlag((flag) => !flag)
  }

  const onPolicyRemove = async (id) => {
    setPoliciesLoader(true)
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl:
        type === 'Member'
          ? process.env.USER_DELETE_EXISTING_POL_GRP_SUBS
          : type === 'Team'
            ? process.env.TEAMS_DELETE_EXISTING_POL_GRP_SUBS
            : type === 'Role' && process.env.ROLES_DELETE_EXISTING_POL_GRP_SUBS,
      value: {
        id,
      },
      apiType: 'delete',
      showSuccessMessage: true,
      spaceId,
    })
    setPoliciesToAdd((current) =>
      current.map((obj) => {
        if (obj.subs_id === id) {
          return { ...obj, subs_id: '' }
        }
        return obj
      }),
    )
    setFlag((flag) => !flag)
  }

  const onPolicyChange = (e, policyData) => {
    if (e.target.checked) {
      setSelectedPolicies([...selectedPolicies, policyData])
      onPolicyAdd(policyData.ac_pol_grp_id)
    } else {
      const updatedSelectedPolicies = selectedPolicies?.filter(
        (item) => item !== policyData,
      )
      onPolicyRemove(policyData.subs_id)
      setSelectedPolicies(updatedSelectedPolicies)
    }
  }

  return (
    <div className="float-left w-full">
      <div className="border-ab-gray-dark custom-h-scroll-primary mt-3.5 w-full overflow-x-auto border">
        <table className="text-ab-black w-full text-left">
          <thead>
            <tr className="bg-ab-gray-light">
              <th className="text-ab-sm whitespace-nowrap p-3 font-medium">
                Policy
              </th>
              <th className="text-ab-sm whitespace-nowrap p-3 font-medium">
                Type
              </th>
              {/* <th className='w-12 min-w-[48px] p-3'></th> */}
            </tr>
          </thead>
          <tbody>
            {policies?.map((policy, index) => (
              <tr
                key={index}
                className="text-ab-sm border-ab-gray-dark border-t"
              >
                <td className="px-3 py-2 text-xs">
                  <div className="flex flex-wrap items-center">
                    <span className="my-1 pr-2">{policy?.name}</span>{' '}
                    <div className="flex flex-wrap items-center">
                      <p className="float-left mr-2">
                        {[...(policy?.teams || []), ...(policy?.roles || [])]
                          ?.slice(0, chipsToDisplay)
                          ?.map((team) => (
                            <span
                              key={team?.team_id || team?.role_id}
                              className={`text-primary m-1 whitespace-nowrap rounded-full py-1 px-2 text-xs font-medium ${team?.team_id ? 'bg-ab-orange' : 'bg-ab-yellow'}`}
                            >
                              {team?.team_name || team?.role_name}
                            </span>
                          ))}
                      </p>
                      {[...(policy?.teams || []), ...(policy?.roles || [])]
                        .length > chipsToDisplay && (
                        <div className="float-left">
                          <Popup
                            arrow={false}
                            keepTooltipInside
                            trigger={() => (
                              <span className="bg-primary float-left my-1 flex h-6 min-w-[24px] cursor-pointer items-center justify-center rounded-full py-1 px-1 text-[10px] font-medium leading-tight text-white">
                                +
                                {[
                                  ...(policy?.teams || []),
                                  ...(policy?.roles || []),
                                ].length - chipsToDisplay}
                              </span>
                            )}
                            position="bottom center"
                            className="ab-dropdown-popup max-w-270"
                          >
                            <div className="border-ab-gray-dark shadow-box dropdownFade border-ab-gray-dark shadow-box float-left mt-4 w-full max-w-[270px] rounded-b-md border bg-white px-2 py-2.5">
                              <div className="float-left w-full bg-white">
                                <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-2">
                                  {[
                                    ...(policy?.teams || []),
                                    ...(policy?.roles || []),
                                  ]
                                    ?.slice(chipsToDisplay)
                                    ?.map((team) => (
                                      <li key={team?.team_id || team?.role_id}>
                                        <span
                                          className={`text-primary inline m-1 whitespace-nowrap rounded-full py-1 px-2 text-xs font-medium ${team?.team_id ? 'bg-ab-orange' : 'bg-ab-yellow'}`}
                                        >
                                          {team?.team_name || team?.role_name}
                                        </span>
                                      </li>
                                    ))}
                                </ul>
                              </div>
                            </div>
                          </Popup>
                        </div>
                      )}
                    </div>
                  </div>
                </td>
                <td className="p-3 text-xs">
                  {policy?.is_predefined ? 'Predefined' : 'Custom'}
                </td>
                {/* <td className='w-12 min-w-[48px] p-3 text-center'>
                    <svg
                      className='fill-ab-black hover:fill-ab-red cursor-pointer'
                      width='20'
                      height='20'
                      viewBox='0 0 20 20'
                      fill='none'
                      xmlns='http://www.w3.org/2000/svg'
                    >
                      <path d='M7.75 13.2083L10.0208 10.9167L12.3125 13.2083L12.8542 12.6667L10.5833 10.375L12.8542 8.0625L12.3125 7.52083L10.0208 9.8125L7.75 7.52083L7.1875 8.0625L9.47917 10.375L7.1875 12.6667L7.75 13.2083ZM6.33333 16.75C5.95833 16.75 5.63195 16.6111 5.35417 16.3333C5.07639 16.0556 4.9375 15.7292 4.9375 15.3542V4.77083H4.10417V4.04167H7.39583V3.4375H12.625V4.04167H15.9167V4.77083H15.0833V15.3542C15.0833 15.7292 14.9479 16.0556 14.6771 16.3333C14.4062 16.6111 14.0764 16.75 13.6875 16.75H6.33333ZM14.3542 4.77083H5.66667V15.3542C5.66667 15.5208 5.73611 15.6736 5.875 15.8125C6.01389 15.9514 6.16667 16.0208 6.33333 16.0208H13.6875C13.8542 16.0208 14.0069 15.9514 14.1458 15.8125C14.2847 15.6736 14.3542 15.5208 14.3542 15.3542V4.77083ZM5.66667 4.77083V16.0208V15.3542V4.77083Z' />
                    </svg>
                  </td> */}
              </tr>
            ))}
            {!policiesLoader && !policies && (
              <tr className="flex justify-center items-center">
                <td className="text-ab-black float-left w-full py-10 text-center text-sm">
                  No Policies Found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      <Pagination
        Padding={0}
        marginTop={1}
        pageCount={Math.ceil(policiesTotalCount / page_limit)}
        handlePageChange={handlePoliciesPageChange}
        selected={policiesSelectedPage}
      />
      {!newPolicy && (
        <button
          onClick={() => setNewPolicy(true)}
          className="btn-primary text-ab-sm mt-3 px-5 py-2 font-bold leading-normal focus:outline-none"
        >
          Add New Policy
        </button>
      )}
      {newPolicy && (
        <AddNewPolicy
          loader={policiesToAddLoader}
          policiesLoader={policiesLoader}
          policies={policies}
          policiesToAdd={policiesToAdd}
          policiesSelectedPage={policiesToAddSelectedPage}
          policiesTotalCount={policiesToAddTotalCount}
          handlePageChange={handlePoliciesToAddPageChange}
          onSearchTextChange={onSearchTextChange}
          onPolicyChange={onPolicyChange}
          selectedPolicies={selectedPolicies}
        />
      )}
    </div>
  )
}

export default PolicyList
