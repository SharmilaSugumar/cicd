package handlers

import (
	"encoding/json"
	"forgeflow/internal/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MetricsHandler struct {
	db *gorm.DB
}

func NewMetricsHandler(db *gorm.DB) *MetricsHandler {
	return &MetricsHandler{db: db}
}

func (h *MetricsHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", h.GetMetrics)
}

type MetricsResponse struct {
	PipelineCount    int64                  `json:"pipeline_count"`
	ProjectCount     int64                  `json:"project_count"`
	OrgCount         int64                  `json:"org_count"`
	SuccessRate      float64                `json:"success_rate"`
	ActivityData     []ActivityMetric       `json:"activity_data"`
	LanguageData     []LanguageMetric       `json:"language_data"`
}

type ActivityMetric struct {
	Name      string `json:"name"`
	Pipelines int64  `json:"pipelines"`
	Success   int64  `json:"success"`
}

type LanguageMetric struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
	Color string `json:"color"`
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	var res MetricsResponse

	h.db.Model(&database.Pipeline{}).Count(&res.PipelineCount)
	h.db.Model(&database.Project{}).Count(&res.ProjectCount)
	h.db.Model(&database.Organization{}).Count(&res.OrgCount)

	var totalRuns int64
	var successRuns int64
	h.db.Model(&database.PipelineRun{}).Count(&totalRuns)
	h.db.Model(&database.PipelineRun{}).Where("status = ?", "COMPLETED").Or("status = ?", "SUCCESS").Count(&successRuns)

	if totalRuns > 0 {
		res.SuccessRate = float64(successRuns) / float64(totalRuns) * 100
	}

	// 7 days activity
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	today := time.Now()
	for i := 6; i >= 0; i-- {
		targetDate := today.AddDate(0, 0, -i)
		startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
		endOfDay := startOfDay.AddDate(0, 0, 1)

		var dailyTotal int64
		var dailySuccess int64
		h.db.Model(&database.PipelineRun{}).Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).Count(&dailyTotal)
		h.db.Model(&database.PipelineRun{}).Where("created_at >= ? AND created_at < ? AND (status = 'COMPLETED' OR status = 'SUCCESS')", startOfDay, endOfDay).Count(&dailySuccess)

		res.ActivityData = append(res.ActivityData, ActivityMetric{
			Name:      days[targetDate.Weekday()],
			Pipelines: dailyTotal,
			Success:   dailySuccess,
		})
	}

	// Language distribution from jobs payload
	var jobs []database.Job
	h.db.Find(&jobs)

	langCounts := make(map[string]int64)
	for _, j := range jobs {
		var payload struct {
			Language string `json:"language"`
		}
		if err := json.Unmarshal([]byte(j.Payload), &payload); err == nil && payload.Language != "" {
			langCounts[payload.Language]++
		}
	}

	colors := []string{"#4ade80", "#60a5fa", "#fcd34d", "#f472b6", "#a78bfa"}
	i := 0
	for lang, count := range langCounts {
		res.LanguageData = append(res.LanguageData, LanguageMetric{
			Name:  lang,
			Count: count,
			Color: colors[i%len(colors)],
		})
		i++
	}

	if res.LanguageData == nil {
		res.LanguageData = []LanguageMetric{}
	}

	c.JSON(http.StatusOK, res)
}
