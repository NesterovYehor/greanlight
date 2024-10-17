package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

// Include method to check if a permission is included in the Permissions slice
func (permissions Permissions) Include(code string) bool {
	for _, p := range permissions {
		if p == code {
			return true
		}
	}
	return false
}

// PermissionModel struct to interact with the database
type PermissionModel struct {
	DB *sql.DB
}

// GetAllForUser method to retrieve all permissions for a specific user
func (model *PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
        SELECT permissions.code
        FROM permissions
        INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
        WHERE users_permissions.user_id = $1
    `

	// Creating a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query the database
	rows, err := model.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	// Iterate through the rows
	for rows.Next() {
		var permission string
		// Scan the permission code into the variable
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		// Append the permission to the Permissions slice
		permissions = append(permissions, permission)
	}

	// Check for any error during row iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (model *PermissionModel) AddForUser(userID int64, codes ...string) error {
	query := `
        INSERT INTO users_permissions 
        SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := model.DB.ExecContext(ctx, query, userID, pq.Array(codes))

	return err
}
