/* eslint-disable react/prop-types */
import React, { useState } from 'react'
import PolicyList from './policy-list'
import PolicyListApp from './policy-list-app'
import Permissions from './permissions'

const Policies = (props) => {
  const {
    hasTabMenu,
    tabMenus,
    getExistingPolicies,
    getPoliciesToAdd,
    current,
    type,
  } = props
  const [AclTabActive, setAclTabActive] = useState(tabMenus[0])

  return (
    <div className="float-left mt-4 w-full">
      <p className="float-left mb-6 w-full text-sm">Available Policies</p>
      {hasTabMenu && (
        <div className="md-h-scroll-primary float-left mb-3 flex w-full overflow-x-auto">
          <div className="border-ab-gray-dark float-left flex w-full space-x-3 border-b">
            {tabMenus.map((item) => (
              <div
                key={item}
                onClick={() => {
                  setAclTabActive(item)
                  // setNewPolicy(false);
                }}
                className={`text-ab-sm relative -bottom-px flex flex-1 cursor-pointer items-center justify-center border-b px-3 py-2.5 font-medium ${
                  AclTabActive === item
                    ? 'text-primary border-primary'
                    : 'text-ab-black hover:text-primary border-transparent'
                }`}
              >
                <p className="whitespace-nowrap">{item}</p>
              </div>
            ))}
          </div>
        </div>
      )}
      {(AclTabActive === 'Member' ||
        AclTabActive === 'Team' ||
        AclTabActive === 'Role') && (
        <PolicyList
          getExistingPolicies={getExistingPolicies}
          getPoliciesToAdd={getPoliciesToAdd}
          current={current}
          type={type}
        />
      )}
      {AclTabActive === 'Entities' && (
        <PolicyListApp current={current} type={type} />
      )}
      {AclTabActive === 'Permissions' && (
        <Permissions current={current} type={type} />
      )}
    </div>
  )
}

export default Policies
