import { useEffect, useState } from "react"
import { Link } from 'react-router-dom'

function Dashboard() {
  const [domains, setDomains] = useState([])

  const fetchDomainData = () => {
    fetch("http://localhost:8080/list", {
      method: 'GET',
      headers: {
        Accept: 'application/json',
      },
    })
      .then(response => {
        return response.json()
      })
      .then(data => {
        setDomains(data)
      })
  }

  useEffect(() => {
    fetchDomainData()
  }, [])

  return (
    <div>
      Domains
      {domains.length >= 0 && (
        <ul>
          {domains.map(domain => (
            <li key={domain.ID}>{domain.domain} <Link to={`/domain/${domain.domain}`}>View</Link></li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default Dashboard;