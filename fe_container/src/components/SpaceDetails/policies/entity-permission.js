/* eslint-disable no-unsafe-optional-chaining */
/* eslint-disable no-unused-expressions */
/* eslint-disable react/prop-types */
import React, { useState, useCallback, useEffect, useContext } from 'react'
import useOnclickOutside from 'react-cool-onclickoutside'
import Popup from 'reactjs-popup'
import { debounce } from 'lodash'
import DownArrow from '../../../assets/img/icons/down-arrow.svg'
import MyContext from '../common/my-context'
import apiHelper from '../common/helpers/apiGetters'

const EntityPermission = (props) => {
  const { type, selectedBn, updateSelectedBn, selection, updateSelection } =
    props
  const chipesToDisplay = 4

  const { spaceId } = useContext(MyContext)
  const [bnDropdown, setBnDropdown] = useState(false)

  const [entityList, setEntityList] = useState([])
  const [loader, setLoader] = useState(false)

  const getUserListEntities = async (entitytype, searchtext = null) => {
    setLoader(true)
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.USER_LIST_AVAILABLE_ENTITIES,
      value: {
        search_keyword: searchtext,
        entity_types: [entitytype],
        page_limit: 100,
        offset: 0,
        sort_column: 'CreatedAt',
        sort_direction: 'ASC',
      },
      spaceId,
    })
    setEntityList(res?.data)
    setLoader(false)
  }

  useEffect(() => {
    getUserListEntities(+type?.id)
  }, [type])

  const bnDropContainer = useOnclickOutside(() => {
    setBnDropdown(false)
  })

  const handleBn = (e) => {
    const updatedSelectedBn = e.target.checked
      ? [...selectedBn, e.target.value]
      : selectedBn.filter((item) => item !== e.target.value)

    updateSelectedBn(updatedSelectedBn)
  }
  const handleBnRemove = (e, val) => {
    e.stopPropagation()
    const updatedSelectedBn = selectedBn.filter((item) => item !== val)
    updateSelectedBn(updatedSelectedBn)
  }

  const handler = useCallback(
    debounce((text) => {
      getUserListEntities(+type?.id, text)
    }, 1000),
    [type],
  )

  const onSearchTextChange = (e) => {
    handler(e.target.value)
  }

  return (
    <>
      <div className="float-left mt-3 mb-3 flex w-full items-center">
        <label className="text-ab-black float-left mr-4 mt-2 flex cursor-pointer items-center text-sm font-medium">
          <input
            type="checkbox"
            checked={selection === 'none'}
            onChange={(e) => {
              e.target.checked && updateSelection('none')
            }}
            className="peer hidden"
            name="block-selection"
          />
          <span className="border-ab-gray-dark peer-checked:border-primary float-left mr-2 h-5 w-5 rounded-full border peer-checked:border-4" />
          None
        </label>
        <label className="text-ab-black float-left mr-4 mt-2 flex cursor-pointer items-center text-sm font-medium">
          <input
            type="checkbox"
            checked={selection === 'all'}
            onChange={(e) => {
              e.target.checked && updateSelection('all')
            }}
            className="peer hidden"
            name="block-selection"
          />
          <span className="border-ab-gray-dark peer-checked:border-primary float-left mr-2 h-5 w-5 rounded-full border peer-checked:border-4" />
          All Entities
        </label>
        <label className="text-ab-black float-left mr-3 mt-2 flex cursor-pointer items-center text-sm font-medium">
          <input
            type="checkbox"
            checked={selection === 'custom'}
            onChange={(e) => {
              e.target.checked && updateSelection('custom')
            }}
            className="peer hidden"
            name="block-selection"
          />
          <span className="border-ab-gray-dark peer-checked:border-primary float-left mr-2 h-5 w-5 rounded-full border peer-checked:border-4" />
          Custom Entities
        </label>
      </div>
      {selection === 'custom' && (
        <div className="relative float-left mt-3 w-full" ref={bnDropContainer}>
          <div
            onClick={() => setBnDropdown(!bnDropdown)}
            className={`text-ab-sm bg-ab-gray-light focus:border-primary float-left flex h-12 w-full cursor-pointer select-none items-center justify-between rounded-md border py-0.5 px-2.5 focus:outline-none ${
              bnDropdown ? 'border-primary' : 'border-ab-gray-light'
            }`}
          >
            <div className="flex-grow overflow-hidden">
              <p className="text-ab-black text-ab-sm flex items-center truncate px-1 text-xs font-medium">
                Select Entities
              </p>
            </div>
            <img
              className={`ml-3 flex-shrink-0 transform transition-transform duration-300 ${
                bnDropdown && 'rotate-180'
              }`}
              src={DownArrow}
              alt=""
            />
          </div>
          <div
            className={`shadow-box dropDownFade border-ab-gray-medium position-inverse-2 absolute top-full left-0 z-10 mt-1 w-full border bg-white p-1 py-3 ${
              bnDropdown ? '' : 'hidden'
            }`}
          >
            <div className="border-ab-gray-dark float-left mb-2 w-full px-2">
              <input
                type="text"
                onChange={onSearchTextChange}
                className="search-input border-ab-gray-dark text-ab-sm h-10 w-full rounded-md border !bg-[length:14px_14px] px-2 pl-9 focus:outline-none"
                placeholder="Search Entities"
              />
            </div>
            {loader ? (
              <div className="flex justify-center items-center w-full pb-1">
                <span className="text-ab-black float-left w-full text-center text-sm">
                  Loading ...
                </span>
              </div>
            ) : entityList?.length ? (
              <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-2">
                {entityList?.map((item) => (
                  <li
                    key={item.entity_id}
                    className="float-left mb-4 w-full last-of-type:mb-0"
                  >
                    <label className="float-left flex max-w-full cursor-pointer items-center leading-normal">
                      <input
                        checked={selectedBn?.includes(item?.entity_id)}
                        onChange={(e) => handleBn(e)}
                        name="team"
                        value={item?.entity_id}
                        className="peer hidden"
                        type="checkbox"
                      />
                      <span className="chkbox-icon border-ab-disabled float-left mr-2 h-5 w-5 flex-shrink-0 cursor-pointer rounded border bg-white" />
                      <p className="text-ab-black truncate text-xs font-medium tracking-tight">
                        {item?.label}
                      </p>
                    </label>
                  </li>
                ))}
              </ul>
            ) : (
              <div className="flex justify-center items-center w-full pb-1">
                <span className="text-ab-black float-left w-full text-center text-sm">
                  No Entities Found
                </span>
              </div>
            )}
          </div>
        </div>
      )}
      {selection === 'custom' && (
        <div className="float-left mt-3 flex w-full flex-wrap items-center">
          {entityList
            ?.filter((item) => selectedBn?.includes(item?.entity_id))
            ?.slice(0, chipesToDisplay)
            ?.map((entity) => (
              <div
                key={entity?.entity_id}
                className="bg-primary/10 float-left my-1 mr-2 inline-flex max-w-full items-center space-x-2 rounded-full py-0.5 px-3"
              >
                <p className="text-primary truncate text-xs font-medium leading-[20px]">
                  {entity?.label}
                </p>
                <svg
                  onClick={(e) => handleBnRemove(e, entity?.entity_id)}
                  width="8"
                  height="8"
                  viewBox="0 0 8 8"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M0.799988 0.799805L7.19999 7.19981M0.799988 7.19981L7.19999 0.799805"
                    stroke="#000"
                    strokeWidth="1.53333"
                  />
                </svg>
              </div>
            ))}
          {entityList?.filter((item) => selectedBn?.includes(item?.entity_id))
            ?.length > chipesToDisplay && (
            <div className="float-left">
              <Popup
                arrow={false}
                keepTooltipInside
                // eslint-disable-next-line react/no-unstable-nested-components
                trigger={() => (
                  <span className="float-left flex h-6 min-w-[24px] cursor-pointer items-center justify-center rounded-full bg-[#9747FF] px-1 text-xs font-medium leading-tight text-white">
                    +{selectedBn?.length - chipesToDisplay}
                  </span>
                )}
                position="bottom right"
                className="ab-dropdown-popup max-w-270"
              >
                <div className="border-ab-gray-dark shadow-box dropdownFade border-ab-gray-dark shadow-box float-left mt-3 w-full max-w-[270px] rounded-b-md border bg-white px-2 py-2.5">
                  <div className="float-left w-full bg-white">
                    <div className="float-left w-full p-2">
                      <p className="text-ab-black border-ab-gray-dark float-left w-full border-b pb-2 text-xs font-medium">
                        Selected {type?.display_name}s
                      </p>
                    </div>
                    <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-2">
                      {entityList
                        ?.filter((item) =>
                          selectedBn?.includes(item?.entity_id),
                        )
                        ?.slice(chipesToDisplay, selectedBn?.length)
                        .map((entity) => (
                          <li
                            key={entity?.entity_id}
                            className="text-ab-black float-left mb-4 w-full truncate text-xs font-medium tracking-tight last-of-type:mb-0"
                          >
                            {entity?.label}
                          </li>
                        ))}
                    </ul>
                  </div>
                </div>
              </Popup>
            </div>
          )}
        </div>
      )}
    </>
  )
}

export default EntityPermission
