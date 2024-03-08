/* eslint-disable no-param-reassign */
/* eslint-disable no-return-assign */
/* eslint-disable no-shadow */
/* eslint-disable array-callback-return */
/* eslint-disable no-unused-expressions */
/* eslint-disable react/prop-types */
import React, { useState, useEffect, useContext } from 'react'
import useOnclickOutside from 'react-cool-onclickoutside'
import EntityPermission from './entity-permission'
import apiHelper from '../common/helpers/apiGetters'
import DownArrow from '../../../assets/img/icons/down-arrow.svg'
import MyContext from '../common/my-context'

const ADD_PERMISSIONS_API_MAPPING = {
  Member: process.env.USER_ADD_PERMISSIONS,
  Team: process.env.TEAMS_ADD_PERMISSIONS,
  Role: process.env.ROLES_ADD_PERMISSIONS,
}

const Entities = (props) => {
  const {
    type,
    current,
    currentPage,
    updateList,
    currentPermission,
    closeAccordian,
  } = props

  const { spaceId, entityTypeList } = useContext(MyContext)

  const payloadType = {
    ...(type === 'Member' && { user_id: current?.user_id }),
    ...(type === 'Team' && { team_id: current?.team_id }),
    ...(type === 'Role' && { role_id: current?.role_id }),
  }

  const showEntityTabs = Object.keys(currentPermission?.entity_types).length

  const getSelectedEntities = (attachedEntities) => {
    let selectedEntities = []
    attachedEntities &&
      Object.keys(attachedEntities)?.map((key) => {
        attachedEntities[key]?.map((item) => {
          selectedEntities = [...selectedEntities, item?.entity_id]
        })
      })
    return selectedEntities
  }

  const [entityTypeDropdown, setEntityTypeDropdown] = useState(false)
  const [entityType, setEntityType] = useState(
    entityTypeList.find(
      (x) => x.id === Object.keys(currentPermission?.entity_types)[0]
    )
  )
  const entityTypeDropContainer = useOnclickOutside(() => {
    setEntityTypeDropdown(false)
  })

  const [selectedBn, setSelectedBn] = useState(
    getSelectedEntities(currentPermission?.attached_entities?.AddedEntities)
  )

  const initialBlockSelectValue = entityTypeList?.reduce(
    (a, c) => ({
      ...a,
      [c.id]: !currentPermission?.attached_entities?.SpaceAccessEntities[c.id]
        ?.length
        ? !currentPermission?.attached_entities?.AddedEntities[c.id]?.length
          ? 'none'
          : 'custom'
        : 'all',
    }),
    {}
  )

  const [blockSelection, setBlockSelection] = useState(initialBlockSelectValue)
  const [blockList, setBlockList] = useState(null)
  const [showSaveButton, setShowSaveButton] = useState(false)

  const filterDataStructure = (type, searchText) => ({
    search_keyword: searchText,
    entity_types: [...type],
    page_limit: 100,
    offset: 0,
    sort_column: 'CreatedAt',
    sort_direction: 'ASC',
  })

  const getUserListEntities = async (type, searchText, updateList) => {
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.USER_LIST_AVAILABLE_ENTITIES,
      value: filterDataStructure(type, searchText),
      spaceId,
    })
    updateList(
      res?.data?.reduce((a, c) => {
        a[c.type] ??= []
        a[c.type].push(c)
        return a
      }, {})
    )
  }

  useEffect(() => {
    getUserListEntities([], null, setBlockList)
  }, [])

  useEffect(() => {
    setSelectedBn(
      getSelectedEntities(currentPermission?.attached_entities?.AddedEntities)
    )
  }, [currentPermission])

  const apiHandler = async (payload) => {
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: ADD_PERMISSIONS_API_MAPPING[type],
      value: payload,
      spaceId,
      showSuccessMessage: true,
    })
    updateList()
    closeAccordian()
  }

  const handleEnitiesChange = async () => {
    // check for added entities
    let addedEntities = []
    selectedBn?.map((item) => {
      entityTypeList?.map((entityType) =>
        blockList[entityType.id]?.some(
          (entity) =>
            entity?.entity_id === item &&
            !currentPermission?.attached_entities?.AddedEntities[
              entityType?.id
            ]?.some((item) => item?.entity_id === entity?.entity_id) &&
            blockSelection[entityType?.id] === 'custom' &&
            (addedEntities = [
              ...addedEntities,
              { id: item, type: Number(entityType?.id) },
            ])
        )
      )
    })

    // check for deleted entities
    let deletedEntities = []
    currentPermission?.attached_entities &&
      Object.keys(currentPermission?.attached_entities?.AddedEntities)?.map(
        (key) =>
          currentPermission?.attached_entities?.AddedEntities[key].map(
            (entity) =>
              (!selectedBn.includes(entity?.entity_id) ||
                blockSelection[key] === ('all' || 'none')) &&
              (deletedEntities = [
                ...deletedEntities,
                { id: entity?.entity_id, type: Number(key) },
              ])
          )
      )

    // check for added space entities
    let addedSpaceEntities = []
    entityTypeList.map(
      (spaceEntity) =>
        blockSelection[spaceEntity?.id] === 'all' &&
        (addedSpaceEntities = [
          ...addedSpaceEntities,
          {
            id: String(spaceEntity?.id),
            type: Number(spaceEntity?.id),
          },
        ])
    )

    // check for deleted space entities
    let deletedSpaceEntities = []
    entityTypeList.map(
      (spaceEntity) =>
        (blockSelection[spaceEntity?.id] === 'custom' ||
          blockSelection[spaceEntity?.id] === 'none') &&
        currentPermission?.attached_entities?.SpaceAccessEntities[
          spaceEntity?.id
        ]?.some((item) => item?.entity_type === Number(spaceEntity?.id)) &&
        (deletedSpaceEntities = [
          ...deletedSpaceEntities,
          {
            id: String(spaceEntity?.id),
            type: Number(spaceEntity?.id),
          },
        ])
    )

    const payload = {
      ...payloadType,
      permissions: [
        {
          permission_id: currentPermission?.permission_id,
          added_entities: addedEntities,
          added_space_access_entities: addedSpaceEntities,
          deleted_space_access_entities: deletedSpaceEntities,
          deleted_entities: deletedEntities,
          is_delete: false,
        },
      ],
    }
    apiHandler(payload)
  }

  function checkEntityIdInSelectedBn(listArray) {
    // Extract entity_id values from the array of objects
    const entityIdValues = listArray?.map((obj) => obj?.entity_id)

    // Check if any of the entityIdValues are present in the selectedBn
    const found = entityIdValues?.some((entityId) =>
      selectedBn?.includes(entityId)
    )
    return found
  }

  const handleDiscard = () => {
    const payload = {
      ...payloadType,
      permissions: [
        {
          permission_id: currentPermission?.permission_id,
          added_entities: [],
          added_space_access_entities: [],
          deleted_entities: [],
          deleted_space_access_entities: [],
          is_delete: true,
        },
      ],
    }
    apiHandler(payload)
  }

  const handleCancel = () => {
    setBlockSelection(initialBlockSelectValue)

    setSelectedBn(
      getSelectedEntities(currentPermission?.attached_entities?.AddedEntities)
    )
    setShowSaveButton(false)
  }

  return (
    <div className="float-left w-full">
      <div className="text-ab-black float-left w-full">
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
              {entityTypeList?.map(
                (entityTyp) =>
                  currentPermission?.entity_types[entityTyp.id] && (
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
                  )
              )}
            </ul>
          </div>
        </div>
      </div>
      {showEntityTabs && (
        <div className="float-left w-full">
          <EntityPermission
            type={entityType}
            selectedBn={selectedBn}
            updateSelectedBn={(arg) => {
              setSelectedBn(arg)
              !showSaveButton && setShowSaveButton(true)
            }}
            selection={blockSelection[entityType?.id]}
            updateSelection={(arg) => {
              !showSaveButton && setShowSaveButton(true)
              setBlockSelection((x) => ({ ...x, [entityType?.id]: arg }))
            }}
          />
        </div>
      )}

      <div className="border-ab-gray-dark float-left mt-3 flex w-full items-center justify-between border-t pt-3.5">
        <div className="float-left">
          {(currentPage === 'all-permissions' || showSaveButton) && (
            <button
              type="button"
              disabled={
                // add env check when integrating
                entityTypeList.every(
                  ({ id }) =>
                    blockSelection[id] === 'custom' &&
                    !checkEntityIdInSelectedBn(blockList[id])
                )
              }
              onClick={handleEnitiesChange}
              className="btn-secondary text-ab-sm disabled:bg-ab-disabled mr-4 rounded px-3 py-2 font-medium leading-normal text-white transition-all focus:outline-none min-w-[100px]"
            >
              {currentPage === 'all-permissions' ? 'Add' : 'Save Changes'}
            </button>
          )}
          {showSaveButton && (
            <button
              type="button"
              onClick={handleCancel}
              className="hover:text-ab-black text-ab-sm rounded px-3 py-1 font-medium leading-tight text-[#979696] focus:outline-none"
            >
              Cancel
            </button>
          )}
        </div>
        {currentPage === 'assigned-permissions' && (
          <div className="float-left">
            <button
              type="button"
              onClick={handleDiscard}
              className="btn-secondary bg-ab-red text-ab-sm disabled:bg-ab-disabled ml-4 rounded px-3 py-2 font-medium leading-normal text-white transition-all focus:outline-none"
            >
              Remove
            </button>
          </div>
        )}
      </div>
    </div>
  )
}

export default Entities
