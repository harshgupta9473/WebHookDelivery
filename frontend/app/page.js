import Link from "next/link";

export default function Home() {
  return (
    <main className="flex flex-col items-center justify-center h-screen">
      <h1 className="text-4xl font-bold mb-4">
        Welcome to the Subscriptions App
      </h1>
      <Link
        href="/subscriptions"
        className="px-6 py-3 bg-blue-500 text-white rounded-md"
      >
        Go to Subscriptions
      </Link>
    </main>
  );
}