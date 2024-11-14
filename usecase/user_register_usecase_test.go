package usecase

import (
	"strings"
	"kaizen/model"
	"testing"
)

// MockUserDAO は、UserDAOInterfaceのモック実装
type MockUserDAO struct {
	RegisterUserFunc   func(user model.User) error
	FindUserByNameFunc func(name string) ([]model.User, error)
}

func (m *MockUserDAO) RegisterUser(user model.User) error {
	return m.RegisterUserFunc(user)
}

func (m *MockUserDAO) FindUserByName(name string) ([]model.User, error) {
	return m.FindUserByNameFunc(name)
}

func TestRegisterUserUseCase_Execute(t *testing.T) {
	mockDAO := &MockUserDAO{
		RegisterUserFunc: func(user model.User) error {
			return nil // 成功時はエラーなし
		},
		FindUserByNameFunc: func(name string) ([]model.User, error) {
			return []model.User{}, nil // ここでは使わないため適当な戻り値
		},
	}

	useCase := NewRegisterUserUseCase(mockDAO)

	tests := []struct {
		id      string
		name    string
		age     int
		wantErr bool
		errMsg  string
	}{
		// 正常な入力
		{"test_id_1", "ValidName", 30, false, ""},
		
		// 名前が空
		{"test_id_2", "", 30, true, "invalid name"},
		
		// 名前が長すぎる
		{"test_id_3", "ThisNameIsWayTooLongToBeValidAndShouldCauseAnError!!!!!!!!!!!!!!!!!", 30, true, "invalid name"},
		
		// 年齢が範囲外（20未満）
		{"test_id_4", "ValidName", 19, true, "invalid age"},
		
		// 年齢が範囲外（80を超える）
		{"test_id_5", "ValidName", 81, true, "invalid age"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			_, err := useCase.Execute(tt.id, tt.name, tt.age)

			// エラーの有無が期待と一致するか確認
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーメッセージの内容を確認
			if err != nil && !containsErrorMessage(err, tt.errMsg) {
				t.Errorf("Execute() error message = %v, wantErr %v", err, tt.errMsg)
			}
		})
	}
}

// containsErrorMessageはエラーメッセージに特定の文字列が含まれているか確認します
func containsErrorMessage(err error, msg string) bool {
	return err != nil && strings.Contains(err.Error(), msg)
}
