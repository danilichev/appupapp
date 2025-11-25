package schemas

import (
	"apps/api/internal/api"

	z "github.com/Oudwins/zog"
)

var folderName = z.String().
	Max(100, z.Message("Should be less than 100 characters"))

var CreateFolderRequestSchema = z.Struct(z.Shape{
	"name":     folderName.Required(z.Message("Folder name is required")),
	"parentId": z.String().Required(z.Message("Parent ID is required")),
})

var MoveFolderItemsRequestSchema = z.Struct(z.Shape{
	"items": z.Slice(z.Struct(z.Shape{
		"id": z.String().Required(z.Message("Item ID is required")),
		"type": z.StringLike[api.FolderItemType]().OneOf([]api.FolderItemType{"folder", "link"}, z.Message("Must be folder or link")).
			Required(z.Message("Item type is required")),
	})).Required(z.Message("Items are required")),
	"targetFolderId": z.String().Required(z.Message("Folder ID is required")),
})

var UpdateFolderRequestSchema = z.Struct(z.Shape{
	"name":     z.Ptr(folderName.Optional()),
	"parentId": z.Ptr(z.String().Optional()),
})
