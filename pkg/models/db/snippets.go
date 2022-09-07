package db

import (
	"gorm.io/gorm"
	"snippet_box/pkg/models"
	"strconv"
	"time"
)

type SnippetModel struct {
	DB *gorm.DB
}

func (snippetModel *SnippetModel) Insert(title, content, expires string) (string, error) {

	expirationDate, err := strconv.Atoi(expires)
	if err != nil {
		return "", err
	}
	snippet := models.Snippets{
		Title:   title,
		Content: content,
		Expires: time.Now().AddDate(0, 0, expirationDate),
	}
	res := snippetModel.DB.Create(&snippet)

	if res.Error != nil {
		return "", res.Error
	}
	return snippet.Uuid, nil
}

func (snippetModel *SnippetModel) Get(uuid string) (*models.Snippets, error) {
	snippet := &models.Snippets{}
	result := snippetModel.DB.First(&snippet, "uuid = ?", uuid)

	if result.Error != nil {
		return nil, result.Error
	}
	return snippet, nil
}

func (snippetModel *SnippetModel) Latest() ([]*models.Snippets, error) {
	snippets := []*models.Snippets{}
	result := snippetModel.DB.Order("created DESC").Limit(10).Where("expires > UTC_TIMESTAMP()").Find(&snippets)
	if result.Error != nil {
		return nil, result.Error
	}
	return snippets, nil
}

func (snippetModel *SnippetModel) Delete(uuid string) {
	snippetModel.DB.Delete(models.Snippets{}, uuid)
}
