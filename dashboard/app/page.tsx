"use client"

import { signIn, useSession } from "next-auth/react";

export default function Home() {
  const { data } = useSession()
  const [id, token] = [data?.id, data?.token]

  return (
    <div>
      <button onClick={() => signIn()}>
        <p>Click me to login</p>
      </button>
    </div>
  );
}
