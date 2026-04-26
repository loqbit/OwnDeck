package main

import (
	"context"

	"OwnDeck/internal/discovery"
	"OwnDeck/internal/repository/config"
	"OwnDeck/internal/service/connectionsvc"
	"OwnDeck/internal/service/discoverysvc"
)

// App is the Wails handler layer. It only validates input and delegates
// to services; all business logic lives in internal/service/*.
type App struct {
	ctx           context.Context
	discoverySvc  *discoverysvc.Service
	connectionSvc *connectionsvc.Service
	store         config.Store
}

func NewApp(discoverySvc *discoverysvc.Service, connectionSvc *connectionsvc.Service, store config.Store) *App {
	return &App{
		discoverySvc:  discoverySvc,
		connectionSvc: connectionSvc,
		store:         store,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) DiscoverMCPServers() ([]discovery.MCPServer, error) {
	return a.discoverySvc.MCPServers(a.ctx)
}

func (a *App) DiscoverClients() []discovery.ClientInfo {
	return a.discoverySvc.Clients()
}

func (a *App) DiscoverMCPServersForClients(clientIDs []string) ([]discovery.MCPServer, error) {
	return a.discoverySvc.MCPServersForClients(a.ctx, clientIDs)
}

func (a *App) DiscoverSkillsForClients(clientIDs []string) ([]discovery.SkillAsset, error) {
	return a.discoverySvc.SkillsForClients(a.ctx, clientIDs)
}

func (a *App) GetConfig() (config.AppConfig, error) {
	return a.connectionSvc.GetConfig()
}

func (a *App) GetConnectedClientIDs() ([]string, error) {
	return a.connectionSvc.ConnectedIDs()
}

func (a *App) ConnectClient(clientID string) (config.AppConfig, error) {
	return a.connectionSvc.Connect(clientID)
}

func (a *App) DisconnectClient(clientID string) (config.AppConfig, error) {
	return a.connectionSvc.Disconnect(clientID)
}

func (a *App) ConfigPath() string {
	return a.connectionSvc.Path()
}

func (a *App) IntrospectMCPServer(server discovery.MCPServer) discovery.MCPServer {
	return a.discoverySvc.IntrospectServer(a.ctx, server)
}

// RescanAgents re-runs the agent filesystem scanner, updates the
// persisted config, and returns the fresh agent list.
func (a *App) RescanAgents() ([]config.AgentConfig, error) {
	agents := discovery.ScanAgents()

	cfg, err := a.store.Load()
	if err != nil {
		return agents, err
	}
	cfg.Agents = agents
	if err := a.store.Save(cfg); err != nil {
		return agents, err
	}
	return agents, nil
}
