/* eslint-disable react/prop-types */
import React, { useState } from 'react'
import 'reactjs-popup/dist/index.css'
import AllPermissions from './all-permissions'
import AssignedPermissions from './assigned-permissions'

const Permissions = (props) => {
  const { current, type } = props
  // const BlockData = [
  //   { id: 1, title: 'AppBlock-1' },
  //   { id: 2, title: 'AppBlock-2' },
  // ];
  // const [accordianID, setAccordianID] = useState(null);
  const [permissionType, setPermissionType] = useState('assigned-permissions')

  return (
    <div className="float-left w-full">
      <div className="md-h-scroll-primary float-left mb-3 mt-4 flex w-full overflow-x-auto">
        <div className="border-ab-gray-dark float-left flex w-full space-x-3 border-b">
          <div
            onClick={() => {
              setPermissionType('assigned-permissions')
            }}
            className={`text-ab-sm relative -bottom-px flex cursor-pointer items-center justify-center border-b-2 px-3 py-2.5 font-medium ${
              permissionType === 'assigned-permissions'
                ? 'text-primary border-primary'
                : 'text-ab-black hover:text-primary border-transparent'
            }`}
          >
            <p className="whitespace-nowrap">Assigned Permissions</p>
          </div>
          <div
            onClick={() => {
              setPermissionType('all-permissions')
            }}
            className={`text-ab-sm relative -bottom-px flex cursor-pointer items-center justify-center border-b-2 px-3 py-2.5 font-medium ${
              permissionType === 'all-permissions'
                ? 'text-primary border-primary'
                : 'text-ab-black hover:text-primary border-transparent'
            }`}
          >
            <p className="whitespace-nowrap">All Permissions</p>
          </div>
        </div>
      </div>
      <div className="float-left w-full">
        {permissionType === 'assigned-permissions' && (
          <AssignedPermissions current={current} type={type} />
        )}
        {permissionType === 'all-permissions' && (
          <AllPermissions current={current} type={type} />
        )}
      </div>
      {/* <Pagination /> */}
    </div>
  )
}

export default Permissions
