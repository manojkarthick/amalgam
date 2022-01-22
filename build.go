package main

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"

	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
)

// application's configuration
type Configuration struct {
	Owner              string
	Repository         string
	Tag                string
	Amd64Substring     string
	Arm64Substring     string
	Compressed         bool
	BinaryPath         string
	GithubToken        string
	Overwrite          bool
	UniversalIdentifer string
}

const (
	EXTRACT_PREFIX  = "extract"
	DOWNLOAD_PREFIX = "download"
	FAT_PREFIX      = "fat"
)

// create universal binary using the given macOS binaries
// for amd64 and arm64
func CreateUniveralBinary(cfg *Configuration) error {
	log.Info("Starting process...")
	log.Debugf("Repository: %s/%s", cfg.Owner, cfg.Repository)
	log.Debugf("Tag: %s", cfg.Tag)
	log.Debugf("Amd64 Regex: %s", cfg.Amd64Substring)
	log.Debugf("Arm64 Regex: %s", cfg.Arm64Substring)
	log.Debugf("Compressed?: %t", cfg.Compressed)
	log.Debugf("Binary path: %s", cfg.BinaryPath)

	release, err := cfg.GetRelease()
	if err != nil {
		return err
	}

	log.Info("Finding relevant Github release")
	amd64Asset, arm64Asset, err := cfg.FilterArtifacts(release.Assets)
	if err != nil {
		return err
	}

	amd64Path, err := cfg.DownloadAndGetPath(amd64Asset)
	if err != nil {
		return err
	}
	arm64Path, err := cfg.DownloadAndGetPath(arm64Asset)
	if err != nil {
		return err
	}
	log.Info("Successfully downloaded release assets")

	log.Debugf("The path to the amd64 binary is: %s", amd64Path)
	log.Debugf("The path to the arm64 binary is: %s", arm64Path)

	fatPath, err := cfg.Combine(amd64Path, arm64Path, *amd64Asset.Name)
	if err != nil {
		return err
	}
	log.Info("Combined assets into universal binary")
	log.Debugf("The path to the fat binary is: %s", fatPath)

	err = cfg.UploadAsset(release, fatPath)
	if err != nil {
		return err
	}
	log.Info("Uploaded universal binary to Github")

	return nil
}

// Get the github release with the given tag name.
// If the given tag is `latest`, return the latest github release
func (cfg *Configuration) GetRelease() (*github.RepositoryRelease, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: cfg.GithubToken,
	})

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var release *github.RepositoryRelease
	var err error
	if cfg.Tag == "latest" {
		log.Info("Finding the latest release")
		release, _, err = client.Repositories.GetLatestRelease(ctx, cfg.Owner, cfg.Repository)
	} else {
		log.Info("Finding the release with tag ", cfg.Tag)
		release, _, err = client.Repositories.GetReleaseByTag(ctx, cfg.Owner, cfg.Repository, cfg.Tag)
	}

	if err != nil {
		return nil, err
	}
	return release, nil
}

// Filter the available assets for the given release to find the macOS specific assets
func (cfg *Configuration) FilterArtifacts(assets []*github.ReleaseAsset) (*github.ReleaseAsset, *github.ReleaseAsset, error) {
	log.Debug("Filtering artifacts")

	var amd64Asset, arm64Asset *github.ReleaseAsset

	for _, asset := range assets {
		assetName := asset.GetName()
		if strings.Contains(assetName, cfg.Amd64Substring) {
			amd64Asset = asset
		}
		if strings.Contains(assetName, cfg.Arm64Substring) {
			arm64Asset = asset
		}
	}

	if amd64Asset == nil || arm64Asset == nil {
		return nil, nil, errors.New("could not find needed artifacts from github release")
	}

	return amd64Asset, arm64Asset, nil
}

