'use client'

import { useEffect, useState } from "react";
import { user } from "../util/apiTypes";
import API from "../util/API";
import Login from "./Authenticate";
import LinkItem from "../Links/LinkItem";
import Authenticate from "./Authenticate";
import Links from "../Links/Links";

export default function AuthState() {
    let [myUser, setMyUser] : [user | null, any] = useState(null); 

    const checkSignedIn = async () => {
      let idString : string | null = localStorage.getItem("current_user_id");
      if (idString == null)
      {
        setMyUser(null);
        return;
      }
  
      setMyUser(await API.getUser(parseInt(idString, 10)));
    }

    const updateUser = (newUser : user) => {
      localStorage.setItem("current_user_id", newUser.id.toString());
      setMyUser(newUser);
    }

    const nullUser = () => {
      setMyUser(null);
    }
  
    useEffect(()=> {
      if (myUser == null)
        checkSignedIn();
    })

    return (
        <div>
            {myUser == null ? <Authenticate setUserCallback={updateUser}/> : <Links myUser={myUser} nullUser={nullUser}/>}
        </div>
    )
}