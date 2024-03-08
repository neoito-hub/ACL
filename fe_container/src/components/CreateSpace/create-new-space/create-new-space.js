import React, { useState, useRef, useContext } from 'react'
import useOnclickOutside from 'react-cool-onclickoutside'
import { Link, useHistory } from 'react-router-dom'
import { Formik } from 'formik'
import axios from 'axios'
import { shield } from '@appblocks/js-sdk'
import Toast from '../toast/Toast'
import DownArrow from '../../../assets/img/icons/down-arrow.svg'
import SpaceValidationSchema from './validation/validation'
import countries from './countries'
import { ACLContext } from '../../../context/ACLContext'
// import EditIcon from "../../assets/img/icons/edit-icon.svg";

const initialValues = {
  name: '',
  belongs: 'personal',
  email: '',
  business_name: '',
  address: '',
  country: '',
  // acceptTerms: false,
}

const CreateNewSpace = () => {
  const formikRef = useRef()
  const history = useHistory()

  const { setSpaceData } = useContext(ACLContext)

  const [countryDropdown, setCountryDropdown] = useState(false)
  const [spaceType, setSpaceType] = useState('personal')
  const [countrySearchText, setCountrySearchText] = useState('')

  const countryDropContainer = useOnclickOutside(() => {
    setCountryDropdown(false)
    setCountrySearchText('')
  })

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
      return data
    } catch (err) {
      console.log('msg', err)
      if (err.response.status === 401) shield.logout()
    }
  }

  const updateSpaceData = async () => {
    const res = await apiHelper(
      process.env.BLOCK_ENV_URL_API_BASE_URL,
      process.env.LIST_SPACES_URL,
      {
        search_keyword: null,
      }
    )
    setSpaceData(res?.data)
  }

  const createSpace = async (values) => {
    const res = await apiHelper(
      process.env.BLOCK_ENV_URL_API_BASE_URL,
      process.env.CREATE_SPACE_URL,
      values
    )
    console.log(res)
    if (res) {
      history.push('/')
      formikRef?.current?.resetForm(initialValues)
      updateSpaceData()
    }
  }

  const onSubmit = (values) => {
    const personalBody = {
      name: values.name,
      type:
        values.belongs === 'personal'
          ? 'P'
          : values.belongs === 'business'
          ? 'B'
          : '',
      description: '',
      developer_portal_access: true,
    }
    const businessBody = {
      name: values.name,
      type:
        values.belongs === 'personal'
          ? 'P'
          : values.belongs === 'business'
          ? 'B'
          : '',
      email: values.email || '',
      country: values.country || '',
      business_name: values.business_name || '',
      address: values.address || '',
      business_category: '',
      description: '',
      developer_portal_access: true,
    }
    // isEdit && Object.assign(personalBody, { space_id: listSpaceId });
    // isEdit && Object.assign(businessBody, { space_id: listSpaceId });

    // createSpace(values.belongs === 'business' ? businessBody : personalBody,)
    // formikRef?.current?.resetForm(initialValues)
    // isEdit
    //   ? updateSpace(values.belongs === 'business' ? businessBody : personalBody)
    //   :
    createSpace(values.belongs === 'business' ? businessBody : personalBody)
  }

  const onChange = (data) => {
    formikRef?.current?.setFieldValue(data.target.name, data.target.value)
  }

  const onCountryChange = (name, country) => {
    formikRef?.current?.setFieldValue(name, country)
    setCountryDropdown(false)
  }

  // const onChangeCheckBox = (data) => {
  //   formikRef?.current?.setFieldValue(data.target.name, data.target.checked);
  // };

  return (
    <>
      <Toast />

      <Formik
        innerRef={formikRef}
        initialValues={initialValues}
        onSubmit={onSubmit}
        validationSchema={SpaceValidationSchema}
        validateOnMount
        validateOnChange
        validateOnBlur
        enableReinitialize
      >
        {({ handleSubmit, values, errors, touched, setFieldValue }) => (
          <div className="float-left w-full py-6 md-lt:px-4">
            <h4 className="text-ab-black mt-1 mb-5 text-xl font-semibold">
              Create New Space
            </h4>
            <div className="float-left mb-6 flex w-full max-w-xs flex-col">
              {/* <div className="bg-secondary relative mb-6 flex h-16 w-16 items-center justify-center rounded-full text-white">
          <span className="text-3xl font-bold">S</span>
          <label className="absolute bottom-0 right-0 flex h-5 w-5 cursor-pointer items-center justify-center rounded-full bg-white shadow-lg">
            <input type="file" className="hidden"></input>
            <img src={EditIcon} alt="Edit"></img>
          </label>
        </div> */}
              <div className="flex flex-col float-left mb-6 w-full">
                <label className="text-ab-sm float-left mb-2 font-medium text-black">
                  Space Name
                </label>
                <input
                  value={values.name}
                  name="name"
                  onChange={onChange}
                  placeholder="Enter Space Name"
                  type="text"
                  className={`${
                    touched.name && errors.name
                      ? 'border-ab-red'
                      : 'border-ab-gray-light focus:border-primary'
                  } text-ab-sm bg-ab-gray-light float-left w-full rounded-md border py-3.5 px-4 focus:outline-none`}
                />
                <p className="text-xs text-ab-red left-0 mt-0.5">
                  {touched.name && errors.name}
                </p>
              </div>
              <div className="float-left mb-6 w-full">
                <label className="text-ab-sm float-left mb-2 font-medium text-black">
                  Belongs to:
                </label>
                <div className="float-left flex w-full">
                  <label className="mt-2 mr-10 flex cursor-pointer items-center">
                    <input
                      name="belongs"
                      onChange={(e) => {
                        setSpaceType('personal')
                        onChange(e)
                      }}
                      checked={spaceType === 'personal'}
                      className="peer hidden"
                      type="radio"
                      value="personal"
                    />
                    <span className="chkbox-icon border-ab-disabled float-left mr-2 h-5 w-5 flex-shrink-0 rounded border bg-white" />
                    <p className="text-xs text-black">Personal</p>
                  </label>
                  <label className="mt-2 flex cursor-pointer items-center">
                    <input
                      name="belongs"
                      onChange={(e) => {
                        setSpaceType('business')
                        onChange(e)
                      }}
                      checked={spaceType === 'business'}
                      className="peer hidden"
                      type="radio"
                      value="business"
                    />
                    <span className="chkbox-icon border-ab-disabled float-left mr-2 h-5 w-5 flex-shrink-0 rounded border bg-white" />
                    <p className="text-xs text-black">Business</p>
                  </label>
                </div>
              </div>
              {spaceType === 'business' && (
                <div className="float-left w-full">
                  <div className="flex flex-col float-left mb-6 w-full">
                    <label className="text-ab-sm float-left mb-2 font-medium text-black">
                      Business Name
                    </label>
                    <input
                      value={values.business_name}
                      onChange={onChange}
                      name="business_name"
                      placeholder="Business Name"
                      type="text"
                      className={`text-ab-sm ${
                        touched.business_name && errors.business_name
                          ? 'border-ab-red'
                          : 'border-ab-gray-light focus:border-primary'
                      } bg-ab-gray-light  float-left w-full rounded-md border py-3.5 px-4 focus:outline-none`}
                    />
                    <p className="text-xs text-ab-red left-0 mt-0.5">
                      {touched.business_name && errors.business_name}
                    </p>
                  </div>
                  <div className="flex flex-col float-left mb-6 w-full">
                    <label className="text-ab-sm float-left mb-2 font-medium text-black">
                      Business e-mail
                    </label>
                    <input
                      value={values.email}
                      onChange={onChange}
                      name="email"
                      placeholder="Business Email"
                      type="text"
                      className={`text-ab-sm ${
                        touched.email && errors.email
                          ? 'border-ab-red'
                          : 'border-ab-gray-light focus:border-primary'
                      } bg-ab-gray-light float-left w-full rounded-md border py-3.5 px-4 focus:outline-none`}
                    />
                    <p className="text-xs text-ab-red left-0 mt-0.5">
                      {touched.email && errors.email}
                    </p>
                  </div>
                  <div className="flex flex-col float-left mb-6 w-full">
                    <label className="text-ab-sm float-left mb-2 font-medium text-black">
                      Business Address
                    </label>
                    <textarea
                      value={values.address}
                      onChange={onChange}
                      name="address"
                      placeholder="Business Address"
                      type="text"
                      className={`text-ab-sm ${
                        touched.business_name && errors.business_name
                          ? 'border-ab-red'
                          : 'border-ab-gray-light focus:border-primary'
                      } bg-ab-gray-light custom-scroll-primary float-left h-32 w-full resize-none rounded-md border py-3.5 px-4 focus:outline-none`}
                    />
                    <p className="text-xs text-ab-red left-0 mt-0.5">
                      {touched.address && errors.address}
                    </p>
                  </div>
                  <div className="flex flex-col float-left mb-6 w-full">
                    <label className="text-ab-sm float-left mb-2 font-medium text-black">
                      Country
                    </label>
                    <div
                      className="relative float-left w-full"
                      ref={countryDropContainer}
                    >
                      <div
                        onClick={() => setCountryDropdown(!countryDropdown)}
                        className={`text-ab-sm bg-ab-gray-light focus:border-primary float-left flex w-full cursor-pointer select-none items-center justify-between rounded-md border py-3.5 px-4 focus:outline-none ${
                          touched.country && errors.country
                            ? 'border-ab-red'
                            : countryDropdown
                            ? 'border-primary'
                            : 'border-ab-gray-light'
                        }`}
                      >
                        <p className="text-ab-black text-ab-sm truncate">
                          {values.country || 'Choose Country'}
                        </p>
                        <img
                          className={`ml-3 flex-shrink-0 transform transition-transform duration-300 ${
                            countryDropdown ? 'rotate-180' : ''
                          }`}
                          src={DownArrow}
                          alt=""
                        />
                      </div>

                      <ul
                        className={`curson-pointer shadow-box dropDownFade border-ab-gray-medium custom-scroll-primary absolute top-12 left-0 z-10 max-h-[150px] w-full overflow-y-auto border bg-white py-2 ${
                          countryDropdown ? '' : 'hidden'
                        }`}
                      >
                        <input
                          type="text"
                          value={countrySearchText}
                          autoComplete="new-password"
                          onChange={(e) => setCountrySearchText(e.target.value)}
                          className="search-input-xs border-ab-gray-dark mb-2 h-8 w-full border-b !bg-[length:12px_12px] px-2 pl-8 text-[13px] focus:outline-none"
                          placeholder="Search countries"
                        />
                        {countries?.map(
                          (country) =>
                            country?.name
                              .toLowerCase()
                              .match(countrySearchText.toLowerCase()) && (
                              <li
                                key={country?.code}
                                className="px-3 py-1 cursor-pointer group"
                                onClick={() => {
                                  onCountryChange('country', country?.name)
                                  setCountrySearchText('')
                                }}
                              >
                                <span className="text-ab-black text-ab-sm group-hover:text-primary">
                                  {country?.name}
                                </span>
                              </li>
                            )
                        )}
                      </ul>
                    </div>
                    <p className="text-xs text-ab-red left-0 mt-0.5">
                      {touched.country && errors.country}
                    </p>
                  </div>
                </div>
              )}
              {/* <div className='flex flex-col float-left mb-6 w-full'>
                <div className='my-2 flex w-full items-center'>
                  <label className='float-left flex items-center'>
                    <input
                      checked={values.acceptTerms}
                      name='acceptTerms'
                      onChange={onChangeCheckBox}
                      className='peer hidden'
                      type='checkbox'
                    />
                    <span
                      className={`${
                        touched.acceptTerms && errors.acceptTerms
                          ? 'border-ab-red'
                          : 'border-ab-disabled'
                      } chkbox-icon float-left mr-2 h-5 w-5 flex-shrink-0 cursor-pointer rounded border bg-white`}
                    ></span>
                  </label>
                  <p className='text-xs tracking-tight text-black'>
                    I Accept the Appblocks{' '}
                    <a className='text-primary underline'>terms of services</a>{' '}
                    of appblocks
                  </p>
                </div>
                <p className='text-xs text-ab-red left-0'>
                  {touched.acceptTerms && errors.acceptTerms}
                </p>
              </div> */}
              <div className="float-left mb-6 flex w-full items-center">
                <button
                  onClick={handleSubmit}
                  type="button"
                  className="btn-primary text-ab-sm mr-4 rounded px-5 py-2.5 font-bold leading-tight text-white transition-all hover:opacity-90"
                >
                  Create Space
                </button>
                <Link
                  to="/"
                  className="text-ab-disabled hover:text-ab-black text-ab-sm rounded px-3 py-1 font-bold leading-tight text-white"
                >
                  Cancel
                </Link>
              </div>
            </div>
          </div>
        )}
      </Formik>
    </>
  )
}

export default CreateNewSpace
