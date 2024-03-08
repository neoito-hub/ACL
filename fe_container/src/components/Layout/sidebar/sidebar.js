/* eslint-disable no-unused-expressions */
import React, { useState, useEffect, useCallback, useContext } from 'react'
import { useLocation, Link } from 'react-router-dom'
import { shield } from '@appblocks/js-sdk'
import axios from 'axios'
import { debounce } from 'lodash'
import { ACLContext } from '../../../context/ACLContext'

const Sidebar = (props) => {
  // eslint-disable-next-line react/prop-types
  const { activePage, handleActivePage } = props
  const location = useLocation()

  const { spaceId, setSpaceId, spaceData, setSpaceData, setSpaceDetails } =
    useContext(ACLContext)

  const [loader, setLoader] = useState(false)
  const [searchText, setSearchText] = useState(null)
  const [displayAllSpaces, setDisplayAllSpaces] = useState(false)

  useEffect(() => {
    const pathName = window?.location?.pathname
    const spaceMenu = [
      { pathname: '/', slug: 'home' },
      { pathname: '/spaces', slug: 'spaces' },
      { pathname: '/profile', slug: 'profile-settings' },
    ]
    handleActivePage(spaceMenu.find((menu) => pathName === menu.pathname)?.slug)
  }, [location])

  // eslint-disable-next-line consistent-return
  const apiHelper = async (baseUrl, subUrl, value = null, apiType = 'post') => {
    const token = shield.tokenStore.getToken()
    try {
      const { data } = await axios({
        method: apiType,
        url: `${baseUrl}${subUrl}`,
        data: value && value,
        headers: token && {
          Authorization: `Bearer ${token}`,
        },
      })
      return data?.data
    } catch (err) {
      console.log('msg', err)
      // if (err.response.status === 401) shield.logout()
    }
  }

  useEffect(async () => {
    // setSpaceData(null);
    setLoader(true)
    const res = await apiHelper(
      process.env.BLOCK_ENV_URL_API_BASE_URL,
      process.env.LIST_SPACES_URL,
      {
        search_keyword: searchText,
      },
    )
    const defSpace = res?.find((space) => space.is_default)
    !spaceId && setSpaceId(defSpace?.space_id)
    setSpaceData(res)
    setLoader(false)
  }, [searchText, spaceId])

  const handleSpaceChange = (id) => {
    if (id !== spaceId) setSpaceDetails(null)
    setSpaceId(id)
  }

  const handler = useCallback(
    debounce((text) => {
      setSearchText(text)
    }, 1000),
    [],
  )

  const onSearchTextChange = (e) => {
    handler(e.target.value)
  }

  return (
    <div className="custom-scroll-primary float-left flex h-full w-full flex-col justify-between overflow-auto pr-2">
      <div className="float-left w-full">
        <div className="float-left mb-2 w-full">
          <Link
            to="/"
            // onClick={() => handleActivePage('home')}
            className={`text-ab-sm hover:text-primary hover:bg-ab-gray-light group flex w-full cursor-pointer items-center space-x-2 py-2 px-3 font-semibold leading-normal transition-all ${
              activePage === 'home'
                ? 'text-primary bg-ab-gray-light'
                : 'text-ab-black'
            }`}
          >
            <svg
              className={`group-hover:fill-primary flex-shrink-0 ${
                activePage === 'home' ? 'fill-primary' : 'fill-[#484848]'
              }`}
              width="20"
              height="20"
              viewBox="0 0 20 20"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path d="M4.16675 16.542V8.77116L2.64591 9.93783L2.22925 9.37533L10.0001 3.41699L17.7917 9.35449L17.3542 9.93783L15.8334 8.75033V16.542H4.16675ZM4.89591 15.8128H15.1042V8.20866L10.0001 4.33366L4.89591 8.22949V15.8128ZM4.89591 15.8128H15.1042H4.89591ZM6.72925 12.3128C6.56258 12.3128 6.41328 12.2469 6.28133 12.1149C6.14939 11.983 6.08341 11.8337 6.08341 11.667C6.08341 11.5003 6.14939 11.351 6.28133 11.2191C6.41328 11.0871 6.56258 11.0212 6.72925 11.0212C6.88203 11.0212 7.02786 11.0871 7.16675 11.2191C7.30564 11.351 7.37508 11.5003 7.37508 11.667C7.37508 11.8337 7.30564 11.983 7.16675 12.1149C7.02786 12.2469 6.88203 12.3128 6.72925 12.3128ZM10.0001 12.3128C9.83341 12.3128 9.68411 12.2469 9.55216 12.1149C9.42022 11.983 9.35425 11.8337 9.35425 11.667C9.35425 11.5003 9.42022 11.351 9.55216 11.2191C9.68411 11.0871 9.83341 11.0212 10.0001 11.0212C10.1529 11.0212 10.2987 11.0871 10.4376 11.2191C10.5765 11.351 10.6459 11.5003 10.6459 11.667C10.6459 11.8337 10.5765 11.983 10.4376 12.1149C10.2987 12.2469 10.1529 12.3128 10.0001 12.3128ZM13.2709 12.3128C13.1042 12.3128 12.9549 12.2469 12.823 12.1149C12.6911 11.983 12.6251 11.8337 12.6251 11.667C12.6251 11.5003 12.6911 11.351 12.823 11.2191C12.9549 11.0871 13.1042 11.0212 13.2709 11.0212C13.4237 11.0212 13.5695 11.0871 13.7084 11.2191C13.8473 11.351 13.9167 11.5003 13.9167 11.667C13.9167 11.8337 13.8473 11.983 13.7084 12.1149C13.5695 12.2469 13.4237 12.3128 13.2709 12.3128Z" />
            </svg>
            <p className="truncate">Home</p>
          </Link>
        </div>
        <div className="float-left w-full">
          <div className="float-left mb-2 w-full">
            <Link
              to="/spaces"
              // onClick={() => handleActivePage('spaces')}
              className={`text-ab-sm hover:text-primary hover:bg-ab-gray-light group flex w-full cursor-pointer items-center space-x-2 py-2 px-3 font-semibold leading-normal transition-all ${
                activePage === 'spaces'
                  ? 'text-primary bg-ab-gray-light'
                  : 'text-ab-black'
              }`}
            >
              <svg
                className={`group-hover:fill-primary flex-shrink-0 ${
                  activePage === 'spaces' ? 'fill-primary' : 'fill-[#484848]'
                }`}
                width="20"
                height="20"
                viewBox="0 0 20 20"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path d="M4.70833 16.25C4.45833 16.25 4.23611 16.1528 4.04167 15.9583C3.84722 15.7639 3.75 15.5417 3.75 15.2917C3.75 15.0278 3.84722 14.8021 4.04167 14.6146C4.23611 14.4271 4.45833 14.3333 4.70833 14.3333C4.97222 14.3333 5.19792 14.4271 5.38542 14.6146C5.57292 14.8021 5.66667 15.0278 5.66667 15.2917C5.66667 15.5417 5.57292 15.7639 5.38542 15.9583C5.19792 16.1528 4.97222 16.25 4.70833 16.25ZM10 16.25C9.75 16.25 9.52778 16.1528 9.33333 15.9583C9.13889 15.7639 9.04167 15.5417 9.04167 15.2917C9.04167 15.0278 9.13889 14.8021 9.33333 14.6146C9.52778 14.4271 9.75 14.3333 10 14.3333C10.2639 14.3333 10.4896 14.4271 10.6771 14.6146C10.8646 14.8021 10.9583 15.0278 10.9583 15.2917C10.9583 15.5417 10.8646 15.7639 10.6771 15.9583C10.4896 16.1528 10.2639 16.25 10 16.25ZM15.2917 16.25C15.0417 16.25 14.8194 16.1528 14.625 15.9583C14.4306 15.7639 14.3333 15.5417 14.3333 15.2917C14.3333 15.0278 14.4306 14.8021 14.625 14.6146C14.8194 14.4271 15.0417 14.3333 15.2917 14.3333C15.5556 14.3333 15.7812 14.4271 15.9688 14.6146C16.1562 14.8021 16.25 15.0278 16.25 15.2917C16.25 15.5417 16.1562 15.7639 15.9688 15.9583C15.7812 16.1528 15.5556 16.25 15.2917 16.25ZM4.70833 10.9583C4.45833 10.9583 4.23611 10.8611 4.04167 10.6667C3.84722 10.4722 3.75 10.25 3.75 10C3.75 9.73611 3.84722 9.51042 4.04167 9.32292C4.23611 9.13542 4.45833 9.04167 4.70833 9.04167C4.97222 9.04167 5.19792 9.13542 5.38542 9.32292C5.57292 9.51042 5.66667 9.73611 5.66667 10C5.66667 10.25 5.57292 10.4722 5.38542 10.6667C5.19792 10.8611 4.97222 10.9583 4.70833 10.9583ZM10 10.9583C9.75 10.9583 9.52778 10.8611 9.33333 10.6667C9.13889 10.4722 9.04167 10.25 9.04167 10C9.04167 9.73611 9.13889 9.51042 9.33333 9.32292C9.52778 9.13542 9.75 9.04167 10 9.04167C10.2639 9.04167 10.4896 9.13542 10.6771 9.32292C10.8646 9.51042 10.9583 9.73611 10.9583 10C10.9583 10.25 10.8646 10.4722 10.6771 10.6667C10.4896 10.8611 10.2639 10.9583 10 10.9583ZM15.2917 10.9583C15.0417 10.9583 14.8194 10.8611 14.625 10.6667C14.4306 10.4722 14.3333 10.25 14.3333 10C14.3333 9.73611 14.4306 9.51042 14.625 9.32292C14.8194 9.13542 15.0417 9.04167 15.2917 9.04167C15.5556 9.04167 15.7812 9.13542 15.9688 9.32292C16.1562 9.51042 16.25 9.73611 16.25 10C16.25 10.25 16.1562 10.4722 15.9688 10.6667C15.7812 10.8611 15.5556 10.9583 15.2917 10.9583ZM4.70833 5.66667C4.45833 5.66667 4.23611 5.56944 4.04167 5.375C3.84722 5.18056 3.75 4.95833 3.75 4.70833C3.75 4.44444 3.84722 4.21875 4.04167 4.03125C4.23611 3.84375 4.45833 3.75 4.70833 3.75C4.97222 3.75 5.19792 3.84375 5.38542 4.03125C5.57292 4.21875 5.66667 4.44444 5.66667 4.70833C5.66667 4.95833 5.57292 5.18056 5.38542 5.375C5.19792 5.56944 4.97222 5.66667 4.70833 5.66667ZM10 5.66667C9.75 5.66667 9.52778 5.56944 9.33333 5.375C9.13889 5.18056 9.04167 4.95833 9.04167 4.70833C9.04167 4.44444 9.13889 4.21875 9.33333 4.03125C9.52778 3.84375 9.75 3.75 10 3.75C10.2639 3.75 10.4896 3.84375 10.6771 4.03125C10.8646 4.21875 10.9583 4.44444 10.9583 4.70833C10.9583 4.95833 10.8646 5.18056 10.6771 5.375C10.4896 5.56944 10.2639 5.66667 10 5.66667ZM15.2917 5.66667C15.0417 5.66667 14.8194 5.56944 14.625 5.375C14.4306 5.18056 14.3333 4.95833 14.3333 4.70833C14.3333 4.44444 14.4306 4.21875 14.625 4.03125C14.8194 3.84375 15.0417 3.75 15.2917 3.75C15.5556 3.75 15.7812 3.84375 15.9688 4.03125C16.1562 4.21875 16.25 4.44444 16.25 4.70833C16.25 4.95833 16.1562 5.18056 15.9688 5.375C15.7812 5.56944 15.5556 5.66667 15.2917 5.66667Z" />
              </svg>
              <p className="truncate">Spaces</p>
            </Link>
          </div>
          <div className="float-left w-full">
            <form autoComplete="off" className="mb-2">
              <input
                id="space_name"
                type="space_name"
                name="space_name"
                onChange={onSearchTextChange}
                className="search-input-xs border-ab-gray-dark focus:border-primary h-8 w-full border-b !bg-[length:12px_12px] px-2 pl-8 text-[13px] focus:outline-none"
                placeholder="Search spaces"
              />
            </form>

            <div className="w-full float-left mb-2 custom-scroll-primary max-h-[144px] overflow-auto">
              {spaceData &&
                spaceData.length > 0 &&
                spaceData
                  .slice(0, displayAllSpaces ? spaceData.length : 2)
                  .map((item) => (
                    <Link
                      to="/spaces/my-entities"
                      key={item.space_id}
                      onClick={() => handleSpaceChange(item.space_id)}
                      className="hover:bg-ab-gray-light float-left w-full"
                    >
                      <label className="text-ab-black flex cursor-pointer items-center overflow-hidden py-2 px-3">
                        {location?.pathname.includes('spaces/') && (
                          <input
                            id="space"
                            name="space"
                            className="peer hidden"
                            checked={item.space_id === spaceId}
                            type="radio"
                            onChange={() => {}}
                          />
                        )}
                        <span className="no-bg-chkbox float-left mr-2 h-5 w-5 flex-shrink-0 bg-white" />
                        <span className="bg-secondary mr-2 flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full text-xs font-bold text-white capitalize">
                          {item.space_name[0]}
                        </span>
                        <p className="text-ab-sm truncate font-semibold">
                          {item.space_name}
                        </p>
                      </label>
                    </Link>
                  ))}
              {!loader && !spaceData && (
                <div className="text-ab-sm font-noraml py-2 px-3">
                  No spaces found
                </div>
              )}
            </div>
            {spaceData && spaceData.length > 2 && (
              <div className="float-left mb-2 w-full">
                <button
                  type="button"
                  onClick={() => {
                    // eslint-disable-next-line no-shadow
                    setDisplayAllSpaces((displayAllSpaces) => !displayAllSpaces)
                  }}
                  className="text-primary float-left mb-3 mt-1 cursor-pointer truncate pl-10 text-xs font-semibold underline focus:outline-none"
                >
                  {displayAllSpaces ? 'See less' : 'See more'}
                </button>
              </div>
            )}
            <div className="float-left mb-2 w-full pl-10">
              <div
                // onClick={() => handleActivePage('create-new-space')}
                className="float-left flex max-w-full cursor-pointer items-center space-x-3"
              >
                <svg
                  className="flex-shrink-0"
                  width="14"
                  height="14"
                  viewBox="0 0 14 14"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M6.99998 0.0664062V13.9331M0.0666504 6.99974H13.9333"
                    stroke="#5E5EDD"
                  />
                </svg>
                <Link
                  to="/create-space"
                  className="text-ab-sm text-primary float-left truncate font-semibold"
                >
                  Create New Space
                </Link>
              </div>
            </div>
          </div>
        </div>

        <div className="float-left mb-2 w-full">
          <Link
            to="/profile"
            // onClick={() => handleActivePage('profile-settings')}
            className={`text-ab-sm hover:text-primary hover:bg-ab-gray-light group flex w-full cursor-pointer items-center space-x-2 py-2 px-3 font-semibold leading-normal transition-all ${
              activePage === 'profile-settings'
                ? 'text-primary bg-ab-gray-light'
                : 'text-ab-black'
            }`}
          >
            <svg
              className={`group-hover:fill-primary flex-shrink-0 ${
                activePage === 'profile-settings'
                  ? 'fill-primary'
                  : 'fill-[#484848]'
              }`}
              width="20"
              height="20"
              viewBox="0 0 20 20"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path d="M5.16667 14.7917C5.93056 14.2778 6.69444 13.8715 7.45833 13.5729C8.22222 13.2743 9.06944 13.125 10 13.125C10.9306 13.125 11.7778 13.2743 12.5417 13.5729C13.3056 13.8715 14.0694 14.2778 14.8333 14.7917C15.4306 14.1667 15.9063 13.4514 16.2604 12.6458C16.6146 11.8403 16.7917 10.9583 16.7917 10C16.7917 8.11111 16.1319 6.50694 14.8125 5.1875C13.4931 3.86806 11.8889 3.20833 10 3.20833C8.11111 3.20833 6.50694 3.86806 5.1875 5.1875C3.86806 6.50694 3.20833 8.11111 3.20833 10C3.20833 10.9583 3.38542 11.8403 3.73958 12.6458C4.09375 13.4514 4.56944 14.1667 5.16667 14.7917ZM10 10.2292C9.33333 10.2292 8.77431 10 8.32292 9.54167C7.87153 9.08333 7.64583 8.52778 7.64583 7.875C7.64583 7.20833 7.875 6.64931 8.33333 6.19792C8.79167 5.74653 9.34722 5.52083 10 5.52083C10.6667 5.52083 11.2257 5.75 11.6771 6.20833C12.1285 6.66667 12.3542 7.22222 12.3542 7.875C12.3542 8.54167 12.125 9.10069 11.6667 9.55208C11.2083 10.0035 10.6528 10.2292 10 10.2292ZM10 17.25C8.98611 17.25 8.04167 17.0625 7.16667 16.6875C6.29167 16.3125 5.52431 15.7951 4.86458 15.1354C4.20486 14.4757 3.6875 13.7083 3.3125 12.8333C2.9375 11.9583 2.75 11.0139 2.75 10C2.75 8.98611 2.9375 8.04167 3.3125 7.16667C3.6875 6.29167 4.20486 5.52431 4.86458 4.86458C5.52431 4.20486 6.29167 3.6875 7.16667 3.3125C8.04167 2.9375 8.98611 2.75 10 2.75C11.0139 2.75 11.9583 2.9375 12.8333 3.3125C13.7083 3.6875 14.4757 4.20486 15.1354 4.86458C15.7951 5.52431 16.3125 6.29167 16.6875 7.16667C17.0625 8.04167 17.25 8.98611 17.25 10C17.25 11.0139 17.0625 11.9583 16.6875 12.8333C16.3125 13.7083 15.7951 14.4757 15.1354 15.1354C14.4757 15.7951 13.7083 16.3125 12.8333 16.6875C11.9583 17.0625 11.0139 17.25 10 17.25ZM10 16.7917C10.7917 16.7917 11.5833 16.6493 12.375 16.3646C13.1667 16.0799 13.8611 15.6667 14.4583 15.125C13.8611 14.6528 13.184 14.2778 12.4271 14C11.6701 13.7222 10.8611 13.5833 10 13.5833C9.13889 13.5833 8.32639 13.7188 7.5625 13.9896C6.79861 14.2604 6.13194 14.6389 5.5625 15.125C6.14583 15.6667 6.83333 16.0799 7.625 16.3646C8.41667 16.6493 9.20833 16.7917 10 16.7917ZM10 9.77083C10.5417 9.77083 10.9931 9.59028 11.3542 9.22917C11.7153 8.86805 11.8958 8.41667 11.8958 7.875C11.8958 7.33333 11.7153 6.88194 11.3542 6.52083C10.9931 6.15972 10.5417 5.97917 10 5.97917C9.45833 5.97917 9.00694 6.15972 8.64583 6.52083C8.28472 6.88194 8.10417 7.33333 8.10417 7.875C8.10417 8.41667 8.28472 8.86805 8.64583 9.22917C9.00694 9.59028 9.45833 9.77083 10 9.77083Z" />
            </svg>
            <span className="truncate">Profile Settings</span>
          </Link>
        </div>
      </div>
    </div>
  )
}

export default Sidebar
