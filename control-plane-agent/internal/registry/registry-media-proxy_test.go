// nolint
package registry

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"control-plane-agent/internal/model"
)

func TestMediaProxyRegistry(t *testing.T) {
	MediaProxyRegistry.Init(MediaProxyRegistryConfig{
		CancelCommandRequestFunc: func(reqId string) {
			// Mock function for testing
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := MediaProxyRegistry.Run(ctx)
		require.NoError(t, err)
	}()

	// Test Add
	proxy := model.MediaProxy{
		Config: &model.MediaProxyConfig{},
		Status: &model.MediaProxyStatus{},
	}

	id, err := MediaProxyRegistry.Add(ctx, proxy)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// Test Get
	retrievedProxy, err := MediaProxyRegistry.Get(ctx, id, false)
	require.NoError(t, err)
	require.Equal(t, id, retrievedProxy.Id)
	require.NotNil(t, retrievedProxy.Status)
	require.Nil(t, retrievedProxy.Config)

	retrievedProxy, err = MediaProxyRegistry.Get(ctx, id, true)
	require.NoError(t, err)
	require.Equal(t, id, retrievedProxy.Id)
	require.NotNil(t, retrievedProxy.Status)
	require.NotNil(t, retrievedProxy.Config)

	// Test List
	proxy2 := model.MediaProxy{
		Config: &model.MediaProxyConfig{},
		Status: &model.MediaProxyStatus{},
	}

	id2, err := MediaProxyRegistry.Add(ctx, proxy2)
	require.NoError(t, err)

	proxies, err := MediaProxyRegistry.List(ctx, nil, false, false)
	require.NoError(t, err)
	require.Len(t, proxies, 2)
	require.Equal(t, id, proxies[0].Id)
	require.Nil(t, proxies[0].Status)
	require.Nil(t, proxies[0].Config)
	require.Equal(t, id2, proxies[1].Id)
	require.Nil(t, proxies[1].Status)
	require.Nil(t, proxies[1].Config)

	proxies, err = MediaProxyRegistry.List(ctx, nil, true, true)
	require.NoError(t, err)
	require.Len(t, proxies, 2)
	require.NotNil(t, proxies[0].Status)
	require.NotNil(t, proxies[0].Config)
	require.NotNil(t, proxies[1].Status)
	require.NotNil(t, proxies[1].Config)

	// Test Update_LinkConn
	connId := uuid.NewString()
	err = MediaProxyRegistry.Update_LinkConn(ctx, id, connId)
	require.NoError(t, err)

	updatedProxy, err := MediaProxyRegistry.Get(ctx, id, false)
	require.NoError(t, err)
	require.Contains(t, updatedProxy.ConnIds, connId)

	// Test Update_UnlinkConn
	err = MediaProxyRegistry.Update_UnlinkConn(ctx, id, connId)
	require.NoError(t, err)

	updatedProxy, err = MediaProxyRegistry.Get(ctx, id, false)
	require.NoError(t, err)
	require.NotContains(t, updatedProxy.ConnIds, connId)

	// Test Update_LinkBridge
	bridgeId := uuid.NewString()
	err = MediaProxyRegistry.Update_LinkBridge(ctx, id, bridgeId)
	require.NoError(t, err)

	updatedProxy, err = MediaProxyRegistry.Get(ctx, id, false)
	require.NoError(t, err)
	require.Contains(t, updatedProxy.BridgeIds, bridgeId)

	// Test Update_UnlinkBridge
	err = MediaProxyRegistry.Update_UnlinkBridge(ctx, id, bridgeId)
	require.NoError(t, err)

	updatedProxy, err = MediaProxyRegistry.Get(ctx, id, false)
	require.NoError(t, err)
	require.NotContains(t, updatedProxy.BridgeIds, bridgeId)

	// Test Delete
	err = MediaProxyRegistry.Delete(ctx, id)
	require.NoError(t, err)

	_, err = MediaProxyRegistry.Get(ctx, id, false)
	require.Error(t, err)

	cancel()
	wg.Wait()
}