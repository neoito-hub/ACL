/* eslint-disable no-sequences */
/* eslint-disable no-unused-expressions */
/* eslint-disable import/no-extraneous-dependencies */
/* eslint-disable consistent-return */
import React, { useState, useEffect, useRef, useContext } from 'react'
import { shield } from '@appblocks/js-sdk'
import { Formik } from 'formik'
import axios from 'axios'
import { toast } from 'react-toastify'
import Skeleton from 'react-loading-skeleton'
import Toast from '../toast/Toast'
import 'react-loading-skeleton/dist/skeleton.css'
import PasswordValidationSchema from '../validation/validation'
// import EditIcon from '../../assets/img/icons/edit-icon.svg';
import GoogleLogo from '../../../assets/img/icons/google.svg'
import LinkedInLogo from '../../../assets/img/icons/linkedin.svg'
import TwitterLogo from '../../../assets/img/icons/twitter.svg'
// import GithubLogo from '../../assets/img/icons/github.svg';
import OpenEye from '../../../assets/img/icons/open-eye.svg'
import CloseEye from '../../../assets/img/icons/close-eye.svg'
import { ACLContext } from '../../../context/ACLContext'

const initialValues = { current_password: '', new_password: '' }

const ProfileSettings = () => {
  const ipName = useRef()
  const ipUserName = useRef()
  const formikRef = useRef(null)

  const { userDetails, setUserDetails } = useContext(ACLContext)

  const [showPwd, setShowPwd] = useState(false)
  const [showNewPwd, setShowNewPwd] = useState(false)
  const [flag, setFlag] = useState(false)
  const [changeName, setChangeName] = useState(false)
  const [changeUserName, setChangeUserName] = useState(false)
  const [name, setName] = useState('')
  const [userName, setUserName] = useState('')
  const [loader, setLoader] = useState(false)
  // const [hasStripeConnected, setHasStripeConnected] = useState(null);
  // const [stripeConnectUrl, setStripeConnectUrl] = useState(null);
  // const [selectedImage, setSelectedImage] = useState(null);

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
      console.log('msg', err.response.data.data)
      toast.error(err.response.data.data)
      // if (err.response.status === 401) shield.logout();
    }
  }

  const onChange = (data) => {
    formikRef?.current?.setFieldValue(data.target.name, data.target.value)
  }

  useEffect(async () => {
    setLoader(true)
    const res = await apiHelper(
      process.env.SHIELD_AUTH_URL,
      process.env.GET_USER_DETAILS_URL,
      null,
      'get'
    )
    setUserDetails(res)
    setName(res.full_name)
    setUserName(res.user_name)
    setLoader(false)
  }, [flag])

  // useEffect(async () => {
  //   const res = await apiHelper(
  //     process.env.PAYMENTS_HOST,
  //     process.env.GET_ACCOUNT_DETAILS,
  //     {}
  //   );
  //   setHasStripeConnected(res.details_submitted);
  //   if (!res.details_submitted) {
  //     const response = await apiHelper(
  //       process.env.PAYMENTS_HOST,
  //       process.env.CREATE_ACCOUNT_LINK,
  //       {}
  //     );
  //     setStripeConnectUrl(response.onboarding_url);
  //   }
  // }, []);

  const onPasswordReset = () => {
    formikRef?.current?.setFieldValue('current_password', '')
    formikRef?.current?.setFieldValue('new_password', '')
    formikRef?.current?.setFieldTouched('current_password', false)
    formikRef?.current?.setFieldTouched('new_password', false)
    setShowNewPwd(false)
    setShowPwd(false)
  }

  const onPasswordChange = async (values) => {
    const requestData = new FormData()
    requestData.append('current_password', values.current_password)
    requestData.append('new_password', values.new_password)
    const res = await apiHelper(
      process.env.SHIELD_AUTH_URL,
      process.env.CHANGE_USER_PASSWORD_URL,
      requestData
    )
    onPasswordReset()
    res && toast.success(res)
  }

  const onDetailsChange = async (type) => {
    const requestData = new FormData()
    if (type === 'name') requestData.append('fullname', name)
    else requestData.append('username', userName)
    const res = await apiHelper(
      process.env.SHIELD_AUTH_URL,
      process.env.UPDATE_USER_PROFILE_URL,
      requestData,
      'patch'
    )
    if (res) {
      toast.success(res)
      setChangeName(false)
      setChangeUserName(false)
      setFlag(!flag)
    }
  }

  const handleResetPassword = async () => {
    const url = await shield.getAuthUrl('password-recovery')
    const res = await apiHelper(url, '', { email: userDetails?.email })
    onPasswordReset()
    res && toast.success('Email Sent Successfully')
  }

  return (
    <>
      <Toast />
      <div className="float-left w-full py-6 md-lt:px-4">
        <h4 className="text-ab-black mt-1 mb-5 text-xl font-semibold">
          Profile Settings
        </h4>
        <div className="float-left mb-6 flex w-full flex-col md:max-w-2xl">
          {/* <div
          className={`relative mb-6 flex h-16 w-16 items-center justify-center rounded-full text-white ${
            selectedImage === null && 'bg-secondary'
          }`}
        >
          {selectedImage === null && (
            <span className='text-3xl font-bold'>S</span>
          )}
          {selectedImage !== null && (
            <img
              className='border-ab-gray-medium h-16 w-16 rounded-full border object-cover'
              src={URL.createObjectURL(selectedImage)}
              alt='avatar'
            />
          )}
          <label className='absolute bottom-0 right-0 flex h-5 w-5 cursor-pointer items-center justify-center rounded-full bg-white shadow-lg'>
            <input
              type='file'
              accept='image/*'
              onChange={(event) => {
                setSelectedImage(event.target.files[0]);
              }}
              className='hidden'
            ></input>
            <img src={EditIcon} alt='Edit'></img>
          </label>
        </div> */}
          <div className="float-left mb-6 grid w-full grid-cols-1 gap-8 md:grid-cols-2">
            <div className="float-left w-full flex flex-col">
              <label className="text-ab-sm float-left mb-2 font-medium text-black">
                Name
              </label>
              {loader ? (
                <div className="flex flex-col gap-1.5 rounded py-1.5">
                  <Skeleton height={8} />
                  <Skeleton height={8} />
                  <Skeleton height={8} />
                </div>
              ) : (
                <input
                  ref={ipName}
                  value={name}
                  type="text"
                  readOnly={!changeName}
                  onChange={(e) => setName(e.target.value)}
                  className="text-ab-sm border-ab-gray-light bg-ab-gray-light focus:border-primary float-left w-full rounded-md border py-3.5 px-4 focus:outline-none"
                />
              )}
              {!changeName ? (
                <span
                  onClick={() => {
                    setChangeName(true), ipName.current.focus()
                  }}
                  className="text-primary float-left mt-1 cursor-pointer text-xs underline"
                >
                  Change name
                </span>
              ) : (
                <div className="float-left mt-4 flex w-full items-center">
                  <button
                    type="button"
                    onClick={() => onDetailsChange('name')}
                    className="btn-primary text-ab-sm mr-4 rounded px-5 py-2 font-bold leading-tight text-white transition-all hover:opacity-90"
                  >
                    Save
                  </button>
                  <button
                    onClick={() => {
                      setName(userDetails?.full_name)
                      setChangeName(false)
                    }}
                    type="button"
                    className="text-ab-black hover:text-ab-black text-ab-sm rounded px-3 py-1 font-bold leading-tight"
                  >
                    Cancel
                  </button>
                </div>
              )}
            </div>
          </div>
          <div className="float-left mb-6 grid w-full grid-cols-1 gap-8 md:grid-cols-2">
            <div className="float-left w-full flex flex-col">
              <label className="text-ab-sm float-left mb-2 font-medium text-black">
                Username
              </label>
              {loader ? (
                <div className="flex flex-col gap-1.5 rounded py-1.5">
                  <Skeleton height={8} />
                  <Skeleton height={8} />
                  <Skeleton height={8} />
                </div>
              ) : (
                <input
                  ref={ipUserName}
                  value={userName}
                  type="text"
                  readOnly={!changeUserName}
                  onChange={(e) => setUserName(e.target.value)}
                  className="text-ab-sm border-ab-gray-light bg-ab-gray-light focus:border-primary float-left w-full rounded-md border py-3.5 px-4 focus:outline-none"
                />
              )}
              {!changeUserName ? (
                <span
                  onClick={() => {
                    setChangeUserName(true), ipUserName.current.focus()
                  }}
                  className="text-primary float-left mt-1 cursor-pointer text-xs underline"
                >
                  Change username
                </span>
              ) : (
                <div className="float-left mt-4 flex w-full items-center">
                  <button
                    type="button"
                    onClick={() => onDetailsChange('username')}
                    className="btn-primary text-ab-sm mr-4 rounded px-5 py-2 font-bold leading-tight text-white transition-all hover:opacity-90"
                  >
                    Save
                  </button>
                  <button
                    onClick={() => {
                      setUserName(userDetails?.user_name)
                      setChangeUserName(false)
                    }}
                    type="button"
                    className="text-ab-black hover:text-ab-black text-ab-sm rounded px-3 py-1 font-bold leading-tight"
                  >
                    Cancel
                  </button>
                </div>
              )}
            </div>
          </div>
          {userDetails?.provider?.find((item) => item === 'Password') && (
            <Formik
              innerRef={formikRef}
              initialValues={initialValues}
              onSubmit={onPasswordChange}
              validationSchema={PasswordValidationSchema()}
              validateOnMount
              validateOnChange
              validateOnBlur
              enableReinitialize
            >
              {({ handleSubmit, values, errors, touched }) => (
                <div className="float-left mb-6 w-full">
                  <p className="border-ab-gray-dark my-1 border-b pb-2.5 text-sm font-semibold text-black">
                    Change Password
                  </p>
                  <div className="float-left mt-3 grid w-full grid-cols-1 gap-x-8 gap-y-4 md:grid-cols-2">
                    <input type="hidden" />
                    <div className="float-left w-full">
                      <label className="text-ab-sm float-left mb-2 font-medium text-black">
                        Current Password
                      </label>
                      <div className="relative flex flex-col float-left w-full">
                        <input
                          id="current_password"
                          name="current_password"
                          placeholder="Current Password"
                          value={values.current_password}
                          onChange={onChange}
                          type={showNewPwd ? 'text' : 'password'}
                          autoComplete="off"
                          className={`${
                            errors?.current_password &&
                            touched?.current_password
                              ? 'border-ab-red'
                              : 'border-ab-gray-light focus:border-primary'
                          } text-ab-sm bg-ab-gray-light float-left w-full rounded-md border py-3.5 px-4 pr-11 focus:outline-none`}
                        />
                        <span className="absolute top-1/2 right-4 -translate-y-1/2 transform cursor-pointer">
                          <img
                            onClick={() => setShowNewPwd(!showNewPwd)}
                            src={showNewPwd ? OpenEye : CloseEye}
                            alt=""
                          />
                        </span>
                        <p className="text-xs text-ab-red left-0 mt-0.5">
                          {errors?.current_password &&
                            touched?.current_password &&
                            errors?.current_password}
                        </p>
                      </div>
                      <p className="float-left mt-2 text-xs text-black">
                        Forgot Password?{' '}
                        <span
                          onClick={handleResetPassword}
                          className="text-primary cursor-pointer"
                        >
                          Reset
                        </span>{' '}
                      </p>
                    </div>
                    <div className="flex flex-col float-left w-full">
                      <label className="text-ab-sm float-left mb-2 font-medium text-black">
                        New Password
                      </label>
                      <div className="relative">
                        <input
                          id="new_password"
                          name="new_password"
                          placeholder="New Password"
                          value={values.new_password}
                          autoComplete="off"
                          onChange={onChange}
                          type={showPwd ? 'text' : 'password'}
                          className={`${
                            errors?.new_password && touched?.new_password
                              ? 'border-ab-red'
                              : 'border-ab-gray-light focus:border-primary'
                          } text-ab-sm bg-ab-gray-light float-left w-full rounded-md border py-3.5 px-4 pr-11 focus:outline-none`}
                        />
                        <span className="absolute top-1/2 right-4 -translate-y-1/2 transform cursor-pointer">
                          <img
                            onClick={() => setShowPwd(!showPwd)}
                            src={showPwd ? OpenEye : CloseEye}
                            alt=""
                          />
                        </span>
                      </div>
                      <p className="text-xs text-ab-red left-0 mt-0.5">
                        {errors?.new_password &&
                          touched?.new_password &&
                          errors?.new_password}
                      </p>
                    </div>
                  </div>
                  {values.current_password !== '' &&
                    values.new_password !== '' && (
                      <div className="float-left mt-4 flex w-full items-center">
                        <button
                          type="button"
                          onClick={handleSubmit}
                          className="btn-primary text-ab-sm mr-4 rounded px-5 py-2 font-bold leading-tight text-white transition-all hover:opacity-90"
                        >
                          Save
                        </button>
                        <button
                          onClick={onPasswordReset}
                          type="button"
                          className="text-ab-disabled hover:text-ab-black text-ab-sm rounded px-3 py-1 font-bold leading-tight text-white"
                        >
                          Cancel
                        </button>
                      </div>
                    )}
                </div>
              )}
            </Formik>
          )}
          <div className="float-left mb-6 w-full">
            <p className="border-ab-gray-dark my-1 border-b pb-2.5 text-sm font-semibold text-black">
              Social Connections
            </p>
            <div className="mt-5 flex w-full flex-col space-y-6">
              <div className="flex w-full items-center space-x-2">
                <div className="float-left flex flex-grow items-center space-x-2 overflow-hidden md:max-w-[234px]">
                  <img
                    className="flex-shrink-0"
                    src={GoogleLogo}
                    alt="Google"
                  />
                  <p className="text-ab-sm truncate text-black">Google</p>
                </div>
                <button
                  type="button"
                  className={`text-ab-sm rounded px-3 py-1.5 font-bold leading-tight text-white transition-all ${
                    userDetails?.provider?.includes('Google')
                      ? 'bg-ab-red hover:bg-opacity-80 hover:opacity-80 '
                      : 'btn-primary hover:opacity-90'
                  }`}
                >
                  {userDetails?.provider?.includes('Google')
                    ? 'Disconnect'
                    : 'Connect'}
                </button>
              </div>
              <div className="flex w-full items-center space-x-2">
                <div className="float-left flex flex-grow items-center space-x-2 overflow-hidden md:max-w-[234px]">
                  <img
                    className="flex-shrink-0"
                    src={LinkedInLogo}
                    alt="LinkedIn"
                  />
                  <p className="text-ab-sm truncate text-black">Linkedin</p>
                </div>

                <button
                  type="button"
                  className={`text-ab-sm rounded px-3 py-1.5 font-bold leading-tight text-white transition-all ${
                    userDetails?.provider?.includes('LinkedIn')
                      ? 'bg-ab-red hover:bg-opacity-80 hover:opacity-80 '
                      : 'btn-primary hover:opacity-90'
                  }`}
                >
                  {userDetails?.provider?.includes('LinkedIn')
                    ? 'Disconnect'
                    : 'Connect'}
                </button>
              </div>
              <div className="flex w-full items-center space-x-2">
                <div className="float-left flex flex-grow items-center space-x-2 overflow-hidden md:max-w-[234px]">
                  <img
                    className="flex-shrink-0"
                    src={TwitterLogo}
                    alt="Twitter"
                  />
                  <p className="text-ab-sm truncate text-black">X (Twitter)</p>
                </div>
                <button
                  type="button"
                  className={`text-ab-sm rounded px-3 py-1.5 font-bold leading-tight text-white transition-all ${
                    userDetails?.provider?.includes('Twitter')
                      ? 'bg-ab-red hover:bg-opacity-80 hover:opacity-80 '
                      : 'btn-primary hover:opacity-90'
                  }`}
                >
                  {userDetails?.provider?.includes('Twitter')
                    ? 'Disconnect'
                    : 'Connect'}
                </button>
              </div>
              {/* <div className='flex w-full items-center space-x-2'>
              <div className='float-left flex flex-grow items-center space-x-2 overflow-hidden md:max-w-[234px]'>
                <img className='flex-shrink-0' src={GithubLogo} alt='Google' />
                <p className='text-ab-sm truncate text-black'>Github</p>
              </div>
              <button
                type='button'
                className='bg-ab-red text-ab-sm rounded px-3 py-1.5 font-bold leading-tight text-white transition-all hover:bg-opacity-80 hover:opacity-80'
              >
                Disconnect
              </button>
            </div> */}
            </div>
          </div>
          {/* <div className='float-left mb-6 w-full'>
            <p className='border-ab-gray-dark my-1 border-b pb-2.5 text-sm font-semibold text-black'>
              Stripe account
            </p>
            <div className='float-left mt-3 w-full'> */}
          {/* <p className="text-ab-sm text-[#636363]">No Account found!</p> */}
          {/* <p className='text-ab-sm font-medium text-black'>ACC091282</p> */}
          {/* <button type="button" className="btn-default bg-secondary rounded px-3 py-2 mt-3 leading-tight text-ab-sm font-bold text-white transition-all hover:opacity-90">
                      Connect your stripe account
                    </button> */}
          {/* <button
                type='button'
                disabled={hasStripeConnected || !stripeConnectUrl}
                onClick={() => window.location.replace(stripeConnectUrl)}
                className={`btn-default ${
                  hasStripeConnected ? 'bg-ab-red' : 'bg-secondary'
                } text-ab-sm mt-3 rounded px-3 py-2 font-bold leading-tight text-white transition-all hover:opacity-90`}
              >
                {hasStripeConnected
                  ? 'Disconnect your stripe account'
                  : 'Connect your stripe account'}
              </button>
            </div>
          </div> */}
        </div>
      </div>
    </>
  )
}

export default ProfileSettings
