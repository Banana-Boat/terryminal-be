package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Banana-Boat/terryminal/main-service/internal/db"
	"github.com/Banana-Boat/terryminal/main-service/internal/pb"
	"github.com/Banana-Boat/terryminal/main-service/internal/pty"
	"github.com/Banana-Boat/terryminal/main-service/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type PtyHandler struct {
	gRPCConn   *grpc.ClientConn        // gRPC连接
	gRPCClient pb.BasePtyClient        // gRPC客户端
	gRPCStream pb.BasePty_RunCmdClient // gRPC数据流
}

// socket连接上下文
type WSContext struct {
	conn          net.Conn
	config        util.Config
	PtyHandlerMap map[string]*PtyHandler
}

/* Websocket连接 */
func (server *Server) handleTermWS(ctx *gin.Context) {
	conn, _, _, err := ws.UpgradeHTTP(ctx.Request, ctx.Writer)
	if err != nil {
		log.Error().Err(err).Msg("cannot upgrade http to websocket")
		return
	}
	log.Info().Msgf("new socket conn from %s", conn.RemoteAddr().String())
	defer conn.Close()

	wsCtx := &WSContext{
		conn:          conn,
		config:        server.config,
		PtyHandlerMap: make(map[string]*PtyHandler),
	}
	for {
		msg, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			if len(wsCtx.PtyHandlerMap) != 0 { // 客户端主动断开连接
				endAll(wsCtx)
			}
			log.Info().Msgf("socket conn closed from %s", conn.RemoteAddr().String())
			return
		}

		/* 解析 message */
		var wsMsg Message
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			log.Error().Err(err).Msg("cannot unmarshal message")
			return
		}

		/* 根据Event字段进行路由 */
		routeByEvent(wsCtx, wsMsg)
	}
}

func routeByEvent(wsCtx *WSContext, wsMsg Message) {
	switch wsMsg.Event {
	case "start":
		startEventHandle(wsCtx, wsMsg.PtyId, wsCtx.config)

	case "end":
		endEventHandle(wsCtx, wsMsg.PtyId)

	case "run-cmd":
		/* 将Data字段解析为对应结构体 */
		var data RunCmdClientData
		if err := mapstructure.Decode(wsMsg.Data, &data); err != nil {
			log.Error().Err(err).Msg("cannot decode data")
			return
		}

		runCmdEventHandle(wsCtx, wsMsg.PtyId, data.Cmd)
	}
}

/* 创建终端实例 */
type createTermReq struct {
	TemplateID int64  `json:"templateID" binding:"required"`
	Remark     string `json:"remark" binding:"omitempty"`
}

func (server *Server) handleCreateTerm(ctx *gin.Context) {
	tokenPayload := ctx.MustGet("token_payload").(*TokenPayload)

	/* 校验参数 */
	var req createTermReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Info().Err(err).Msg("invalid request body")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	/* 获得template */
	template, err := server.store.GetTerminalTemplateById(ctx, req.TemplateID)
	if err != nil {
		log.Error().Err(err).Msg("cannot get terminal template")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "模版不存在", nil))
		return
	}

	/* 创建pty */
	// 容器名格式: <用户ID>-<终端模版ID>-<时间戳>
	containerName := fmt.Sprintf("%d-%d-%d", tokenPayload.ID, template.ID, time.Now().Unix())
	ptyId, err := pty.NewPty(template.ImageName, containerName, server.config.PtyNetwork, nil)
	if err != nil {
		log.Error().Err(err).Msg("cannot create pty")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "创建失败", nil))
		return
	}

	/* 更新数据库 */
	args := db.CreateTerminalParams{
		ID:         ptyId,
		Name:       containerName,
		Size:       template.Size,
		Remark:     req.Remark,
		OwnerID:    tokenPayload.ID,
		TemplateID: req.TemplateID,
	}
	_, err = server.store.CreateTerminal(ctx, args)
	if err != nil {
		log.Error().Err(err).Msg("cannot create terminal")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "创建失败", nil))
		return
	}

	/* 查询新增终端实例 */
	term, err := server.store.GetTerminalById(ctx, ptyId)
	if err != nil {
		log.Error().Err(err).Msg("cannot get terminal")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "创建失败", nil))
		return
	}

	log.Info().Msg("create pty success")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{
		"terminal": term,
	}))
}

/* 销毁终端实例 */
func (server *Server) handleDestroyTerm(ctx *gin.Context) {
	/* 校验参数 */
	terminalID := ctx.Query("terminalID")
	if terminalID == "" {
		log.Info().Msg("invalid params")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	/* 销毁 */
	err := pty.RemovePty(terminalID)
	if err != nil {
		log.Error().Err(err).Msg("cannot destroy pty")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "销毁失败", nil))
		return
	}

	/* 更新数据库 */
	if err = server.store.DeleteTerminal(ctx, terminalID); err != nil {
		log.Error().Err(err).Msg("cannot delete terminal")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "销毁失败", nil))
		return
	}

	log.Info().Msg("destroy pty success")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{"isOk": true}))
}

/* 获取终端模版列表 */
func (server *Server) handleGetTermTemplates(ctx *gin.Context) {
	templates, err := server.store.GetTerminalTemplates(ctx)
	if err != nil {
		log.Error().Err(err).Msg("cannot get terminal templates")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "获取失败", nil))
		return
	}

	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{
		"templates": templates,
	}))
}

/* 获取某个用户的终端实例列表 */
func (server *Server) handleGetUserTerms(ctx *gin.Context) {
	tokenPayload := ctx.MustGet("token_payload").(*TokenPayload)

	terms, err := server.store.GetTerminalByOwnId(ctx, tokenPayload.ID)
	if err != nil {
		log.Error().Err(err).Msg("cannot get terminals")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "获取失败", nil))
		return
	}

	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{
		"terminals": terms,
	}))
}

/* 修改终端实例信息 */
type updateTermInfoReq struct {
	TerminalID string `json:"terminalID" binding:"required"`
	Remark     string `json:"remark" binding:"omitempty"`
}

func (server *Server) handleUpdateTermInfo(ctx *gin.Context) {
	/* 校验参数 */
	var req updateTermInfoReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Info().Err(err).Msg("invalid request body")
		ctx.JSON(http.StatusBadRequest, wrapResponse(false, "参数不合法", nil))
		return
	}

	/* 获取终端实例 */
	term, err := server.store.GetTerminalById(ctx, req.TerminalID)
	if err != nil {
		log.Error().Err(err).Msg("cannot get terminal")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "更新失败", nil))
		return
	}

	/* 更新数据库 */
	args := db.UpdateTerminalInfoParams{
		ID:            req.TerminalID,
		Size:          term.Size,
		Remark:        req.Remark,
		TotalDuration: term.TotalDuration,
		UpdatedAt:     time.Now(),
	}
	if err := server.store.UpdateTerminalInfo(ctx, args); err != nil {
		log.Error().Err(err).Msg("cannot update terminal")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "更新失败", nil))
		return
	}

	/* 查询新增终端实例 */
	term, err = server.store.GetTerminalById(ctx, req.TerminalID)
	if err != nil {
		log.Error().Err(err).Msg("cannot get terminal")
		ctx.JSON(http.StatusInternalServerError, wrapResponse(false, "更新失败", nil))
		return
	}

	log.Info().Msg("create pty success")
	ctx.JSON(http.StatusOK, wrapResponse(true, "", gin.H{
		"terminal": term,
	}))
}
