import { useState } from 'react'
import reactLogo from './assets/react.svg'
import goLogo from './assets/Go-Logo_Blue.svg'
import './App.css'
import PasskeyRegisterAuth from './PasskeyRegisterAuthenticate'
import UserDashboard from './UserDashboard'

function App() {
  const [username, setUsername] = useState('')
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  const handleSuccessfulAuth = () => {
    setIsAuthenticated(true)
  }

  const handleLogout = () => {
    setIsAuthenticated(false)
    // Optionally clear username if you want users to re-enter it
    // setUsername('')
  }

  return (
    <>
      <div>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
        <a href="https://go.dev/" target="_blank">
          <img src={goLogo} className="logo" alt="Go logo" />
        </a>
      </div>
      <h1>Passkey Registration and Authentication</h1>
      
      {!isAuthenticated ? (
        <div>
          <input
            type="text"
            placeholder="Enter username"
            className="username-input"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <PasskeyRegisterAuth 
            username={username} 
            onAuthSuccess={handleSuccessfulAuth}
          />
        </div>
      ) : (
        <UserDashboard
          username={username}
          onLogout={handleLogout}
        />
      )}
    </>
  )
}

export default App