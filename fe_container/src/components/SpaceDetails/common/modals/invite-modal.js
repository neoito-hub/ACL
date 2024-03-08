/* eslint-disable react/button-has-type */
/* eslint-disable no-unused-expressions */
/* eslint-disable consistent-return */
/* eslint-disable import/no-named-as-default */
/* eslint-disable react/prop-types */
import React, {
  useRef,
  useState,
  useEffect,
  useContext,
  useCallback,
} from 'react'
import useOnclickOutside from 'react-cool-onclickoutside'
import { debounce } from 'lodash'
import { Tooltip } from 'react-tooltip'
import 'react-tooltip/dist/react-tooltip.css'
import apiHelper from '../helpers/apiGetters'
import MyContext from '../my-context'
import CreateNewIcon from '../../../../assets/img/icons/user-icon.gif'
import TickBlack from '../../../../assets/img/icons/tick-black.svg'

const INVITE_LINK_API_MAPPING = {
  Space: process.env.CREATE_INVITE_LINK_URL,
  Team: process.env.TEAMS_CREATE_INVITE_LINK_URL,
  Role: process.env.ROLES_CREATE_INVITE_LINK_URL,
}

const SEND_USER_INVITE_EMAIL_API_MAPPING = {
  Space: process.env.SEND_USER_INVITE_EMAIL_URL,
  Team: process.env.TEAMS_SEND_USER_INVITE_EMAIL_URL,
  Role: process.env.ROLES_SEND_USER_INVITE_EMAIL_URL,
}

const SEARCH_USER_API_MAPPING = {
  Team: process.env.TEAMS_SEARCH_USER_URL,
  Role: process.env.ROLES_SEARCH_USER_URL,
}

