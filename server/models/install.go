package models


var adminInit = func() error {
	return CreateAdminRow("admin", "admin", "admin", "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif", "admin introduce", "admin")
}

var initProject *Project

var projectInit = func() error {
	initProject = &Project{}
	initProject.Name = "ExampleProject"
	initProject.Desc = "Example project"
	return initProject.Save("admin")
}

var envInit = func() error {
	return CreateEnvironmentRow(initProject, "Development", "Development Environment", "admin")
}

var installScripts = []func() error{adminInit, projectInit, envInit}

func Install() error {
	for _, f := range installScripts {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
