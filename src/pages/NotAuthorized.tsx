import { Link } from "react-router-dom"

export default function NotAuthorized() {
  return (
    <div className="min-h-screen bg-base-200 flex items-center justify-center p-4">
      <div className="card w-full max-w-md bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl font-bold text-center mb-2">
            Not Authorized
          </h2>
          <p className="text-center text-base-content/70 mb-6">
            You do not have permission to access this page.
          </p>
          <div className="flex justify-center">
            <Link to="/login" className="btn btn-primary">
              Go to Login
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}