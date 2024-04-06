package main

import (
	"context"
	"fmt"
	user "shorturl/kitex_gen/short/user"
	"shorturl/model"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}


// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	// TODO: Your code here...
	resp=user.NewUserResponse()
	username:=req.Username
	
	password:=req.Password

	u:=model.User{}
	result:=model.DB.Where("username=?",username).First(&u)
	if result.RecordNotFound(){
		resp.Message="该用户不存在"
		return 
	}else{
		if u.Password==password{
			resp.Message="登录成功"
		}else{
			resp.Message="密码错误"
		}
	}
	return resp,nil
}



// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	// TODO: Your code here...
	resp=user.NewUserResponse()
	username:=req.Username
	password:=req.Password

	u:=model.User{}
	result:=model.DB.Where("username=?",username).First(&u)
	if result.RowsAffected==1{
		resp.Message="该用户已存在"
	}else{
		u=model.User{Username: username,Password: password}
		if err:=model.DB.Create(&u).Error;err!=nil{
			return nil, fmt.Errorf("注册失败：%v", err)
		}else{
			resp.Message="注册成功"
		}
	}
	
	
	return resp,nil
}
