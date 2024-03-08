/* eslint-disable no-unused-expressions */
/* eslint-disable consistent-return */
/* eslint-disable camelcase */
import React, { useState, useEffect } from 'react'
import { shield } from '@appblocks/js-sdk'
import { useHistory, useLocation } from 'react-router-dom'
import Axios from 'axios'
import queryString from 'query-string'

const Invitation = () => {
  const history = useHistory()
  const location = useLocation()
  const { invite_id } = queryString.parse(location.search)

  const [inviteData, setInviteData] = useState(null)

  const getData = async (baseUrl, subUrl, value = null) => {
    const token = shield.tokenStore.getToken()
    try {
      const { data } = await Axios.post(
        `${baseUrl}${subUrl}`,
        value && value,
        token && {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )
      return data?.data
    } catch (error) {
      console.log(error)
    }
  }

  useEffect(async () => {
    const res = await getData(
      process.env.BLOCK_ENV_URL_API_BASE_URL,
      process.env.GET_INVITE_BY_ID_URL,
      { invite_id }
    )
    res && setInviteData(res)
  }, [])

  const onAcceptClick = async () => {
    const res = await getData(
      process.env.BLOCK_ENV_URL_API_BASE_URL,
      process.env.ACCEPT_INVITE_URL,
      { invite_id }
    )
    res && history.push('/')
  }

  const onDeclineClick = async () => {
    inviteData?.invite_type !== 2 &&
      (await getData(
        process.env.BLOCK_ENV_URL_API_BASE_URL,
        process.env.CANCEL_INVITE_URL,
        { invite_id }
      ))
    history.push('/')
  }

  return (
    <div className="flex min-h-[calc(100vh-64px)] w-full items-center justify-center">
      {inviteData && (
        <div className="border-ab-gray-medium relative float-left flex w-full max-w-[480px] flex-col rounded-md border bg-white p-6 md:px-16 md:pt-12 md:pb-10">
          {!inviteData?.expired &&
          ((inviteData?.invite_type === 2 && inviteData?.status === 0) ||
            (inviteData?.invite_type === 1 && inviteData?.status === 1)) ? (
            <div className="bg-secondary mb-6 flex h-16 w-16 items-center justify-center rounded-full text-3xl font-medium leading-tight text-white capitalize">
              {inviteData?.invite_details?.length &&
                inviteData?.invite_details[0].space_name[0]}
            </div>
          ) : (
            <div className="bg-secondary mb-6 flex h-16 w-16 items-center justify-center rounded-full text-3xl font-medium leading-tight text-white capitalize">
              !
            </div>
          )}
          <div className="flex flex-grow flex-col overflow-hidden">
            <h5 className="text-ab-black mb-5 text-xl font-semibold md:text-[23px]">
              {inviteData?.msg}
            </h5>
            {/* {inviteData?.expired ||
            inviteData?.status === 0 ||
            (inviteData?.invite_type === 1 && inviteData?.status !== 1) ( */}
            {!inviteData?.expired &&
            ((inviteData?.invite_type === 2 && inviteData?.status === 0) ||
              (inviteData?.invite_type === 1 && inviteData?.status === 1)) ? (
              <>
                <button
                  type="button"
                  onClick={onAcceptClick}
                  className="btn-primary text-ab-base !focus:outline-none mb-4 !py-3 font-semibold hover:no-underline"
                >
                  Accept
                </button>
                <button
                  type="button"
                  onClick={onDeclineClick}
                  className="text-ab-red text-ab-base !focus:outline-none mb-4 !py-3 font-semibold hover:no-underline hover:opacity-70"
                >
                  Decline
                </button>
              </>
            ) : (
              <button
                type="button"
                onClick={() => history.push('/')}
                className="btn-primary text-ab-base !focus:outline-none mb-4 !py-3 mt-4 font-semibold hover:no-underline"
              >
                Return To Home
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

export default Invitation
