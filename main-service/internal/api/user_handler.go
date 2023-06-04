package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Banana-Boat/terryminal/main-service/internal/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/rs/zerolog/log"
)

/* 向前端返回去除敏感信息的用户信息 */
type userOfResponse struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Nickname     string    `json:"nickname"`
	ChatbotToken int32     `json:"chatbotToken"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func newUserOfResponse(user db.User) userOfResponse {
	return userOfResponse{
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
type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) registerHandle(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Info().Err(err).Msg("参数不合法")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	/* 判断邮箱是否存在 */
	isExistUser, _ := server.store.IsUserExisted(ctx, req.Email)
	if isExistUser {
		log.Info().Msg("邮箱已存在")
		ctx.JSON(http.StatusConflict, wrapResponse(false, "邮箱已存在", nil))
		return
	}

	/* 创建用户 */
	hashedPassword, err := hashPassword(req.Password) // 对密码加密
	if err != nil {
		log.Error().Err(err).Msg("注册失败")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "注册失败", nil))
		return
	}

	arg := db.CreateUserParams{
		Email:    req.Email,
		Nickname: req.Nickname,
		Password: hashedPassword,
	}
	res, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("注册失败")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "注册失败", nil))
		return
	}

	/* 查询新增用户 */
	id, _ := res.LastInsertId()
	user, err := server.store.GetUserById(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("注册失败")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "注册失败", nil))
		return
	}

	/* 返回结果 */
	log.Info().Msg("注册成功")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{"user": newUserOfResponse(user)}))
}

/* 登录 */
type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) loginHandle(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Info().Err(err).Msg("参数不合法")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	/* 判断邮箱是否存在 */
	isExistUser, _ := server.store.IsUserExisted(ctx, req.Email)
	if !isExistUser {
		log.Info().Msg("用户不存在")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "用户不存在", nil))
		return
	}

	/* 获取用户信息 */
	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("登录失败")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "登录失败", nil))
		return
	}

	/* 校验密码 */
	err = checkPassword(req.Password, user.Password)
	if err != nil {
		log.Info().Err(err).Msg("密码错误")
		ctx.JSON(http.StatusUnauthorized, wrapResponse(false, "密码错误", nil))
		return
	}

	/* 颁发Token */
	token, err := server.tokenMaker.CreateToken(user.ID, user.Email, server.config.AccessTokenDuration)
	if err != nil {
		log.Error().Err(err).Msg("登录失败")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "登录失败", nil))
		return
	}

	log.Info().Msg("登录成功")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{
		"token": token,
		"user":  newUserOfResponse(user),
	}))
}

/* 修改用户信息 */
type updateInfoRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func (server *Server) updateInfoHandle(ctx *gin.Context) {
	tokenPayload := ctx.MustGet("token_payload").(*TokenPayload)

	var req updateInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Info().Err(err).Msg("参数不合法")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	user, err := server.store.GetUserById(ctx, tokenPayload.ID)
	if err != nil {
		log.Error().Err(err).Msg("修改失败")
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
			log.Error().Err(err).Msg("修改失败")
			ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "修改失败", nil))
			return
		}
		user.Password = hashedPassword
	}

	arg := db.UpdateUserParams{
		ID:           user.ID,
		Nickname:     user.Nickname,
		Password:     user.Password,
		ChatbotToken: user.ChatbotToken,
		UpdatedAt:    time.Now(),
	}
	err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("修改失败")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "修改失败", nil))
		return
	}

	/* 查询更新后的用户信息 */
	_user, err := server.store.GetUserById(ctx, tokenPayload.ID)
	if err != nil {
		log.Error().Err(err).Msg("修改失败")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "修改失败", nil))
		return
	}

	log.Info().Msg("修改成功")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{
		"user": newUserOfResponse(_user),
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
