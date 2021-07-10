import React, {
    useRef,
    useEffect,
    useState,
    createContext,
    useCallback,
} from "react"

import Join from "./components/Join"
import Chat from "./components/Chat"

export const WebsocketContext = createContext()

export default function App() {
    const [joined, setJoined] = useState({ joined: false, err: null })
    const [connection, setConnection] = useState(false)
    const [messages, setMessages] = useState([])

    const setMessagesValue = useCallback((newValue) => {
        setMessages((prev) => [...prev, newValue])
    }, [])

    const joinedValue = useRef(null)
    const ws = useRef(null)
    joinedValue.current = joined

    useEffect(() => {
        let wsProtocol = "ws"

        if (location.protocol === "https:") {
            wsProtocol = "wss"
        }

        ws.current = new WebSocket(
            `${wsProtocol}://${window.location.hostname}:8000/ws`
        )

        ws.current.onopen = () => setConnection(true)
        ws.current.onclose = () => setConnection(false)
        ws.current.onmessage = (event) => {
            const data = JSON.parse(event.data)
            switch (data.event) {
                case "joined":
                    setJoined({ joined: true, err: null })
                    break
                case "alreadyExists":
                    if (joinedValue.current.joined) {
                        setJoined({ joined: true, err: data.content })
                    } else {
                        setJoined({ joined: false, err: data.content })
                    }
                    break
                case "message":
                    setMessagesValue(data.content)
                default:
                    break
            }
        }

        return () => {
            ws.current.close()
        }
    }, [])

    return (
        <WebsocketContext.Provider value={{ ws, joined }}>
            {connection && joined.joined ? (
                <Chat messages={messages} />
            ) : (
                <Join joined={joined} />
            )}
        </WebsocketContext.Provider>
    )
}
