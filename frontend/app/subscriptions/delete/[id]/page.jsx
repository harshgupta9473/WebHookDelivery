"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { toast } from "react-toastify";
import { useParams } from "next/navigation";

export default function DeleteSubscription() {
  const { id } = useParams();
  const router = useRouter();

  const [subscription, setSubscription] = useState(null);

  useEffect(() => {
    const fetchSubscription = async () => {
      if (id) {
        try {
          const response = await fetch(
            ` http://13.51.170.153:8080/subscriptions/${id}`
          );
          const result = await response.json();

          if (result && result.data) {
            setSubscription(result.data);
          }
        } catch (error) {
          console.error("Error fetching subscription:", error);
          alert("Failed to fetch subscription data");
        }
      }
    };

    fetchSubscription();
  }, [id]);

  const handleDelete = async () => {
    if (!id) {
      alert("Subscription ID is missing!");
      return;
    }

    try {
      const response = await fetch(
        `http://13.51.170.153:8080/subscriptions/${id}`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      const result = await response.json();
      console.log(result);

      if (result.status === "success") {
        toast.success("Subscription deleted successfully!");
        router.push("/subscriptions");
      } else {
        toast.error("Failed to delete subscription");
      }
    } catch (error) {
      console.error("Error deleting subscription:", error);
      toast.error("Error deleting subscription");
    }
  };

  if (!subscription) {
    return <div>Loading...</div>;
  }

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Delete Subscription</h1>
      <p>Are you sure you want to delete the subscription with ID: {id}?</p>
      <div className="mb-4">
        <p>
          <strong>Target URL:</strong> {subscription.target_url}
        </p>
        <p>
          <strong>Event Types:</strong> {subscription.event_types.join(", ")}
        </p>
      </div>

      <button
        onClick={handleDelete}
        className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
      >
        Delete Subscription
      </button>
    </div>
  );
}