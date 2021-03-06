package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-services/comment-service/bean"
	"github.com/ninjadotorg/handshake-services/comment-service/request_obj"
	"github.com/ninjadotorg/handshake-services/comment-service/response_obj"
	"github.com/ninjadotorg/handshake-services/comment-service/service"
	"github.com/ninjadotorg/handshake-services/comment-service/utils"
)

var hookService = new(service.HookService)

type Api struct {
}

func (api Api) Init(router *gin.Engine) *gin.Engine {
	router.GET("/list", func(context *gin.Context) {
		api.GetComments(context)
	})
	router.POST("/", func(context *gin.Context) {
		api.CreateComment(context)
	})
	router.GET("/count", func(context *gin.Context) {
		api.GetCommentCount(context)
	})
	return router
}

func (api Api) CreateComment(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}

	request := new(request_obj.CommentRequest)

	requestJson := context.Request.PostFormValue("request")
	err := json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	sourceFile, sourceFileHeader, err := context.Request.FormFile("image")
	comment, appErr := commentService.CreateComment(userId.(int64), *request, &sourceFile, sourceFileHeader)
	if appErr != nil {
		log.Print(appErr.OrgError)
		result.SetStatus(bean.UnexpectedError)
		result.Error = appErr.OrgError.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	data := response_obj.MakeCommentResponse(comment)
	count, err := commentService.CountCommentByObjectId(comment.ObjectId)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	hookService.CommentCountHooks(comment.ObjectId, count)
	result.Data = data
	result.Status = 1
	result.Message = ""

	context.JSON(http.StatusOK, result)
	return
}

func (api Api) GetComments(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}

	pageSizeStr := context.Query("page_size")
	if len(pageSizeStr) == 0 {
		pageSizeStr = utils.DEFAULT_PAGE_SIZE
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	pageStr := context.Query("page")
	if len(pageStr) == 0 {
		pageStr = utils.DEFAULT_PAGE
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	objectId := context.Query("object_id")

	var pagination *bean.Pagination
	pagination = &bean.Pagination{PageSize: pageSize, Page: page}

	pagination, err = commentService.GetCommentPagination(0, objectId, pagination)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	data := response_obj.MakePaginationCommentResponse(pagination)

	result.Data = data
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (api Api) GetCommentCount(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}

	objectId := context.Query("object_id")

	count, err := commentService.CountCommentByObjectId(objectId)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Data = count
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}
