import React, { useContext, useState, useEffect } from "react"

import { WebsocketContext } from "../App"
import Rename from "./Rename"

export default function Chat({ messages }) {
    const [message, setMessage] = useState("")

    const { ws } = useContext(WebsocketContext)

    const handleMessageSent = () => {
        ws.current.send(JSON.stringify({ event: "message", content: message }))
        setMessage("")
    }

    return (
        <div>
            {messages.map((message) => (
                <p>{message}</p>
            ))}
            <input
                type="text"
                placeholder="Type a message"
                onChange={(event) => setMessage(event.target.value)}
                value={message}
            ></input>
            <button onClick={handleMessageSent}>Send</button>
            <Rename />
        </div>
    )
}
