package util

import (
	"context"
	"database/sql"

	"github.com/Banana-Boat/terryminal/main-service/internal/db"
)

var terminalTemplateDict = []db.TerminalTemplate{
	{
		Name:        "Bash",
		Size:        "529MB",
		ImageName:   "tiangexiang/terryminal-base-pty:0.1.0",
		Description: sql.NullString{String: "基于Alpine Linux，仅可使用Bash", Valid: true},
	},
}

// 初始化终端模版字典表
func InitTermTemplates(store *db.Store) error {
	tmps, err := store.GetTerminalTemplates(context.Background())
	if err != nil {
		return err
	}

	// 如果表为空，则进行初始化
	if len(tmps) == 0 {
		for _, tmp := range terminalTemplateDict {
			if _, err = store.CreateTerminalTemplate(context.Background(),
				db.CreateTerminalTemplateParams{
					Name:        tmp.Name,
					Size:        tmp.Size,
					ImageName:   tmp.ImageName,
					Description: tmp.Description,
				},
			); err != nil {
				return err
			}
		}
	}

	return nil
}
