package controllers

import (
	"errandify/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
}

// create new task
func (t *TaskController) CreateTask(ctx *gin.Context) {
	task := models.Task{}
	errBindJson := ctx.ShouldBindJSON(&task)
	if errBindJson != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errBindJson.Error()})
		return
	}

	errDB := t.DB.Create(&task).Error
	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, task)

}

// delete task
func (t *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	task := models.Task{}

	err := t.DB.First(&task, id).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	errDB := t.DB.Delete(&models.Task{}, id).Error
	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
	}

	if task.Attachment != "" {
		os.Remove("attachments/" + task.Attachment)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})

}

// submit task
func (t *TaskController) SubmitTask(ctx *gin.Context) {
	task := models.Task{}
	id := ctx.Param("id")
	submitDate := ctx.PostForm("submitDate")

	file, errFile := ctx.FormFile("attachment")
	err := t.DB.First(&task, id).Error

	//check error file
	if errFile != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errFile.Error()})
		return
	}

	//check error task
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	//remove old attachment (Tambahkan pengecekan jika nama file tidak kosong)
	attachment := task.Attachment
	if attachment != "" {
		fileInfo, _ := os.Stat("attachments/" + attachment)
		if fileInfo != nil {
			//old attachment found
			os.Remove("attachments/" + attachment)
		}
	}

	//create new attachment (save new attachment)
	newAttachmentName := file.Filename
	errSave := ctx.SaveUploadedFile(file, "attachments/"+newAttachmentName)
	if errSave != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file: " + errSave.Error()})
		return
	}

	//patch new task
	errDB := t.DB.Where("id=?", id).Updates(models.Task{
		Status:     "Review",
		SubmitDate: submitDate,
		Attachment: newAttachmentName, // use new file name
	}).Error

	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	//if no error/success
	ctx.JSON(http.StatusOK, gin.H{"message": "Submitted to review!"})
}

// reject task
func (t *TaskController) RejectTask(ctx *gin.Context) {
	task := models.Task{}
	id := ctx.Param("id")
	rejectedDate := ctx.PostForm("rejectedDate")
	reason := ctx.PostForm("reason")

	//check error task
	err := t.DB.First(&task, id).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	//patch new task
	errDB := t.DB.Where("id=?", id).Updates(models.Task{
		Status:       "Rejected",
		Reason:       reason,
		RejectedDate: rejectedDate, // use reason for rejected task
	}).Error

	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	//if no error/success
	ctx.JSON(http.StatusOK, gin.H{"message": "Rejected!"})
}
