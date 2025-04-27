"use client";


import { useEffect, useState } from 'react';
import axios from 'axios';
import { useRouter } from 'next/router';
import { useParams } from 'next/navigation';

const WebhookLogsPage = () => {
  const { webhookID } = useParams(); // Get the webhookID from URL parameters

  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Fetch logs when the webhookID changes
  useEffect(() => {
    if (webhookID) {
      axios
        .get(` http://13.51.170.153:8080/logs/webhook/${webhookID}`)
        .then((response) => {
          if (response.data.status === "success") {
            setLogs(response.data.data);
          } else {
            setError(response.data.message);
          }
          setLoading(false);
        })
        .catch((err) => {
          setError(err.response.data.message);
          setLoading(false);
        });
    }
  }, [webhookID]);

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
      <h1>Webhook Logs for Webhook ID: {webhookID}</h1>
      <table>
        <thead>
          <tr>
            <th>Webhook ID</th>
            <th>Subscription ID</th>
            <th>Target URL</th>
            <th>Attempt Number</th>
            <th>Status</th>
            <th>HTTP Status Code</th>
            <th>Error Details</th>
            <th>Timestamp</th>
          </tr>
        </thead>
        <tbody>
          {logs.map((log) => (
            <tr key={log.id}>
              <td>{log.webhook_id}</td>
              <td>{log.subscription_id}</td>
              <td>{log.target_url}</td>
              <td>{log.attempt_number}</td>
              <td>{log.status}</td>
              <td>{log.http_status_code}</td>
              <td>{log.error_details || 'N/A'}</td>
              <td>{log.timestamp}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default WebhookLogsPage;
