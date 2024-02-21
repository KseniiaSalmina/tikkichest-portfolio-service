package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uptrace/bunrouter"
	"net/http"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/api/response_errors"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/api/validation"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/models"
)

// @Summary Get portfolios
// @Tags portfolios
// @Description get portfolios (all, by user id or by category id)
// @Produce json
// @Param page query int false "page number"
// @Param limit query int false "limit records by page"
// @Param id query int false "user or category id"
// @Param filter query string false "filtered by" Enums(ByProfileID, ByCategoryID)
// @Success 200 {object} models.PortfoliosPage
// @Success 204
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios [get]
func (s *Server) getPortfoliosHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := validation.ID(idStr)
	if err != nil && !errors.Is(err, response_errors.ErrMissingID) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filterType := r.FormValue("filter")
	filter, err := validation.PortfoliosFilter(filterType, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := s.getPageInfo(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect page info: %s", err.Error()), http.StatusBadRequest)
		return
	}

	portfolios, pagesAmount, err := s.databaseConnector.GetAllPortfolios(r.Context(), page.limit, page.offset, *filter)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, true)
		return
	}

	response := models.PortfoliosPage{Portfolios: portfolios, PageNo: page.number, Limit: page.limit, PagesAmount: pagesAmount}

	json.NewEncoder(w).Encode(response)
}

// @Summary Get portfolio
// @Tags portfolios
// @Description get portfolio by its id
// @Produce json
// @Param id path int true "portfolio id"
// @Success 200 {object} models.Portfolio
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id} [get]
func (s *Server) getPortfolioByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("id")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	portfolio, err := s.databaseConnector.GetPortfolioByID(r.Context(), id)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	json.NewEncoder(w).Encode(portfolio)
}

// @Summary Post portfolio
// @Tags portfolios
// @Description create new portfolio, return its id
// @Accept json
// @Produce json
// @Param portfolio body models.Portfolio true "portfolio without crafts, profile id is required"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios [post]
func (s *Server) postPortfolioHandler(w http.ResponseWriter, r *http.Request) {
	var portfolio models.Portfolio
	if err := json.NewDecoder(r.Body).Decode(&portfolio); err != nil {
		http.Error(w, fmt.Sprintf("incorrect portfolio data: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if portfolio.ProfileID <= 0 {
		http.Error(w, response_errors.ErrIncorrectID.Error(), http.StatusBadRequest)
		return
	}

	portfolioID, err := s.databaseConnector.CreatePortfolio(r.Context(), portfolio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(portfolioID)
}

// @Summary Patch portfolio
// @Tags portfolios
// @Description update portfolio by its id
// @Accept json
// @Param id path int true "portfolio id"
// @Param portfolio body models.Portfolio true "updated portfolio, info without changes is also required"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id} [patch]
func (s *Server) patchPortfolioHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("id")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var portfolio models.Portfolio
	if err = json.NewDecoder(r.Body).Decode(&portfolio); err != nil {
		http.Error(w, fmt.Sprintf("incorrect portfolio data: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	portfolio.ID = id
	if err = s.databaseConnector.PatchPortfolio(r.Context(), portfolio); err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

// @Summary Delete portfolio
// @Tags portfolios
// @Description delete portfolio by its id
// @Param page query int false "page number"
// @Param id path int true "portfolio id"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id} [delete]
func (s *Server) deletePortfolioHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("id")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = s.databaseConnector.DeletePortfolio(r.Context(), id); err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Post category
// @Tags categories
// @Description create new category, return its id
// @Accept json
// @Produce json
// @Param category body models.Category true "category, name required"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /categories [post]
func (s *Server) postCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, fmt.Sprintf("incorrect category data: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if category.Name == "" {
		http.Error(w, "category name is required", http.StatusBadRequest)
		return
	}

	id, err := s.databaseConnector.CreateCategory(r.Context(), category.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(id)
}

// @Summary Delete category
// @Tags categories
// @Description delete category by its id
// @Param id path int true "category id"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /categories/{id} [delete]
func (s *Server) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("id")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = s.databaseConnector.DeleteCategory(r.Context(), id); err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get categories
// @Tags categories
// @Description get all categories
// @Produce json
// @Param page query int false "page number"
// @Param limit query int false "limit records by page"
// @Success 200 {object} models.CategoriesPage
// @Success 204
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /categories [get]
func (s *Server) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	page, err := s.getPageInfo(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect page info: %s", err.Error()), http.StatusBadRequest)
		return
	}

	categories, pagesAmount, err := s.databaseConnector.GetAllCategories(r.Context(), page.limit, page.offset)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, true)
		return
	}

	response := models.CategoriesPage{Categories: categories, PageNo: page.number, Limit: page.limit, PagesAmount: pagesAmount}

	json.NewEncoder(w).Encode(response)
}

// @Summary Get crafts by portfolio id
// @Tags crafts
// @Description get all crafts by portfolio id
// @Produce json
// @Param page query int false "page number"
// @Param limit query int false "limit records by page"
// @Param id path int true "portfolio id"
// @Success 200 {object} models.CraftsPage
// @Success 204
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id}/crafts [get]
func (s *Server) getCraftsByPortfolioIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("id")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := s.getPageInfo(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect page info: %s", err.Error()), http.StatusBadRequest)
		return
	}

	crafts, pagesAmount, err := s.databaseConnector.GetAllCraftsByPortfolioID(r.Context(), id, page.limit, page.offset)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, true)
		return
	}

	response := models.CraftsPage{Crafts: crafts, PageNo: page.number, Limit: page.limit, PagesAmount: pagesAmount}

	json.NewEncoder(w).Encode(response)
}