const InviteModal = (props) => {
  const { hasInviteModal, handleInviteModal, type, updateList, current } = props
  const { spaceId } = useContext(MyContext)

  const emailInputRef = useRef()
  const [emailDropdown, setEmailDropdown] = useState(false)
  const [inviteLink, setInviteLink] = useState(null)
  const [users, setUsers] = useState(null)
  const [emails, setEmails] = useState([])
  const [loader, setLoader] = useState(false)

  const handleEmailDropdown = () => {
    setEmailDropdown(true)
  }
  const emailDropContainer = useOnclickOutside(() => {
    setEmailDropdown(false)
  })

  useEffect(async () => {
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: INVITE_LINK_API_MAPPING[type],
      value: {
        data: [
          {
            space_id: spaceId,
            // ...(type === 'Space' && { team_ids: [] }),
            ...(type === 'Role' && { role_ids: [current] }),
            ...(type === 'Team' && { team_ids: [current] }),
          },
        ],
      },
      spaceId,
    })
    setInviteLink(res.invite_link)
  }, [])

  const handleEmailLabel = (e) => {
    if (
      /^[\w-\\.]+@([\w-]+\.)+[\w-]{2,4}$/.test(e.target.value) &&
      !emails.includes(e.target.value)
    ) {
      setEmails([...emails, e.target.value])
      e.target.value = ''
      e.target.focus()
    } else {
      return false
    }
  }

  const handleRemoveEmail = (email) => {
    setEmails(emails.filter((item) => item !== email))
  }

  const handleInviteSubmit = async () => {
    setLoader(true)
    if (emails?.length) {
      await apiHelper({
        baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
        subUrl: SEND_USER_INVITE_EMAIL_API_MAPPING[type],
        value: {
          data: [
            {
              space_id: spaceId,
              // ...(type === 'Space' && { team_ids: [] }),
              ...(type === 'Role' && { role_ids: [current] }),
              ...(type === 'Team' && { team_ids: [current] }),
            },
          ],
          email: emails,
        },
        showSuccessMessage: true,
        spaceId,
      })
      setEmails([])
      handleInviteModal()
      updateList()
      setLoader(false)
    }
  }

  const onCopyToClipboard = async () => {
    navigator.clipboard.writeText(inviteLink)
  }

  const onEmailChange = async (searchText) => {
    setUsers(null)
    const res =
      type !== 'Member' &&
      (await apiHelper({
        baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
        subUrl: SEARCH_USER_API_MAPPING[type],
        value: {
          search_string: searchText,
          space_id: spaceId,
          ...(type === 'Team' && { team_id: current }),
          ...(type === 'Role' && { role_id: current }),
        },

        spaceId,
      }))
    setUsers(res?.filter((item) => emails.indexOf(item.email) === -1))
    setLoader(false)
  }

  const handler = useCallback(
    debounce((text) => {
      onEmailChange(text)
    }, 1000),
    [],
  )

  const onEmailSelect = (email) => {
    setEmails((prevArray) => [...prevArray, email])
    setUsers((prevUsers) => prevUsers.filter((user) => user.email !== email))
    document.getElementById('searchInput').value = ''
    setEmailDropdown(false)
  }

  return (
    <div
      className={`fixed left-0 top-0 z-[10001] h-screen w-full ${
        hasInviteModal ? 'fadeIn' : 'hidden'
      }`}
    >
      <div
        onClick={() => {
          handleInviteModal()
        }}
        className="fixed left-0 top-0 z-[10001] h-full w-full bg-black/40"
      />
      <div
        className={`absolute top-1/2 left-1/2 z-[10002] w-full max-w-[620px] -translate-x-1/2 -translate-y-1/2 transform px-4 ${
          hasInviteModal ? '' : 'hidden'
        }`}
      >
        <div className="relative float-left flex w-full rounded-md bg-white p-6 md:space-x-10 md:p-[60px] md-lt:flex-col">
          <img
            src={CreateNewIcon}
            alt="Create New"
            className="h-20 w-20 flex-shrink-0 rounded-full md-lt:mx-2 md-lt:mb-4"
          />
          <div className="flex w-full flex-col md:max-w-[calc(100%-120px)]">
            <h5 className="mb-3 text-lg font-semibold text-black">
              Invite members to {type}
            </h5>
            <div className="float-left mb-4 w-full">
              <div className="float-left mb-2 w-full">
                <label className="float-left mb-2 text-xs font-medium text-[#24292E]">
                  Invite through e-mail
                </label>
                <div
                  className="relative float-left w-full"
                  ref={emailDropContainer}
                >
                  <div
                    onClick={() => emailInputRef.current.focus()}
                    className="border-ab-gray-light bg-ab-gray-light focus-within:border-primary float-left w-full select-none rounded-md border py-0.5 px-1 text-xs focus-within:outline-none "
                  >
                    <div className="w-full float-left max-h-[100px] overflow-auto custom-scroll-primary px-1.5">
                      {emails &&
                        emails.map((email) => (
                          <div
                            key={email}
                            className="bg-ab-disabled-yellow float-left my-1 mr-2 inline-flex max-w-full items-center space-x-2 rounded-full py-0.5 px-3"
                          >
                            <p className="text-ab-black truncate text-xs font-medium leading-normal">
                              {email}
                            </p>
                            <svg
                              className="flex-shrink-0 cursor-pointer"
                              onClick={() => handleRemoveEmail(email)}
                              width="8"
                              height="8"
                              viewBox="0 0 8 8"
                              fill="none"
                              xmlns="http://www.w3.org/2000/svg"
                            >
                              <path
                                d="M0.799988 0.799805L7.19999 7.19981M0.799988 7.19981L7.19999 0.799805"
                                stroke="#484848"
                                strokeWidth="1.53333"
                              />
                            </svg>
                          </div>
                        ))}
                      <input
                        type="text"
                        id="searchInput"
                        ref={emailInputRef}
                        autoComplete="off"
                        onChange={(e) => handler(e.target.value)}
                        className="text-ab-black inline-block w-auto border-0 bg-transparent py-2 align-baseline text-xs outline-none focus:border-none focus:shadow-none focus:outline-none"
                        placeholder="Emails"
                        onBlur={(e) => handleEmailLabel(e)}
                        onFocus={(e) => {
                          type !== 'Space' && handleEmailDropdown()
                          type !== 'Space' && onEmailChange(e.target.value)
                        }}
                        onKeyPress={(e) =>
                          e.code === 'Enter' && handleEmailLabel(e)
                        }
                      />
                    </div>
                  </div>
                  <div
                    className={`shadow-box dropDownFade border-ab-gray-medium absolute top-full left-0 z-10 mt-1 w-full border bg-white p-1 py-3 ${
                      emailDropdown && users?.length ? '' : 'hidden'
                    }`}
                  >
                    <div className="float-left w-full p-2">
                      <p className="text-ab-black float-left w-full text-xs font-medium">
                        Select a person
                      </p>
                    </div>
                    <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-1">
                      {users?.map((user) => (
                        <li
                          key={user?.email}
                          className="hover:bg-ab-gray-light float-left mb-1.5 w-full py-1.5 last-of-type:mb-0"
                        >
                          <label
                            onClick={(e) => {
                              onEmailSelect(user?.email)
                              e.preventDefault()
                            }}
                            className="float-left flex w-full cursor-pointer items-center space-x-2 leading-normal"
                          >
                            <input
                              name="email"
                              className="peer hidden"
                              type="checkbox"
                            />
                            <img
                              className="invisible flex-shrink-0 peer-checked:visible"
                              src={TickBlack}
                              alt=""
                            />
                            <span className="bg-primary flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-full text-xs font-bold text-white">
                              {user?.email[0].toUpperCase()}
                            </span>
                            <p className="text-ab-black truncate text-xs font-medium tracking-tight">
                              {user.email}
                            </p>
                          </label>
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </div>
              <div className="float-left mt-2 flex w-full items-center">
                <button
                  type="button"
                  onClick={handleInviteSubmit}
                  className="btn-primary disabled:bg-ab-disabled rounded-md px-4 py-2 text-xs font-bold leading-tight text-white transition-all"
                  disabled={!emails?.length || loader}
                >
                  Invite
                </button>
              </div>
            </div>
            <div className="float-left mb-2 w-full">
              <label className="float-left mb-2 text-xs font-medium text-[#24292E]">
                Invite Via link
              </label>
              <div className="border-ab-gray-light bg-ab-gray-light focus:border-primary relative float-left flex w-full items-center space-x-2 rounded-md border py-2.5 px-4">
                <input
                  readOnly
                  defaultValue={inviteLink && inviteLink}
                  type="text"
                  className="float-left flex-grow bg-transparent text-xs focus:outline-none"
                />
                <button
                  id="copied"
                  data-tooltip-content="Copied"
                  data-tooltip-delay-hide={1200}
                  className="text-primary flex-shrink-0 cursor-pointer text-xs underline-offset-2 hover:underline"
                  onClick={() => onCopyToClipboard()}
                >
                  Copy Link
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <Tooltip
        anchorId="copied"
        effect="solid"
        padding="4px 8px"
        className="ab-tooltip-copied"
        arrowColor="transparent"
        events={['click']}
        noArrow
      />
    </div>
  )
}

export default InviteModal
