import { idStruct, link, linkRequest, user, userRequest } from "./apiTypes";

export default class API
{
    static API_URL : string = ""; 
    static FRONTEND_URL : string = ""; 
    static init() {
        if(this.API_URL !== "" && this.FRONTEND_URL !== "")
            return


        let apiURL : string | undefined = process.env.NEXT_PUBLIC_API_URL
        if (!apiURL)
            return

        this.API_URL = apiURL;

        console.log("API_URL=" + this.API_URL)

        let frontendURL : string | undefined = process.env.NEXT_PUBLIC_FRONTEND_URL
        if (!frontendURL)
            return

        this.FRONTEND_URL = frontendURL;

        console.log("FRONTEND_URL=" + this.FRONTEND_URL)

    }

    static async createUser(email : string, password : string) : Promise<user | number> {
        this.init();
        let body : userRequest = {email: email, password: password}
        const result = await fetch(this.API_URL + "/accounts", {
            method: "POST", 
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body : JSON.stringify(body)
        });

        if (!result.ok) {
            console.log(result.body);
            return result.status;
        }
        return await result.json() as user;
    }

    static async login(email : string, password : string) : Promise<user | number> {
        this.init();

        let body : userRequest = {email: email, password: password}

        const result = await fetch(this.API_URL + "/login", {
            method: "POST", 
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body : JSON.stringify(body)
        });

        if (!result.ok) {
            console.log(result.body);
            return result.status;
        }

        return await result.json() as user;
    }

    static async getUser(id : number) : Promise<user | null> {
        this.init();
        const result = await fetch(this.API_URL + "/accounts/" + id, {
            method: "GET", 
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
        });

        if (!result.ok) {
            return null
        }

        return await result.json() as user;
    }

    static async createLink(id : number, url_redirect : string) : Promise<link | number> {
        this.init();

        let body : linkRequest = {id: id, url_redirect: url_redirect}
        const result = await fetch(this.API_URL + "/links/" + id, {
            method: "POST", 
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body : JSON.stringify(body)
        });

        if (!result.ok) {
            console.log(result.body);
            return result.status;
        }
        return await result.json() as link;
    }

    static async deleteLink(user_id : number, link_id: number) : Promise<link | number> {
        this.init();

        let body : idStruct = {id: user_id, link_id: link_id}
        const result = await fetch(this.API_URL + "/links/" + user_id, {
            method: "DELETE", 
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body : JSON.stringify(body)
        });

        if (!result.ok) {
            console.log(result.body);
            return result.status;
        }
        return await result.json() as link;
    }


    static async getLinks(id : number) : Promise<link[] | null> {
        this.init();

        const result = await fetch(this.API_URL + "/links/" + id, {
            method: "GET", 
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
        });

        if (!result.ok) {
            console.log(result.body);
            return null
        }

        return await result.json() as link[];
    }

    static async logout() : Promise<boolean> {
        this.init();

        const result = await fetch(this.API_URL + "/logout", {
            method: "POST", 
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
        });

        if (!result.ok) {
            console.log(result.body);
            return false;
        }

        return true;
    }
}