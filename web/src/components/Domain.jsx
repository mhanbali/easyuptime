import { useEffect, useState } from "react"
import { useParams } from 'react-router-dom'

function Domain() {
    const [domains, setDomains] = useState([])
    const { domain } = useParams()
    const fetchDomainData = () => {
        fetch(`http://localhost:8080/view/${domain}`, {
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
            <h1>Status for {domain}</h1>
            <table>
                <thead>
                    <tr>
                        <th>Domain</th>
                        <th>Response Code</th>
                        <th>Response Time</th>
                        <th>Checked At</th>
                    </tr>
                </thead>
                <tbody>
                    {domains.map(domain => (
                        <tr key={domain.ID}>
                            <td>{domain.domain}</td>
                            <td>{domain.response}</td>
                            <td>{domain.duration} ms</td>
                            <td>{formatTimestamp(domain.CreatedAt)}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

function formatTimestamp(timestamp) {
    const date = new Date(timestamp);
    const dateSlice = date.toLocaleDateString();
    const timeSlice = date.toLocaleTimeString();
    const formattedTimestamp = `${dateSlice} at ${timeSlice}`;

    return formattedTimestamp;
}

export default Domain; 