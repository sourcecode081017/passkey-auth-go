import { useState } from 'react'
import { get } from '@github/webauthn-json'

interface AuthenticateProps {
  username: string
}

function Authenticate({ username }: AuthenticateProps) {
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [isAuthenticating, setIsAuthenticating] = useState(false)

  const handleAuthenticate = async () => {
    if (!username.trim()) {
      setError('Please enter a username')
      return
    }

    setError('')
    setMessage('Starting authentication...')
    setIsAuthenticating(true)

    try {
      // Step 1: Call your Go backend to initiate authentication
      const response = await fetch(`http://localhost:8080/passkey-auth/authenticate-begin/${encodeURIComponent(username)}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`Server responded with status: ${response.status}`)
      }

      // Step 2: Get the challenge data from the server
      const publicKeyOptions = await response.json()

      // Step 3: Get credentials using the webauthn-json library
      const credential = await get(publicKeyOptions)

      // Step 4: Send the credential to your server for verification
      const verifyResponse = await fetch('http://localhost:8080/passkey-auth/authenticate-complete', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credential),
      })

      if (!verifyResponse.ok) {
        throw new Error(`Verification failed with status: ${verifyResponse.status}`)
      }

      const result = await verifyResponse.json()
      setMessage(`Authentication successful! ${result.message || ''}`)
    } catch (error) {
      console.error('Authentication error:', error)
      setError(`Authentication failed: ${error instanceof Error ? error.message : 'Unknown error'}`)
      setMessage('')
    } finally {
      setIsAuthenticating(false)
    }
  }

  return (
    <div>
      <button 
        onClick={handleAuthenticate} 
        disabled={isAuthenticating || !username.trim()}
      >
        {isAuthenticating ? 'Authenticating...' : 'Authenticate'}
      </button>
      {error && <p className="error">{error}</p>}
      {message && <p className="message">{message}</p>}
    </div>
  )
}

export default Authenticate