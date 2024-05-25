import { idStruct, link, linkRequest, user, userRequest } from "./apiTypes";

export default class API
{
    static API_URL : string = "http://localhost:8080"; //this needs to be changed in production

    static async createUser(email : string, password : string) : Promise<user | number> {
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
        const result = await fetch(this.API_URL + "/accounts/" + id, {
            method: "GET", 
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
        });

        if (!result.ok) {
            console.log(result.body);
            return null
        }

        return await result.json() as user;
    }

    static async createLink(id : number, url_redirect : string) : Promise<link | number> {
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