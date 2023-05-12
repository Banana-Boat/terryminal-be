package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Banana-Boat/terryminal/main-service/internal/db"
	"github.com/Banana-Boat/terryminal/main-service/internal/pb"
	"github.com/Banana-Boat/terryminal/main-service/internal/util"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type registerRequest struct {
	Username string         `json:"username" binding:"required"`
	Password string         `json:"password" binding:"required"`
	Email    string         `json:"email" binding:"required"`
	Gender   db.UsersGender `json:"gender" binding:"required"`
	Age      int32          `json:"age" binding:"required"`
}

type userOfResponse struct {
	ID       int32          `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Gender   db.UsersGender `json:"gender"`
	Age      int32          `json:"age"`
}

func newUserOfResponse(user db.User) userOfResponse {
	return userOfResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Gender:   user.Gender,
		Age:      user.Age,
	}
}

func (server *Server) register(ctx *gin.Context) {
	var req registerRequest

	/* 通过gin的binding校验参数合法性 */
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, err.Error(), nil))
		return
	}

	/* 判断用户名是否存在 */
	isExistUser, _ := server.store.IsExistUser(ctx, req.Username)
	if isExistUser {
		ctx.JSON(http.StatusConflict, wrapResponse(false, "用户名已经存在", nil))
		return
	}

	/* 创建用户 */
	hashedPassword, err := util.HashPassword(req.Password) // 对密码加密
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, err.Error(), nil))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Gender:   req.Gender,
		Age:      req.Age,
	}
	res, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, err.Error(), nil))
		return
	}

	/* 查询新增用户 */
	id, _ := res.LastInsertId()
	user, err := server.store.GetUserById(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, err.Error(), nil))
		return
	}
	_user := newUserOfResponse(user)

	/* 返回结果 */
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{"user": _user}))
}

type listUserRequest struct {
	PageIdx  int32 `form:"pageIdx" binding:"min=0"`
	PageSize int32 `form:"pageSize" binding:"required,min=5,max=10"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, err.Error(), nil))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: req.PageIdx * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{"userList": users}))
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, err.Error(), nil))
		return
	}

	/* 判断用户名是否存在 */
	isExistUser, _ := server.store.IsExistUser(ctx, req.Username)
	if !isExistUser {
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "用户不存在", nil))
		return
	}

	/* 获取用户信息 */
	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, err.Error(), nil))
		return
	}

	/* 校验密码 */
	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, wrapResponse(false, err.Error(), nil))
		return
	}

	/* 颁发Token */
	token, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, err.Error(), nil))
		return
	}

	/* 通过gRPC调用mail-service发送登录邮件 */
	err = sendMailViaGRPC(
		fmt.Sprintf("%s:%s", server.config.TerminalServiceHost, server.config.TerminalServicePort),
		user.Email,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{
		"token": token,
		"user":  newUserOfResponse(user),
	}))
}

// 通过gRPC调用mail-service发送登录邮件
func sendMailViaGRPC(serverAddr string, destAddr string) error {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := pb.NewMailServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.SendMail(ctx, &pb.SendMailRequest{
		DestAddr: destAddr,
		Content:  "hello",
	})
	if err != nil {
		return err
	}
	return nil
}
