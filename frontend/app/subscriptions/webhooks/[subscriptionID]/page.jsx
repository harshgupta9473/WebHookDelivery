"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation"; // Import useParams from next/navigation
import axios from 'axios'; // Import axios

export default function Webhooks() {
  const [webhooks, setWebhooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null); // Add error state for more detailed error handling
  const { subscriptionID } = useParams(); // Get the subscription ID using useParams

  useEffect(() => {
    // Check if the subscriptionID is available
    if (!subscriptionID) {
      setError("No subscription ID available.");
      setLoading(false); // Stop loading if no subscription ID is found
      return;
    }

    // Define the fetch function
    const fetchData = async () => {
      try {
        setLoading(true); // Ensure loading is true when starting to fetch data
        console.log("Fetching data for subscription ID:", subscriptionID);
        
        const response = await axios.get(` http://13.51.170.153:8080/subscriptions/webhooks/${subscriptionID}`);
        console.log("Fetched webhooks:", response.data);

        // Check if data is present and then update state
        if (response.data && response.data.data) {
          setWebhooks(response.data.data);
        } else {
          setError("No webhooks found for this subscription.");
        }
      } catch (err) {
        console.error("Error fetching data:", err);
        
        // Handle different types of errors based on the Axios error
        if (err.response) {
          setError(`Error: ${err.response.data.message || 'Unknown error occurred'}`);
        } else if (err.request) {
          setError("No response received from the server.");
        } else {
          setError(`Request error: ${err.message}`);
        }
      } finally {
        setLoading(false); // Ensure loading is false after data is fetched or error occurs
      }
    };

    // Call the fetchData function
    fetchData();
  }, [subscriptionID]); // Dependency on subscriptionID, re-fetch when it's available

  // Handle the loading, error, and data display
  if (loading) return <div>Loading...</div>; // Show loading message while fetching data

  // Handle error display
  if (error) return <div>Error: {error}</div>; // Display error message if an error occurred

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Webhooks for Subscription: {subscriptionID}</h1>

      {webhooks.length === 0 ? (
        <p>No webhooks found for this subscription.</p>
      ) : (
        <table className="table-auto w-full border mt-4">
          <thead>
            <tr className="bg-gray-100">
              <th className="border px-4 py-2">Webhook ID</th>
              <th className="border px-4 py-2">Event Type</th>
              <th className="border px-4 py-2">Status</th>
              <th className="border px-4 py-2">Delivered</th>
              <th className="border px-4 py-2">Created At</th>
              <th className="border px-4 py-2">Retries</th>
              <th className="border px-4 py-2">Logs</th> {/* New column for logs */}
            </tr>
          </thead>
          <tbody>
            {webhooks.map((webhook) => (
              <tr key={webhook.id}>
                <td className="border px-4 py-2">{webhook.id}</td>
                <td className="border px-4 py-2">{webhook.event_type}</td>
                <td className="border px-4 py-2">{webhook.status}</td>
                <td className="border px-4 py-2">
                  {webhook.delivered ? "Yes" : "No"}
                </td>
                <td className="border px-4 py-2">{webhook.created_at}</td>
                <td className="border px-4 py-2">{webhook.retries}</td>
                <td className="border px-4 py-2">
                  {/* Button to view logs */}
                  <a
                    href={`/subscriptions/webhooks/wlogs/${webhook.id}`} // Dynamic link to webhook logs
                    className="text-blue-500 hover:text-blue-700"
                  >
                    View Logs
                  </a>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
