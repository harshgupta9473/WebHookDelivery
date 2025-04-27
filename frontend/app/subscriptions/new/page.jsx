"use client";

import { useState } from "react";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

const CreateSubscriptionPage = () => {
  const [targetUrl, setTargetUrl] = useState("");
  const [eventTypes, setEventTypes] = useState("");
  const [loading, setLoading] = useState(false);
  const handleTargetUrlChange = (e) => {
    setTargetUrl(e.target.value);
  };

  const handleEventTypesChange = (e) => {
    setEventTypes(e.target.value);
  };

  const handleCreateSubscription = async (e) => {
    e.preventDefault();
    setLoading(true);
    const eventTypesArray = eventTypes.split(",").map((item) => item.trim());

    try {
      const response = await fetch(" http://13.51.170.153:8080/subscriptions", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          target_url: targetUrl,
          secret: "",
          event_types: eventTypesArray,
        }),
      });

      const result = await response.json();

      if (response.ok && result.status === "success") {
        toast.success("Subscription added successfully!");
        setTargetUrl("");
        setEventTypes("");
      } else {
        toast.error("Failed to add subscription!");
      }
    } catch (error) {
      console.error(error);
      toast.error("Something went wrong!");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-4">
      <h1 className="text-2xl font-bold mb-6">Create New Subscription</h1>
      <form
        onSubmit={handleCreateSubscription}
        className="space-y-4 w-full max-w-md bg-white p-6 rounded-lg shadow-md"
      >
        <div>
          <label
            htmlFor="target_url"
            className="block text-sm font-medium text-gray-700"
          >
            Target URL
          </label>
          <input
            id="target_url"
            type="url"
            value={targetUrl}
            onChange={handleTargetUrlChange}
            required
            className="mt-1 p-2 w-full border border-gray-300 rounded-md"
            placeholder="Enter target URL"
          />
        </div>

        <div>
          <label
            htmlFor="event_types"
            className="block text-sm font-medium text-gray-700"
          >
            Event Types (comma separated)
          </label>
          <input
            id="event_types"
            type="text"
            value={eventTypes}
            onChange={handleEventTypesChange}
            required
            className="mt-1 p-2 w-full border border-gray-300 rounded-md"
            placeholder="Enter event types (comma separated)"
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full mt-4 px-6 py-3 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400"
        >
          {loading ? "Creating..." : "Create Subscription"}
        </button>
      </form>
    </div>
  );
};

export default CreateSubscriptionPage;