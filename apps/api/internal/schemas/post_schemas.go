package schemas

import z "github.com/Oudwins/zog"

var postContent = z.String().
	Max(5000, z.Message("Should be less than 5000 characters"))

var postTitle = z.String().
	Max(100, z.Message("Should be less than 100 characters"))

var CreatePostRequestSchema = z.Struct(z.Shape{
	"content": postContent.Required(z.Message("Content is required")),
	"title":   postTitle.Required(z.Message("Title is required")),
})

var GetPostsParamsSchema = z.Struct(z.Shape{
	"limit": z.Ptr(
		z.Int().GTE(
			1, z.Message("Limit must be 1 or greater"),
		).LTE(100, z.Message("Limit must be less or equal 100")).Optional()),
	"offset": z.Ptr(
		z.Int().GTE(0, z.Message("Offset must be 0 or greater")).Optional(),
	),
})

var UpdatePostRequestSchema = z.Struct(z.Shape{
	"content": z.Ptr(postContent.Optional()),
	"title":   z.Ptr(postTitle.Optional()),
})
