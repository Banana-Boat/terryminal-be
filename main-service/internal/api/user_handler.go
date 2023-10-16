package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Banana-Boat/terryminal/main-service/internal/db"
	"github.com/Banana-Boat/terryminal/main-service/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"

	"github.com/rs/zerolog/log"
)

/* 向前端返回去除敏感信息的用户信息 */
type userOfResp struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Nickname     string    `json:"nickname"`
	ChatbotToken int32     `json:"chatbotToken"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func newUserOfResp(user db.User) userOfResp {
	return userOfResp{
		ID:           user.ID,
		Email:        user.Email,
		Nickname:     user.Nickname,
		ChatbotToken: user.ChatbotToken,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

/* 密码的加密与验证 */
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func checkPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

/* 注册 */
type registerReq struct {
	Email    string `json:"email" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) handleRegister(ctx *gin.Context) {
	var req registerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Info().Err(err).Msg("invalid request body")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	/* 判断邮箱是否存在 */
	isExistUser, _ := server.store.IsUserExisted(ctx, req.Email)
	if isExistUser {
		log.Info().Msg("email already exists")
		ctx.JSON(http.StatusOK, wrapResponse(false, "邮箱已存在", nil))
		return
	}

	/* 创建用户 */
	hashedPassword, err := hashPassword(req.Password) // 对密码加密
	if err != nil {
		log.Error().Err(err).Msg("register failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "注册失败", nil))
		return
	}

	arg := db.CreateUserParams{
		Email:    req.Email,
		Nickname: req.Nickname,
		Password: hashedPassword,
	}
	_, err = server.store.CreateUser(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("register failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "注册失败", nil))
		return
	}

	/* 返回结果 */
	log.Info().Msg("register success")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{"isOk": true}))
}

/* 登录 */
type loginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) handleLogin(ctx *gin.Context) {
	var req loginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Info().Err(err).Msg("invalid request body")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	/* 判断邮箱是否存在 */
	isExistUser, _ := server.store.IsUserExisted(ctx, req.Email)
	if !isExistUser {
		log.Info().Msg("user not found")
		ctx.JSON(http.StatusOK, wrapResponse(false, "用户不存在", nil))
		return
	}

	/* 获取用户信息 */
	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("login failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "登录失败", nil))
		return
	}

	/* 校验密码 */
	err = checkPassword(req.Password, user.Password)
	if err != nil {
		log.Info().Err(err).Msg("password incorrect")
		ctx.JSON(http.StatusOK, wrapResponse(false, "密码错误", nil))
		return
	}

	/* 颁发Token */
	token, err := server.tokenMaker.CreateToken(user.ID, user.Email, server.config.AccessTokenDuration)
	if err != nil {
		log.Error().Err(err).Msg("login failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "登录失败", nil))
		return
	}

	log.Info().Msg("login success")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{
		"token": token,
		"user":  newUserOfResp(user),
	}))
}

/* 获取用户信息 */
func (server *Server) handleGetUserInfo(ctx *gin.Context) {
	tokenPayload := ctx.MustGet("token_payload").(*TokenPayload)
	user, err := server.store.GetUserById(ctx, tokenPayload.ID)
	if err != nil {
		log.Error().Err(err).Msg("get user info failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "获取用户信息失败", nil))
		return
	}

	log.Info().Msg("get user info success")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{"user": newUserOfResp(user)}))
}

