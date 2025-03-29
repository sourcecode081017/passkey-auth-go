import { useState } from 'react'
import { create } from '@github/webauthn-json'

interface RegisterProps {
  username: string
}

function Register({ username }: RegisterProps) {
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [isRegistering, setIsRegistering] = useState(false)

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
      // log the publicKeyOptions to see the challenge and other parameters
        console.log('PublicKeyOptions:', publicKeyOptions)

      // Step 3: Create credentials using the webauthn-json library
      const credential = await create(publicKeyOptions)

      // Step 4: Send the credential to your server for verification
      const verifyResponse = await fetch('http://localhost:8080/passkey-auth/register-complete', {
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

  return (
    <div>
      <button 
        onClick={handleRegister} 
        disabled={isRegistering || !username.trim()}
      >
        {isRegistering ? 'Registering...' : 'Register'}
      </button>
      {error && <p className="error">{error}</p>}
      {message && <p className="message">{message}</p>}
    </div>
  )
}

export default Register