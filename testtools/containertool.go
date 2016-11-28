package testtools

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/msproject/relive/dbmodel"
	// Used for the mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/samalba/dockerclient"
)

// RunMySQL starts a docker instance of mysql, connect to it, and
// create the associated tables.
func RunMySQL(setupSQL string) (mysqlID, dbconnect string, err error) {
	docker, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		return "", "", fmt.Errorf("error creating docker client %s", err.Error())
	}

	// If a container already exists use that
	var createSQL []string
	containers, err := docker.ListContainers(true, false, "")
	if err != nil {
		return "", "", fmt.Errorf("error listing docker containers %s", err.Error())
	}
Outer:
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/relive-mysql" {
				mysqlID = c.Id
				createSQL = []string{"drop database if exists relive"}
				break Outer
			}
		}
	}

	mysqlID, err = createAndStartMySQL(mysqlID, setupSQL)
	if err != nil {
		return "", "", fmt.Errorf("could not create or start container %s", err.Error())
	}

	ip, err := ContainerIP(mysqlID, docker)
	if err != nil {
		return "", "", fmt.Errorf("error fetching mysql ip %s", err)
	}

	// Open the db connections
	dbaddr := fmt.Sprintf("root:@tcp(%s:3306)/", ip)
	db, err := sql.Open("mysql", dbaddr)
	if err != nil {
		return "", "", fmt.Errorf("could not connect to mysql using %s %s", dbaddr, err.Error())
	}

	// And run the sql to remove and create DB (so as to verify the container is UP)
	for _, stmt := range append(createSQL, dbmodel.TableCreateSQL[:2]...) {
		// Sometimes mysql is a bit slower up, retry 5 times with a 10second sleep.
		for i := 0; i < 5; i++ {
			fmt.Println("Issue stmt:", stmt)
			_, err = db.Exec(stmt)
			if err != nil {
				err = fmt.Errorf("error issuing sql statement %s %s sleep 10 and retry", stmt, err.Error())
				fmt.Println(err.Error())
				time.Sleep(10 * time.Second)
			} else {
				break
			}
		}
		// If we still have an err after our retries, return err
		if err != nil {
			// If we still have an err after our retries, maybe the
			// docker container is not good anymore. we will try to
			// remove it, so hopefully next time will pass.
			_ = docker.RemoveContainer("relive-mysql", true, true)

			fmt.Println("DOCKER CONTAINER FOR mySQL MIGHT BE IN BAD SHAPE, REMOVED AND PLEASE RE-RUN TEST!!!")
			return "", "", err
		}
	}

	// return dbaddr to include DB name and default options
	dbaddr += "relive?interpolateParams=true&parseTime=true"

	return mysqlID, dbaddr, nil
}

// CreateAndStartMySQL - Create and Start MemSql container.
//
// Parameters:
// mysqlID is the name of existing container.
//          - If this value is empty then cName container will be deleted if present and recreated.
//          - If this is non empty then this API shall just start the existing cName container.
func createAndStartMySQL(mysqlID, cName string) (string, error) {
	if cName == "" {
		cName = "relive-mysql"
	}
	docker, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		return "", fmt.Errorf("error creating docker client %s", err.Error())
	}

	if mysqlID == "" {

		// If this API is called with mysqlID as empty, this means with
		_ = docker.RemoveContainer(cName, true, true)

		fmt.Printf("DOCKER CONTAINER FOR %s NOT FOUND!!!!!!! Creating!!!!\n", cName)
		// Create a container
		containerConfig := &dockerclient.ContainerConfig{
			Image: "mysql:latest",
		}
		mysqlID, err = docker.CreateContainer(containerConfig, cName, nil)
		if err != nil {
			return "", err
		}

	}

	// Start the container
	hostConfig := &dockerclient.HostConfig{
		PortBindings: map[string][]dockerclient.PortBinding{
			"3306/tcp": {{HostIp: "", HostPort: "3306"}},
		},
	}

	err = docker.StartContainer(mysqlID, hostConfig)
	if err != nil {
		return "", fmt.Errorf("could not start container %s", err.Error())
	}

	// Have to wait until mysql has bound ports
	time.Sleep(10 * time.Second)

	return mysqlID, nil
}

// ContainerIP - get container ip address
func ContainerIP(ID string, dc *dockerclient.DockerClient) (string, error) {
	res, err := dc.InspectContainer(ID)
	if err != nil {
		return "", fmt.Errorf("could not inspect docker container %s", err)
	}
	return res.NetworkSettings.IPAddress, nil
}

//StopContainer - Stop a running container identified by the id.
func StopContainer(id string) error {
	docker, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		return err
	}
	if err := docker.StopContainer(id, 2); err != nil {
		return err
	}
	return nil
}

// CleanUpTables - clean up tables, remove all data
func CleanUpTables(db *sql.DB) error {

	sqlStrs := []string{
		`Delete From Account`,
		`Delete From MediaType`,
		`Delete From Payment`,
		`Delete From PaymentHistory`,
		`Delete From Product`,
		`Delete From Subscription`,
		`Delete From SubscriptionAccount`,
	}

	for _, sqlStr := range sqlStrs {

		_, err := db.Exec(sqlStr)
		if err != nil {
			return err
		}
	}

	return nil
}
