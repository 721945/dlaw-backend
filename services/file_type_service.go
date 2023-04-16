package services

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
)

type FileTypeService struct {
	logger             *libs.Logger
	fileTypeRepository repositories.FileTypeRepository
}

func NewFileTypeService(logger *libs.Logger, r repositories.FileTypeRepository) FileTypeService {
	return FileTypeService{logger: logger, fileTypeRepository: r}
}

func (s *FileTypeService) GetFileTypes() ([]dtos.FileTypeDto, error) {
	fileTypes, err := s.fileTypeRepository.GetFileTypes()

	if err != nil {
		return nil, err
	}

	types := dtos.ToFileTypes(fileTypes)

	return types, nil
}

func (s *FileTypeService) GetFileType(id uuid.UUID) (*dtos.FileTypeDto, error) {
	fileTypeModel, err := s.fileTypeRepository.GetFileType(id)

	if err != nil {
		return nil, err
	}

	fileType := dtos.ToFileType(fileTypeModel)

	return fileType, nil
}

func (s *FileTypeService) CreateFileType(dto dtos.CreateFileTypeDto) (string, error) {
	fileType, err := s.fileTypeRepository.CreateFileType(dto.ToFileType())

	if err != nil {
		return "", err
	}

	return fileType.ID.String(), nil
}

func (s *FileTypeService) UpdateFileType(id uuid.UUID, dto dtos.UpdateFileTypeDto) error {
	fileType := dto.ToFileType()
	return s.fileTypeRepository.UpdateFileType(id, fileType)
}

func (s *FileTypeService) DeleteFileType(id uuid.UUID) error {
	return s.fileTypeRepository.DeleteFileType(id)
}
