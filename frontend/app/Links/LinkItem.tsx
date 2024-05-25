import API from "../util/API";
import { link } from "../util/apiTypes";

export default function LinkItem({key, myLink, deleteCallback} : {key : number, myLink : link, deleteCallback : (link : link) => void}) {
    return (
        <div className="flex items-center justify-between full h-20 bg-gray-300 rounded-md m-2 border-4 border-black">
            <a href={API.API_URL + "/l/" + myLink.encoded_id} className="text-left align-middle text-black" style={{paddingInline: "1vw", flex: "2"}}>
                {API.API_URL + "/l/" + myLink.encoded_id + " ↗"}
            </a>

            <a href={myLink.url_redirect} className="text-left align-middle text-black" 
            style={{paddingInline: "1vw", maxHeight: "3rem", maxWidth: "62.5%", flex: "5", textOverflow: "ellipsis", overflow: "hidden",
                overflowWrap: "anywhere"
            }}>
                {myLink.url_redirect + " ↗"}
            </a>
           
            <button onClick={()=>{deleteCallback(myLink)}} className="text-left align-middle text-black font-bold"
                style={{paddingInline: "1vw", flex: "1"}}>Delete</button>
        </div>
    )
}