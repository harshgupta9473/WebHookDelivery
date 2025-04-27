"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

export default function SubscriptionList() {
  const [subscriptions, setSubscriptions] = useState([]);

  useEffect(() => {
    // Define an async function inside the useEffect to handle async operations
    const fetchSubscriptions = async () => {
      try {
        const res = await fetch("http://13.51.170.153:8080/subscriptions");
        const data = await res.json();
        console.log("Fetched subscriptions:", data);
        setSubscriptions(data.data);
      } catch (err) {
        console.error("Error fetching subscriptions:", err);
      }
    };

    fetchSubscriptions(); // Call the async function
  }, []); // Empty dependency array means it runs only once after the initial render

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Subscriptions</h1>

      <Link href="/subscriptions/new" className="text-blue-500 mb-4 inline-block">
        Create New Subscription
      </Link>

      <table className="table-auto w-full border mt-4">
        <thead>
          <tr className="bg-gray-100">
            <th className="border px-4 py-2">ID</th>
            <th className="border px-4 py-2">Target URL</th>
            <th className="border px-4 py-2">Event Types</th>
            <th className="border px-4 py-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {subscriptions.length === 0 ? (
            <tr>
              <td colSpan="4" className="text-center">No subscriptions available.</td>
            </tr>
          ) : (
            subscriptions.map((sub) => (
              <tr key={sub.id}>
                <td className="border px-4 py-2">{sub.id}</td>
                <td className="border px-4 py-2">{sub.target_url}</td>
                <td className="border px-4 py-2">{sub.event_types.join(", ")}</td>
                <td className="border px-4 py-2">
                  <Link href={`/subscriptions/update/${sub.id}`} className="text-green-600 mr-4">
                    Update
                  </Link>
                  <Link href={`/subscriptions/delete/${sub.id}`} className="text-red-500 mr-4">
                    Delete
                  </Link>
                  <Link href={`/subscriptions/webhooks/${sub.id}`} className="text-purple-600 mr-4">
                    Webhooks
                  </Link>
                  <Link href={`/subscriptions/sublogs/${sub.id}`} className="text-blue-600">
                    View Logs
                  </Link>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
