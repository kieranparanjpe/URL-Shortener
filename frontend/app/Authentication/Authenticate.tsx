'use client'

import { useState } from "react"
import { user } from "../util/apiTypes"
import Login from "./Login";
import SignUp from "./SignUp";

export default function Authenticate({setUserCallback} : {setUserCallback: (user: user) => void}) {

    let [loginMethod, setLoginMethod] : [boolean, any] = useState(true);

    const toggleMethod = () => {
        setLoginMethod(!loginMethod);
    }

    return (
        <div>
            {loginMethod ? <Login setUserCallback={setUserCallback} toggleMethod={toggleMethod}/> : <SignUp setUserCallback={setUserCallback} toggleMethod={toggleMethod}/>}
        </div>
    )
}

