import * as Yup from 'yup'
import { shield } from '@appblocks/js-sdk'
import Axios from 'axios'

function checkDuplicate(valueToCheck, apiUrl) {
  const token = shield.tokenStore.getToken()
  return new Promise((resolve) => {
    let isDuplicateExists
    Axios.post(
      `${process.env.BLOCK_ENV_URL_API_BASE_URL}${apiUrl}/invoke`,
      {
        name: valueToCheck,
      },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    )
      .then((valueFromAPIResponse) => {
        isDuplicateExists = valueFromAPIResponse.data.data.exists // boolean: true or false
        resolve(isDuplicateExists)
      })
      .catch(() => {
        isDuplicateExists = false
        resolve(isDuplicateExists)
      })
  })
}
const SpaceValidationSchema = () =>
  Yup.object().shape({
    name: Yup.string()
      .min(2, 'Please enter a name more than 1 character')
      .required('Name is required')
      .matches(/^[A-Za-z0-9_-]+$/, 'Only alphanumeric characters are allowed')
      .test(
        'checkDuplicateSpaceName',
        'Space name already exists',
        async (value) => {
          if (value) {
            const isDuplicateExists = await checkDuplicate(
              value,
              process.env.CHECK_SPACE_NAME_URL
            )
            return !isDuplicateExists
          }
          return true
        }
      ),
    belongs: Yup.string().required('Please select a type'),
    business_name: Yup.string().when('belongs', {
      is: 'business',
      then: () =>
        Yup.string()
          .min(2, 'Please enter a name more than 1 character')
          .required('Name is required')
          .test(
            'checkDuplicateBusinessName',
            'Business name already exists',
            async (value) => {
              if (value) {
                const isDuplicateExists = await checkDuplicate(
                  value,
                  process.env.CHECK_BUSINESS_NAME_URL
                )
                return !isDuplicateExists
              }
              return true
            }
          ),
    }),
    email: Yup.string().when('belongs', {
      is: 'business',
      then: () =>
        Yup.string()
          .matches(
            /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/,
            'Enter a valid email ID'
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
    // acceptTerms: Yup.boolean().oneOf([true], 'Required'),
  })

export default SpaceValidationSchema
