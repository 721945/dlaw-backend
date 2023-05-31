package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type FileViewLogRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewFileViewLogRepository(logger *libs.Logger, db libs.Database) FileViewLogRepository {
	return FileViewLogRepository{logger: logger, db: db}
}

func (r *FileViewLogRepository) GetFileViewLogs() (logs []models.FileViewLog, err error) {
	return logs, r.db.DB.Find(&logs).Error
}

func (r *FileViewLogRepository) GetFileViewLog(id uuid.UUID) (permissionLog *models.FileViewLog, err error) {
	return permissionLog, r.db.DB.First(&permissionLog, id).Error
}

func (r *FileViewLogRepository) CreateFileViewLog(log models.FileViewLog) (models.FileViewLog, error) {
	var fileViewLog models.FileViewLog
	fileViewLog = log
	res := r.db.DB.Where("file_id = ? AND user_id = ?", log.FileId, log.UserId).FirstOrCreate(&fileViewLog)

	if res.Error != nil {
		return log, res.Error
	}

	return log, r.db.DB.Model(&models.FileViewLog{}).Where("id = ?", fileViewLog.ID).Updates(log).Error
}

func (r *FileViewLogRepository) UpdateFileViewLog(id uuid.UUID, log models.FileViewLog) error {
	return r.db.DB.Model(&models.FileViewLog{}).Where("id = ?", id).Updates(log).Error
}

func (r *FileViewLogRepository) DeleteFileViewLog(id uuid.UUID) error {
	return r.db.DB.Select(clause.Associations).Delete(&models.FileViewLog{}, id).Error
}

func (r *FileViewLogRepository) GetFileViewLogsForUser(userId uuid.UUID) (logs []models.FileViewLog, err error) {
	return logs, r.db.DB.Preload("File.FileType").Preload("File.Tags").Preload("File", "deleted_at IS NULL").Order("updated_at DESC").Distinct().Limit(10).Where("user_id = ?", userId).Find(&logs).Error
}
