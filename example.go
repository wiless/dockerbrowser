package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fsouza/go-dockerclient"
)

func images(rw http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {

		rw.Header().Add("Access-Control-Allow-Origin", "*")
		rw.WriteHeader(http.StatusOK)

		endpoint := "unix:///var/run/docker.sock"
		client, _ := docker.NewClient(endpoint)

		imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
		log.Printf("Received request ")

		data, jerr := json.Marshal(imgs)
		if jerr == nil {
			rw.Write(data)
		}

		for _, img := range imgs {
			// fmt.Println("ID: ", img.ID)
			fmt.Println("RepoTags: ", img.RepoTags)
			fmt.Println("Created: ", img.Created)
			fmt.Println("Size: ", img.Size)
			// fmt.Println("VirtualSize: ", img.VirtualSize)
			fmt.Println("ParentId: ", img.ParentID)

		}
	} else {
		/// START the image

		rw.Header().Add("Access-Control-Allow-Origin", "*")
		rw.WriteHeader(http.StatusOK)

		endpoint := "unix:///var/run/docker.sock"
		client, _ := docker.NewClient(endpoint)

		req.ParseForm()

		image := req.FormValue("image")
		log.Printf("Requested for Image %v", image)
		var opt docker.CreateContainerOptions
		opt.Name = image
		opt.Config = new(docker.Config)
		opt.HostConfig = new(docker.HostConfig)

		opt.Config.Cmd = []string{"ls"}
		opt.HostConfig.RestartPolicy = docker.AlwaysRestart()
		info, derr := client.CreateContainer(opt)
		if derr != nil {
			log.Printf("CREATE CONTAINER ERROR  ?? %#v ", derr)
		} else {
			log.Printf("CREATE CONTAINER OUTPUT :  %#v", info)

			data, jerr := json.Marshal(info)
			if jerr == nil {
				rw.Write(data)
			}
		}

	}

}

func containers(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	if req.Method == "GET" {

		endpoint := "unix:///var/run/docker.sock"
		client, _ := docker.NewClient(endpoint)

		imgs, _ := client.ListContainers(docker.ListContainersOptions{All: true})
		log.Printf("Received request ")

		rw.WriteHeader(http.StatusOK)
		data, jerr := json.Marshal(imgs)
		if jerr == nil {
			rw.Write(data)
		}

	} else {

		endpoint := "unix:///var/run/docker.sock"
		client, _ := docker.NewClient(endpoint)

		req.ParseForm()

		id := req.FormValue("id")
		log.Printf("Requested for %v", id)
		info, _ := client.InspectContainer(id)
		rw.WriteHeader(http.StatusOK)
		data, jerr := json.Marshal(info)
		if jerr == nil {
			rw.Write(data)
		}
	}
}

func main() {
	http.HandleFunc("/images", images)
	http.HandleFunc("/containers", containers)
	http.ListenAndServe(":8080", nil)
}
