'use client'

import { use, useEffect, useState } from "react";
import { link, user } from "../util/apiTypes";
import API from "../util/API";
import LinkItem from "./LinkItem";

export default function Links({myUser}: {myUser : user}) {

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

    return (
        <div>
            <div style={{height: '500px', overflowY: 'scroll'}}>
                {links.map((_value, _index)=>{
                    return <LinkItem key={_index} myLink={_value} deleteCallback={deleteCallback}/>
                })}
            </div>

            <form onSubmit={handleSubmit} className=" w-3/4 m-auto text-left mt-8">
                <label className=" w-full">
                    <span className=" text-2xl">New Long URL</span>
                    <input
                        required 
                        type="text"
                        onChange={(e) => setNewURL(e.target.value)}
                        value={newURL}
                        placeholder="https://google.com"
                        className=" border-2 rounded-sm border-black w-full text-xl"
                    />
                </label>
                <button className=" m-auto self-center btn-primary h-12 w-full text-white text-2xl bg-green-600 rounded-lg" >Add Link</button>
            </form>
        </div>
    )


}