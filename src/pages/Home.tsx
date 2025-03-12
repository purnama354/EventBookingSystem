import { Link } from "react-router-dom"

export default function Home() {
  return (
    <>
      <div className="min-h-screen bg-base-200 flex items-center justify-center p-4">
        <div className="card w-full max-w-md bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title text-2xl font-bold text-center mb-2">
              Welcome to Event Booking System
            </h2>
            <p className="text-center text-base-content/70 mb-6">
              Find and book your favorite events with ease!
            </p>

            <div className="space-y-4">
              <p>
                Our platform allows you to discover a wide range of events, from
                concerts and workshops to conferences and festivals.
              </p>
              <p>
                Create an account to browse events, book tickets, and manage
                your bookings.
              </p>
            </div>

            <div className="divider my-6">Get Started</div>

            <div className="flex flex-col gap-3">
              <Link to="/register" className="btn btn-primary">
                Sign Up
              </Link>
              <Link to="/login" className="btn btn-outline">
                Login
              </Link>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
