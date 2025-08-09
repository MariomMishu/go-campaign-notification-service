package consts

const (
	RoleIdAdmin    = iota + 1
	RoleIdManager  = 2
	RoleIdCustomer = 3

	RoleAdmin    = "ADMIN"
	RoleManager  = "MANAGER"
	RoleCustomer = "CUSTOMER"

	DefaultPageSize = 10
	DefaultPage     = 1

	PermissionUserCreate = "user.create" // Permission to create a new user
	PermissionUserUpdate = "user.update" // Permission to update an existing user's information
	PermissionUserFetch  = "user.fetch"  // Permission to fetch a specific user's data
	PermissionUserList   = "user.list"   // Permission to list all users
	PermissionUserDelete = "user.delete" // Permission to delete a user

	PermissionCampaignCreate        = "campaign.create"         // Permission to create a new campaign
	PermissionCampaignUpdate        = "campaign.update"         // Permission to update an existing campaign's information
	PermissionCampaignFetch         = "campaign.fetch"          // Permission to fetch a specific campaign's data
	PermissionCampaignList          = "campaign.list"           // Permission to list all campaign
	PermissionCampaignDelete        = "campaign.delete"         // Permission to delete a campaign
	PermissionCampaignApproveReject = "campaign.approve_reject" // Permission to delete a campaign
)

var RoleMap = map[int]string{
	RoleIdAdmin:    RoleAdmin,
	RoleIdManager:  RoleManager,
	RoleIdCustomer: RoleCustomer,
}
