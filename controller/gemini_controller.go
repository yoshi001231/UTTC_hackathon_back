package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/usecase"

	"github.com/gorilla/mux"
)

// InstructionRequest リクエストボディの構造体
type InstructionRequest struct {
	Instruction *string `json:"instruction"`
}

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

	// JSONリクエストのパース
	var req InstructionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[gemini_controller.go] リクエストデコード失敗 (auth_id: %s): %v", authID, err)
		http.Error(w, "リクエスト形式が不正です", http.StatusBadRequest)
		return
	}

	// `instruction` が null の場合は空文字列を代入
	instruction := ""
	if req.Instruction != nil {
		instruction = *req.Instruction
	}

	part, err := c.geminiUseCase.GenerateBio(authID, instruction)
	if err != nil {
		log.Printf("[gemini_controller.go] 自己紹介生成失敗 (auth_id: %s): %v", authID, err)
		http.Error(w, "自己紹介の生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(part); err != nil {
		log.Printf("[gemini_controller.go] jsonエンコード失敗 (auth_id: %s): %v", authID, err)
		http.Error(w, "レスポンスの生成に失敗しました", http.StatusInternalServerError)
	}
}
