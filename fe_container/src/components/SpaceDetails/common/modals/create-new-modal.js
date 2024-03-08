/* eslint-disable no-unused-expressions */
/* eslint-disable react/prop-types */
import React, { useRef, useState, useEffect, useContext } from 'react'
import { Formik } from 'formik'
import useOnclickOutside from 'react-cool-onclickoutside'
import CreateNewIcon from '../../../../assets/img/icons/user-icon.gif'
import {
  CreateNewAppValidation,
  CreateNewValidation,
} from '../validation/validation'
import MyContext from '../my-context'
import DownArrow from '../../../../assets/img/icons/down-arrow.svg'

const initialValues = { name: '' }

const CreateNewModal = (props) => {
  const {
    hasCreateNewModal,
    handleCreateNewModal,
    handleClick,
    type,
    placeHolder,
    current,
    spaceId,
  } = props
  const formikRef = useRef(null)
  const { entityTypeList } = useContext(MyContext)

  const [entityTypeDropdown, setEntityTypeDropdown] = useState(false)
  const [entityType, setEntityType] = useState(entityTypeList[0])
  const entityTypeDropContainer = useOnclickOutside(() => {
    setEntityTypeDropdown(false)
  })

  useEffect(() => {
    current && formikRef?.current?.setFieldValue('name', current?.name)
  })

  const onChange = (data) => {
    formikRef?.current?.setFieldValue(data.target.name, data.target.value)
  }

  return (
    <Formik
      innerRef={formikRef}
      initialValues={initialValues}
      onSubmit={(values) => handleClick(values?.name, Number(entityType.id))}
      validationSchema={
        type === 'Entity'
          ? CreateNewAppValidation(spaceId, entityType.id)
          : CreateNewValidation(type, spaceId)
      }
      validateOnMount
      validateOnChange={false}
      validateOnBlur={false}
      enableReinitialize
    >
      {({ handleSubmit, values, errors, touched }) => (
        <div
          className={`fixed left-0 top-0 z-[10001] h-screen w-full ${
            hasCreateNewModal ? 'fadeIn' : 'hidden'
          }`}
        >
          <div
            onClick={() => {
              handleCreateNewModal()
            }}
            className="fixed left-0 top-0 z-[10001] h-full w-full bg-black/40"
          />
          <div
            className={`absolute top-1/2 left-1/2 z-[10002] w-full max-w-[620px] -translate-x-1/2 -translate-y-1/2 transform px-4 ${
              hasCreateNewModal ? '' : 'hidden'
            }`}
          >
            <div className="relative float-left flex w-full rounded-md bg-white p-6 md:space-x-10 md:p-[60px] md-lt:flex-col">
              <img
                src={CreateNewIcon}
                alt="Create New"
                className="h-20 w-20 flex-shrink-0 rounded-full md-lt:mx-2 md-lt:mb-4"
              />
              <div className="flex flex-grow flex-col overflow-hidden">
                <h5 className="mb-3 text-lg font-semibold text-black">
                  {!current ? 'Create new ' : 'Edit '}
                  {type}
                </h5>
                {type === 'Entity' && (
                  <div className="text-ab-black float-left mb-3 w-full">
                    <label className="float-left text-xs font-medium">
                      Entity Type
                    </label>
                    <div
                      className="relative float-left my-1.5 w-full"
                      ref={entityTypeDropContainer}
                    >
                      <div
                        onClick={() => setEntityTypeDropdown((x) => !x)}
                        className={`text-ab-sm bg-ab-gray-light focus:border-primary float-left flex w-full cursor-pointer select-none items-center justify-between rounded-md border py-0.5 px-2.5 focus:outline-none ${
                          entityTypeDropdown
                            ? 'border-primary'
                            : 'border-ab-gray-light'
                        }`}
                      >
                        <div className="text-ab-black text-ab-sm flex-grow overflow-hidden">
                          <p className="flex h-[32px] items-center px-1 text-xs font-medium">
                            {entityType?.name || 'Select Entity Type'}
                          </p>
                        </div>
                        <img
                          className={`ml-3 flex-shrink-0 transform transition-transform duration-300 ${
                            entityTypeDropdown && 'rotate-180'
                          }`}
                          src={DownArrow}
                          alt=""
                        />
                      </div>
                      <div
                        className={`shadow-box dropDownFade border-ab-gray-medium position-inverse-2 absolute top-full left-0 z-10 mt-1 w-full border bg-white p-1 py-3 ${
                          entityTypeDropdown ? '' : 'hidden'
                        }`}
                      >
                        <ul className="custom-scroll-primary float-left max-h-[150px] w-full overflow-y-auto p-2">
                          {entityTypeList?.map((entityTyp) => (
                            <li
                              key={entityTyp.id}
                              onClick={() => {
                                setEntityType(entityTyp)
                                setEntityTypeDropdown((x) => !x)
                              }}
                              className={`float-left mb-4 w-full last-of-type:mb-0 cursor-pointer items-center leading-normal truncate text-xs font-medium tracking-tight hover:text-primary ${
                                entityTyp.id === entityType.id
                                  ? 'text-primary'
                                  : 'text-ab-black'
                              }`}
                            >
                              {entityTyp.name}
                            </li>
                          ))}
                        </ul>
                      </div>
                    </div>
                  </div>
                )}
                <label className="text-ab-black float-left text-xs font-medium">
                  {type} Name
                </label>

                <div className="float-left mb-2 w-full">
                  <input
                    value={values.name}
                    onChange={onChange}
                    placeholder={placeHolder}
                    type="text"
                    name="name"
                    className={`${
                      touched.name && errors.name
                        ? 'border-ab-red'
                        : 'border-ab-gray-light focus:border-primary'
                    } bg-ab-gray-light float-left w-full rounded-md border py-2.5 px-4 text-xs focus:outline-none`}
                  />
                  <p className="text-xs text-ab-red left-0 pt-10">
                    {touched.name && errors.name}
                  </p>
                </div>
                <div className="float-left mt-3 flex w-full items-center">
                  <button
                    type="button"
                    onClick={handleSubmit}
                    className="btn-primary text-ab-sm disabled:bg-ab-disabled mr-4 rounded px-5 py-2.5 font-bold leading-tight text-white transition-all"
                  >
                    {!current ? `Create new ${type}` : 'Save'}
                  </button>
                  <button
                    onClick={() => {
                      handleCreateNewModal()
                    }}
                    type="button"
                    className="text-ab-disabled hover:text-ab-black rounded px-3 py-1 text-sm font-bold leading-tight text-white"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </Formik>
  )
}

export default CreateNewModal
