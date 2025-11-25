package schemas

import z "github.com/Oudwins/zog"

var linkDescription = z.String().
	Max(500, z.Message("Should be less than 500 characters"))

var linkName = z.String().
	Max(100, z.Message("Should be less than 100 characters"))

var CreateLinkRequestSchema = z.Struct(z.Shape{
	"description": z.Ptr(linkDescription.Optional()),
	"folderId":    z.String().Required(z.Message("Folder ID is required")),
	"name":        z.Ptr(linkName.Optional()),
	"tags":        z.Ptr(z.Slice(z.String()).Optional()),
	"url":         z.String().URL().Required(z.Message("URL is required")),
})

var GetLinksParamsSchema = z.Struct(z.Shape{
	"folderId": z.Ptr(z.String().Optional()),
	"limit": z.Ptr(
		z.Int().GTE(
			1, z.Message("Limit must be 1 or greater"),
		).LTE(100, z.Message("Limit must be less or equal 100")).Optional()),
	"offset": z.Ptr(
		z.Int().GTE(0, z.Message("Offset must be 0 or greater")).Optional(),
	),
})

var ParseHtmlRequestSchema = z.Struct(z.Shape{
	"url": z.String().URL().Required(z.Message("URL is required")),
})

var UpdateLinkRequestSchema = z.Struct(z.Shape{
	"description": z.Ptr(linkDescription.Optional()),
	"folderId":    z.Ptr(z.String().Optional()),
	"name":        z.Ptr(linkName.Optional()),
	"tags":        z.Ptr(z.Slice(z.String()).Optional()),
	"url":         z.Ptr(z.String().URL().Optional()),
})
