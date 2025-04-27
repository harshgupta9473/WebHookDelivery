"use client";


import { useEffect, useState } from 'react';
import axios from 'axios';

import { useParams } from 'next/navigation';

const LogsPage = () => {

  const { subscriptionID} =useParams() // Get the subscriptionId from URL parameters

  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Fetch logs when the subscriptionId changes
  useEffect(() => {
    if (subscriptionID) {
      axios
        .get(` http://13.51.170.153:8080/logs/subscription/${subscriptionID}`)
        .then((response) => {
          setLogs(response.data.data);
          setLoading(false);
        })
        .catch((error) => {
          setError(error.response.data.message);
          setLoading(false);
        });
    }
  }, [subscriptionID]);

  // Loading state
  if (loading) {
    return <div>Loading logs...</div>;
  }

  // Error state
  if (error) {
    return <div>{error}</div>;
  }

  // Logs table
  return (
    <div>
      <h1>Webhook Logs for Subscription {subscriptionID}</h1>
      <table>
        <thead>
          <tr>
            <th>Webhook ID</th>
            <th>Target URL</th>
            <th>Attempt Number</th>
            <th>Status</th>
            <th>HTTP Status Code</th>
            <th>Timestamp</th>
          </tr>
        </thead>
        <tbody>
          {logs.map((log) => (
            <tr key={log.id}>
              <td>{log.webhook_id}</td>
              <td>{log.target_url}</td>
              <td>{log.attempt_number}</td>
              <td>{log.status}</td>
              <td>{log.http_status_code}</td>
              <td>{log.timestamp}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default LogsPage;
