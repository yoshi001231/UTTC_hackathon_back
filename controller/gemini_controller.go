package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/usecase"

	"github.com/gorilla/mux"
)

// GeminiController Gemini関連エンドポイントのコントローラ
type GeminiController struct {
	geminiUseCase *usecase.GeminiUseCase
}

// NewGeminiController コントローラの初期化
func NewGeminiController(useCase *usecase.GeminiUseCase) *GeminiController {
	return &GeminiController{geminiUseCase: useCase}
}

// HandleGenerateBio 自己紹介生成エンドポイントのハンドラ
func (c *GeminiController) HandleGenerateBio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authID := vars["auth_id"]
	instruction := vars["instruction"]

	part, err := c.geminiUseCase.GenerateBio(authID, instruction)
	if err != nil {
		log.Printf("[gemini_controller.go] 自己紹介生成失敗 (auth_id: %s): %v", authID, err)
		http.Error(w, "自己紹介の生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(part); err != nil {
		log.Printf("[gemini_controller.go] jsonエンコード失敗 (auth_id: %s, instruction): %v", authID, instruction, err)
	}
}