// @Summary Get craft
// @Tags crafts
// @Description get craft by its id
// @Produce json
// @Param craftID query int true "craft id"
// @Success 200 {object} models.Craft
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router //portfolios/{id}/crafts/{craftID} [get]
func (s *Server) getCraftHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("craftID")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	craft, err := s.databaseConnector.GetCraftByID(r.Context(), id)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	json.NewEncoder(w).Encode(craft)
}

// @Summary Post craft
// @Tags crafts
// @Description create new craft, return its id
// @Accept json
// @Produce json
// @Param id query int true "portfolio id"
// @Param craft body models.Craft true "craft without contents"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id}/crafts [post]
func (s *Server) postCraftHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("id")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var craft models.Craft
	if err = json.NewDecoder(r.Body).Decode(&craft); err != nil {
		http.Error(w, fmt.Sprintf("incorrect craft data: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	craftID, err := s.databaseConnector.CreateCraft(r.Context(), id, craft)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(craftID)
}

// @Summary Post tag patch craft
// @Tags crafts
// @Description add tag to the craft
// @Param craftID query int true "craft id"
// @Param tagID query int true "tag id"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id}/crafts/{craftID}/tag/{tagID} [post]
func (s *Server) postTagPatchCraftHandler(w http.ResponseWriter, r *http.Request) {
	craftIdStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("craftID")
	craftID, err := validation.ID(craftIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tagIdStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("tagID")
	tagID, err := validation.ID(tagIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = s.databaseConnector.AddTagToCraft(r.Context(), craftID, tagID); err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete tag patch craft
// @Tags crafts
// @Description delete tag from the craft
// @Param craftID query int true "craft id"
// @Param tagID query int true "tag id"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id}/crafts/{craftID}/tag/{tagID} [delete]
func (s *Server) deleteTagPatchCraftHandler(w http.ResponseWriter, r *http.Request) {
	craftIdStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("craftID")
	craftID, err := validation.ID(craftIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tagIdStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("tagID")
	tagID, err := validation.ID(tagIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = s.databaseConnector.DeleteTagFromCraft(r.Context(), craftID, tagID); err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Patch craft
// @Tags crafts
// @Description update craft by its id
// @Accept json
// @Param craftID query int true "craft id"
// @Param craft query models.Craft true "updated craft, info without changes is also required"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id}/crafts/{craftID} [patch]
func (s *Server) patchCraftHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("craftID")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var craft models.Craft
	if err = json.NewDecoder(r.Body).Decode(&craft); err != nil {
		http.Error(w, fmt.Sprintf("incorrect craft data: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	craft.ID = id

	if err = s.databaseConnector.PatchCraft(r.Context(), craft); err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete craft
// @Tags crafts
// @Description delete craft by its id
// @Param craftID query int true "craft id"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id}/crafts/{craftID} [delete]
func (s *Server) deleteCraftHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("craftID")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = s.databaseConnector.DeleteCraft(r.Context(), id); err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get crafts by tag
// @Tags crafts
// @Description get all crafts by tag id
// @Produce json
// @Param id query int true "tag id"
// @Success 200 {object} models.CraftsPage
// @Success 204
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /tags/{id}/crafts [get]
func (s *Server) getCraftsByTagIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("id")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := s.getPageInfo(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect page info: %s", err.Error()), http.StatusBadRequest)
		return
	}

	crafts, pagesAmount, err := s.databaseConnector.GetAllCraftsByTagID(r.Context(), id, page.limit, page.offset)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, true)
		return
	}

	response := models.CraftsPage{Crafts: crafts, PageNo: page.number, Limit: page.limit, PagesAmount: pagesAmount}

	json.NewEncoder(w).Encode(response)
}

// @Summary Get tags
// @Tags tags
// @Description get all tags
// @Produce json
// @Param page query int false "page number"
// @Param limit query int false "limit records by page"
// @Success 200 {object} models.TagsPage
// @Success 204
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /tags [get]
func (s *Server) getTagsHandler(w http.ResponseWriter, r *http.Request) {
	page, err := s.getPageInfo(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect page info: %s", err.Error()), http.StatusBadRequest)
		return
	}

	tags, pagesAmount, err := s.databaseConnector.GetAllTags(r.Context(), page.limit, page.offset)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, true)
		return
	}

	response := models.TagsPage{Tags: tags, PageNo: page.number, Limit: page.limit, PagesAmount: pagesAmount}

	json.NewEncoder(w).Encode(response)
}

// @Summary Post tag
// @Tags tags
// @Description create new tag, return its id
// @Accept json
// @Produce json
// @Param tag body models.Tag true "tag name required"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /tags [post]
func (s *Server) postTagHandler(w http.ResponseWriter, r *http.Request) {
	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		http.Error(w, fmt.Sprintf("incorrect tag data: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	id, err := s.databaseConnector.CreateTag(r.Context(), tag.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(id)
}

// @Summary Delete tag
// @Tags tags
// @Description delete tag by its id
// @Param id query int true "tag id"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /tags/{id} [delete]
func (s *Server) deleteTagHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("id")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = s.databaseConnector.DeleteTag(r.Context(), id); err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Post content
// @Tags contents
// @Description create new content, return its id
// @Accept json
// @Produce json
// @Param craftID query int true "craft id"
// @Param content body models.Content true "content"
// @Success 200 {object} models.PortfoliosPage
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id}/crafts/{craftID}/contents [post]
func (s *Server) postContentHandler(w http.ResponseWriter, r *http.Request) {
	craftIdStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("craftID")
	craftId, err := validation.ID(craftIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var content models.Content
	if err = json.NewDecoder(r.Body).Decode(&content); err != nil {
		http.Error(w, fmt.Sprintf("incorrect content data: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(content.Data) == 0 {
		http.Error(w, "content data is required", http.StatusBadRequest)
		return
	}

	id, err := s.databaseConnector.CreateContent(r.Context(), craftId, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(id)
}

// @Summary Delete content
// @Tags contents
// @Description delete content by its id
// @Param contentID query int true "content id"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id}/crafts/{craftID}/contents/{contentID} [delete]
func (s *Server) deleteContentHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("contentID")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = s.databaseConnector.DeleteContent(r.Context(), id); err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Patch content
// @Tags contents
// @Description update content by its id
// @Accept json
// @Param contentID query int true "content id"
// @Param content body models.Content true "updated content, info without changes is also required"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /portfolios/{id}/crafts/{craftID}/contents/{contentID} [patch]
func (s *Server) patchContentHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := bunrouter.ParamsFromContext(r.Context()).Get("contentID")
	id, err := validation.ID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var content models.Content
	if err = json.NewDecoder(r.Body).Decode(&content); err != nil {
		http.Error(w, fmt.Sprintf("incorrect content data: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(content.Data) == 0 {
		http.Error(w, "content data is required", http.StatusBadRequest)
		return
	}

	content.ID = id

	if err = s.databaseConnector.PatchContent(r.Context(), content); err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, false)
		return
	}

	w.WriteHeader(http.StatusOK)
}
