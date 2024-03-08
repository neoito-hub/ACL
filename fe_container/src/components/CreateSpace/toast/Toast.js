import React from 'react'
import { ToastContainer, Flip } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
// import './ToastStyles.scss';

const Toast = () => (
  <div>
    <ToastContainer
      className="ab-toast"
      position="top-center"
      transition={Flip}
      hideProgressBar
      autoClose={5000}
      newestOnTop={false}
      closeOnClick
      rtl={false}
      pauseOnFocusLoss
      draggable
      pauseOnHover
    />
  </div>
)

export default Toast
