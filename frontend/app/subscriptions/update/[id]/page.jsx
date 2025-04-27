"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useParams } from "next/navigation";

export default function UpdateSubscription() {
  const { id } = useParams();
  const router = useRouter();

  const [targetUrl, setTargetUrl] = useState("");
  const [eventTypes, setEventTypes] = useState([]);

  useEffect(() => {
    const fetchSubscription = async () => {
      if (id) {
        try {
          const response = await fetch(
            `http://13.51.170.153:8080/subscriptions/${id}`
          );
          const result = await response.json();

          if (result && result.data) {
            const subscription = result.data;
            setTargetUrl(subscription.target_url);
            setEventTypes(subscription.event_types);
          }
        } catch (error) {
          console.error("Error fetching subscription:", error);
          alert("Failed to fetch subscription data");
        }
      }
    };

    fetchSubscription();
  }, [id]);

  const handleUpdate = async () => {
    if (!id) {
      alert("Subscription ID is missing!");
      return;
    }

    const cleanedEventTypes = eventTypes
      .filter((event) => event.trim() !== "")
      .map((event) => event.trim());

    const updatedData = {
      target_url: targetUrl,
      secret: "",
      event_types: cleanedEventTypes,
    };

    try {
      const response = await fetch(
        `http://13.51.170.153:8080/subscriptions/${id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(updatedData),
        }
      );

      const result = await response.json();
      console.log(result);

      if (result.status === "success") {
        alert("Subscription updated successfully!");
        router.push("/subscriptions");
      } else {
        alert("Failed to update subscription");
      }
    } catch (error) {
      console.error("Error updating subscription:", error);
      alert("Error updating subscription");
    }
  };

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Update Subscription</h1>

      <div className="mb-4">
        <label className="block mb-1 font-semibold">Target URL:</label>
        <input
          type="text"
          value={targetUrl}
          onChange={(e) => setTargetUrl(e.target.value)}
          className="border p-2 w-full"
          placeholder="Enter target URL"
        />
      </div>

      <div className="mb-4">
        <label className="block mb-1 font-semibold">
          Event Types (comma separated):
        </label>
        <input
          type="text"
          value={eventTypes.join(", ")}
          onChange={(e) =>
            setEventTypes(e.target.value.split(",").map((type) => type.trim()))
          }
          className="border p-2 w-full"
          placeholder="Ex: user.account.logout"
        />
      </div>

      <button
        onClick={handleUpdate}
        className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700"
      >
        Update Subscription
      </button>
    </div>
  );
}