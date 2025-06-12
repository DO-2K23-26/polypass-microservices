// filepath: gateway/application/searchhttp/http.go
package searchhttp

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/DO-2K23-26/polypass-microservices/gateway/core"
	pb "github.com/DO-2K23-26/polypass-microservices/gateway/proto/search"
)

// Controller provides REST endpoints for SearchService
type Controller struct {
	search core.SearchService
}

// NewController creates a new search HTTP controller
func NewController(s core.SearchService) *Controller {
	return &Controller{search: s}
}

// Register mounts routes on the Fiber app
func (c *Controller) Register(app *fiber.App) {
	app.Get("/search/folders", c.handleFolders)
	app.Get("/search/tags", c.handleTags)
	app.Get("/search/credentials", c.handleCredentials)
}

// handleFolders GET /search/folders?search_query=&limit=&page=&user_id=
func (c *Controller) handleFolders(ctx *fiber.Ctx) error {
	q := ctx.Query("search_query")
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	req := &pb.SearchFoldersRequest{SearchQuery: q, Limit: int32(limit), Page: int32(page), UserId: ctx.Query("user_id")}
	folders, total, err := c.search.SearchFolders(ctx.Context(), req)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"folders": folders, "total": total})
}

func (c *Controller) handleTags(ctx *fiber.Ctx) error {
	q := ctx.Query("search_query")
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	req := &pb.SearchTagsRequest{SearchQuery: q, FolderId: ctx.Query("folder_id"), Limit: int32(limit), Page: int32(page), UserId: ctx.Query("user_id")}
	tags, total, err := c.search.SearchTags(ctx.Context(), req)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"tags": tags, "total": total})
}

func (c *Controller) handleCredentials(ctx *fiber.Ctx) error {
	q := ctx.Query("search_query")
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	// parse tag_ids as comma-separated
	tagIdsParam := ctx.Query("tag_ids")
	var tagIds []string
	if tagIdsParam != "" {
		tagIds = strings.Split(tagIdsParam, ",")
	}
	req := &pb.SearchCredentialsRequest{SearchQuery: q, FolderId: ctx.Query("folder_id"), TagIds: tagIds, Limit: int32(limit), Page: int32(page), UserId: ctx.Query("user_id")}
	creds, total, err := c.search.SearchCredentials(ctx.Context(), req)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"credentials": creds, "total": total})
}
