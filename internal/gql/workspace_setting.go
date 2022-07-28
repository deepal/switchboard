package gql

import (
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/graphql-go/graphql"
)

var ScenarioConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ScenarioConfig",
	Fields: graphql.Fields{
		"scenario": &graphql.Field{
			Type: ScenarioGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				scenarioID := p.Source.(models.ScenarioConfig).ScenarioID
				scenario, err := db.GetScenarioByID(scenarioID)
				if err != nil {
					return nil, err
				}
				return scenario, nil
			},
		},
		"isActive": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var EndpointConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EndpointConfig",
	Fields: graphql.Fields{
		"endpoint": &graphql.Field{
			Type:    EndpointGqlType,
			Resolve: EndpointResolver,
		},
		"scenarioConfigs": &graphql.Field{
			Type: graphql.NewList(ScenarioConfigGqlType),
		},
		"responseDelay": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var WorkspaceSettingGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WorkspaceSetting",
	Fields: graphql.Fields{
		"workspace": &graphql.Field{
			Type: WorkspaceGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				workspaceID := p.Source.(*models.WorkspaceSetting).WorkspaceID
				currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
				ws, err := db.GetUserWorkspaceByID(currentUser.ID, workspaceID)
				if err != nil {
					return nil, err
				}
				return ws, nil
			},
		},
		"mockService": &graphql.Field{
			Type: MockServiceGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mockServiceID := p.Source.(*models.WorkspaceSetting).MockServiceID
				ms, err := db.GetMockServiceByID(mockServiceID)
				if err != nil {
					return nil, err
				}
				return *ms, nil
			},
		},
		"config": &graphql.Field{
			Type: GlobalMockServiceConfigGqlType,
		},
		"endpointConfigs": &graphql.Field{
			Type: graphql.NewList(EndpointConfigGqlType),
		},
	},
})

var EndpointResolver = func(p graphql.ResolveParams) (interface{}, error) {
	endpointID := p.Source.(models.EndpointConfig).EndpointID
	ep, err := db.GetEndpointByID(endpointID)
	if err != nil {
		return nil, err
	}
	return *ep, nil
}

var WorkspaceSettingsResolver = func(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Source.(models.Workspace).ID
	wss, err := db.GetWorkspaceSettings(workspaceID)
	if err != nil {
		return make([]models.WorkspaceSetting, 0), err
	}

	return wss, nil
}

var WorkspaceMockServicesResolver = func(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Source.(models.Workspace).ID
	wss, errWs := db.GetWorkspaceSettings(workspaceID)
	if errWs != nil {
		return make([]models.MockService, 0), errWs
	}
	mockServiceIds := make([]string, 0)
	for _, ws := range *wss {
		mockServiceIds = append(mockServiceIds, ws.MockServiceID)
	}

	ms, errMs := db.GetMockServicesByIds(mockServiceIds)
	if errMs != nil {
		return nil, errMs
	}

	return *ms, nil
}

var WorkspaceSettingResolver = func(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Args["workspaceId"].(string)
	mockServiceID := p.Args["mockServiceId"].(string)
	wss, err := db.GetWorkspaceSetting(workspaceID, mockServiceID)
	if err != nil {
		return nil, err
	}

	return wss, nil
}