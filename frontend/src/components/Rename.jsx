import React, { useState, useContext } from "react"

import { WebsocketContext } from "../App"

export default function Rename() {
    const [name, setName] = useState("")

    const { ws } = useContext(WebsocketContext)

    const changeNameHandler = () => {
        const trimmedName = name.trim()
        if (trimmedName == "") return
        ws.current.send(JSON.stringify({ event: "join", content: trimmedName }))
        setName("")
    }

    return (
        <div>
            <input
                type="text"
                onChange={(event) => setName(event.target.value)}
                value={name}
            />
            <button onClick={changeNameHandler}>Change Name</button>
        </div>
    )
}
