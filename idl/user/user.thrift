namespace go short.user

struct UserRequest{
    1:string Username
    2:string Password
}

struct UserResponse{
   1:string message
}

service UserService{

    UserResponse Login(1:UserRequest req)
    UserResponse Register(1:UserRequest req)
    
}

