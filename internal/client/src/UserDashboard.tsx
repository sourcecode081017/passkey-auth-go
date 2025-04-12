import { useState, useEffect } from "react";

interface PasskeyInfo {
  Credential: {
    id: string;
    transport: string[];
    authenticator: {
      AAGUID: string;
    };
  };
  createdAt: string;
  lastUsedAt: string;
}

interface UserDashboardProps {
  username: string;
  onLogout: () => void;
}

function UserDashboard({ username, onLogout }: UserDashboardProps) {
  const [passkeys, setPasskeys] = useState<PasskeyInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [deleteStatus, setDeleteStatus] = useState("");

  useEffect(() => {
    fetchPasskeys();
  }, [username]);

  const fetchPasskeys = async () => {
    setLoading(true);
    setError("");
    
    try {
      const response = await fetch(
        `http://localhost:8080/passkey-auth/${encodeURIComponent(username)}/keys`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to fetch passkeys: ${response.status}`);
      }

      const data = await response.json();
      setPasskeys(data.keys || []);
    } catch (error) {
      console.error("Error fetching passkeys:", error);
      setError(`Failed to load passkeys: ${(error as Error).message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleDeletePasskey = async (passkeyId: string) => {
    setDeleteStatus("");
    
    try {
      const response = await fetch(
        `http://localhost:8080/passkey-auth/${encodeURIComponent(username)}/keys`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ passkeyId }),
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to delete passkey: ${response.status}`);
      }

    setDeleteStatus("Passkey deleted successfully");
    setTimeout(() => setDeleteStatus(""), 2000);
      // Refresh the passkeys list
      fetchPasskeys();
    } catch (error) {
      console.error("Error deleting passkey:", error);
      setDeleteStatus(`Failed to delete passkey: ${(error as Error).message}`);
    }
  };

  // Format Base64 ID for better display
  const formatId = (id: string) => {
    if (id.length > 10) {
      return `${id.substring(0, 10)}...`;
    }
    return id;
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      timeZoneName: "short",
    }).format(date);
  };

  return (
    <div className="user-dashboard">
      <h2>Welcome, {username}!</h2>
      <h3>Your Registered Passkeys</h3>
      
      {loading && <p>Loading your passkeys...</p>}
      {error && <p className="error">{error}</p>}
      {deleteStatus && <p className="status">{deleteStatus}</p>}
      
      {!loading && passkeys.length === 0 && (
        <p>You don't have any registered passkeys yet.</p>
      )}
      
      {passkeys.length > 0 && (
        <div className="passkeys-list">
          <table>
            <thead>
              <tr>
                <th>Credential ID</th>
                <th>Transports</th>
                <th>AAGUID</th>
                <th>Created At</th>
                <th>Last Used At</th>
              </tr>
            </thead>
            <tbody>
              {passkeys.map((passkey) => (
                <tr key={passkey.Credential.id}>
                  <td title={passkey.Credential.id}>{passkey.Credential.id}</td>
                  <td>{passkey.Credential.transport.join(", ")}</td>
                  <td>{passkey.Credential.authenticator.AAGUID}</td>
                  <td>{formatDate(passkey.createdAt)}</td>
                  <td>{passkey.lastUsedAt ? formatDate(passkey.lastUsedAt) : "Never Used"}</td>
                  <td>
                    <button
                      className="btn-delete"
                      onClick={() => handleDeletePasskey(passkey.Credential.id)}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      
      <div className="dashboard-actions">
        <button className="btn-logout" onClick={onLogout}>
          Logout
        </button>
      </div>
    </div>
  );
}

export default UserDashboard;