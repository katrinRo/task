package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todolist/model"
)

// MockRepository для тестирования сервиса
type MockRepository struct {
	CreateFunc func(task model.TaskRequest) error
	ChekIdFunc func(id int) (bool, error)
	GetFunc    func(id int) (model.TaskResponse, error)
	UpdateFunc func(task model.TaskRequest) error
	DeleteFunc func(id int) error
	ListFunc   func(req model.ListRequest) ([]model.TaskResponse, error)
}

func (m *MockRepository) Create(task model.TaskRequest) error {
	return m.CreateFunc(task)
}

func (m *MockRepository) ChekId(id int) (bool, error) {
	return m.ChekIdFunc(id)
}

func (m *MockRepository) Get(id int) (model.TaskResponse, error) {
	return m.GetFunc(id)
}

func (m *MockRepository) Update(task model.TaskRequest) error {
	return m.UpdateFunc(task)
}

func (m *MockRepository) Delete(id int) error {
	return m.DeleteFunc(id)
}

func (m *MockRepository) List(req model.ListRequest) ([]model.TaskResponse, error) {
	return m.ListFunc(req)
}

func TestTodoService_Get(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(repo *MockRepository)
		wantResp  model.TaskResponse
		wantErr   bool
	}{
		{
			name: "Успешное получение",
			setupRepo: func(repo *MockRepository) {
				repo.ChekIdFunc = func(id int) (bool, error) { return true, nil }
				repo.GetFunc = func(id int) (model.TaskResponse, error) {
					return model.TaskResponse{
						ID:          1,
						Title:       "Test Task",
						Description: "Test Description",
						Date:        "2023-01-01",
					}, nil
				}
			},
			wantResp: model.TaskResponse{
				ID:          1,
				Title:       "Test Task",
				Description: "Test Description",
				Date:        "2023-01-01",
			},
			wantErr: false,
		},
		{
			name: "id не найден",
			setupRepo: func(repo *MockRepository) {
				repo.ChekIdFunc = func(id int) (bool, error) { return false, nil }
			},
			wantResp: model.TaskResponse{},
			wantErr:  true,
		},
		{
			name: "Ошибка в базе",
			setupRepo: func(repo *MockRepository) {
				repo.ChekIdFunc = func(id int) (bool, error) { return true, nil }
				repo.GetFunc = func(id int) (model.TaskResponse, error) { return model.TaskResponse{}, errors.New("Ошибка") }
			},
			wantResp: model.TaskResponse{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{}
			tt.setupRepo(repo)
			service := NewTodoService(repo)

			params := httprouter.Params{
				httprouter.Param{Key: "id", Value: "1"},
			}
			resp, errResp := service.Get(params)

			if (errResp != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", errResp, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("Get() resp = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestTodoService_Create(t *testing.T) {
	repo := &MockRepository{}
	service := NewTodoService(repo)

	tasks := []struct {
		name      string
		task      model.TaskRequest
		wantError bool
		setupRepo func()
	}{
		{
			name: "Успех",
			task: model.TaskRequest{
				Title:       "Test Task",
				Description: "Test Description",
				Date:        "2023-01-01",
				Status:      true,
			},
			wantError: false,
			setupRepo: func() {
				repo.CreateFunc = func(task model.TaskRequest) error { return nil }
			},
		},
		{
			name: "Отсутсвие статуса",
			task: model.TaskRequest{
				Title:       "Test Task",
				Description: "Test Description",
				Date:        "2023-01-01",
			},
			wantError: true,
			setupRepo: func() {
				repo.CreateFunc = func(task model.TaskRequest) error { return errors.New("mock error") }
			},
		},
		{
			name: "Отсутвует заголовка",
			task: model.TaskRequest{
				Description: "Test Description",
				Date:        "2023-01-01",
				Status:      true,
			},
			wantError: true,
			setupRepo: func() {
				repo.CreateFunc = func(task model.TaskRequest) error { return nil }
			},
		},
		{
			name: "Отсутсвие описания",
			task: model.TaskRequest{
				Title:  "Test Task",
				Date:   "2023-01-01",
				Status: true,
			},
			wantError: false,
			setupRepo: func() {
				repo.CreateFunc = func(task model.TaskRequest) error { return nil }
			},
		},
		{
			name: "Не верная дата",
			task: model.TaskRequest{
				Title:       "Test Task",
				Description: "Test Description",
				Date:        "invalid-date",
			},
			wantError: true,
			setupRepo: func() {
				repo.CreateFunc = func(task model.TaskRequest) error { return nil }
			},
		},
	}

	for _, tt := range tasks {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupRepo()

			jsonTask, err := json.Marshal(tt.task)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonTask))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := service.Create(req)

			if tt.wantError {
				require.NotNil(t, resp)
			} else {
				require.Nil(t, resp)
			}
		})
	}
}

