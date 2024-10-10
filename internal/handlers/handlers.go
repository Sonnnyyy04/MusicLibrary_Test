package handlers

import (
	"MusicLibrary_Test/internal/models"
	"MusicLibrary_Test/internal/repository"
	_ "MusicLibrary_Test/internal/response"
	"MusicLibrary_Test/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type SongHandler struct {
	Repo       *repository.SongRepository
	APIService *service.APIService
}

type Input struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

func NewSongHandler(repo *repository.SongRepository, apiService *service.APIService) *SongHandler {
	return &SongHandler{
		Repo:       repo,
		APIService: apiService,
	}
}

// @Summary Get songs
// @Description Get songs with filtering and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Group name"
// @Param song query string false "Song name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} models.Song
// @Failure 500 {object} response.ErrorResponse
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {
	filters := map[string]string{}
	if group := c.Query("group"); group != "" {
		filters["group"] = group
	}
	if song := c.Query("song"); song != "" {
		filters["song"] = song
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	songs, err := h.Repo.GetSongs(c.Request.Context(), filters, page, limit)
	if err != nil {
		logrus.Error("Error fetching songs: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, songs)
}

// @Summary Get song text
// @Description Get song text with pagination by verses
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number"
// @Param limit query int false "Verses per page"
// @Success 200 {string} string
// @Failure 500 {object} response.ErrorResponse
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongsText(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "4"))
	text, err := h.Repo.GetSongText(c.Request.Context(), id, page, limit)
	if err != nil {
		logrus.Error("Error fetching song text: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, text)
}

// @Summary Add new song
// @Description Add new song and enrich data from external API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song data"
// @Success 201
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /songs [post]
func (h *SongHandler) AddSong(c *gin.Context) {
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	song, err := h.APIService.GetSongDetails(input.Group, input.Song)
	if err != nil {
		logrus.Error("Error fetching song details: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song details"})
		return
	}

	err = h.Repo.AddSong(c.Request.Context(), *song)
	if err != nil {
		logrus.Error("Error add song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song"})
		return
	}
	c.Status(http.StatusCreated)
}

// @Summary Update song
// @Description Update song data
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Song data"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := h.Repo.UpdateSong(c.Request.Context(), id, song)
	if err != nil {
		logrus.Error("Error update song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song"})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Delete song
// @Description Delete song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 204
// @Failure 500 {object} response.ErrorResponse
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Repo.DeleteSong(c.Request.Context(), id)
	if err != nil {
		logrus.Error("Error delete song: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}
	c.Status(http.StatusNoContent)
}
