/* eslint-disable no-shadow */
/* eslint-disable no-unused-expressions */
/* eslint-disable no-unused-vars */
import React, { useState, useEffect, useContext } from 'react'
import useOnclickOutside from 'react-cool-onclickoutside'
import * as dayjs from 'dayjs'
import CreateNewModal from '../common/modals/create-new-modal'
import MyContext from '../common/my-context'
import apiHelper from '../common/helpers/apiGetters'
import HorizontalDotIcon from '../../../assets/img/icons/horizontal-dot.svg'
import PlusPrimaryIcon from '../../../assets/img/icons/plus-primary.svg'
import DownArrow from '../../../assets/img/icons/down-arrow.svg'
import { ACLContext } from '../../../context/ACLContext'

const MyEntities = () => {
  // const urlparams = new URLSearchParams(location.search);
  // const params = Object.fromEntries(urlparams.entries());
  const { spaceId, entityTypeList } = useContext(MyContext)
  const [actionId, setActionId] = useState(null)
  const [isCreateNewEntity, setCreateNewEntity] = useState(false)
  const [myEntities, setMyEntities] = useState(null)
  const [current, setCurrent] = useState(null)
  const [flag, setFlag] = useState(false)
  const [loader, setLoader] = useState(true)
  const actionDropContainer = useOnclickOutside(() => {
    setActionId(null)
  })

  const entityTypeListMod = [
    {
      id: 'all',
      name: 'All Entities',
      display_name: 'All Entities',
    },
    ...entityTypeList,
  ]

  const [entityTypeDropdown, setEntityTypeDropdown] = useState(false)
  const [entityType, setEntityType] = useState(entityTypeListMod[0])
  const entityTypeDropContainer = useOnclickOutside(() => {
    setEntityTypeDropdown(false)
  })

  const filterData = () => ({
    search_keyword: null,
    entity_types: entityType.id === 'all' ? [] : [+entityType.id],
    page_limit: 100,
    offset: 0,
    sort_column: 'CreatedAt',
    sort_direction: 'ASC',
  })

  useEffect(async () => {
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.USER_LIST_AVAILABLE_ENTITIES_URL,
      value: filterData(),
      spaceId,
    })
    res && setMyEntities(res?.data)
    setLoader(false)
  }, [flag, entityType])

  const handleActionId = (id) => {
    if (actionId === id) {
      setActionId(null)
    } else {
      setActionId(id)
    }
  }

  const onCreateNewEntity = async (name, type_id) => {
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: !current
        ? process.env.CREATE_ENTITY_URL
        : process.env.UPDATE_APP_URL,
      value: !current
        ? {
            name,
            space_id: spaceId,
            type_id,
          }
        : {
            id: current?.id,
            name,
          },
      showSuccessMessage: true,
      spaceId,
    })
    setCurrent(null)
    setCreateNewEntity(false)
    setFlag((flag) => !flag)
  }

  return (
    <div className="float-left w-full">
      <div className="border-ab-gray-dark mb-2 flex items-center justify-between border-b py-2.5">
        <p className="text-lg font-medium text-black">My Entities</p>
        <div>
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
                {entityTypeListMod?.map((entityTyp) => (
                  <li
                    key={entityTyp.id}
                    onClick={() => {
                      setEntityType(entityTyp)
                      setEntityTypeDropdown((x) => !x)
                    }}
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
      </div>
      <div className="float-left mt-1 flex w-full flex-wrap">
        {!loader &&
          myEntities?.map((entities) => (
            <div
              key={entities.entity_id}
              className="border-ab-gray-dark hover:bg-ab-gray-light group relative float-left mt-2 mr-3 flex w-32 flex-col items-center border p-4"
            >
              <div className="bg-primary/20 flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full">
                <span className="text-primary text-2xl font-semi/bold capitalize">
                  {entities.label[0]}
                </span>
              </div>
              {/* <img src={} className="w-12 h-12 flex-shrink-0 rounded-full border-ab-gray-medium border object-cover" alt=""/> */}
              <p className="text-ab-sm text-ab-black mt-3 w-full truncate text-center font-medium">
                {entities.label}
              </p>
              <p className="text-ab-black mt-3 w-full text-center text-xs">
                Created on
              </p>
              <p className="text-ab-sm text-ab-black w-full text-center">
                {dayjs(entities.created_at).format('DD/MM/YYYY')}
              </p>
              {/* <div
                className={`absolute top-1 right-1 group-hover:visible ${
                  actionId !== entities.entity_id && 'invisible'
                }`}
              >
                <div ref={actionDropContainer} className="relative float-left">
                  <div
                    onClick={() => handleActionId(entities.entity_id)}
                    className={`float-left flex h-5 w-5 cursor-pointer items-center justify-center rounded-full hover:bg-[#E2E2E2] ${
                      actionId === entities.entity_id ? 'bg-[#E2E2E2]' : ''
                    }`}
                  >
                    <img src={HorizontalDotIcon} alt="" />
                  </div>
                  {actionId === entities.entity_id && (
                    <div className="shadow-box border-ab-gray-dark absolute top-full left-0 z-20 mt-1 min-w-[130px] max-w-[160px] rounded-md border bg-white py-2">
                      <span
                        onClick={() => {
                          setCurrent(entities)
                          setCreateNewEntity(true)
                        }}
                        className="text-ab-black hover:text-primary float-left w-full cursor-pointer truncate px-3 py-1 text-xs font-medium focus:outline-none"
                      >
                        Edit
                      </span>
                      <span className='text-ab-black hover:text-primary float-left w-full cursor-pointer truncate px-3 py-1 text-xs font-medium focus:outline-none'>
                        Delete
                      </span>
                    </div>
                  )}
                </div>
              </div> */}
            </div>
          ))}
        {!loader && !myEntities && (
          <p className="text-ab-black float-left w-full py-10 text-center text-sm">
            No Entities Found!
          </p>
        )}
        <div
          onClick={() => setCreateNewEntity(true)}
          className="border-ab-gray-medium trasnsition-all hover:bg-primary hover:border-primary group mt-2 flex w-32 cursor-pointer flex-col items-center justify-center border p-4 duration-300"
        >
          <span className="bg-primary-light trasnsition-all text-primary float-left flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full text-2xl font-semibold duration-200 group-hover:bg-white">
            <img src={PlusPrimaryIcon} alt="Plus Icon" />
          </span>
          <p className="text-ab-sm text-primary mt-2.5 w-full text-center group-hover:text-white">
            Create New Entity
          </p>
        </div>
      </div>
      {isCreateNewEntity && (
        <CreateNewModal
          hasCreateNewModal={isCreateNewEntity}
          handleCreateNewModal={() => {
            setCreateNewEntity(false)
            setCurrent(null)
          }}
          current={current}
          handleClick={onCreateNewEntity}
          type="Entity"
          placeHolder="Enter Here"
          spaceId={spaceId}
        />
      )}
    </div>
  )
}

export default MyEntities
