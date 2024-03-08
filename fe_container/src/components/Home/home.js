import React, { useState, useEffect, useContext } from 'react'
import { shield } from '@appblocks/js-sdk'
import { Link } from 'react-router-dom'
import axios from 'axios'
import SpacesIcon from '../../assets/img/icons/spaces-icon.svg'
import BlogIcon from '../../assets/img/icons/blog.svg'
import DocIcon from '../../assets/img/icons/paper-file.svg'
import HomeSpaceListingSkeletonLoader from './homeSpaceListingSkeletonLoader'
import { ACLContext } from '../../context/ACLContext'
// import Checked from "../../assets/img/icons/check-line.svg";

const Home = () => {
  const { setSpaceId, setSpaceDetails, userDetails } = useContext(ACLContext)

  const loaderData = Array(7).fill(null)
  // const [warningBox, setWarningBox] = useState(true);
  const [loader, setLoader] = useState(false)
  const [spaceData, setSpaceData] = useState(null)
  const [numberOfDataToDisplay, setNumberOfDataToDisplay] = useState(10)
  const handleSpaceData = () => {
    setNumberOfDataToDisplay((x) => x + 10)
  }

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
      // if (err.response.status === 401) shield.logout();
    }
  }

  useEffect(async () => {
    setLoader(true)
    const res = await apiHelper(
      process.env.BLOCK_ENV_URL_API_BASE_URL,
      process.env.LIST_SPACES_URL
    )
    setSpaceData(res)
    setLoader(false)
  }, [])

  const handleSpaceChange = (spaceId) => {
    setSpaceDetails(null)
    setSpaceId(spaceId)
  }

  return (
    <div className="w-full float-left flex flex-col lg:flex-row lg:space-x-6 xl:space-x-10">
      <div className="float-left w-full max-w-3xl py-6 px-4 lg:px-0">
        <h4 className="text-ab-black md:mt-3 text-2xl font-medium">
          {userDetails &&
            `Welcome${
              userDetails?.full_name ? ` ${userDetails?.full_name}` : ''
            }!`}
        </h4>
        <p className="text-ab-black mt-2.5 text-base font-medium">
          Manage your account
        </p>
        <div className="float-left mt-4 grid w-full grid-cols-2 gap-7">
          <div className="border-ab-gray-medium col-span-2 flex w-full flex-col border rounded p-6">
            <span className="bg-primary-light float-left flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full">
              <img src={SpacesIcon} alt="Manage Block" />
            </span>
            <div className="float-left mt-3 flex w-full flex-col">
              <h3 className="text-ab-black text-xl font-semibold">
                Your Spaces {spaceData && `(`}
                {spaceData && spaceData?.length}
                {spaceData && `)`}
              </h3>
              {/* <p className="text-ab-black text-sm font-medium mt-3">
                A simply dummy text of the printing and typesetting industry.{" "}
              </p> */}
              <div className="mt-6 grid grid-cols-2 gap-3 lg:grid-cols-2 xl:grid-cols-4">
                {spaceData &&
                  spaceData?.length > 0 &&
                  spaceData?.slice(0, numberOfDataToDisplay).map((item) => (
                    <Link
                      to="/spaces/my-entities"
                      key={item.space_id}
                      onClick={() => handleSpaceChange(item.space_id)}
                      className="border-ab-gray-medium flex w-full flex-col items-center border px-6 py-8 cursor-pointer"
                    >
                      <span className="bg-primary-light text-primary float-left flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full text-2xl font-semibold capitalize">
                        {item.space_name[0]}
                      </span>
                      <p className="text-base font-medium text-ab-black mt-3.5 w-full truncate text-center">
                        {item.space_name}
                      </p>
                    </Link>
                  ))}
                {loader &&
                  !spaceData &&
                  loaderData.map((_, index) => (
                    // eslint-disable-next-line react/no-array-index-key
                    <HomeSpaceListingSkeletonLoader key={index} />
                  ))}

                <Link
                  to="/create-space"
                  className="border-ab-gray-medium trasnsition-all hover:bg-[#F5F5FF] group flex w-full cursor-pointer flex-col items-center border px-6 py-8 duration-300"
                >
                  <span className="bg-primary-light trasnsition-all textstroke-width-primary float-left flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full text-2xl font-semibold duration-200 group-hover:bg-primary">
                    <svg
                      className="stroke-primary group-hover:stroke-white"
                      width="21"
                      height="21"
                      viewBox="0 0 21 21"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M10.4997 1.83337V19.1667M1.83301 10.5H19.1663"
                        strokeWidth="1.33333"
                      />
                    </svg>
                  </span>
                  <p className="text-base font-medium text-primary mt-3.5 w-full truncate text-center group-hover:text-primary">
                    Create New
                  </p>
                </Link>
              </div>
              {spaceData && spaceData.length > numberOfDataToDisplay && (
                <span
                  onClick={() => handleSpaceData()}
                  className="text-primary text-ab-sm mt-5 inline-block cursor-pointer font-semibold underline"
                >
                  Show more
                </span>
              )}
              {/* <a className="text-primary text-ab-sm mt-5 inline-block cursor-pointer font-semibold underline">
                Manage Spaces
              </a> */}
            </div>
          </div>
        </div>
      </div>
      <div className="w-full float-left lg:max-w-[280px] mb-10 lg:mt-10 flex flex-col space-y-4 px-4 lg:px-0">
        <div className="w-full float-left flex flex-col items-start space-y-4 bg-white border outline outline-1 outline-transparent border-ab-gray-medium rounded-lg p-6 group">
          <div className="float-left">
            <img
              className="w-full max-w-[64px] float-left"
              src={BlogIcon}
              alt=""
            />
          </div>
          <h4 className="text-xl text-ab-black font-bold">Visit Our Blog</h4>
          <p className="text-sm text-ab-black font-medium">
            We love to share ideas! Visit our blog if you&apos;re looking for
            great articles.
          </p>
          <a
            href="https://discord.gg/b7YSVvHp2x"
            target="_blank"
            className="flex items-center justify-center py-2.5 px-6 bg-secondary rounded-full text-white text-sm font-semibold hover:opacity-90"
            rel="noreferrer"
          >
            Go to Blogs
          </a>
        </div>
        <div className="w-full float-left flex flex-col items-start space-y-4 bg-white border outline outline-1 outline-transparent border-ab-gray-medium rounded-lg p-6 group">
          <div className="float-left">
            <img
              className="w-full max-w-[64px] float-left"
              src={DocIcon}
              alt=""
            />
          </div>
          <h4 className="text-xl text-ab-black font-bold">
            Read our Documentation
          </h4>
          <p className="text-sm text-ab-black font-medium">
            Learn more about Spaces
          </p>
          <a
            href="https://docs.appblocks.com/"
            target="_blank"
            className="flex items-center justify-center py-2.5 px-6 bg-secondary rounded-full text-white text-sm font-semibold hover:opacity-90"
            rel="noreferrer"
          >
            Check Docs
          </a>
        </div>
      </div>
    </div>
  )
}

export default Home
