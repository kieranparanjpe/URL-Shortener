'use client'

import { use, useEffect, useState } from "react";
import { link, user } from "../util/apiTypes";
import API from "../util/API";
import LinkItem from "./LinkItem";

export default function Links({myUser, nullUser}: {myUser : user, nullUser : ()=>void}) {

    let [links, setLinks] : [link[], any] = useState([]);
    let [newURL, setNewURL] : [string, any] = useState('');

    useEffect(()=>{
        const getLinks = async () => {
            let ln : link[] | null = await API.getLinks(myUser.id);

            if (ln == null) {
                return;
            }

            setLinks(ln);
        }

        getLinks();
    }, [myUser, newURL])

    const handleSubmit = async (e : any) => {
        e.preventDefault();
        let ln : link | number = await API.createLink(myUser.id, newURL);
        if (typeof(ln) == 'number') {
            console.log("could not create link");
            return;
        }

        setLinks([...links, ln]);
    }

    const deleteCallback = async (delLink : link) => {
        let success : link | number = await API.deleteLink(myUser.id, delLink.id);

        if (typeof(success) == 'number') {
            return;
        }
        let l : link[] = links.splice(links.indexOf(delLink), 1); //this is very inefficient but also links array will never be large
        setLinks([...links]);
    }

    const signoutButton = async () => {
        let ok : boolean = await API.logout();
        if (ok)
            nullUser();
    }

    return (
        <div>
            <div className="flex align-middle justify-between ml-2 mr-10 mt-2">
                <h2 className=" text-left text-2xl font-semibold">Logged in as: {myUser.email}</h2>
                <button className=" w-40 text-white text-xl bg-red-600 rounded-lg" onClick={signoutButton}>Logout</button>
            </div>

            <div style={{height: '500px', overflowY: 'scroll'}} className=" bg-gray-50 ml-2 mr-2 mt-4 rounded-md shadow-md">
                {links.map((_value, _index)=>{
                    return <LinkItem key={_index} myLink={_value} deleteCallback={deleteCallback}/>
                })}
            </div>

            <form onSubmit={handleSubmit} className=" w-3/4 m-auto text-left mt-2">
                <label className=" w-full">
                    <span className=" text-2xl">New Long URL</span>
                    <input
                        required 
                        type="text"
                        onChange={(e) => setNewURL(e.target.value)}
                        value={newURL}
                        placeholder="https://google.com"
                        className=" border-2 rounded-sm border-black w-full text-xl mt-2"
                    />
                </label>
                <button className=" mt-2 mb-2 m-auto self-center btn-primary h-12 w-full text-white text-2xl bg-green-600 rounded-lg" >Add Link</button>
            </form>
        </div>
    )


}