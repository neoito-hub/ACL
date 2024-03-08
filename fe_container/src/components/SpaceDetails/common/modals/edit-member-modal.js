/* eslint-disable react/prop-types */
import React, { useRef, useState, useEffect, useContext } from 'react'
import useOnclickOutside from 'react-cool-onclickoutside'
import dayjs from 'dayjs'
import { Formik } from 'formik'
import DownArrow from '../../../../assets/img/icons/down-arrow.svg'
import MyContext from '../my-context'
import apiHelper from '../helpers/apiGetters'
import { MemberValidationSchema } from '../validation/validation'
// import Avatar from '../../../../assets/img/icons/github.svg';

const initialValues = {
  user_id: '',
  name: '',
  email: '',
  roles: [],
  teams: [],
}

const EditMemberModal = (props) => {
  const {
    hasEditMemberModal,
    handleEditMemberModal,
    selectedUser,
    updateList,
  } = props
  const formikRef = useRef()
  const { spaceId } = useContext(MyContext)

  const [teamList, setTeamList] = useState(null)
  const [roleList, setRoleList] = useState(null)
  // const [userDetails, setUserDetails] = useState(null);
  const [roleDropdown, setRoleDropdown] = useState(false)
  const [teamDropdown, setTeamDropdown] = useState(false)
  const [loader, setLoader] = useState(false)
  const roleDropContainer = useOnclickOutside(() => {
    setRoleDropdown(false)
  })
  const teamDropContainer = useOnclickOutside(() => {
    setTeamDropdown(null)
  })

  useEffect(async () => {
    setLoader(true)
    formikRef?.current?.setFieldValue('user_id', selectedUser?.user_id)
    formikRef?.current?.setFieldValue('name', selectedUser?.full_name)
    formikRef?.current?.setFieldValue('email', selectedUser?.email)
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.GET_USER_BY_ID_URL,
      value: {
        user_id: selectedUser.user_id,
        space_id: spaceId,
      },
      spaceId,
    })
    setRoleList(res.roles)
    setTeamList(res.teams)
    // setUserDetails(res);
    const filteredRoles = res.roles ? res.roles.map((role) => role.id) : []
    formikRef?.current?.setFieldValue('roles', filteredRoles)
    const filteredTeams = res?.teams
      ? res.teams?.map((team) => team.team_id)
      : []
    formikRef?.current?.setFieldValue('teams', filteredTeams)
    setLoader(false)
  }, [selectedUser])

  const handleTeam = (e) => {
    let updatedSelectedTeam = formikRef?.current?.values?.teams
    if (e.target.checked) {
      updatedSelectedTeam = [...updatedSelectedTeam, e.target.value]
    } else {
      updatedSelectedTeam = updatedSelectedTeam.filter(
        (item) => item !== e.target.value
      )
    }
    formikRef?.current?.setFieldValue('teams', updatedSelectedTeam)
  }
  const handleTeamRemove = (e, team) => {
    e.stopPropagation()
    const updatedSelectedTeam = formikRef?.current?.values?.teams?.filter(
      (item) => item !== team?.team_id
    )
    formikRef?.current?.setFieldValue('teams', updatedSelectedTeam)
  }
  const handleRole = (e) => {
    let updatedSelectedRole = formikRef?.current?.values?.roles
    if (e.target.checked) {
      updatedSelectedRole = [...updatedSelectedRole, e.target.value]
    } else {
      updatedSelectedRole = updatedSelectedRole.filter(
        (item) => item !== e.target.value
      )
    }
    formikRef?.current?.setFieldValue('roles', updatedSelectedRole)
  }

  const handleRoleRemove = (e, role) => {
    e.stopPropagation()
    const updatedSelectedRole = formikRef?.current?.values?.roles?.filter(
      (item) => item !== role?.id
    )
    formikRef?.current?.setFieldValue('roles', updatedSelectedRole)
  }

  const onSubmit = async (values) => {
    setLoader(true)
    const deletedRoles = roleList
      ? roleList
          .filter((role) => !values.roles.includes(role.id))
          .map((item) => item.id)
      : []
    const deletedTeams = teamList
      ? teamList
          .filter((team) => !values.teams.includes(team.team_id))
          .map((item) => item.team_id)
      : []
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.UPDATE_USER_URL,
      value: {
        user_id: values.user_id,
        deleted_roles: deletedRoles,
        deleted_teams: deletedTeams,
        user_details: {
          full_name: values.name,
        },
      },
      showSuccessMessage: true,
      spaceId,
    })
    setLoader(false)
    handleEditMemberModal()
    updateList()
  }

  return (
    <Formik
      innerRef={formikRef}
      initialValues={initialValues}
      onSubmit={onSubmit}
      validationSchema={MemberValidationSchema()}
      validateOnMount
      validateOnChange
      validateOnBlur
      enableReinitialize
    >
      {({ handleSubmit, values }) => (
        <div
          className={`fixed left-0 top-0 z-[999] h-screen w-full ${
            hasEditMemberModal ? 'fadeIn' : 'hidden'
          }`}
        >
          <div
            onClick={() => {
              handleEditMemberModal()
            }}
            className="fixed left-0 top-0 z-[1000] h-full w-full bg-black/40"
          />
          <div
            className={`absolute top-1/2 left-1/2 z-[1001] w-full max-w-[620px] -translate-x-1/2 -translate-y-1/2 transform px-4 ${
              hasEditMemberModal ? '' : 'hidden'
            }`}
          >
            <div className="relative float-left flex w-full rounded-md bg-white p-6 md:space-x-10 md:p-[60px] md-lt:flex-col">
              <div className="bg-primary/20 flex h-20 w-20 flex-shrink-0 items-center justify-center rounded-full md-lt:mb-3">
                <span className="text-primary text-3xl font-semibold capitalize">
                  {selectedUser?.full_name[0] || selectedUser?.email[0]}
                </span>
              </div>
              {/* <img src={Avatar} className="h-20 w-20 flex-shrink-0 rounded-full md-lt:mx-2 md-lt:mb-3 border-ab-gray-medium border object-cover" alt=""/> */}
              <div className="flex w-full flex-col md:max-w-[calc(100%-120px)]">
                <h5 className="mb-4 text-lg font-semibold text-black">
                  Edit Member
                </h5>
                <div className="text-ab-black float-left mb-3 w-full">
                  <label className="float-left text-xs font-medium">Name</label>
                  <div className="float-left flex w-full items-center justify-between space-x-3">
                    <p className="text-ab-sm truncate py-1.5">
                      {selectedUser?.full_name}
                    </p>
                  </div>
                </div>
                <div className="text-ab-black float-left mb-3 w-full">
                  <label className="float-left text-xs font-medium">
                    E-mail
                  </label>
                  <div className="float-left flex w-full items-center justify-between space-x-3">
                    <p className="text-ab-sm truncate py-1.5">
                      {values?.email}
                    </p>
                    <span className="bg-ab-green/10 text-ab-green flex-shrink-0 rounded-full py-1 px-2.5 text-xs font-medium">
                      Active
                    </span>
                  </div>
                </div>
                <div className="text-ab-black float-left mb-3 w-full">
                  <label className="float-left text-xs font-medium">
                    Member Since
                  </label>
                  <div className="float-left flex w-full items-center justify-between space-x-3">
                    <p className="text-ab-sm truncate py-1.5 font-medium text-gray-400">
                      {selectedUser?.created_date &&
                        dayjs(selectedUser.created_date).format('DD/MM/YYYY')}
                    </p>
                  </div>
                </div>
                <div className="text-ab-black float-left mb-3 w-full">
                  <label className="float-left text-xs font-medium">Role</label>
                  <div
                    className="relative float-left my-1.5 w-full"
                    ref={roleDropContainer}
                  >
                    <div
                      onClick={() => setRoleDropdown(!roleDropdown)}
                      className={`text-ab-sm bg-ab-gray-light focus:border-primary float-left flex w-full cursor-pointer select-none items-center justify-between rounded-md border py-0.5 px-2.5 focus:outline-none ${
                        roleDropdown ? 'border-primary' : 'border-ab-gray-light'
                      }`}
                    >
                      <div className="text-ab-black text-ab-sm flex-grow overflow-hidden">
                        {values?.roles?.length ? (
                          roleList?.map(
                            (role) =>
                              values?.roles?.includes(role?.id) && (
                                <div
                                  key={role?.id}
                                  className="bg-ab-disabled-yellow float-left my-1 mr-2 inline-flex max-w-full items-center space-x-2 rounded-full py-0.5 px-3"
                                >
                                  <p className="truncate text-xs font-medium leading-[20px]">
                                    {role?.name}
                                  </p>
                                  <svg
                                    className="flex-shrink-0 cursor-pointer"
                                    onClick={(e) => handleRoleRemove(e, role)}
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
                              )
                          )
                        ) : (
                          <p className="flex h-[32px] items-center px-1 text-xs font-medium text-gray-400">
                            Select
                          </p>
                        )}
                      </div>
                      <img
                        className={`ml-3 flex-shrink-0 transform transition-transform duration-300 ${
                          roleDropdown && 'rotate-180'
                        }`}
                        src={DownArrow}
                        alt=""
                      />
                    </div>
                    <div
                      className={`shadow-box dropDownFade border-ab-gray-medium position-inverse-2 absolute top-full left-0 z-10 mt-1 w-full border bg-white p-1 py-3 ${
                        roleDropdown ? '' : 'hidden'
                      }`}
                    >
                      {roleList && (
                        <div className="float-left w-full p-2">
                          <p className="text-ab-black border-ab-gray-dark float-left w-full border-b pb-3 text-xs font-medium">
                            Select All that Apply
                          </p>
                        </div>
                      )}
                      <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-2">
                        {roleList?.map((role) => (
                          <li
                            key={role.id}
                            className="float-left mb-4 w-full last-of-type:mb-0"
                          >
                            <label className="float-left flex max-w-full cursor-pointer items-center leading-normal">
                              <input
                                checked={values?.roles?.includes(role?.id)}
                                onChange={(e) => handleRole(e)}
                                name="team"
                                value={role.id}
                                className="peer hidden"
                                type="checkbox"
                              />
                              <span className="chkbox-icon border-ab-disabled float-left mr-2 h-5 w-5 flex-shrink-0 cursor-pointer rounded border bg-white" />
                              <p className="text-ab-black truncate text-xs font-medium tracking-tight">
                                {role.name}
                              </p>
                            </label>
                          </li>
                        ))}
                        {!roleList && (
                          <li className="float-left mb-4 w-full last-of-type:mb-0">
                            <label className="float-left flex max-w-full cursor-pointer items-center leading-normal">
                              <p className="text-ab-black truncate text-xs font-medium tracking-tight">
                                No Roles to display
                              </p>
                            </label>
                          </li>
                        )}
                      </ul>
                    </div>
                  </div>
                </div>
                <div className="text-ab-black float-left mb-3 w-full">
                  <label className="float-left text-xs font-medium">Team</label>
                  <div
                    className="relative float-left my-1.5 w-full"
                    ref={teamDropContainer}
                  >
                    <div
                      onClick={() => setTeamDropdown(!teamDropdown)}
                      className={`text-ab-sm bg-ab-gray-light focus:border-primary float-left flex w-full cursor-pointer select-none items-center justify-between rounded-md border py-0.5 px-2.5 focus:outline-none ${
                        teamDropdown ? 'border-primary' : 'border-ab-gray-light'
                      }`}
                    >
                      <div className="text-ab-black text-ab-sm flex-grow overflow-hidden">
                        {values?.teams?.length ? (
                          teamList?.map(
                            (team) =>
                              values?.teams?.includes(team?.team_id) && (
                                <div
                                  key={team?.team_id}
                                  className="bg-ab-disabled-yellow float-left my-1 mr-2 inline-flex max-w-full items-center space-x-2 rounded-full py-0.5 px-3"
                                >
                                  <p className="truncate text-xs font-medium leading-[20px]">
                                    {team?.name}
                                  </p>
                                  <svg
                                    onClick={(e) => handleTeamRemove(e, team)}
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
                              )
                          )
                        ) : (
                          <p className="flex h-[32px] items-center px-1 text-xs font-medium text-gray-400">
                            Select
                          </p>
                        )}
                      </div>
                      <img
                        className={`ml-3 flex-shrink-0 transform transition-transform duration-300 ${
                          teamDropdown && 'rotate-180'
                        }`}
                        src={DownArrow}
                        alt=""
                      />
                    </div>
                    <div
                      className={`shadow-box dropDownFade border-ab-gray-medium position-inverse-2 absolute top-full left-0 z-10 mt-1 w-full border bg-white p-1 py-3 ${
                        teamDropdown ? '' : 'hidden'
                      }`}
                    >
                      {teamList && (
                        <div className="float-left w-full p-2">
                          <p className="text-ab-black border-ab-gray-dark float-left w-full border-b pb-3 text-xs font-medium">
                            Select All that Apply
                          </p>
                        </div>
                      )}
                      <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-2">
                        {teamList?.map((team) => (
                          <li
                            key={team?.team_id}
                            className="float-left mb-4 w-full last-of-type:mb-0"
                          >
                            <label className="float-left flex max-w-full cursor-pointer items-center leading-normal">
                              <input
                                checked={values?.teams?.includes(team?.team_id)}
                                onChange={(e) => handleTeam(e)}
                                name="team"
                                value={team.team_id}
                                className="peer hidden"
                                type="checkbox"
                              />
                              <span className="chkbox-icon border-ab-disabled float-left mr-2 h-5 w-5 flex-shrink-0 cursor-pointer rounded border bg-white" />
                              <p className="text-ab-black truncate text-xs font-medium tracking-tight">
                                {team?.name}
                              </p>
                            </label>
                          </li>
                        ))}
                        {!teamList && (
                          <li className="float-left mb-4 w-full last-of-type:mb-0">
                            <label className="float-left flex max-w-full cursor-pointer items-center leading-normal">
                              <p className="text-ab-black truncate text-xs font-medium tracking-tight">
                                No Teams to display
                              </p>
                            </label>
                          </li>
                        )}
                      </ul>
                    </div>
                  </div>
                </div>
                <div className="float-left mt-2 flex w-full items-center">
                  <button
                    type="button"
                    disabled={loader}
                    onClick={handleSubmit}
                    className="btn-primary text-ab-sm disabled:bg-ab-disabled mr-4 rounded px-5 py-2.5 font-bold leading-normal text-white transition-all focus:outline-none"
                  >
                    Save Changes
                  </button>
                  <button
                    onClick={() => {
                      handleEditMemberModal()
                    }}
                    type="button"
                    className="text-ab-disabled hover:text-ab-black text-ab-sm rounded px-3 py-1 font-bold leading-tight text-white focus:outline-none"
                  >
                    Cancel
                  </button>
                </div>
                {/* <div className='float-left mt-5 flex w-full items-center'>
                  <a className='text-ab-red float-left cursor-pointer text-xs font-medium underline focus:outline-none'>
                    Delete Member
                  </a>
                </div> */}
              </div>
            </div>
          </div>
        </div>
      )}
    </Formik>
  )
}

export default EditMemberModal
