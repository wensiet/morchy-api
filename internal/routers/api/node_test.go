package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/wensiet/morchy-api/internal/routers/api"
	"github.com/wensiet/morchy-api/internal/usecase"
	"github.com/wensiet/morchy-api/pkg/entity"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockedNodes = []*entity.Node{
	{
		ID:     uuid.MustParse("96222e5e-159e-46fd-8d43-32346436758c"),
		Status: entity.NewNodeStatus,
	},
	{
		ID:     uuid.MustParse("e3f50db7-5d0f-425e-98ac-0d575950033c"),
		Status: entity.RunningNodeStatus,
	},
	{
		ID:     uuid.MustParse("09a183f9-dde0-47f4-94f0-f5220d1d10bf"),
		Status: entity.FailedNodeStatus,
	},
}

type mockService struct {
	dbPool *pgxpool.Pool
}

func newMockService(dbPool *pgxpool.Pool) *mockService {
	return &mockService{
		dbPool: dbPool,
	}
}

func (m mockService) GetNode(_ context.Context, id uuid.UUID) (*entity.Node, error) {
	for _, node := range mockedNodes {
		if node.ID.String() == id.String() {
			return node, nil
		}
	}
	return nil, usecase.NodeNotFoundErr
}

func (m mockService) ListNodes(_ context.Context) ([]*entity.Node, error) {
	return mockedNodes, nil
}

func (m mockService) AddNode(_ context.Context) (*entity.Node, error) {
	newNode := entity.NewNode()
	mockedNodes = append(mockedNodes, newNode)
	return newNode, nil
}

func (m mockService) UpdateNode(_ context.Context, nodeModel *entity.Node) error {
	for _, node := range mockedNodes {
		if node.ID == nodeModel.ID {
			err := nodeModel.Status.Validate()
			if err != nil {
				return err
			}
			node.Status = nodeModel.Status
		}
	}
	return nil
}

func (m mockService) DeleteNode(_ context.Context, id uuid.UUID) error {
	var index *int
	for idx, node := range mockedNodes {
		if node.ID == id {
			index = &idx
		}
	}
	if index == nil {
		return usecase.NodeNotFoundErr
	}
	mockedNodes = append(mockedNodes[:*index], mockedNodes[*index+1:]...)
	return nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := newMockService(nil)
	nr := api.NewNodeRouter(mockService)

	r.GET("/node/:resource_id", nr.GetNode)
	r.GET("/nodes", nr.ListNodes)
	r.POST("/node", nr.AddNode)
	r.PUT("/node", nr.UpdateNode)
	r.DELETE("/node/:resource_id", nr.DeleteNode)

	return r
}

func TestNodeRouter_GetNode(t *testing.T) {
	r := setupRouter()
	testNode := mockedNodes[0]
	id := testNode.ID.String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/node/"+id, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseNode entity.Node
	err := json.Unmarshal(w.Body.Bytes(), &responseNode)
	assert.NoError(t, err)
	assert.Equal(t, testNode.ID, responseNode.ID)
	assert.Equal(t, testNode.Status, responseNode.Status)
}

func TestNodeRouter_ListNodes(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nodes", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNodeRouter_AddNode(t *testing.T) {
	r := setupRouter()
	oldLen := len(mockedNodes)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/node", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var responseNode entity.Node
	err := json.Unmarshal(w.Body.Bytes(), &responseNode)
	assert.NoError(t, err)
	assert.Equal(t, len(mockedNodes), oldLen+1)
	gotNode := false
	for _, nodeObj := range mockedNodes {
		if nodeObj.ID == responseNode.ID {
			assert.Equal(t, responseNode.Status, responseNode.Status)
			gotNode = true
		}
	}
	assert.True(t, gotNode)
}

func TestNodeRouter_UpdateNode(t *testing.T) {
	r := setupRouter()
	testNodeUpdate := entity.Node{
		ID:     mockedNodes[0].ID,
		Status: entity.FailedNodeStatus,
	}
	w := httptest.NewRecorder()
	updateNodeBody, err := json.Marshal(testNodeUpdate)
	req, _ := http.NewRequest("PUT", "/node", bytes.NewBuffer(updateNodeBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.NoError(t, err)
	assert.Equal(t, mockedNodes[0].Status, entity.FailedNodeStatus)
}

func TestNodeRouter_DeleteNode(t *testing.T) {
	r := setupRouter()
	testNode := mockedNodes[0]
	oldLen := len(mockedNodes)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/node/"+testNode.ID.String(), nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, len(mockedNodes), oldLen-1)
	for _, nodeObj := range mockedNodes {
		assert.NotEqual(t, testNode.ID, nodeObj.ID)
	}
}
