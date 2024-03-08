/* eslint-disable react/prop-types */
import React, { useState, useEffect, useRef, useContext } from 'react'
import { Formik } from 'formik'
import useOnclickOutside from 'react-cool-onclickoutside'
import apiHelper from '../common/helpers/apiGetters'
import countries from './countries'
import { SpaceNameUpdateSchema } from '../common/validation/validation'
import MyContext from '../common/my-context'
import DownArrow from '../../../assets/img/icons/down-arrow.svg'
import { ACLContext } from '../../../context/ACLContext'
// import { set } from 'lodash'
// import EditIcon from '../../../../assets/img/icons/edit-icon.svg';

const initialValues = {
  current_name: '',
  name: '',
  belongs: 'P',
  email: '',
  current_business_name: '',
  business_name: '',
  address: '',
  country: '',
}

const Settings = (props) => {
  const { spaceDetails, onUpdateSpace } = props
  const formikRef = useRef()
  const { spaceId } = useContext(MyContext)

  const { setSpaceData } = useContext(ACLContext)

  // const [selectedImage, setSelectedImage] = useState(null);
  const [countryDropdown, setCountryDropdown] = useState(false)
  const [loader, setLoader] = useState(false)
  const [flag, setFlag] = useState(false)
  const [countrySearchText, setCountrySearchText] = useState('')

  const countryDropContainer = useOnclickOutside(() => {
    setCountryDropdown(false)
  })

  const setFormikField = (name, value) => {
    formikRef?.current?.setFieldValue(name, value)
  }

  const setFormikFieldTouched = (name, value) => {
    formikRef?.current?.setFieldTouched(name, value)
  }

  const updateSpaceData = async () => {
    const res = await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.LIST_SPACES_URL,
      value: {
        search_keyword: null,
      },
    })
    setSpaceData(res)
  }

  useEffect(async () => {
    setFormikField('current_name', spaceDetails?.name)
    setFormikField('name', spaceDetails?.name)
    setFormikField('belongs', spaceDetails?.type)
    setFormikField('email', spaceDetails?.email)
    setFormikField('country', spaceDetails?.country)
    setFormikField('current_business_name', spaceDetails?.business_name)
    setFormikField('business_name', spaceDetails?.business_name)
    setFormikField('address', spaceDetails?.address)
  }, [spaceId, spaceDetails, flag])

  const setTouchedToFalse = () => {
    setFormikFieldTouched('current_name', false)
    setFormikFieldTouched('name', false)
    setFormikFieldTouched('belongs', false)
    setFormikFieldTouched('email', false)
    setFormikFieldTouched('country', false)
    setFormikFieldTouched('current_business_name', false)
    setFormikFieldTouched('business_name', false)
    setFormikFieldTouched('address', false)
  }

  const onSubmit = async (values) => {
    setLoader(true)
    const personalBody = spaceDetails?.type === 'P' && {
      name: values?.name,
      type: 'P',
      description: '',
      developer_portal_access: spaceDetails?.developer_portal_access,
      space_id: spaceId,
    }
    const businessBody = spaceDetails?.type === 'B' && {
      name: values?.name,
      type: 'B',
      email: values?.email || '',
      country: values?.country || '',
      business_name: values?.business_name || '',
      address: values?.address || '',
      business_category: '',
      description: '',
      developer_portal_access: values?.developer_portal_access,
      space_id: spaceId,
    }
    await apiHelper({
      baseUrl: process.env.BLOCK_ENV_URL_API_BASE_URL,
      subUrl: process.env.UPDATE_SPACE_URL,
      value: spaceDetails?.type === 'B' ? businessBody : personalBody,
      apiType: 'put',
      showSuccessMessage: true,
      spaceId,
    })
    updateSpaceData()
    onUpdateSpace()
    setLoader(false)
    setTouchedToFalse()
  }

  const onChange = (data) => {
    setFormikField(data.target.name, data.target.value)
  }

  const onCancel = () => {
    setFlag((flg) => !flg)
    setTouchedToFalse()
  }

  const onCountryChange = (name, country) => {
    setFormikField(name, country)
    setCountryDropdown(false)
  }

  // const [spaceType, setSpaceType] = useState('business');
  return (
    <Formik
      innerRef={formikRef}
      initialValues={initialValues}
      onSubmit={onSubmit}
      validationSchema={SpaceNameUpdateSchema}
      validateOnMount
      validateOnChange
      validateOnBlur
      enableReinitialize
    >
      {({ handleSubmit, values, errors, touched }) => (
        <div className="float-left w-full">
          <div className="border-ab-gray-dark mb-2 flex items-center justify-between border-b py-2.5">
            <p className="text-lg font-medium text-black">Space Settings</p>
          </div>
          <div className="float-left mt-4 flex w-full flex-wrap">
            <div className="float-left mb-6 flex w-full max-w-2xl flex-col">
              <div className="flex flex-col float-left mb-6 w-80">
                <label className="text-ab-sm float-left mb-2 font-medium text-black">
                  Space Name
                </label>
                <input
                  value={values.name}
                  name="name"
                  onChange={onChange}
                  placeholder="Enter Space Name"
                  type="text"
                  className={`text-ab-sm ${
                    touched.name && errors.name
                      ? 'border-ab-red'
                      : 'border-ab-gray-light focus:border-primary'
                  } bg-ab-gray-light float-left w-full rounded-md border py-3.5 px-4 focus:outline-none`}
                />
                <p className="text-xs text-ab-red left-0 mt-0.5">
                  {touched.name && errors.name}
                </p>
              </div>
              {values?.belongs === 'B' && (
                <div className="float-left w-full">
                  <div className="flex justify-between flex-wrap">
                    <div className="flex flex-col float-left mb-6 w-80">
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
                    <div className="flex flex-col float-left mb-6 w-80">
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
                  </div>
                  <div className="flex justify-between flex-wrap">
                    <div className="flex flex-col float-left mb-6 w-80">
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
                    <div className="flex flex-col float-left mb-6 w-80">
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
                            onChange={(e) =>
                              setCountrySearchText(e.target.value)
                            }
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
                </div>
              )}
              <div className="float-left mb-10 mt-1 flex w-full items-center">
                <button
                  type="button"
                  disabled={loader}
                  onClick={handleSubmit}
                  className="btn-primary text-ab-sm mr-4 rounded px-5 py-2.5 font-bold leading-tight text-white transition-all hover:opacity-90"
                >
                  Save Changes
                </button>
                <button
                  type="button"
                  onClick={onCancel}
                  className="text-ab-disabled hover:text-ab-black text-ab-sm rounded px-3 py-1 font-bold leading-tight text-white"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </Formik>
  )
}

export default Settings
