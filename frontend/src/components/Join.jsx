import React, { useContext, useState, memo } from "react"

import { WebsocketContext } from "../App"

export default memo(function Join() {
    const [name, setName] = useState("")
    const [validName, setValidName] = useState(true)

    const { ws, joined } = useContext(WebsocketContext)

    const handleInput = (event) => {
        const value = event.target.value.trim()
        setValidName(!(value == ""))
        setName(value)
    }

    const handleJoin = () => {
        const trimmedName = name.trim()
        if (trimmedName == "") return
        ws.current.send(JSON.stringify({ event: "join", content: trimmedName }))
        setName("")
    }

    return (
        <div>
            {joined.err}
            <button onClick={handleJoin}>Join</button>
            <input
                required
                onChange={handleInput}
                value={name}
                className={!validName ? "unvalid" : ""}
            ></input>
        </div>
    )
})
