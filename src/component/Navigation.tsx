import { Link, useNavigate } from "react-router-dom"

const Navigation = () => {
  const navigate = useNavigate()

  const handleLogout = () => {
    localStorage.removeItem("token")
    navigate("/")
  }

  return (
    <div className="navbar bg-base-100">
      <div className="flex-1">
        <Link to="/events" className="btn btn-ghost text-xl">
          Event Booking System
        </Link>
      </div>
      <div className="flex-none">
        <button className="btn btn-ghost" onClick={handleLogout}>
          Logout
        </button>
      </div>
    </div>
  )
}

export default Navigation