/* 发送邮箱验证码 */
func (server *Server) handleSendCodeByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		log.Info().Msg("invalid params")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	/* 判断邮箱是否存在 */
	isExistUser, _ := server.store.IsUserExisted(ctx, email)
	if !isExistUser {
		log.Info().Msg("user not found")
		ctx.JSON(http.StatusOK, wrapResponse(false, "用户不存在", nil))
		return
	}

	/* 生成6位随机验证码 */
	var code string
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		randNum := rand.Intn(10)
		code += fmt.Sprintf("%d", randNum)
	}

	/* 更新数据库，验证码过期时间为5分钟 */
	arg := db.UpdateVerificationCodeParams{
		Email:            email,
		VerificationCode: sql.NullString{String: code, Valid: true},
		ExpiredAt:        sql.NullTime{Time: time.Now().Add(time.Minute * 5), Valid: true},
		UpdatedAt:        time.Now(),
	}
	err := server.store.UpdateVerificationCode(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("verification code generate failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "验证码生成失败", nil))
		return
	}

	/* 将发送任务入队 */
	payload := worker.PayloadSendMail{
		To:      email,
		Subject: "邮箱校验",
		Html:    fmt.Sprintf("您的验证码为：<b>%s</b>，5分钟内有效。", code),
	}
	_payload, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("json marshal failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "验证码生成失败", nil))
		return
	}
	server.taskDistributor.DistributeTask(
		ctx, worker.TaskSendMail,
		_payload,
		asynq.ProcessIn(time.Second),
	)

	ctx.JSON(http.StatusOK, wrapResponse(true, "", nil))
}

/* 修改用户密码 */
type updateUserPwdReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

func (server *Server) handleUpdateUserPwd(ctx *gin.Context) {
	var req updateUserPwdReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Info().Err(err).Msg("invalid request body")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("user not found")
		ctx.JSON(http.StatusOK, wrapResponse(false, "用户不存在", nil))
		return
	}

	/* 校验验证码 */
	if user.VerificationCode.String != req.Code {
		log.Info().Msg("code incorrect")
		ctx.JSON(http.StatusOK, wrapResponse(false, "验证码错误", nil))
		return
	}
	if user.ExpiredAt.Time.Before(time.Now()) {
		log.Info().Msg("code expired")
		ctx.JSON(http.StatusOK, wrapResponse(false, "验证码已过期", nil))
		return
	}

	/* 更新用户信息 */
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("update failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "修改失败", nil))
		return
	}
	user.Password = hashedPassword

	arg := db.UpdateUserInfoParams{
		ID:           user.ID,
		Nickname:     user.Nickname,
		Password:     user.Password,
		ChatbotToken: user.ChatbotToken,
		UpdatedAt:    time.Now(),
	}
	err = server.store.UpdateUserInfo(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("update failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "修改失败", nil))
		return
	}

	log.Info().Msg("update success")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{"isOk": true}))
}

/* 修改用户信息 */
type updateUserInfoReq struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func (server *Server) handleUpdateUserInfo(ctx *gin.Context) {
	tokenPayload := ctx.MustGet("token_payload").(*TokenPayload)

	var req updateUserInfoReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Info().Err(err).Msg("invalid request body")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	user, err := server.store.GetUserById(ctx, tokenPayload.ID)
	if err != nil {
		log.Error().Err(err).Msg("update failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "修改失败", nil))
		return
	}

	/* 更新用户信息 */
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			log.Error().Err(err).Msg("update failed")
			ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "修改失败", nil))
			return
		}
		user.Password = hashedPassword
	}

	arg := db.UpdateUserInfoParams{
		ID:           user.ID,
		Nickname:     user.Nickname,
		Password:     user.Password,
		ChatbotToken: user.ChatbotToken,
		UpdatedAt:    time.Now(),
	}
	err = server.store.UpdateUserInfo(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("update failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "修改失败", nil))
		return
	}

	/* 查询更新后的用户信息 */
	_user, err := server.store.GetUserById(ctx, tokenPayload.ID)
	if err != nil {
		log.Error().Err(err).Msg("update failed")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "修改失败", nil))
		return
	}

	log.Info().Msg("update success")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{
		"user": newUserOfResp(_user),
	}))
}

/* type listUserRequest struct {
	PageIdx  int32 `form:"pageIdx" binding:"min=0"`
	PageSize int32 `form:"pageSize" binding:"required,min=5,max=10"`
}

func (server *Server) listUsersHandle(ctx *gin.Context) {
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
} */
