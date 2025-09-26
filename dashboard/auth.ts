import NextAuth from "next-auth"
import { JWT } from "next-auth/jwt"
import Credentials from "next-auth/providers/credentials"

declare module "next-auth" {
  interface User {
    id: string 
    token: string
  }

  interface Session {
    id: string
    token: string
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    id: string
    token: string
  }
}

export const { handlers, signIn, signOut, auth } = NextAuth({
  session: {
    strategy: "jwt"
  },
  providers: [
    Credentials({
      credentials: {
        username: { label: "Username" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        const username = credentials.username
		    const password = credentials.password

        const res = await fetch("http://localhost:8080/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            username: username, 
            password: password
          })
        })

        const body = await res.json()
        if(res.ok) return body
        return null
      },
    }),
  ],
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.id = user.id
        token.token = user.token
      }
      return token
    },
    async session({ session, token }) {
      session.id = token.id
      session.token = token.token
      return session
    }
  }
})