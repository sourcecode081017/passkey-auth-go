import { useState } from 'react'
import reactLogo from './assets/react.svg'
import goLogo from './assets/Go-Logo_Blue.svg'
import './App.css'
import Register from './Register'
import Authenticate from './Authenticate'

function App() {
  const [username, setUsername] = useState('')

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
      <h2>Passkey Registration and Authentication</h2>
      <div>
        <input
          type="text"
          placeholder="Enter username"
          className="username-input"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <div className="button-group">
          <Register username={username} /> &nbsp; &nbsp; &nbsp; &nbsp;
          <Authenticate username={username} />
        </div>
      </div>
    </>
  )
}

export default App