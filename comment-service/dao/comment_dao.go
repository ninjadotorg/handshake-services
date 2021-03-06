package dao

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ninjadotorg/handshake-services/comment-service/bean"
	"github.com/ninjadotorg/handshake-services/comment-service/models"
)

type CommentDao struct {
}

func (commentDao CommentDao) GetById(id int64) models.Comment {
	dto := models.Comment{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (commentDao CommentDao) Create(dto models.Comment, tx *gorm.DB) (models.Comment, error) {
	if tx == nil {
		tx = models.Database()
	}
	dto.DateCreated = time.Now()
	dto.DateModified = dto.DateCreated
	err := tx.Create(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (commentDao CommentDao) Update(dto models.Comment, tx *gorm.DB) (models.Comment, error) {
	if tx == nil {
		tx = models.Database()
	}
	dto.DateModified = time.Now()
	err := tx.Save(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (commentDao CommentDao) Delete(dto models.Comment, tx *gorm.DB) (models.Comment, error) {
	if tx == nil {
		tx = models.Database()
	}
	err := tx.Delete(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (commentDao CommentDao) GetCommentPagination(userId int64, objectId string, pagination *bean.Pagination) (*bean.Pagination, error) {
	dtos := []models.Comment{}
	db := models.Database()
	if pagination != nil {
		db = db.Limit(pagination.PageSize)
		db = db.Offset(pagination.PageSize * (pagination.Page - 1))
	}
	if objectId != "" {
		db = db.Where("object_id = ?", objectId)
	}

	err := db.Order("date_created desc").Find(&dtos).Error
	if err != nil {
		log.Print(err)
		return pagination, err
	}
	pagination.Items = dtos
	total := 0
	if pagination.Page == 1 && len(dtos) < pagination.PageSize {
		total = len(dtos)
	} else {
		err := db.Find(&dtos).Count(&total).Error
		if err != nil {
			log.Print(err)
			return pagination, err
		}
	}
	pagination.Total = total
	return pagination, nil
}

func (commentDao CommentDao) CountByObjectId(objectId string) (int, error) {
	var count int
	db := models.Database()
	rows, err := db.Raw("SELECT count(1) FROM comment WHERE object_id = ?", objectId).Rows()
	if err != nil {
		log.Print(err)
		return count, err
	}
	for rows.Next() {
		rows.Scan(&count)
	}
	return count, nil
}
