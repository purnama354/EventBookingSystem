import React, { useState, useEffect } from "react"
import { Event } from "../types/types"
import Navigation from "./Navigation"
import { useNavigate } from "react-router-dom"
import toast from "react-hot-toast"

const EventList: React.FC = () => {
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const token = localStorage.getItem("token")
        if (!token) {
          toast.error("Missing authentication token")
          navigate("/login")
          return
        }

        const response = await fetch("http://localhost:8080/api/events", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        })

        if (response.status === 401) {
          toast.error("Your session has expired. Please login again.")
          localStorage.removeItem("token")
          navigate("/login")
          return
        }

        if (!response.ok) {
          const errorText = await response.text()
          throw new Error(errorText || `HTTP error! status: ${response.status}`)
        }

        const data: Event[] = await response.json()
        setEvents(data)
        setLoading(false)
      } catch (e: any) {
        console.error("Error fetching events:", e)
        setError(e.message || "Failed to fetch events")
        setLoading(false)
      }
    }

    fetchEvents()
  }, [navigate])

  if (loading) {
    return <div className="text-center p-10">Loading events...</div>
  }

  if (error) {
    return (
      <div className="p-10">
        <div className="alert alert-error">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>{error}</span>
        </div>
      </div>
    )
  }

  return (
    <>
      <Navigation />
      <div className="container mx-auto py-8">
        <h1 className="text-2xl font-bold mb-4">Upcoming Events</h1>
        {events.length === 0 ? (
          <div className="text-center p-10">
            <p>No events found. Check back later!</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {events.map((event) => (
              <div key={event.id} className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h2 className="card-title">{event.title}</h2>
                  <p>{event.description}</p>
                  <p>Date: {new Date(event.date).toLocaleDateString()}</p>
                  <p>Location: {event.location}</p>
                  <p>Capacity: {event.capacity}</p>
                  <div className="card-actions justify-end">
                    <button className="btn btn-primary">Book Now</button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  )
}

export default EventList
