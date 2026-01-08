package main

func main() {
	
	// Entry point for the budget API v2 service
	// Create a Gin router with default middleware (logger and recovery)
	engine := server.WebServerProvisioner{}

	_ := engine.ConfigureEngine()
	

	engine.StartEngine()
}

