import * as Yup from 'yup'
import { shield } from '@appblocks/js-sdk'
import Axios from 'axios'
import { debounce } from 'lodash'

const typeURLMapping = {
  Team: process.env.CHECK_TEAM_NAME_URL,
  Role: process.env.CHECK_ROLE_NAME_URL,
}

const checkDuplicate = async (payload, apiUrl, existsExist = true) => {
  const token = shield.tokenStore.getToken()
  const headers = {
    Authorization: `Bearer ${token}`,
    space_id: payload?.space_id,
  }
  return new Promise((resolve) => {
    let isDuplicateExists
    Axios.post(
      `${process.env.BLOCK_ENV_URL_API_BASE_URL}${apiUrl}/invoke`,
      payload,
      {
        headers,
      },
    )
      .then((valueFromAPIResponse) => {
        isDuplicateExists = existsExist
          ? valueFromAPIResponse.data.data.exists
          : valueFromAPIResponse.data.data // boolean: true or false
        resolve(isDuplicateExists)
      })
      .catch(() => {
        isDuplicateExists = false
        resolve(isDuplicateExists)
      })
  })
}

const handler = debounce((payload, apiUrl, existsExist) => {
  checkDuplicate(payload, apiUrl, existsExist)
}, 1000)

export const SpaceNameUpdateSchema = () =>
  Yup.object().shape({
    current_name: Yup.string().required(),
    name: Yup.string()
      .min(2, 'Please enter a name more than 1 character')
      .required('Name is required')
      .matches(/^[A-Za-z0-9_]+$/, 'Only alphanumeric characters are allowed')
      .test(
        'checkDuplicateSpaceName',
        'Space name already exists',
        async (value, currentValues) => {
          if (value && value !== currentValues?.parent?.current_name) {
            const isDuplicateExists = await handler(
              {
                name: value,
              },
              process.env.CHECK_SPACE_NAME_URL,
            )
            return !isDuplicateExists
          }
          return true
        },
      ),
    belongs: Yup.string().required('Please select a type'),
    current_business_name: Yup.string().optional(),
    business_name: Yup.string().when('belongs', {
      is: 'B',
      then: () =>
        Yup.string()
          .min(2, 'Please enter a name more than 1 character')
          .required('Name is required')
          // .matches(/^[A-Za-z0-9]+$/, 'Only alphanumeric characters are allowed')
          .test(
            'checkDuplicateBusinessName',
            'Business name already exists',
            async (value, currentValues) => {
              if (
                value &&
                value !== currentValues?.parent?.current_business_name
              ) {
                const isDuplicateExists = await checkDuplicate(
                  {
                    name: value,
                  },
                  process.env.CHECK_BUSINESS_NAME_URL,
                )
                return !isDuplicateExists
              }
              return true
            },
          ),
    }),
    email: Yup.string().when('belongs', {
      is: 'business',
      then: () =>
        Yup.string()
          .matches(
            /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/,
            'Enter a valid email ID',
          )
          .required('Email is required'),
    }),
    address: Yup.string().when('belongs', {
      is: 'business',
      then: () => Yup.string().required('Address is required'),
    }),
    country: Yup.string().when('belongs', {
      is: 'business',
      then: () => Yup.string().required('Country is required'),
    }),
  })

export const CreateNewValidation = (type, spaceId) =>
  Yup.object().shape({
    name: Yup.string()
      .required('This field cannot be blank')
      .matches(/^[A-Za-z0-9_-]+$/, 'Only alphanumeric characters are allowed')
      .test(
        `checkDuplicate${type}Name`,
        `${type} name already exists`,
        async (value) => {
          if (value) {
            console.log(spaceId)
            const isDuplicateExists = await checkDuplicate(
              {
                name: value,
                space_id: spaceId,
              },
              typeURLMapping[type],
              true,
            )
            return !isDuplicateExists
          }
          return true
        },
      ),
  })

export const CreateNewAppValidation = (spaceId, type_id) =>
  Yup.object().shape({
    name: Yup.string()
      .required('This field cannot be blank')
      .matches(/^[A-Za-z0-9_-]+$/, 'Only alphanumeric characters are allowed')
      .test(
        'checkDuplicateAppName',
        'App name already exists',
        async (value) => {
          if (value) {
            const isDuplicateExists = await checkDuplicate(
              {
                name: value,
                space_id: spaceId,
                type_id,
              },
              process.env.CHECK_ENTITY_NAME_URL,
              true,
            )
            return !isDuplicateExists
          }
          return true
        },
      ),
  })

export const MemberValidationSchema = () =>
  Yup.object().shape({
    name: Yup.string(),
    role: Yup.array(),
    team: Yup.array(),
  })
