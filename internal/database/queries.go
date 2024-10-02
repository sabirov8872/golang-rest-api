package database

import _ "embed"

var (

	//go:embed queries/get_all_users.sql
	getAllUsersQuery string

	//go:embed queries/get_user_by_id.sql
	getUserByIdQuery string

	//go:embed queries/create_user.sql
	createUserQuery string

	//go:embed queries/update_user.sql
	updateUserQuery string

	//go:embed queries/delete_user.sql
	deleteUserQuery string

	//go:embed queries/get_user_by_user.sql
	signInQuery string
)
