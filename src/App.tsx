import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import Register from "./pages/Register"
import Login from "./pages/Login"
import Home from "./pages/Home"
import { Toaster } from "react-hot-toast"
import EventList from "./component/EventList"
import AdminEventList from "./component/AdminEvenList"
import NotAuthorized from "./pages/NotAuthorized"
import { useEffect, useState } from "react"
import { isAuthenticated, isAdmin } from "./utils/auth"

function App() {

  

  return (
    <>
      <BrowserRouter>
        <Routes>

        </Routes>
      </BrowserRouter>
      <Toaster />
    </>
  )
}

export default App
