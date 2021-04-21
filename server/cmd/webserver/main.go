package main

import (
  "flag"
  "github.com/gophercloud/gophercloud"
  "github.com/gophercloud/gophercloud/openstack"
  "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
  "github.com/joho/godotenv"
  "github.com/m1ckswagger/inenp/server/internal/server"
  "log"
  "net/http"
  "os"
)

var l *log.Logger

var openrcPath string

func main() {
  l = log.New(os.Stdout, "server-api ", log.LstdFlags|log.Lmsgprefix)

  flag.StringVar(&openrcPath, "openrc", ".openrc", "openrc file in .env format")
  flag.Parse()

  err := godotenv.Load(openrcPath)
  if err != nil {
    l.Fatalf("Could not open %q", openrcPath)
  }

	// read config from env
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		l.Fatalf("error reading opts from Env: %v", err)
	}
	// create provider from options
	provider, err := openstack.AuthenticatedClient(opts)
  client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")})
	if err != nil {
		l.Fatalf("error creating client from options: %v", err)
	}
	l.Printf("authenticated as: %s\n", opts.Username)

  createOpts, err := server.MakeVMConfig(client, server.VMRequest{})
  if err != nil {
    l.Fatal(err)
  }

	_, err = bootfromvolume.Create(client, createOpts).Extract()
	if err != nil {
		l.Fatalf("error creating server: %v", err)
	}

	l.Printf("successfully created server")

	loginUsers := map[string]string{
		"administrator": "administrator",
	}
	store := &server.UserMemoryStore{DB: loginUsers}

	if err := http.ListenAndServe(":5000", server.NewInfraServer(store)); err != nil {
		log.Fatal(err)
	}
}


