package data

import (
	"context"
	"database/sql"
	"time"
)

type Permissions []string

var (
	GET_ALL_FOR_USERS_QUERY = `
	SELECT permissions.code
	FROM permissions
	INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
	INNER JOIN users ON users_permissions.user_id = users.id
	WHERE users.id = $1
	`
)

func (permissions Permissions) Include(code string) bool {
	for i := range permissions {
		if permissions[i] == code {
			return true
		}
	}

	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (permissionModel PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := permissionModel.DB.QueryContext(ctx, GET_ALL_FOR_USERS_QUERY, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
