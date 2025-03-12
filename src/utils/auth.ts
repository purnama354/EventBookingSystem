/**
 * Parses a JWT token and returns the payload as an object
 */
export function parseJwt(token: string): any {
    try {
      const base64Url = token.split('.')[1]
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      )
      return JSON.parse(jsonPayload)
    } catch (e) {
      console.error("Failed to parse JWT", e)
      return null
    }
  }
  
  /**
   * Checks if the current user is authenticated based on token in localStorage
   */
  export function isAuthenticated(): boolean {
    return localStorage.getItem("token") !== null
  }
  
  /**
   * Checks if the current user has admin privileges
   */
  export function isAdmin(): boolean {
    const token = localStorage.getItem("token")
    if (!token) return false
    
    const decoded = parseJwt(token)
    return decoded && decoded.isAdmin === true
  }