/*
 *  Copyright IBM Corporation 2021
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/konveyor/move2kube/common"
	"github.com/konveyor/move2kube/filesystem"
)

// CheckAndCopyCustomizations checks if the customizations path is an existing directory and copies to assets
func CheckAndCopyCustomizations(customizationsPath string) error {
	if customizationsPath == "" {
		return nil
	}
	customizationsPath, err := filepath.Abs(customizationsPath)
	if err != nil {
		return fmt.Errorf("failed to make the customizations directory path '%s' absolute. Error: %w", customizationsPath, err)
	}
	fi, err := os.Stat(customizationsPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("the given customizations directory '%s' does not exist. Error: %w", customizationsPath, err)
	}
	if err != nil {
		return fmt.Errorf("failed to stat the given customizations directory '%s' Error: %w", customizationsPath, err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("the given customizations path '%s' is a file. Expected a directory", customizationsPath)
	}
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get the current working directory. Error: %w", err)
	}
	if common.IsParent(pwd, customizationsPath) {
		return fmt.Errorf("the given customizations directory '%s' is a parent of the current working directory", customizationsPath)
	}
	if err = CopyCustomizationsAssetsData(customizationsPath); err != nil {
		return fmt.Errorf("failed to copy the customizations data from the directory '%s' . Error: %w", customizationsPath, err)
	}
	return nil
}

// CopyCustomizationsAssetsData copies an customizations to the assets directory
func CopyCustomizationsAssetsData(customizationsPath string) (err error) {
	if customizationsPath == "" {
		return nil
	}
	assetsPath, err := filepath.Abs(common.AssetsPath)
	if err != nil {
		return fmt.Errorf("failed to make the assets path '%s' absolute. Error: %w", assetsPath, err)
	}
	customizationsAssetsPath := filepath.Join(assetsPath, "custom")

	// Create the subdirectory and copy the assets into it.
	if err = os.MkdirAll(customizationsAssetsPath, common.DefaultDirectoryPermission); err != nil {
		return fmt.Errorf("failed to create the customization assets directory at path '%s' . Error: %w", customizationsAssetsPath, err)
	}
	if err = filesystem.Replicate(customizationsPath, customizationsAssetsPath); err != nil {
		return fmt.Errorf("failed to copy the customizations from '%s' to the directory '%s' . Error: %w", customizationsPath, customizationsAssetsPath, err)
	}
	return nil
}
