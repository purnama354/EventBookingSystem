import React, { useState, useEffect } from "react"
import { Event } from "../types/types"
import Navigation from "./Navigation"
import { useNavigate } from "react-router-dom"
import toast from "react-hot-toast"

const AdminEventList: React.FC = () => {
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const token = localStorage.getItem("token")
        if (!token) {
          setError("Unauthorized: Missing token")
          setLoading(false)
          return
        }

        const response = await fetch("http://localhost:8080/api/admin/events", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        })

        if (response.status === 401) {
          navigate("/not-authorized")
          return
        }

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }

        const data: Event[] = await response.json()
        setEvents(data)
        setLoading(false)
      } catch (e: any) {
        setError(e.message || "Failed to fetch events")
        setLoading(false)
      }
    }

    fetchEvents()
  }, [navigate])

  const handleDeleteEvent = async (eventId: string) => {
    try {
      const token = localStorage.getItem("token")
      if (!token) return

      const response = await fetch(`http://localhost:8080/api/admin/events/${eventId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (response.ok) {
        setEvents(events.filter(event => event.id !== eventId))
        toast.success("Event deleted successfully")
      } else {
        toast.error("Failed to delete event")
      }
    } catch (e) {
      toast.error("An error occurred while deleting the event")
    }
  }

  if (loading) {
    return <div className="text-center">Loading events...</div>
  }

  if (error) {
    return <div className="alert alert-error">{error}</div>
  }

  return (
    <>
      <Navigation />
      <div className="container mx-auto py-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-bold">Admin - Manage Events</h1>
          <button 
            className="btn btn-primary" 
            onClick={() => navigate("/admin/events/new")}
          >
            Create New Event
          </button>
        </div>

        <div className="overflow-x-auto">
          <table className="table w-full">
            <thead>
              <tr>
                <th>Title</th>
                <th>Date</th>
                <th>Location</th>
                <th>Capacity</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {events.map((event) => (
                <tr key={event.id}>
                  <td>{event.title}</td>
                  <td>{new Date(event.date).toLocaleDateString()}</td>
                  <td>{event.location}</td>
                  <td>{event.capacity}</td>
                  <td className="flex gap-2">
                    <button 
                      className="btn btn-sm btn-outline btn-info"
                      onClick={() => navigate(`/admin/events/edit/${event.id}`)}
                    >
                      Edit
                    </button>
                    <button 
                      className="btn btn-sm btn-outline btn-error"
                      onClick={() => handleDeleteEvent(event.id)}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
              {events.length === 0 && (
                <tr>
                  <td colSpan={5} className="text-center py-4">
                    No events found. Create some events to get started!
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </>
  )
}

export default AdminEventList