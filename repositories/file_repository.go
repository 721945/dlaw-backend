package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/gorm/clause"
)

type FileRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewFileRepository(logger *libs.Logger, db libs.Database) FileRepository {
	return FileRepository{logger: logger, db: db}
}

func (r *FileRepository) GetFiles() (files []models.File, err error) {
	return files, r.db.DB.Find(&files).Error
}

func (r *FileRepository) GetFile(id uuid.UUID) (file *models.File, err error) {
	return file, r.db.DB.First(&file, id).Error
}

func (r *FileRepository) GetFileByName(name string, folderId uuid.UUID) (file *models.File, err error) {
	return file, r.db.DB.Where("name = ? AND folder_id = ?", name, folderId).First(&file).Error
}

func (r *FileRepository) GetFileContent(id uuid.UUID) (file *models.File, err error) {
	return file, r.db.DB.Preload("Tags").Preload("FileType").First(&file, id).Error
}

func (r *FileRepository) GetFilesByFolderId(folderId uuid.UUID) (files []models.File, err error) {
	return files, r.db.DB.Preload(clause.Associations).Where("folder_id = ?", folderId).Find(&files).Error
}

func (r *FileRepository) GetFilesWithTagByFolderId(folderId uuid.UUID) (files []models.File, err error) {
	return files, r.db.DB.Preload("Tags").Where("folder_id = ?", folderId).Find(&files).Error
}

func (r *FileRepository) CreateFile(file models.File) (models.File, error) {
	return file, r.db.DB.Create(&file).Error
}

func (r *FileRepository) UpdateFile(id uuid.UUID, file models.File) error {
	return r.db.DB.Model(&models.File{}).Where("id = ?", id).Updates(file).Error
}

func (r *FileRepository) DeleteFile(id uuid.UUID) error {
	return r.db.DB.Delete(&models.File{}, id).Error
}

func (r *FileRepository) GetFilesByCaseIds(caseIds []uuid.UUID) (files []models.File, err error) {
	return files, r.db.DB.Where("case_id IN ?", caseIds).Find(&files).Error
}

func (r *FileRepository) CreateFileDocument(file models.MeiliFile) (resp *meilisearch.TaskInfo, err error) {
	return r.db.Meili.Index("files").AddDocuments([]models.MeiliFile{file}, "id")
}

func (r *FileRepository) UpdateFileDocument(file models.MeiliFile) (resp *meilisearch.TaskInfo, err error) {
	return r.db.Meili.Index("files").UpdateDocuments(file, "id")
}

func (r *FileRepository) DeleteFileDocument(id []string) (resp *meilisearch.TaskInfo, err error) {
	return r.db.Meili.Index("files").DeleteDocuments(id)
}

func (r *FileRepository) SearchFiles(query string) (resp *meilisearch.SearchResponse, err error) {
	return r.db.Meili.Index("files").Search(query, nil)
}

func (r *FileRepository) SearchFilesByCaseId(query string, caseId uuid.UUID) (resp *meilisearch.SearchResponse, err error) {
	return r.db.Meili.Index("files").Search(query, &meilisearch.SearchRequest{
		Filter: "case_id = " + caseId.String(),
	})
}

func (r *FileRepository) SearchFilesByFolderId(query string, folderId uuid.UUID) (resp *meilisearch.SearchResponse, err error) {
	return r.db.Meili.Index("files").Search(query, &meilisearch.SearchRequest{
		Filter: "folder_id = " + folderId.String(),
	})
}

// SearchFilesByTag : Search by tag name
func (r *FileRepository) SearchFilesByTag(query string, tag string) (resp *meilisearch.SearchResponse, err error) {
	return r.db.Meili.Index("files").Search(query, &meilisearch.SearchRequest{
		Filter: "tags = " + tag,
	})
}

func (r *FileRepository) SearchFilesByCaseIdAndTag(query string, caseId uuid.UUID, tag string) (resp *meilisearch.SearchResponse, err error) {
	return r.db.Meili.Index("files").Search(query, &meilisearch.SearchRequest{
		Filter: "case_id = " + caseId.String() + " AND tags = " + tag,
	})
}

func (r *FileRepository) SearchFilesByFolderIdAndTag(query string, folderId uuid.UUID, tag string) (resp *meilisearch.SearchResponse, err error) {
	return r.db.Meili.Index("files").Search(query, &meilisearch.SearchRequest{
		Filter: "folder_id = " + folderId.String() + " AND tags = " + tag,
	})
}

func (r *FileRepository) SearchFilesByCaseIdAndFolderId(query string, caseId uuid.UUID, folderId uuid.UUID) (resp *meilisearch.SearchResponse, err error) {
	return r.db.Meili.Index("files").Search(query, &meilisearch.SearchRequest{
		Filter: "case_id = " + caseId.String() + " AND folder_id = " + folderId.String(),
	})
}

func (r *FileRepository) SearchFilesByCaseIdAndFolderIdAndTag(query string, caseId uuid.UUID, folderId uuid.UUID, tag string) (resp *meilisearch.SearchResponse, err error) {
	return r.db.Meili.Index("files").Search(query, &meilisearch.SearchRequest{
		Filter: "case_id = " + caseId.String() + " AND folder_id = " + folderId.String() + " AND tags = " + tag,
	})
}