func TestTodoService_Update(t *testing.T) {
	// Mock репозитория
	repo := &MockRepository{}
	service := NewTodoService(repo)

	tasks := []struct {
		name      string
		task      model.TaskRequest
		wantError bool
		setupRepo func()
	}{
		{
			name: "Успех",
			task: model.TaskRequest{
				ID:          1,
				Title:       "Test Task",
				Description: "Test Description",
				Date:        "2023-01-01",
				Status:      true,
			},
			wantError: false,
			setupRepo: func() {
				repo.ChekIdFunc = func(id int) (bool, error) { return true, nil }
				repo.UpdateFunc = func(task model.TaskRequest) error { return nil }
			},
		},
		{
			name: "Отсутсвие id",
			task: model.TaskRequest{
				Title:       "Test Task",
				Description: "Test Description",
				Date:        "2023-01-01",
			},
			wantError: true,
			setupRepo: func() {
				repo.ChekIdFunc = func(id int) (bool, error) { return false, nil }
				repo.UpdateFunc = func(task model.TaskRequest) error { return nil }
			},
		},
		{
			name: "Отсутсвие статуса",
			task: model.TaskRequest{
				ID:          1,
				Title:       "Test Task",
				Description: "Test Description",
				Date:        "2023-01-01",
			},
			wantError: true,
			setupRepo: func() {
				repo.ChekIdFunc = func(id int) (bool, error) { return true, nil }
				repo.UpdateFunc = func(task model.TaskRequest) error { return errors.New("mock error") }
			},
		},
		{
			name: "Отсутвует заголовка",
			task: model.TaskRequest{
				ID:          1,
				Description: "Test Description",
				Date:        "2023-01-01",
				Status:      true,
			},
			wantError: true,
			setupRepo: func() {
				repo.ChekIdFunc = func(id int) (bool, error) { return true, nil }
				repo.UpdateFunc = func(task model.TaskRequest) error { return nil }
			},
		},
		{
			name: "Отсутсвие описания",
			task: model.TaskRequest{
				ID:     1,
				Title:  "Test Task",
				Date:   "2023-01-01",
				Status: true,
			},
			wantError: false,
			setupRepo: func() {
				repo.UpdateFunc = func(task model.TaskRequest) error { return nil }
			},
		},
		{
			name: "Не верная дата",
			task: model.TaskRequest{
				ID:          1,
				Title:       "Test Task",
				Description: "Test Description",
				Date:        "invalid-date",
			},
			wantError: true,
			setupRepo: func() {
				repo.ChekIdFunc = func(id int) (bool, error) { return true, nil }
				repo.UpdateFunc = func(task model.TaskRequest) error { return nil }
			},
		},
	}

	for _, tt := range tasks {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupRepo()

			jsonTask, err := json.Marshal(tt.task)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonTask))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := service.Update(req)

			if tt.wantError {
				require.NotNil(t, resp)
			} else {
				require.Nil(t, resp)
			}
		})
	}
}

func TestTodoService_Delete(t *testing.T) {
	// Mock репозитория
	repo := &MockRepository{}
	service := NewTodoService(repo)

	tests := []struct {
		name       string
		checkId    func(id int) (bool, error)
		deleteFunc func(id int) error
		wantErr    bool
	}{
		{
			name:       "ID существует",
			checkId:    func(id int) (bool, error) { return true, nil },
			deleteFunc: func(id int) error { return nil },
			wantErr:    false,
		},
		{
			name:       "ID не существует",
			checkId:    func(id int) (bool, error) { return false, nil },
			deleteFunc: func(id int) error { return nil },
			wantErr:    true,
		},
		{
			name:       "База данных недоступна при проверке ID",
			checkId:    func(id int) (bool, error) { return true, errors.New("database error") },
			deleteFunc: func(id int) error { return nil },
			wantErr:    true,
		},
		{
			name:       "База данных недоступна при удалении",
			checkId:    func(id int) (bool, error) { return true, nil },
			deleteFunc: func(id int) error { return errors.New("database error") },
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.ChekIdFunc = tt.checkId
			repo.DeleteFunc = tt.deleteFunc
			params := httprouter.Params{
				httprouter.Param{Key: "id", Value: "1"},
			}
			errResp := service.Delete(params)

			if (errResp != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", errResp, tt.wantErr)
			}
		})
	}
}

func TestTodoService_List(t *testing.T) {
	// Mock репозитория
	repo := &MockRepository{}
	service := NewTodoService(repo)

	// Тестовые данные
	tasks := []model.TaskResponse{
		{ID: 1, Title: "Task 1", Description: "Desc 1", Date: "2023-01-01"},
		{ID: 2, Title: "Task 2", Description: "Desc 2", Date: "2023-01-02"},
	}

	// Сценарии тестирования
	testCases := []struct {
		name        string
		listReq     model.ListRequest
		expectError bool
		expectLen   int
	}{
		{
			name: "Успешное получение",
			listReq: model.ListRequest{
				Limit:  intPtr(10),
				Offset: intPtr(0),
			},
			expectError: false,
			expectLen:   2,
		},
		{
			name: "Отсутсвует limit",
			listReq: model.ListRequest{
				Offset: intPtr(0),
			},
			expectError: true,
			expectLen:   0,
		},
		{
			name: "Отсутсвует offset",
			listReq: model.ListRequest{
				Limit: intPtr(10),
			},
			expectError: true,
			expectLen:   0,
		},
		{
			name: "Некорректная дата",
			listReq: model.ListRequest{
				Limit:  intPtr(10),
				Offset: intPtr(0),
				Date:   "invalid-date",
			},
			expectError: true,
			expectLen:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo.ListFunc = func(req model.ListRequest) ([]model.TaskResponse, error) {
				if tc.expectError {
					return nil, errors.New("mock error")
				}
				return tasks, nil
			}

			jsonReq, err := json.Marshal(tc.listReq)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/tasks/list", bytes.NewBuffer(jsonReq))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, errResp := service.List(req)

			if tc.expectError {
				assert.NotNil(t, errResp)
			} else {
				require.Nil(t, errResp)
				assert.Len(t, resp, tc.expectLen)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
