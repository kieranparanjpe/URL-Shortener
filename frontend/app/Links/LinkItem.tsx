import API from "../util/API";
import { link } from "../util/apiTypes";

export default function LinkItem({key, myLink, deleteCallback} : {key : number, myLink : link, deleteCallback : (link : link) => void}) {
    return (
        <div className="flex items-center justify-between full h-10 bg-gray-300">
            <p className="text-left align-middle text-black" style={{paddingInline: "1vw", flex: "5"}}>{myLink.url_redirect}</p>
            <p className="text-left align-middle text-black" style={{paddingInline: "1vw", flex: "2"}}>{API.API_URL + "/l/" + myLink.encoded_id}</p>
           
           
            <button onClick={()=>{deleteCallback(myLink)}} className="text-left align-middle text-black"
                style={{paddingInline: "1vw", flex: "1"}}>Delete</button>
        </div>
    )
}