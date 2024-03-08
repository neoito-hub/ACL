/* eslint-disable no-shadow */
/* eslint-disable react/jsx-no-constructed-context-values */
/* eslint-disable no-unused-expressions */
import React, { useState, useEffect, useContext } from 'react'
import { useLocation } from 'react-router-dom'
import MyContext from '../common/my-context'
import apiHelper from '../common/helpers/apiGetters'
import Toast from '../common/toast/Toast'
import 'reactjs-popup/dist/index.css'
import MyEntities from '../my-entities/my-entities'
import Members from '../members/members'
import Roles from '../roles/roles'
import Teams from '../teams/teams'
import Settings from '../settings/settings'
import Invites from '../invites/Invites'
import { ACLContext } from '../../../context/ACLContext'

const tabs = () => [
  {
    name: 'My Entities',
    url: '/spaces/my-entities',
  },
  {
    name: 'Members',
    url: '/spaces/members',
  },
  {
    name: 'Role Management',
    url: '/spaces/roles',
  },
  {
    name: 'Teams',
    url: '/spaces/teams',
  },
  {
    name: 'Invites',
    url: '/spaces/invites',
  },
  {
    name: 'Settings',
    url: '/spaces/settings',
  },
]

const SpaceDetails = () => {
  const location = useLocation()

  const { spaceId, spaceDetails, setSpaceDetails } = useContext(ACLContext)

  const [updateFlag, setUpdateFlag] = useState(true)
  const [entityTypeList, setEntityTypeList] = useState([])

  const [activeTab, setActiveTab] = useState(
    tabs().find((tab) => location?.pathname.includes(tab?.url))
  )

  useEffect(() => {
    setActiveTab(tabs().find((tab) => location?.pathname.includes(tab?.url)))
  }, [location])

  const handleClick = (tab) => {
    setActiveTab(tab)
    window.history.pushState(null, '', tab.url)
  }

  const getSpaceDetails = async () => {
    const res =
      spaceId &&
      (await apiHelper({
        baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
        subUrl: process.env.GET_SPACE_BY_ID_URL,
        value: { space_id: spaceId },
        spaceId,
      }))
    res && setSpaceDetails(res)
  }

  const getEntityType = async () => {
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.LIST_ENTITY_DEFINITION,
    })
    res && setEntityTypeList(res)
  }

  useEffect(() => {
    getEntityType()
  }, [])

  useEffect(async () => {
    spaceId && getSpaceDetails()
  }, [spaceId, updateFlag])

  return (
    <MyContext.Provider value={{ spaceId, entityTypeList }}>
      <Toast />
      <div className="float-left w-full max-w-5xl py-6 md-lt:px-4">
        {spaceDetails && (
          <div className="float-left mt-3 flex w-full flex-col items-start">
            <div className="bg-secondary float-left flex h-16 w-16 flex-shrink-0 items-center justify-center overflow-hidden rounded-full">
              <span className="text-3xl font-bold text-white capitalize">
                {spaceDetails?.name[0]}
              </span>
              {/* <img className="border-ab-gray-medium h-full w-full rounded-full border object-cover" src={} alt=""></img> */}
            </div>
            <p className="text-ab-3xl mt-5 font-semibold text-black/80">
              {spaceDetails?.name}
            </p>
          </div>
        )}
        {spaceDetails && (
          <div className="float-left mt-7 w-full">
            <div className="md-h-scroll-primary float-left flex w-full overflow-x-auto">
              <div className="border-ab-gray-dark float-left flex w-full space-x-3 border-b">
                {tabs()?.map((menu) => (
                  <button
                    type="button"
                    key={menu.url}
                    onClick={() => handleClick(menu)}
                    className={`text-ab-sm relative -bottom-px flex cursor-pointer items-center justify-center border-b px-3 py-2.5 font-medium ${
                      activeTab?.url === menu.url
                        ? 'text-primary border-primary'
                        : 'text-ab-black hover:text-primary border-transparent'
                    }`}
                  >
                    <p className="whitespace-nowrap">{menu.name}</p>
                  </button>
                ))}
              </div>
            </div>
            <div className="float-left w-full overflow-x-hidden py-6">
              {activeTab?.url === '/spaces/my-entities' && <MyEntities />}
              {activeTab?.url === '/spaces/members' && <Members />}
              {activeTab?.url === '/spaces/roles' && <Roles />}
              {activeTab?.url === '/spaces/teams' && <Teams />}
              {activeTab?.url === '/spaces/invites' && <Invites />}
              {activeTab?.url === '/spaces/settings' && (
                <Settings
                  spaceDetails={spaceDetails}
                  onUpdateSpace={() =>
                    setUpdateFlag((updateFlag) => !updateFlag)
                  }
                />
              )}
            </div>
          </div>
        )}
      </div>
    </MyContext.Provider>
  )
}

export default SpaceDetails
