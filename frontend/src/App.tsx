import {useEffect, useState} from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

type User = {
	Id: string
	AvatarUrl: string
	Email: string
	Name: string
}

function App() {
	const [count, setCount] = useState(0)
	const [msg, setMsg] = useState("")
	const [user, setUser] = useState<User | null>(null)

	async function fetchFromServer() {
		const resp = await fetch("http://localhost:6969/test")
		const data = await resp.json();
		setMsg(data.message)
	}

	async function oauthDiscord() {
		window.location.replace("http://localhost:6969/auth/discord")
	}

	async function oauthGoogle() {
		window.location.replace("http://localhost:6969/auth/google")
	}

	async function logout() {
		window.location.replace("http://localhost:6969/auth/logout")
	}

	async function fetchUser() {
		const resp = await fetch("http://localhost:6969/auth/me", {
			credentials: "include"
		})

		if (!resp.ok) {
			throw new Error("ahhhh")
		}

		const data: User = await resp.json()
		setUser(data)
	}

	useEffect(() => {
		fetchUser()
	}, [])

	return (
		<>
			{user ? (
				<>
					<h1>User Id: {user.Id}</h1>
					<img src={user?.AvatarUrl} />
					<button onClick={() => logout()}>Logout</button>
				</>
			) : <h1>Not Logged In</h1>}
			<div className="flex justify-between">
				<a href="https://vite.dev" target="_blank">
					<img src={viteLogo} className="logo" alt="Vite logo"/>
				</a>
				<a href="https://react.dev" target="_blank">
					<img src={reactLogo} className="logo react" alt="React logo"/>
				</a>
			</div>
			<h1>Vite + React</h1>
			<div className="card">
				<button onClick={() => setCount((count) => count + 1)}>
					count is {count}
				</button>
				<p>
					Edit <code>src/App.tsx</code> and save to test HMR
				</p>
				<button onClick={() => fetchFromServer()}>Fetch From Server</button>
				<p>
					Message: {msg}
				</p>
				<button onClick={() => oauthDiscord()}>
					Login With Discord
				</button>
				<button onClick={() => oauthGoogle()}>
					Login With Google
				</button>
			</div>
			<p className="read-the-docs">
				Click on the Vite and React logos to learn more
			</p>
		</>
	)
}

export default App
