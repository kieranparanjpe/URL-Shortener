'use client'

import { useState } from "react"
import API from "../util/API"
import { user } from "../util/apiTypes"

export default function SignUp({setUserCallback, toggleMethod} : {setUserCallback: (user: user) => void, toggleMethod:() => void}) {

    const [email, setEmail] : [string, any] = useState('');
    const [password, setPassword] : [string, any] = useState('');
    const [passwordC, setPasswordC] : [string, any] = useState('');
    const [wrongPassword, setWrongPassword] : [boolean, any] = useState(false);
    const [errorSigningIn, setErrorSigningIn] : [boolean, any] = useState(false);
    const [emailTaken, setEmailTaken] : [boolean, any] = useState(false);

    const handleSubmit = async (e : any)  => {
        e.preventDefault();

        if (passwordC !== password) {
            setWrongPassword(true);
            setErrorSigningIn(false);
            return;
        }
    
        let myUser : user | number = await API.createUser(email, password);

        if(typeof(myUser) == "number"){
            if(myUser === 409){
                setEmailTaken(true);
                setErrorSigningIn(false);
                setWrongPassword(false);
            }
            else {
                setErrorSigningIn(true);
                setWrongPassword(false);
                setEmailTaken(false);
            }
            return;
        }

        setUserCallback(myUser);
    }

    return (
        <div className={"content-start text-center justify-center bg-white shadow-xl shadow-gray-700"}
        style={{height: "60vh", width: "25vw", marginLeft: "auto", marginRight: "auto", marginTop:"5vh", borderRadius: "20px"}}>
            <div className="flex items-center justify-items-center w-full h-16">
                <h2 className="text-black text-center font-bold text-4xl" style={{paddingInline: "1vw", flex: "1"}}>Sign Up</h2>
            </div>

            <div className="flex items-center justify-items-center w-full h-1 bg-black"></div>

            {wrongPassword && <p className="text-red-600 mt-2">passwords do not match</p>}
            {errorSigningIn && <p className="text-red-600 mt-2">internal error signing in</p>}
            {emailTaken && <p className="text-red-600 mt-2">email already taken</p>}

            <form onSubmit={handleSubmit} className=" w-3/4 m-auto text-left mt-8">
                <label className=" w-full">
                    <span className=" text-2xl">Email</span>
                    <br/>
                    <input
                        required 
                        type="text"
                        onChange={(e) => setEmail(e.target.value)}
                        value={email}
                        placeholder="hello@mail.com"
                        className=" border-2 rounded-sm border-black w-full text-xl"
                    />
                </label>
                <br/>
                <br/>
                <label className=" w-full">
                    <span className=" text-2xl">Password</span>
                    <br/>
                    <input
                        required 
                        type="password"
                        onChange={(e) => setPassword(e.target.value)}
                        value={password}
                        placeholder="enter password"
                        className=" border-2 rounded-sm border-black w-full text-xl"
                    />
                </label>
                <br/>
                <br/>
                <label className=" w-full">
                    <span className=" text-2xl">Confirm Password</span>
                    <br/>
                    <input
                        required 
                        type="password"
                        onChange={(e) => setPasswordC(e.target.value)}
                        value={passwordC}
                        placeholder="confirm password"
                        className=" border-2 rounded-sm border-black w-full text-xl"
                    />
                </label>
                <br/>
                <br/>
                <br/>
                <button className=" m-auto self-center btn-primary h-12 w-full text-white text-2xl bg-green-600 rounded-lg" >Sign Up</button>
            </form>
            <button onClick={toggleMethod} className=" m-auto mt-10 self-center h-10 w-3/4 text-white text-xl bg-blue-600 rounded-lg" >
                already have an account? log in
            </button>

        </div>
    )
}