// Download the Github asset to the local filesystem and return the path for the binary.
// If the asset is compressed, uncompress and return the specific path to the binary
func (cfg *Configuration) DownloadAndGetPath(asset *github.ReleaseAsset) (string, error) {
	name := asset.GetName()
	url := asset.GetBrowserDownloadURL()

	log.Debugf("Asset Name: %s", name)
	log.Debugf("Asset Download URL: %s", url)

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	downloadDir, err := ioutil.TempDir(currentWorkingDir, DOWNLOAD_PREFIX)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(downloadDir)

	downloadLocation := filepath.Join(downloadDir, name)
	err = downloadFile(downloadLocation, url)
	if err != nil {
		return "", err
	}

	var binary string
	if cfg.Compressed {
		extractDir, err := ioutil.TempDir(currentWorkingDir, EXTRACT_PREFIX)
		if err != nil {
			return "", err
		}

		err = archiver.Unarchive(downloadLocation, extractDir)
		if err != nil {
			return "", err
		}
		binary = filepath.Join(extractDir, cfg.BinaryPath)
	} else {
		binary = downloadLocation
	}

	return binary, err
}

// Combine amd64 and arm64 binary into a macOS universal binary
func (cfg *Configuration) Combine(amd64Path string, arm64Path string, amd64AssetName string) (string, error) {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir, err := ioutil.TempDir(currentWorkingDir, FAT_PREFIX)
	if err != nil {
		return "", err
	}

	fileName := filepath.Base(amd64Path)
	target := filepath.Join(dir, fileName)
	universalAssetName := generateUniversalAssetName(amd64AssetName, cfg.UniversalIdentifer)

	err = MakeFatBinary(amd64Path, arm64Path, target)
	if err != nil {
		return "", err
	}

	var universalAssetPath string
	if cfg.Compressed {
		archiver.Archive([]string{target}, universalAssetName)
		universalAssetPath = filepath.Join(currentWorkingDir, universalAssetName)
	} else {
		renamedPath := filepath.Join(dir, universalAssetName)
		err := os.Rename(target, renamedPath)
		if err != nil {
			return "", err
		}
	}
	return universalAssetPath, nil
}

// Upload the univeral binary to Github as a release asset
func (cfg *Configuration) UploadAsset(release *github.RepositoryRelease, assetPath string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: cfg.GithubToken,
	})

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	assetNameToUpload := filepath.Base(assetPath)
	uploadOptions := github.UploadOptions{
		Name: assetNameToUpload,
	}

	file, err := os.Open(assetPath)
	if err != nil {
		return err
	}

	existingAssets := release.Assets
	for _, asset := range existingAssets {
		// If the asset to upload is already present in the release, delete the asset so it can be reuploaded
		if asset.GetName() == assetNameToUpload {
			if cfg.Overwrite {
				resp, err := client.Repositories.DeleteReleaseAsset(ctx, cfg.Owner, cfg.Repository, asset.GetID())
				if err != nil {
					return err
				}
				log.Info("Deleted pre-existing asset with same name")
				log.Debugf("Asset ID: %d", asset.GetID())
				log.Debugf("Asset Name: %s", asset.GetName())
				log.Debugf("Asset URL: %s", asset.GetURL())
				log.Debugf("Response Status: %s", resp.Status)
				log.Debugf("Response Status Code: %d", resp.StatusCode)
			} else {
				return errors.New("found existing asset with same name")
			}
		}
	}

	uploadedAsset, resp, err := client.Repositories.UploadReleaseAsset(ctx, cfg.Owner, cfg.Repository, *release.ID, &uploadOptions, file)
	if err != nil {
		log.Error("Could not upload asset!")
		return err
	}
	defer file.Close()

	log.Debug("Successfully uploaded asset")
	log.Debugf("Asset ID: %d", uploadedAsset.GetID())
	log.Debugf("Asset Name: %s", uploadedAsset.GetName())
	log.Debugf("Asset URL: %s", uploadedAsset.GetURL())
	log.Debugf("Response Status: %s", resp.Status)
	log.Debugf("Response Status Code: %d", resp.StatusCode)

	return nil
}
