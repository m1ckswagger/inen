package server

import (
  "fmt"
  "github.com/gophercloud/gophercloud"
  "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
  "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
  "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/networks"
  "github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
  "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
  "github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
  "github.com/m1ckswagger/inenp/server/internal/handlers"
  "os"
)

func MakeVMConfig(client *gophercloud.ServiceClient, request handlers.VMRequest) (servers.CreateOptsBuilder, error) {
  flav, err := getFlavor(client, request.Flavor)
  if err != nil {
    return nil, err
  }
  img, err := getImage(client, request.OS)
  if err != nil {
    return nil, err
  }

  var srvNets []servers.Network
  for _, net := range request.Networks {
    n, err := getNetwork(client, net)
    if err != nil {
      return nil, err
    }
    srvNets = append(srvNets, servers.Network{UUID: n.ID})
  }

  blkDevices := []bootfromvolume.BlockDevice{
    {
      DeleteOnTermination: true,
      DestinationType:     bootfromvolume.DestinationVolume,
      SourceType:          bootfromvolume.SourceImage,
      UUID:                img.ID,
      VolumeSize:          request.HDD,
    },
  }
  srvCreateOpts := servers.CreateOpts{
    Name:             request.Hostname,
    AvailabilityZone: "nova",
    ImageRef:         img.ID,
    FlavorRef:        flav.ID,
    Networks:         srvNets,
  }

  blkCreateOpts := bootfromvolume.CreateOptsExt{
    CreateOptsBuilder: srvCreateOpts,
    BlockDevice:       blkDevices,
  }

  createOpts := keypairs.CreateOptsExt{
    CreateOptsBuilder: blkCreateOpts,
    KeyName:           os.Getenv("OS_KEY_NAME"),
  }

  return createOpts, nil
}

func getNetwork(client *gophercloud.ServiceClient, name string) (*networks.Network, error) {
  pages, err := networks.List(client).AllPages()
  if err != nil {
    return nil, fmt.Errorf("error listing networks: %s", err)
  }
  nets, err := networks.ExtractNetworks(pages)
  if err != nil {
    return nil, err
  }
  for _, net := range nets {
    if net.Label == name {
      return &net, nil
    }
  }
  return nil, fmt.Errorf("could not find network named %q", name)
}

func getImage(client *gophercloud.ServiceClient, imageName string) (images.Image, error) {
  listOpts := images.ListOpts{Name: imageName}
  pages, err := images.List(client, listOpts).AllPages()
  if err != nil {
    return images.Image{}, err
  }
  allImages, err := images.ExtractImages(pages)
  if err != nil {
    return images.Image{}, err
  }
  for _, image := range allImages {
    if image.Name == imageName {
      return image, nil
    }
  }
  return images.Image{}, fmt.Errorf("could not find image named %q", imageName)

}

func getFlavor(client *gophercloud.ServiceClient, flavorName string) (flavors.Flavor, error) {
  listOpts := flavors.ListOpts{AccessType: flavors.PublicAccess}
  allPages, err := flavors.ListDetail(client, listOpts).AllPages()
  if err != nil {
    return flavors.Flavor{}, err
  }
  allFlavors, err := flavors.ExtractFlavors(allPages)
  if err != nil {
    return flavors.Flavor{}, err
  }

  for _, flavor := range allFlavors {
    if flavor.Name == flavorName {
      return flavor, nil
    }
  }
  return flavors.Flavor{}, fmt.Errorf("could not find flavor %q", flavorName)
}
