import * as Yup from 'yup'

const PasswordValidationSchema = () =>
  Yup.object().shape({
    current_password: Yup.string().required('Password cannot be blank'),
    new_password: Yup.string()
      .required('Password cannot be blank')
      .matches(
        /^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=[\]{};':"\\|,.<>\\/?])[^\s]{8,}$/,
        'Password must have minimum eight characters, at least one Uppercase letter, one Lowercase letter, one Number and one Special Character'
      ),
  })

export default PasswordValidationSchema
