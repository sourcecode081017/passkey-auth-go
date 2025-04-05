import { useState } from 'react'
import { create, get, parseCreationOptionsFromJSON, parseRequestOptionsFromJSON } from '@github/webauthn-json/browser-ponyfill'

interface RegisterAuthProps {
  username: string
}

function PasskeyRegisterAuthenticate({ username }: RegisterAuthProps) {
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [isRegistering, setIsRegistering] = useState(false)
  const [isAuthenticating, setIsAuthenticating] = useState(false)

  const handleRegister = async () => {
    if (!username.trim()) {
      setError('Please enter a username')
      return
    }

    setError('')
    setMessage('Starting registration...')
    setIsRegistering(true)

    try {
      // Step 1: Call your Go backend to initiate registration
      const response = await fetch(`http://localhost:8080/passkey-auth/register-initiate/${encodeURIComponent(username)}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`Server responded with status: ${response.status}`)
      }

      // Step 2: Get the challenge data from the server
      const publicKeyOptions = await response.json()
      const parsedOptions = parseCreationOptionsFromJSON(publicKeyOptions.options)

      // Step 3: Create credentials using the webauthn-json library
    const credential = await create(parsedOptions)

      // Step 4: Send the credential to your server for verification
      const verifyResponse = await fetch(`http://localhost:8080/passkey-auth/register-complete/${encodeURIComponent(username)}`, {
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
      setMessage(`Registration successful! ${result.message || ''}`)
    } catch (error) {
      console.error('Registration error:', error)
      setError(`Registration failed: ${(error as Error).message}`)
      setMessage('')
    } finally {
      setIsRegistering(false)
    }
  }

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
      const response = await fetch(`http://localhost:8080/passkey-auth/auth-initiate/${encodeURIComponent(username)}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`Server responded with status: ${response.status}`)
      }

      // Step 2: Get the challenge data from the server
      const publicKeyOptions = await response.json()
      const parsedOptions = parseRequestOptionsFromJSON(publicKeyOptions.options)

      // Step 3: Get credentials using the webauthn-json library
      const credential = await get(parsedOptions)

      // Step 4: Send the credential to your server for verification
      const verifyResponse = await fetch(`http://localhost:8080/passkey-auth/auth-complete/${encodeURIComponent(username)}`, {
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
        className="btn-register"
        onClick={handleRegister} 
        disabled={isRegistering || !username.trim()}
      >
        {isRegistering ? 'Registering...' : 'Register'}
      </button> &nbsp; &nbsp; &nbsp;
      <button
        className="btn-authenticate"
        onClick={handleAuthenticate} 
        disabled={isAuthenticating || !username.trim()}
      >
        {isRegistering ? 'Authenticating...' : 'Authenticate'}
      </button>
      {error && <p className="error">{error}</p>}
      {message && <p className="message">{message}</p>}
    </div>
  )
}

export default PasskeyRegisterAuthenticate