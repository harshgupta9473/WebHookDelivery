"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

export default function SubscriptionList() {
  const [subscriptions, setSubscriptions] = useState([]);

  useEffect(() => {
    fetch(" http://13.51.170.153:8080/subscriptions")
      .then((res) => res.json())
      .then((data) => {
        console.log("Fetched subscriptions:", data);
        setSubscriptions(data.data);
      })
      .catch((err) => console.error(err));
  }, []);

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Subscriptions</h1>

      <Link
        href="/subscriptions/new"
        className="text-blue-500 mb-4 inline-block"
      >
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
          {subscriptions.map((sub) => (
            <tr key={sub.id}>
              <td className="border px-4 py-2">{sub.id}</td>
              <td className="border px-4 py-2">{sub.target_url}</td>
              <td className="border px-4 py-2">{sub.event_types.join(", ")}</td>
              <td className="border px-4 py-2">
                <Link
                  href={`/subscriptions/update/${sub.id}`}
                  className="text-green-600 mr-4"
                >
                  Update
                </Link>
                <Link
                  href={`/subscriptions/delete/${sub.id}`}
                  className="text-red-500 mr-4"
                >
                  Delete
                </Link>
                <Link
                  href={`/subscriptions/webhooks/${sub.id}`}
                  className="text-purple-600 mr-4"
                >
                  Webhooks
                </Link>
                {/* View Logs Button */}
                <Link
                  href={`/subscriptions/sublogs/${sub.id}`} // Link to subscription logs page
                  className="text-blue-600"
                >
                  View Logs
                </Link>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
