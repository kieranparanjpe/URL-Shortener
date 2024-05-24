export interface idStruct {
	id : number 
	link_id : number 
}

export interface adminStruct {
    admin_password : string
}

export interface userRequest {
    email : string
    password : string
}

export interface user {
    id : number
    email : string
    hash_password : string
}

export interface linkRequest {
    url_redirect : string
    id : number
}

export interface link {
    id : number
    url_redirect : string
    user_id : number
    encoded_id : string
}