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

	portfolios, err := s.databaseConnector.GetAllPortfolios(r.Context(), page.limit, page.offset, *filter)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, true)
		return
	}

	json.NewEncoder(w).Encode(portfolios)
}

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

func (s *Server) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	page, err := s.getPageInfo(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect page info: %s", err.Error()), http.StatusBadRequest)
		return
	}

	categories, err := s.databaseConnector.GetAllCategories(r.Context(), page.limit, page.offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(categories)
}

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

	crafts, err := s.databaseConnector.GetAllCraftsByPortfolioID(r.Context(), id, page.limit, page.offset)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, true)
		return
	}

	json.NewEncoder(w).Encode(crafts)
}

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

	crafts, err := s.databaseConnector.GetAllCraftsByTagID(r.Context(), id, page.limit, page.offset)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, true)
		return
	}

	json.NewEncoder(w).Encode(crafts)
}

func (s *Server) getTagsHandler(w http.ResponseWriter, r *http.Request) {
	page, err := s.getPageInfo(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect page info: %s", err.Error()), http.StatusBadRequest)
		return
	}

	tags, err := s.databaseConnector.GetAllTags(r.Context(), page.limit, page.offset)
	if err != nil {
		response_errors.StatusCodeByErrorWriter(err, w, true)
		return
	}

	json.NewEncoder(w).Encode(tags)
}

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
