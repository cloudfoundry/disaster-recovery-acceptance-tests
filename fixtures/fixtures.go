// Package fixtures embeds and provides helper methods
// for test apps. It makes it so this package is importable.
// Go "embed" does not allow you to import files from another module.
// It classifies the files in credhub-test-app as a separate module
// because it has go.mod and go.sum.
//
// Run `go generate` after changing any of the fixtures. Then commit the
// test_data directory to the repo.
//
package fixtures

import (
	"archive/tar"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	//go:embed test_data
	testData embed.FS
)

//go:generate tar cf test_data/test_app.tar test_app
//go:generate tar cf test_data/credhub-test-app.tar credhub-test-app

// WriteFixturesToTemporaryDirectory should be passed an already created temporary
// directory and a name for one of the test fixtures. See test_data for examples.
func WriteFixturesToTemporaryDirectory(tempDir, fixture string) (retErr error) {
	targetDir := filepath.Join(tempDir, fixture)

	mkdirErr := os.MkdirAll(targetDir, 0700)
	if mkdirErr != nil {
		return fmt.Errorf("could not make target directory for fixture %q: %w", fixture, mkdirErr)
	}

	root := filepath.Join("test_data", fixture+".tar")

	fixtureFile, testOpenErr := testData.Open(root)
	if testOpenErr != nil {
		return fmt.Errorf("could not open fixture %q: %w", fixture, testOpenErr)
	}
	defer closeHelper(fixtureFile, &retErr)

	fixtureReader := tar.NewReader(fixtureFile)

	for {
		header, nextErr := fixtureReader.Next()
		if nextErr != nil {
			if nextErr == io.EOF {
				break
			}
			return fmt.Errorf("error when reading fixture tarball %q: %w", fixture, nextErr)
		}

		targetPath := filepath.Join(targetDir, strings.TrimPrefix(header.Name, fixture+string(filepath.Separator)))

		if header.FileInfo().IsDir() {
			err := os.MkdirAll(targetPath, 0700)
			if err != nil {
				return err
			}
			continue
		}

		createErr := createAndWriteFile(targetPath, fixtureReader)
		if createErr != nil {
			return createErr
		}
	}

	return nil
}

func createAndWriteFile(dstPath string, src io.Reader) (retErr error) {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer closeHelper(dst, &retErr)

	_, copyErr := io.Copy(dst, src)
	if copyErr != nil {
		return copyErr
	}

	return nil
}

func closeHelper(c io.Closer, retErr *error) {
	closeErr := c.Close()
	if closeErr != nil && *retErr == nil {
		*retErr = closeErr
	}
}
